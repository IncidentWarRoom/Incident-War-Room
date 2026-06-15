# Edge Cases

## EC-1 — Duplicate Incident Creation

| Field | Description |
|-------|-------------|
| **Situation** | `/incident create` when an active incident already exists in the chat |
| **Behavior** | ⚠️ Active incident #15 already exists.<br>Use `/incident close` to close it first. |
| **Why It Matters** | Without this check, the database would have two active incidents with the same `chat_id` — the timeline would break |

---

## EC-2 — Commands Without Active Incident

| Field | Description |
|-------|-------------|
| **Situation** | `/incident close`, `/timeline`, or `/incident <msg>` when no active incident exists |
| **Behavior** | ⚠️ No active incident found in this chat.<br>Use `/incident create <title>` to start one. |
| **Why It Matters** | Without validation, the Go service would receive `null` and crash |

---

## EC-3 — Report Service Unavailable

| Field | Description |
|-------|-------------|
| **Situation** | Python service is not responding during `/incident close` |
| **Behavior** | Incident is still closed in the database.<br>Bot responds: ✅ Incident closed. ⚠️ Report generation failed.<br>Try `/report` later (if implemented) |
| **Why It Matters** | Report service unavailability should not block incident closure |

---

## EC-4 — Empty Timeline

| Field | Description |
|-------|-------------|
| **Situation** | `/incident create` → immediately `/incident close` without adding any events |
| **Behavior** | Bot closes the incident.<br>PDF contains only system events (created, closed).<br>Duration: X min, Participants: 1, Timeline events: 2 |
| **Why It Matters** | Zero timeline should not break PDF generation |

---

## EC-5 — Very Long Timeline Message

| Field | Description |
|-------|-------------|
| **Situation** | `/incident <text longer than Telegram limit>` |
| **Behavior** | Bot truncates to 500 characters and adds `[truncated]`.<br>Saves full text to the database. |
| **Why It Matters** | Formatting timeline in a single message would hit Telegram's limit (4096 characters) if there are many events with long texts |