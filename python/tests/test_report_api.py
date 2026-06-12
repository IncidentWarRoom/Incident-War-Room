from fastapi.testclient import TestClient

from app.main import app


client = TestClient(app)


def test_health_check_returns_ok():
    response = client.get("/health")

    assert response.status_code == 200
    assert response.json() == {"status": "ok"}


def test_generate_report_returns_pdf():
    payload = {
        "incident": {
            "id": "15",
            "title": "Payment Service Down",
            "createdAt": "2026-06-01T14:03:00Z",
            "closedAt": "2026-06-01T14:18:00Z",
            "severity": "HIGH",
            "status": "CLOSED"
        },
        "participants": [
            {
                "userId": 1,
                "username": "rolan"
            },
            {
                "userId": 2,
                "username": "bogdan"
            }
        ],
        "timeline": [
            {
                "timestamp": "2026-06-01T14:05:00Z",
                "username": "rolan",
                "message": "Started investigating database issues"
            },
            {
                "timestamp": "2026-06-01T14:09:00Z",
                "username": "bogdan",
                "message": "Found connection pool errors"
            }
        ]
    }

    response = client.post("/api/v1/reports/generate", json=payload)

    assert response.status_code == 200
    assert response.headers["content-type"] == "application/pdf"
    assert response.content.startswith(b"%PDF")


def test_generate_report_without_required_field_returns_422():
    payload = {
        "incident": {
            "id": "15",
            "createdAt": "2026-06-01T14:03:00Z"
        },
        "participants": [],
        "timeline": []
    }

    response = client.post("/api/v1/reports/generate", json=payload)

    assert response.status_code == 422