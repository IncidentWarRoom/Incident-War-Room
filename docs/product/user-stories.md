# User Stories

## US-1 — Create Incident

**As a** on-call engineer,  
**I want to** create an incident using the `/incident create Payment Service Down` command directly in the working Telegram chat,  
**So that** I can immediately log the issue and notify the team without switching to other tools.

---

## US-2 — Add Events to Timeline

**As a** incident participant,  
**I want to** add messages using the `/incident Started investigating DB` command during the investigation,  
**So that** I can document progress and my actions are visible to the entire team in chronological order.

---

## US-3 — View Timeline

**As a** team lead,  
**I want to** see the chronology of all events using the `/timeline` command,  
**So that** I can understand the current incident picture without reading the entire chat from the beginning.

---

## US-4 — Change Severity

**As a** on-call engineer,  
**I want to** change incident severity via buttons (LOW / MEDIUM / HIGH),  
**So that** the team understands the real scale of the problem as we investigate further.

---

## US-5 — Close Incident and Get Report

**As a** on-call engineer,  
**I want to** close the incident using the `/incident close` command and automatically receive a PDF report in the chat,  
**So that** I don't waste time manually writing a post-mortem immediately after resolving the issue.