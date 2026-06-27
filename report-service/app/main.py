import os
from datetime import datetime, timezone

from dotenv import load_dotenv
from fastapi import FastAPI, HTTPException, Response

from .schemas import ReportRequest, ReportResponse
from .pdf_generator import generate_incident_report
from .s3_storage import S3Storage


load_dotenv()

app = FastAPI(
    title="Incident War Room Report Service",
    description="Service for generating PDF incident reports",
    version="1.0.0",
)


REQUIRED_S3_ENV_VARS = [
    "S3_ENDPOINT_URL",
    "S3_REGION",
    "S3_BUCKET_NAME",
    "S3_ACCESS_KEY",
    "S3_SECRET_KEY",
    "S3_PUBLIC_BASE_URL",
]


def is_s3_configured() -> bool:
    return all(os.getenv(env_name) for env_name in REQUIRED_S3_ENV_VARS)


def should_upload_to_s3() -> bool:
    s3_enabled = os.getenv("S3_ENABLED", "auto").lower()

    if s3_enabled == "false":
        return False

    return is_s3_configured()


def build_report_object_name(request: ReportRequest) -> str:
    timestamp = datetime.now(timezone.utc).strftime("%Y%m%d%H%M%S")
    return f"reports/{request.incident.id}-{timestamp}.pdf"


def upload_report_to_s3(pdf_bytes: bytes, object_name: str) -> str:
    try:
        storage = S3Storage()

        return storage.upload_pdf(
            pdf_bytes=pdf_bytes,
            object_name=object_name,
        )

    except Exception as error:
        raise HTTPException(
            status_code=503,
            detail=f"Failed to upload report to S3: {error}",
        )


@app.get("/health")
def health_check():
    return {"status": "ok"}


@app.post("/api/v1/reports/generate-url", response_model=ReportResponse)
def generate_report_url(request: ReportRequest):
    pdf_bytes = generate_incident_report(request)
    object_name = build_report_object_name(request)

    report_url = upload_report_to_s3(
        pdf_bytes=pdf_bytes,
        object_name=object_name,
    )

    return ReportResponse(reportUrl=report_url)


@app.post("/api/v1/reports/generate-inline")
def generate_report_inline(request: ReportRequest):
    pdf_bytes = generate_incident_report(request)

    filename = f"incident-report-{request.incident.id}.pdf"

    return Response(
        content=pdf_bytes,
        media_type="application/pdf",
        headers={
            "Content-Disposition": f'attachment; filename="{filename}"'
        },
    )


@app.post("/api/v1/reports/generate")
def generate_report(request: ReportRequest):
    if should_upload_to_s3():
        return generate_report_url(request)

    return generate_report_inline(request)