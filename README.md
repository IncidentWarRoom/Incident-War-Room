# Incident War Room

Incident War Room is a self-hosted incident management platform that integrates with Telegram and monitoring systems to simplify production incident response.

The platform automatically creates dedicated Telegram Topics for incidents, records investigation updates as a structured timeline, generates AI-powered incident reports, and stores reports in S3-compatible object storage. It can also integrate with Prometheus Alertmanager to automatically create incidents when monitoring detects a failure.

The platform consists of two core services:

* **Incident Service (Go)** — handles Telegram integration, incident lifecycle management, monitoring integrations, timeline tracking, data persistence, and communication with external services.
* **Report Service (Python)** — generates AI-enhanced PDF reports, integrates with object storage, and supports both S3 and inline report delivery modes.

## Key Features

* Telegram Topics for isolated incident discussions
* Automatic incident creation from Prometheus Alertmanager
* Manual incident creation directly from Telegram
* Chronological incident timeline
* AI-generated incident title and summary
* PDF report generation
* S3-compatible report storage with automatic fallback
* Incident history and web dashboard
* PostgreSQL-based persistence
* Docker-based deployment

## Example Workflow

1. An incident is created manually from Telegram or automatically by Alertmanager.
2. A dedicated Telegram Topic is created for the incident.
3. Engineers discuss the incident inside the Topic while the system records the timeline.
4. AI generates an incident title and summary.
5. The incident is resolved and the Topic is automatically removed.
6. A PDF report is generated and delivered to Telegram.
7. If S3 is configured, the report is stored in object storage; otherwise, it is delivered directly by the Report Service.
8. The incident remains available through the web dashboard for future investigation and analysis.

The platform is designed for engineering teams that already use Telegram as their primary communication tool and want a lightweight, self-hosted incident management solution that can be deployed in their own infrastructure and integrated with existing monitoring systems.
