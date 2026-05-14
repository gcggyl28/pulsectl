# pulsectl

A minimal health-check orchestrator that polls endpoints and reports degraded services via webhook.

---

## Installation

```bash
go install github.com/yourusername/pulsectl@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/pulsectl.git && cd pulsectl && go build -o pulsectl .
```

---

## Usage

Define your endpoints in a `config.yaml` file:

```yaml
interval: 30s
webhook: "https://hooks.example.com/notify"

endpoints:
  - name: api-service
    url: "https://api.example.com/health"
    timeout: 5s
  - name: auth-service
    url: "https://auth.example.com/ping"
    timeout: 3s
```

Then run:

```bash
pulsectl --config config.yaml
```

pulsectl will poll each endpoint at the defined interval. If a service returns a non-2xx status or times out, a POST request is sent to the configured webhook with a JSON payload describing the degraded service.

**Example webhook payload:**

```json
{
  "service": "auth-service",
  "url": "https://auth.example.com/ping",
  "status": "degraded",
  "reason": "timeout after 3s",
  "timestamp": "2024-05-10T14:32:00Z"
}
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `config.yaml` | Path to configuration file |
| `--verbose` | `false` | Enable verbose logging |

---

## License

MIT © 2024 yourusername