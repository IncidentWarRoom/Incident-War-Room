# Incident War Room

Incident War Room is a Telegram-native platform for managing production incidents and preserving incident knowledge in a structured form.

The system allows engineering teams to create incidents directly from Telegram group chats, record important investigation updates, track incident severity, maintain a chronological timeline of events, and automatically generate post-incident reports.

The platform consists of two core services:

* **Incident Service (Go)** — handles Telegram integration, incident management, timeline tracking, data persistence, and report generation requests.
* **Report Service (Python)** — generates PDF incident reports based on collected incident data.

## Key Features

* Incident creation and management directly from Telegram
* Timeline-based incident tracking
* Severity management (Low, Medium, High)
* Automatic participant collection
* Incident report generation in PDF format
* PostgreSQL-based persistence
* Containerized deployment using Docker

## Example Workflow

1. Create an incident using `/incident create`.
2. Add important investigation updates using `/incident <message>`.
3. View the current timeline with `/timeline`.
4. Adjust incident severity when required.
5. Close the incident using `/incident close`.
6. Receive an automatically generated PDF report in Telegram.

The project is designed as an MVP for engineering teams that already use Telegram as their primary communication platform and need a lightweight incident management solution without introducing additional tools into their workflow.
