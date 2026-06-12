# Morpheus OPSEC Considerations

## Detection Surface

- Go binaries have distinct PE characteristics (Go runtime strings, .gopclntab section)
- Default WinHTTP user-agent may be signatured
- Static linking avoids DLL import table entropy (no kernel32.dll imports for common ops)

## Evasion Features

Morpheus uses no built-in evasion (ETW/AMSI patching, unhook, indirect syscalls) by default.
For evasion, pair with Atreus as a wrapper which provides:

- Early Bird APC injection
- Thread hijacking
- RC4/XOR encryption layers
- Hash-based API resolution
- ntdll unhook
- ETW/AMSI patching
- Sandbox detection

## Network OPSEC

- Use HTTPS with valid certificates where possible
- Randomize User-Agent strings
- Configure callback jitter to avoid beacon detection
- Avoid predictable callback intervals
