from mythic_container.C2ProfileBuilder import *

class HTTP(C2Profile):
    name = "http"
    description = "HTTP/HTTPS C2 profile for Morpheus"
    author = "@jeremiznoo (https://github.com/Jeremiznoo/Morpheus/)"
    is_p2p = False
    is_server_routed = False
    mythic_encrypts = True
    parameters = [
        C2ProfileParameter(
            name="callback_host",
            description="Callback host (IP or domain)",
            default_value="https://127.0.0.1:443",
            parameter_type=ParameterType.String,
            required=True,
        ),
        C2ProfileParameter(
            name="headers",
            description="Custom HTTP headers (JSON)",
            default_value='{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"}',
            parameter_type=ParameterType.Dictionary,
            required=False,
        ),
        C2ProfileParameter(
            name="user_agent",
            description="HTTP User-Agent header",
            default_value="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
            parameter_type=ParameterType.String,
            required=False,
        ),
        C2ProfileParameter(
            name="callback_interval",
            description="Callback interval in seconds",
            default_value=5,
            parameter_type=ParameterType.Number,
            required=False,
        ),
        C2ProfileParameter(
            name="callback_jitter",
            description="Callback jitter percentage",
            default_value=10,
            parameter_type=ParameterType.Number,
            required=False,
        ),
        C2ProfileParameter(
            name="killdate",
            description="Agent kill date (YYYY-MM-DD)",
            default_value="",
            parameter_type=ParameterType.String,
            required=False,
        ),
    ]
