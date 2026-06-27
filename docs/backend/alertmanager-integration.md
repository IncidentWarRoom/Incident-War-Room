# Alertmanager Integration (`internal/alert`)

The incident-service can open incidents automatically from an existing
Prometheus / Alertmanager deployment. Alertmanager posts its webhook
notifications to the service; every **firing** alert opens a new incident
(topic + start of the timeline) in the configured Telegram forum supergroup.

The service consumes alerts — it does **not** scrape Prometheus or evaluate
rules. Alerting rules and grouping stay on your Prometheus/Alertmanager side.

## Flow

```
Prometheus ──fires rule──▶ Alertmanager ──webhook (HTTP POST)──▶ incident-service
                                                                      │
                                                          for each alert with
                                                          status == "firing"
                                                                      ▼
                                                   OpenIncidentFromAlert(title, severity)
                                                                      ▼
                                            new Telegram topic + incident in ALERT_CHAT_ID
```

`resolved` alerts are ignored (only `firing` opens an incident). Closing an
incident stays a manual action in the bot.

## Endpoint

| | |
|---|---|
| Method / path | `POST /webhooks/alertmanager` |
| Host / port | the HTTP API address, `HTTP_ADDR` (default `:8080`) |
| Auth | optional `Authorization: Bearer <token>` (see `ALERTMANAGER_WEBHOOK_TOKEN`) |
| Body | Alertmanager webhook payload, schema `version: 4` |

Responses:

| Status | Meaning |
|---|---|
| `200 OK` | payload accepted (returned even if some alerts failed to open — failures are logged, not surfaced) |
| `400 Bad Request` | body is not valid JSON |
| `401 Unauthorized` | a token is configured and the `Authorization` header is missing/wrong |

## Configuration

| Env var | Required | Description |
|---|---|---|
| `ALERT_CHAT_ID` | **yes** | ID of the Telegram **forum supergroup** where alert incidents are created. Without it `OpenIncidentFromAlert` is rejected (`KindUnavailable`, logged) and no incident is opened. |
| `ALERTMANAGER_WEBHOOK_TOKEN` | no | Shared secret. When set, requests must carry `Authorization: Bearer <token>`; when empty, the endpoint is unauthenticated. |
| `HTTP_ADDR` | no | Listen address of the HTTP API (default `:8080`). The webhook is mounted on the same server. |

`ALERT_CHAT_ID` must be a forum supergroup where the bot is an admin with the
**Manage Topics** right (each incident becomes its own topic).

### `.env` example

```dotenv
# Telegram forum supergroup that receives alert-driven incidents
ALERT_CHAT_ID=-1001234567890
# Shared secret for the webhook (optional but recommended)
ALERTMANAGER_WEBHOOK_TOKEN=change-me-to-a-long-random-string
```

## Severity mapping

The incoming `severity` label (case-insensitive) is mapped onto the incident
severity:

| Alert `severity` label | Incident severity |
|---|---|
| `critical`, `page`, `emergency`, `high` | `HIGH` |
| `info`, `information`, `low` | `LOW` |
| anything else / missing | `MEDIUM` |

## Incident title

The title is resolved from the alert in this order (first non-empty wins):

1. `annotations.summary`
2. `annotations.title`
3. `annotations.description`
4. `labels.alertname`
5. fallback: `"Prometheus alert"`

## Alertmanager configuration example

Add a webhook receiver and route firing alerts to it. The Bearer token is
passed via `http_config.authorization`.

```yaml
# alertmanager.yml
route:
  receiver: incident-war-room
  group_by: ['alertname']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 4h

receivers:
  - name: incident-war-room
    webhook_configs:
      - url: 'http://incident-service:8080/webhooks/alertmanager'
        send_resolved: false
        http_config:
          authorization:
            type: Bearer
            credentials: 'change-me-to-a-long-random-string'
```

Notes:

- `credentials` must match `ALERTMANAGER_WEBHOOK_TOKEN`. Drop the whole
  `http_config` block if you left the token empty.
- `send_resolved: false` — the service ignores `resolved` alerts anyway, so
  this just avoids useless traffic.
- Use the address at which Alertmanager can reach the service. Inside the same
  `docker compose` network that is the service name (e.g.
  `http://incident-service:8080`); from outside use the published host/port.

### Example alerting rule (Prometheus)

`severity` and `summary` drive the incident's severity and title:

```yaml
# prometheus rules
groups:
  - name: example
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High 5xx error rate on {{ $labels.service }}"
```

## Smoke test (curl)

Simulate an Alertmanager firing notification:

```bash
curl -i -X POST http://localhost:8080/webhooks/alertmanager \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer change-me-to-a-long-random-string' \
  -d '{
    "version": "4",
    "status": "firing",
    "alerts": [
      {
        "status": "firing",
        "labels": { "alertname": "HighErrorRate", "severity": "critical" },
        "annotations": { "summary": "High 5xx error rate on checkout" }
      }
    ]
  }'
```

Expected: `200 OK`, and a new incident topic opens in `ALERT_CHAT_ID` with a
`HIGH` severity titled "High 5xx error rate on checkout".
