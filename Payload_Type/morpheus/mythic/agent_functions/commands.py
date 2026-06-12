from mythic_container.MythicCommandBase import *
from mythic_container.MythicRPC import *


class ShellCommand(CommandBase):
    cmd = "shell"
    description = "Execute a command via cmd.exe /c"
    version = 1
    supported_ui_features = []
    parameters = [
        CommandParameter(
            name="command",
            type=ParameterType.String,
            description="Command to execute",
            parameter_group_info=[ParameterGroupInfo(ui_position=1)],
        )
    ]

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class RunCommand(CommandBase):
    cmd = "run"
    description = "Execute a binary with arguments"
    version = 1
    parameters = [
        CommandParameter(name="path", type=ParameterType.String, description="Binary path"),
        CommandParameter(name="args", type=ParameterType.Array, description="Arguments"),
    ]

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class CdCommand(CommandBase):
    cmd = "cd"
    description = "Change working directory"
    version = 1
    parameters = [
        CommandParameter(name="path", type=ParameterType.String, description="Directory path"),
    ]

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class PwdCommand(CommandBase):
    cmd = "pwd"
    description = "Print working directory"
    version = 1
    parameters = []

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class LsCommand(CommandBase):
    cmd = "ls"
    description = "List directory contents"
    version = 1
    parameters = [
        CommandParameter(name="path", type=ParameterType.String, description="Directory path", default_value="."),
    ]

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class CatCommand(CommandBase):
    cmd = "cat"
    description = "Read a file"
    version = 1
    parameters = [
        CommandParameter(name="path", type=ParameterType.String, description="File path"),
    ]

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class CpCommand(CommandBase):
    cmd = "cp"
    description = "Copy a file"
    version = 1
    parameters = [
        CommandParameter(name="source", type=ParameterType.String, description="Source path"),
        CommandParameter(name="dest", type=ParameterType.String, description="Destination path"),
    ]

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class MvCommand(CommandBase):
    cmd = "mv"
    description = "Move or rename a file"
    version = 1
    parameters = [
        CommandParameter(name="source", type=ParameterType.String, description="Source path"),
        CommandParameter(name="dest", type=ParameterType.String, description="Destination path"),
    ]

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class RmCommand(CommandBase):
    cmd = "rm"
    description = "Delete a file"
    version = 1
    parameters = [
        CommandParameter(name="path", type=ParameterType.String, description="File path"),
    ]

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class MkdirCommand(CommandBase):
    cmd = "mkdir"
    description = "Create a directory"
    version = 1
    parameters = [
        CommandParameter(name="path", type=ParameterType.String, description="Directory path"),
    ]

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class DownloadCommand(CommandBase):
    cmd = "download"
    description = "Download a file from the target"
    version = 1
    parameters = [
        CommandParameter(name="path", type=ParameterType.String, description="File path"),
    ]

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class UploadCommand(CommandBase):
    cmd = "upload"
    description = "Upload a file to the target"
    version = 1
    parameters = [
        CommandParameter(name="path", type=ParameterType.String, description="Destination path"),
        CommandParameter(name="content", type=ParameterType.String, description="File content (base64)"),
    ]

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class PsCommand(CommandBase):
    cmd = "ps"
    description = "List running processes"
    version = 1
    parameters = []

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class KillCommand(CommandBase):
    cmd = "kill"
    description = "Kill a process by PID"
    version = 1
    parameters = [
        CommandParameter(name="pid", type=ParameterType.Number, description="Process ID"),
    ]

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class GetuidCommand(CommandBase):
    cmd = "getuid"
    description = "Get current username"
    version = 1
    parameters = []

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class WhoamiCommand(CommandBase):
    cmd = "whoami"
    description = "Get user, domain and privileges"
    version = 1
    parameters = [
        CommandParameter(name="privileges", type=ParameterType.Boolean, description="Show privileges"),
    ]

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class SleepCommand(CommandBase):
    cmd = "sleep"
    description = "Get or set callback interval/jitter"
    version = 1
    parameters = [
        CommandParameter(name="interval", type=ParameterType.Number, description="Callback interval (seconds)"),
        CommandParameter(name="jitter", type=ParameterType.Number, description="Jitter percentage"),
    ]

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class IfconfigCommand(CommandBase):
    cmd = "ifconfig"
    description = "List network interfaces"
    version = 1
    parameters = []

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class MakeTokenCommand(CommandBase):
    cmd = "make_token"
    description = "Create impersonation token"
    version = 1
    parameters = [
        CommandParameter(name="username", type=ParameterType.String, description="Username"),
        CommandParameter(name="domain", type=ParameterType.String, description="Domain"),
        CommandParameter(name="password", type=ParameterType.String, description="Password"),
    ]

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class StealTokenCommand(CommandBase):
    cmd = "steal_token"
    description = "Steal token from another process"
    version = 1
    parameters = [
        CommandParameter(name="pid", type=ParameterType.Number, description="Process ID"),
    ]

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class Rev2selfCommand(CommandBase):
    cmd = "rev2self"
    description = "Revert to original token"
    version = 1
    parameters = []

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class RunasCommand(CommandBase):
    cmd = "runas"
    description = "Run a command as another user"
    version = 1
    parameters = [
        CommandParameter(name="username", type=ParameterType.String, description="Username"),
        CommandParameter(name="domain", type=ParameterType.String, description="Domain"),
        CommandParameter(name="password", type=ParameterType.String, description="Password"),
        CommandParameter(name="command", type=ParameterType.String, description="Command to execute"),
    ]

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class SpawnCommand(CommandBase):
    cmd = "spawn"
    description = "Spawn shellcode in a process"
    version = 1
    parameters = [
        CommandParameter(name="pid", type=ParameterType.Number, description="Target PID (0 for current)", default_value=0),
        CommandParameter(name="shellcode", type=ParameterType.String, description="Shellcode (base64)"),
    ]

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class SpawntoCommand(CommandBase):
    cmd = "spawnto"
    description = "Set sacrificial process for spawn"
    version = 1
    parameters = [
        CommandParameter(name="path", type=ParameterType.String, description="Process path"),
    ]

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class ExecuteAssemblyCommand(CommandBase):
    cmd = "execute_assembly"
    description = "Execute a .NET assembly"
    version = 1
    parameters = [
        CommandParameter(name="assembly_id", type=ParameterType.String, description="Assembly ID from Mythic"),
        CommandParameter(name="arguments", type=ParameterType.String, description="CLI arguments"),
    ]

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class BlockDllsCommand(CommandBase):
    cmd = "blockdlls"
    description = "Enable/disable BlockDLLs policy"
    version = 1
    parameters = [
        CommandParameter(name="action", type=ParameterType.String, description="on/off", default_value="on"),
    ]

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class SocksCommand(CommandBase):
    cmd = "socks"
    description = "SOCKS5 proxy management"
    version = 1
    parameters = [
        CommandParameter(name="action", type=ParameterType.String, description="start/stop", default_value="start"),
        CommandParameter(name="port", type=ParameterType.Number, description="Local SOCKS port", default_value=1080),
    ]

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class RportfwdCommand(CommandBase):
    cmd = "rportfwd"
    description = "Reverse port forward"
    version = 1
    parameters = [
        CommandParameter(name="action", type=ParameterType.String, description="start/stop"),
        CommandParameter(name="local_port", type=ParameterType.Number, description="Local bind port"),
        CommandParameter(name="remote_port", type=ParameterType.Number, description="Remote forward port"),
    ]

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class LigoloStartCommand(CommandBase):
    cmd = "ligolo_start"
    description = "Start Ligolo-ng tunnel"
    version = 1
    parameters = [
        CommandParameter(name="tunnel", type=ParameterType.String, description="Tunnel name", default_value="default"),
        CommandParameter(name="listener", type=ParameterType.String, description="Listener address", default_value="0.0.0.0:11601"),
    ]

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class LigoloStopCommand(CommandBase):
    cmd = "ligolo_stop"
    description = "Stop Ligolo-ng tunnel"
    version = 1
    parameters = []

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class LigoloStatusCommand(CommandBase):
    cmd = "ligolo_status"
    description = "List active Ligolo-ng tunnels"
    version = 1
    parameters = []

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass


class ExitCommand(CommandBase):
    cmd = "exit"
    description = "Terminate the agent"
    version = 1
    parameters = []

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        return task

    async def process_response(self, response: AgentResponse):
        pass
