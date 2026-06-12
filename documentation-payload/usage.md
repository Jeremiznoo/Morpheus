# Morpheus Usage

## Basic Commands

| Command | Description |
|---------|-------------|
| `shell` | Execute command via cmd.exe /c |
| `run` | Execute a binary with arguments |
| `cd` | Change working directory |
| `pwd` | Print working directory |
| `ls` | List directory contents |
| `cat` | Read a file |
| `cp` | Copy a file |
| `mv` | Move/rename a file |
| `rm` | Delete a file |
| `mkdir` | Create a directory |
| `download` | Download a file from target |
| `upload` | Upload a file to target |
| `ps` | List processes |
| `kill` | Kill a process by PID |
| `getuid` | Get current username |
| `whoami` | Get user/domain/privileges |
| `sleep` | Get/set callback interval/jitter |
| `ifconfig` | List network interfaces |

## Advanced Commands

| Command | Description |
|---------|-------------|
| `make_token` | Create impersonation token |
| `steal_token` | Steal token from another process |
| `rev2self` | Revert to original token |
| `runas` | Run as another user |
| `spawn` | Inject shellcode into a process |
| `spawnto` | Set sacrificial process path |
| `execute_assembly` | Run .NET assembly |
| `blockdlls` | Toggle BlockDLLs policy |
| `socks` | SOCKS5 proxy start/stop |
| `rportfwd` | Reverse port forward |
| `ligolo_start` | Start Ligolo-ng tunnel |
| `ligolo_stop` | Stop Ligolo-ng tunnel |
| `ligolo_status` | List Ligolo-ng tunnels |
| `exit` | Terminate agent |

## Usage with Atreus

1. Build Morpheus with `output_format = shellcode`
2. Select Atreus as the wrapper
3. Configure injection technique and evasion
4. Deploy the resulting executable

## Building from source outside Mythic

```bash
cd agent/
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build \
  -trimpath \
  -ldflags="-X main.C2Url=https://your-server:443 -X main.CallbackInterval=5 -s -w" \
  -o morpheus.exe .
```
