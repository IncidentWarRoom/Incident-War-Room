from datetime import datetime, timezone

from fastapi import FastAPI

from .schemas import ReportRequest, ReportResponse
from .pdf_generator import generate_incident_report
from .s3_storage import S3Storage


app = FastAPI(
    title="Incident War Room Report Service",
    description="Service for generating PDF incident reports",
    version="1.0.0",
)

storage = S3Storage()

@app.get("/health")
def health_check():
    return {"status": "ok"}


@app.post("/api/v1/reports/generate", response_model=ReportResponse)
def generate_report(request: ReportRequest):
    pdf_bytes = generate_incident_report(request)

    timestamp = datetime.now(timezone.utc).strftime("%Y%m%d%H%M%S")
    object_name = f"reports/{request.incident.id}-{timestamp}.pdf"

    report_url = storage.upload_pdf(
        pdf_bytes=pdf_bytes,
        object_name=object_name,
    )

    return ReportResponse(reportUrl=report_url)