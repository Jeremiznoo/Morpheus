import os
import asyncio
import subprocess
import pathlib
from mythic_container.PayloadBuilder import *
from mythic_container.MythicCommandBase import *


class Morpheus(PayloadType):
    name = "morpheus"
    file_extension = "exe"
    author = "@jeremiznoo (https://github.com/Jeremiznoo/Morpheus/)"
    supported_os = [SupportedOS.Windows]
    wrapper = False
    wrapped_payloads = []
    note = "Go-based Windows C2 agent for Mythic"
    supports_dynamic_loading = False
    build_parameters = [
        BuildParameter(
            name="c2_url",
            parameter_type=BuildParameterType.String,
            description="C2 callback URL (e.g. https://192.168.1.100:443)",
            default_value="https://127.0.0.1:443",
            required=True,
        ),
        BuildParameter(
            name="callback_interval",
            parameter_type=BuildParameterType.Number,
            description="Callback interval in seconds",
            default_value=5,
            required=True,
        ),
        BuildParameter(
            name="callback_jitter",
            parameter_type=BuildParameterType.Number,
            description="Callback jitter percentage (0-100)",
            default_value=10,
            required=True,
        ),
        BuildParameter(
            name="agent_uuid",
            parameter_type=BuildParameterType.String,
            description="Agent UUID (auto-generated if empty)",
            default_value="",
            required=False,
        ),
        BuildParameter(
            name="encryption_key",
            parameter_type=BuildParameterType.String,
            description="AES-256-GCM key (base64, auto-generated if empty)",
            default_value="",
            required=False,
        ),
        BuildParameter(
            name="debug",
            parameter_type=BuildParameterType.Boolean,
            description="Enable debug output",
            default_value=False,
            required=False,
        ),
    ]
    c2_profiles = ["http"]
    agent_path = pathlib.Path(".")
    agent_code_path = agent_path / "agent_code"
    agent_icon_path = pathlib.Path(".") / "mythic" / "agent_functions" / "morpheus.svg"
    build_steps = [
        BuildStep(step_name="Morpheus Build", step_description="Cross-compile Go agent for Windows x64"),
    ]

    async def build(self) -> BuildResponse:
        resp = BuildResponse(status=BuildStatus.Error)

        go_source = self.agent_code_path.resolve()
        if not (go_source / "go.mod").exists():
            resp.build_message = f"Go source not found at {go_source}"
            return resp

        c2_url = self.get_parameter("c2_url") or "https://127.0.0.1:443"
        interval = str(self.get_parameter("callback_interval") or 5)
        jitter = str(self.get_parameter("callback_jitter") or 10)
        uuid = self.get_parameter("agent_uuid") or ""
        enc_key = self.get_parameter("encryption_key") or ""

        ldflags = (
            f"-X main.C2Url={c2_url} "
            f"-X main.CallbackInterval={interval} "
            f"-X main.JitterStr={jitter}"
        )
        if uuid:
            ldflags += f" -X main.AgentUUID={uuid}"
        if enc_key:
            ldflags += f" -X main.EncKey={enc_key}"
        ldflags += " -s -w"

        env = os.environ.copy()
        env["GOOS"] = "windows"
        env["GOARCH"] = "amd64"
        env["CGO_ENABLED"] = "0"

        output_path = go_source / "morpheus.exe"
        cmd = ["go", "build", "-trimpath", f"-ldflags={ldflags}", "-o", str(output_path), "."]

        proc = await asyncio.create_subprocess_exec(
            *cmd,
            cwd=go_source,
            env=env,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        )

        stdout, stderr = await proc.communicate()

        if proc.returncode != 0:
            resp.build_stderr = stderr.decode() if stderr else "build failed"
            resp.build_message = "Go build failed"
            return resp

        resp.payload = output_path.read_bytes()

        try:
            output_path.unlink(missing_ok=True)
        except OSError:
            pass

        resp.status = BuildStatus.Success
        resp.build_message = "Build successful"
        resp.build_stdout = stdout.decode() if stdout else "build complete"
        resp.build_stderr = stderr.decode() if stderr else ""
        return resp
