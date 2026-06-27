from .schemas import ReportRequest
from .timeline_preprocessor import build_timeline_context


def safe_value(value: str | int | None, fallback: str = "Unknown") -> str:
    if value is None:
        return fallback

    value_as_text = str(value).strip()

    if not value_as_text:
        return fallback

    return value_as_text


def build_incident_context(data: ReportRequest) -> str:
    incident = data.incident

    return "\n".join(
        [
            f"Incident ID: {safe_value(incident.id)}",
            f"Title: {safe_value(incident.title)}",
            f"Severity: {safe_value(incident.severity)}",
            f"Status: {safe_value(incident.status)}",
            f"Created at: {safe_value(incident.createdAt)}",
            f"Closed at: {safe_value(incident.closedAt, 'Not closed or not provided')}",
        ]
    )


def build_participants_context(data: ReportRequest) -> str:
    if not data.participants:
        return "No participants provided."

    lines = []

    for participant in data.participants:
        lines.append(
            f"- {participant.username} (Telegram user ID: {participant.userId})"
        )

    return "\n".join(lines)


def build_ai_summary_prompt(data: ReportRequest) -> str:
    incident_context = build_incident_context(data)
    participants_context = build_participants_context(data)
    timeline_context = build_timeline_context(data.timeline)

    return f"""
You are an incident report assistant.

You will receive structured information about one completed incident.
Your task is to generate a clear, detailed, and useful incident summary.

Important rules:
- Do not invent facts.
- Use only the provided incident data, participants, and timeline.
- If some information is missing or unclear, explicitly say that it is unclear.
- Ignore obvious noise messages, but preserve meaningful confirmations.
- If a user accepted an assignment, mention that they took responsibility.
- Write in a professional but simple style.
- The answer must be in English.

The summary should include:
1. Brief overview of what happened
2. Incident severity and status
3. Who participated
4. Who accepted responsibility or took tasks
5. Important investigation steps
6. Actions taken to resolve the incident
7. Probable root cause, if it is clear from the timeline
8. Final resolution
9. Recommendations or possible improvements

Incident:
{incident_context}

Participants:
{participants_context}

Timeline:
{timeline_context}

Generate the final incident summary.
""".strip()