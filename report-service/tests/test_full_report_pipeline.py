import os
import uuid
from pathlib import Path

import pytest
import requests
from dotenv import load_dotenv

from app.pdf_generator import generate_incident_report
from app.s3_storage import S3Storage
from app.schemas import ReportRequest


@pytest.mark.integration
def test_full_report_pipeline_uploads_pdf_to_s3(tmp_path):
    project_root = Path(__file__).resolve().parents[1]
    env_path = project_root / ".env"

    load_dotenv(env_path)

    required_env_vars = [
        "S3_ENDPOINT_URL",
        "S3_REGION",
        "S3_BUCKET_NAME",
        "S3_ACCESS_KEY",
        "S3_SECRET_KEY",
        "S3_PUBLIC_BASE_URL",
    ]

    missing_vars = [
        var_name for var_name in required_env_vars
        if not os.getenv(var_name)
    ]

    assert not missing_vars, f"Missing required env variables: {missing_vars}"

    request = ReportRequest(
        incident={
            "id": "inc-2026-06-27-042",
            "title": "Issue in production",
            "createdAt": "2026-06-27T09:14:00Z",
            "closedAt": "2026-06-27T10:08:00Z",
            "severity": "HIGH",
            "status": "resolved",
        },
        participants=[
            {
                "userId": 1001,
                "username": "teamlead",
            },
            {
                "userId": 1002,
                "username": "backend_dev",
            },
            {
                "userId": 1003,
                "username": "devops_engineer",
            },
            {
                "userId": 1004,
                "username": "qa_engineer",
            },
        ],
        timeline=[
            {
                "timestamp": "2026-06-27T09:14:00Z",
                "username": "teamlead",
                "message": "Alertmanager fired: high 5xx error rate on report-service in production",
            },
            {
                "timestamp": "2026-06-27T09:15:00Z",
                "username": "teamlead",
                "message": "Users cannot generate incident PDF reports from Telegram. Go service receives 500 from report-service",
            },
            {
                "timestamp": "2026-06-27T09:16:00Z",
                "username": "backend_dev",
                "message": "I will check report-service logs and recent changes",
            },
            {
                "timestamp": "2026-06-27T09:18:00Z",
                "username": "devops_engineer",
                "message": "Pods are running, no restarts. CPU and memory look normal",
            },
            {
                "timestamp": "2026-06-27T09:21:00Z",
                "username": "backend_dev",
                "message": "Found repeated errors in report-service logs: botocore validate_bucket_name TypeError, expected string or bytes-like object, got NoneType",
            },
            {
                "timestamp": "2026-06-27T09:23:00Z",
                "username": "teamlead",
                "message": "Could this be related to S3 configuration?",
            },
            {
                "timestamp": "2026-06-27T09:25:00Z",
                "username": "devops_engineer",
                "message": "Checking environment variables inside the report-service container",
            },
            {
                "timestamp": "2026-06-27T09:28:00Z",
                "username": "devops_engineer",
                "message": "Confirmed S3_BUCKET_NAME is missing in production environment",
            },
            {
                "timestamp": "2026-06-27T09:31:00Z",
                "username": "backend_dev",
                "message": "The service tries to upload generated PDF to S3 because S3_ENABLED is true, but bucket name is None",
            },
            {
                "timestamp": "2026-06-27T09:34:00Z",
                "username": "teamlead",
                "message": "backend_dev, please prepare a fallback so report-service returns PDF inline when S3 config is incomplete",
            },
            {
                "timestamp": "2026-06-27T09:35:00Z",
                "username": "backend_dev",
                "message": "ok, I will handle the fallback logic",
            },
            {
                "timestamp": "2026-06-27T09:42:00Z",
                "username": "backend_dev",
                "message": "Implemented S3 config validation. If required S3 env vars are missing, service returns PDF inline instead of failing",
            },
            {
                "timestamp": "2026-06-27T09:46:00Z",
                "username": "qa_engineer",
                "message": "Testing report generation without S3_BUCKET_NAME",
            },
            {
                "timestamp": "2026-06-27T09:49:00Z",
                "username": "qa_engineer",
                "message": "PDF is now returned successfully inline. Response content type is application/pdf",
            },
            {
                "timestamp": "2026-06-27T09:53:00Z",
                "username": "devops_engineer",
                "message": "Added missing S3_BUCKET_NAME variable to production secret and restarted report-service",
            },
            {
                "timestamp": "2026-06-27T09:58:00Z",
                "username": "qa_engineer",
                "message": "Testing report generation with S3 enabled again",
            },
            {
                "timestamp": "2026-06-27T10:02:00Z",
                "username": "qa_engineer",
                "message": "Report successfully generated, uploaded to S3, and public reportUrl is returned to Go service",
            },
            {
                "timestamp": "2026-06-27T10:05:00Z",
                "username": "teamlead",
                "message": "Impact is gone. Telegram users can generate and download reports again",
            },
            {
                "timestamp": "2026-06-27T10:08:00Z",
                "username": "teamlead",
                "message": "Incident resolved. Root cause was missing S3_BUCKET_NAME environment variable in report-service production config",
            },
        ],
    )

    pdf_bytes = generate_incident_report(request)

    assert isinstance(pdf_bytes, bytes), "PDF generator did not return bytes"
    assert len(pdf_bytes) > 0, "Generated PDF is empty"
    assert pdf_bytes.startswith(b"%PDF"), "Generated file is not a valid PDF"

    output_pdf_path = tmp_path / "test_incident_report.pdf"
    output_pdf_path.write_bytes(pdf_bytes)

    assert output_pdf_path.exists(), "PDF file was not saved locally"
    assert output_pdf_path.stat().st_size > 0, "Saved PDF file is empty"

    storage = S3Storage()

    object_name = f"reports/test-report-{uuid.uuid4()}.pdf"

    report_url = storage.upload_pdf(
        pdf_bytes=pdf_bytes,
        object_name=object_name,
    )

    assert report_url is not None, "S3 did not return report URL"
    assert report_url.startswith("http"), f"Invalid report URL: {report_url}"

    response = requests.get(report_url, timeout=30)

    assert response.status_code == 200, (
        f"Report URL is not accessible. "
        f"Status code: {response.status_code}. "
        f"URL: {report_url}"
    )

    content_type = response.headers.get("Content-Type", "")

    assert (
        "application/pdf" in content_type
        or "application/octet-stream" in content_type
        or "binary/octet-stream" in content_type
    ), f"Unexpected Content-Type: {content_type}"

    assert response.content.startswith(b"%PDF"), "Downloaded file is not a valid PDF"

    print("\nReport uploaded successfully.")
    print(f"S3 object name: {object_name}")
    print(f"Report URL: {report_url}")