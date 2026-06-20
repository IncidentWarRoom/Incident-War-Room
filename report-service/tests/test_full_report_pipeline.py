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
            "id": "test-incident-001",
            "title": "Test Incident Report",
            "createdAt": "2026-06-19T12:00:00Z",
            "closedAt": "2026-06-19T13:00:00Z",
            "severity": "HIGH",
            "status": "resolved",
        },
        participants=[
            {
                "userId": 1001,
                "username": "alice",
            },
            {
                "userId": 1002,
                "username": "bob",
            },
        ],
        timeline=[
            {
                "timestamp": "2026-06-19T12:00:00Z",
                "username": "alice",
                "message": "Incident created",
            },
            {
                "timestamp": "2026-06-19T12:30:00Z",
                "username": "bob",
                "message": "Investigation started",
            },
            {
                "timestamp": "2026-06-19T13:00:00Z",
                "username": "alice",
                "message": "Incident resolved",
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