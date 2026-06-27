from .schemas import ReportRequest
from .timeline_preprocessor import build_timeline_context


def build_title_incident_context(data: ReportRequest) -> str:
    incident = data.incident

    parts = [
        f"Incident ID: {incident.id}",
        f"Created at: {incident.createdAt}",
    ]

    if incident.closedAt:
        parts.append(f"Closed at: {incident.closedAt}")

    if incident.severity:
        parts.append(f"Severity: {incident.severity}")

    if incident.status:
        parts.append(f"Status: {incident.status}")

    return "\n".join(parts)


def build_ai_title_prompt(data: ReportRequest) -> str:
    incident_context = build_title_incident_context(data)
    timeline_context = build_timeline_context(data.timeline)

    return f"""
You are an incident management assistant.

Generate a NEW short technical title for this incident.

Important rules:
- Do not copy the original incident title.
- The original title is intentionally not provided.
- Use the timeline as the main source of truth.
- Use only facts from the incident metadata and timeline.
- The title must be in English.
- No more than 8 words.
- Do not use quotation marks.
- Do not put a period at the end.
- Prefer a specific technical formulation.
- Return only the final title.

Good examples:
- Missing S3 bucket configuration
- Report service PDF generation failure
- High 5xx errors in report service
- S3 upload configuration failure
- Telegram report download outage

Bad examples:
- problem_4
- Issue in production
- Summary of the incident
- Something went wrong
- There is a problem

Incident metadata:
{incident_context}

Timeline:
{timeline_context}

Generate only the new incident title.
""".strip()