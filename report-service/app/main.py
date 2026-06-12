from fastapi import FastAPI
from fastapi.responses import Response

from .schemas import ReportRequest
from .pdf_generator import generate_incident_report


app = FastAPI(
    title="Incident War Room Report Service",
    description="Service for generating PDF incident reports",
    version="1.0.0",
)


@app.get("/health")
def health_check():
    return {"status": "ok"}


@app.post("/api/v1/reports/generate")
def generate_report(request: ReportRequest):
    pdf_bytes = generate_incident_report(request)

    return Response(
        content=pdf_bytes,
        media_type="application/pdf",
        headers={
            "Content-Disposition": "attachment; filename=incident_report.pdf"
        },
    )