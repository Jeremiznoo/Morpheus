# Morpheus Installation

## Prerequisites

- Mythic 3.x
- Go 1.22+ (for compilation)

## Install via Mythic CLI

```bash
sudo ./mythic-cli install folder /path/to/Morpheus
```

## Build Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `c2_url` | String | `https://127.0.0.1:443` | C2 callback URL |
| `callback_interval` | Number | 5 | Seconds between callbacks |
| `callback_jitter` | Number | 10 | Jitter percentage (0-100) |
| `agent_uuid` | String | auto | Agent UUID |
| `encryption_key` | String | auto | AES-256-GCM key (base64) |
| `debug` | Bool | false | Enable debug output |
