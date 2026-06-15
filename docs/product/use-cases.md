# Use Cases

## UC-1 — Create Incident

| Field | Description |
|-------|-------------|
| **Actor** | On-call Engineer |
| **Precondition** | Bot is added to chat, no active incident exists |
| **Trigger** | `/incident create Payment Service Down` |
| **Main Flow** | 1. Engineer enters command → 2. Bot parses title → 3. Creates record in `incidents` with status ACTIVE → 4. Creates INCIDENT_CREATED event in `incident_events` → 5. Sends card with buttons |
| **Alternative Flow** | Active incident already exists → bot responds ⚠️ Active incident already exists |
| **Postcondition** | Incident created, team notified |

---

## UC-2 — Add Event to Timeline

| Field | Description |
|-------|-------------|
| **Actor** | Any chat participant |
| **Precondition** | Active incident exists in this chat |
| **Trigger** | `/incident Found connection pool errors` |
| **Main Flow** | 1. Participant enters command → 2. Bot identifies active incident by `chat_id` → 3. Creates record in `incident_events` with type COMMENT_ADDED → 4. Confirms addition |
| **Alternative Flow** | No active incident → ⚠️ No active incident. Use `/incident create` first. |
| **Postcondition** | Event saved and will appear in timeline and report |

---

## UC-3 — View Timeline

| Field | Description |
|-------|-------------|
| **Actor** | Any chat participant |
| **Precondition** | Active incident exists |
| **Trigger** | `/timeline` or button "Show Timeline" |
| **Main Flow** | 1. Bot retrieves all events by `incident_id` → 2. Sorts by `created_at` → 3. Formats and sends message |
| **Alternative Flow** | No events → 📋 No events yet for this incident. |
| **Postcondition** | Team sees chronology |

---

## UC-4 — Change Severity

| Field | Description |
|-------|-------------|
| **Actor** | On-call Engineer |
| **Precondition** | Active incident exists |
| **Trigger** | Button "Change Severity" |
| **Main Flow** | 1. Engineer clicks button → 2. Bot shows popup with options → 3. Engineer selects → 4. Bot updates severity in `incidents` → 5. Logs SEVERITY_CHANGED event → 6. Updates card |
| **Postcondition** | Severity updated, event in timeline |

---

## UC-5 — Close Incident

| Field | Description |
|-------|-------------|
| **Actor** | On-call Engineer |
| **Precondition** | Active incident exists |
| **Trigger** | `/incident close` or button "Close Incident" |
| **Main Flow** | 1. Bot records `closed_at` → 2. Changes status to CLOSED → 3. Creates INCIDENT_CLOSED event → 4. Calculates duration → 5. Builds DTO → 6. POST `/api/v1/reports/generate` → 7. Receives PDF → 8. Sends to chat |
| **Alternative Flow** | Report Service unavailable → bot responds ⚠️ Report generation failed. Incident is closed. |
| **Postcondition** | Incident closed, PDF sent to chat |