from io import BytesIO

from reportlab.lib.pagesizes import A4
from reportlab.platypus import SimpleDocTemplate, Paragraph, Spacer
from reportlab.lib.styles import getSampleStyleSheet

from .schemas import ReportRequest


def generate_incident_report(data: ReportRequest) -> bytes:
    buffer = BytesIO()

    document = SimpleDocTemplate(buffer, pagesize=A4)
    styles = getSampleStyleSheet()
    elements = []

    elements.append(Paragraph("Incident Report", styles["Title"]))
    elements.append(Spacer(1, 12))

    elements.append(Paragraph("Incident Metadata", styles["Heading2"]))
    elements.append(Paragraph(f"ID: {data.incident.id}", styles["Normal"]))
    elements.append(Paragraph(f"Title: {data.incident.title}", styles["Normal"]))
    elements.append(Paragraph(f"Created at: {data.incident.createdAt}", styles["Normal"]))

    if data.incident.closedAt:
        elements.append(Paragraph(f"Closed at: {data.incident.closedAt}", styles["Normal"]))

    if data.incident.severity:
        elements.append(Paragraph(f"Severity: {data.incident.severity}", styles["Normal"]))

    if data.incident.status:
        elements.append(Paragraph(f"Status: {data.incident.status}", styles["Normal"]))

    elements.append(Spacer(1, 12))

    elements.append(Paragraph("Participants", styles["Heading2"]))

    if data.participants:
        for participant in data.participants:
            elements.append(
                Paragraph(
                    f"- {participant.username} ({participant.userId})",
                    styles["Normal"]
                )
            )
    else:
        elements.append(Paragraph("No participants provided.", styles["Normal"]))

    elements.append(Spacer(1, 12))

    elements.append(Paragraph("Timeline", styles["Heading2"]))

    if data.timeline:
        for event in data.timeline:
            elements.append(
                Paragraph(
                    f"[{event.timestamp}] {event.username}: {event.message}",
                    styles["Normal"]
                )
            )
    else:
        elements.append(Paragraph("No timeline events provided.", styles["Normal"]))

    elements.append(Spacer(1, 12))

    elements.append(Paragraph("Summary", styles["Heading2"]))
    elements.append(
        Paragraph(
            f"Incident '{data.incident.title}' involved "
            f"{len(data.participants)} participants and "
            f"{len(data.timeline)} timeline events.",
            styles["Normal"]
        )
    )

    document.build(elements)

    pdf_bytes = buffer.getvalue()
    buffer.close()

    return pdf_bytes