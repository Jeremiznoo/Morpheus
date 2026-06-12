# Morpheus Wrapper Documentation

Morpheus supports Atreus as a wrapper for shellcode injection.

## Wrapper Usage

1. Build Morpheus with `output_format = shellcode`
2. Select Atreus as the wrapper payload
3. Configure Atreus injection technique and evasion
4. Deploy

## Supported Wrappers

| Wrapper | Description |
|---------|-------------|
| Atreus | Shellcode loader with APC/thread hijacking, ETW/AMSI patching, sandbox detection |
