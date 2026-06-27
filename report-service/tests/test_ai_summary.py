from app.ai_summary import (
    build_ai_summary_prompt,
    build_incident_context,
    build_participants_context,
)
from app.schemas import IncidentInfo, Participant, ReportRequest, TimelineEvent


def test_build_incident_context_contains_real_schema_fields():
    data = ReportRequest(
        incident=IncidentInfo(
            id="15",
            title="Database outage",
            createdAt="2026-06-25T12:00:00Z",
            closedAt="2026-06-25T13:00:00Z",
            severity="HIGH",
            status="CLOSED",
        ),
        participants=[],
        timeline=[],
    )

    context = build_incident_context(data)

    assert "Incident ID: 15" in context
    assert "Title: Database outage" in context
    assert "Severity: HIGH" in context
    assert "Status: CLOSED" in context
    assert "Created at: 2026-06-25T12:00:00Z" in context
    assert "Closed at: 2026-06-25T13:00:00Z" in context


def test_build_incident_context_handles_optional_fields():
    data = ReportRequest(
        incident=IncidentInfo(
            id="15",
            title="Database outage",
            createdAt="2026-06-25T12:00:00Z",
            closedAt=None,
            severity=None,
            status=None,
        ),
        participants=[],
        timeline=[],
    )

    context = build_incident_context(data)

    assert "Severity: Unknown" in context
    assert "Status: Unknown" in context
    assert "Closed at: Not closed or not provided" in context


def test_build_participants_context_contains_username_and_user_id():
    data = ReportRequest(
        incident=IncidentInfo(
            id="15",
            title="Database outage",
            createdAt="2026-06-25T12:00:00Z",
        ),
        participants=[
            Participant(
                userId=123,
                username="rolan",
            )
        ],
        timeline=[],
    )

    context = build_participants_context(data)

    assert "rolan" in context
    assert "123" in context


def test_build_ai_summary_prompt_contains_cleaned_timeline():
    data = ReportRequest(
        incident=IncidentInfo(
            id="15",
            title="Database outage",
            createdAt="2026-06-25T12:00:00Z",
            closedAt="2026-06-25T13:00:00Z",
            severity="HIGH",
            status="CLOSED",
        ),
        participants=[
            Participant(
                userId=123,
                username="rolan",
            )
        ],
        timeline=[
            TimelineEvent(
                timestamp="2026-06-25T12:00:00Z",
                username="teamlead",
                message="rolan, take this task",
            ),
            TimelineEvent(
                timestamp="2026-06-25T12:01:00Z",
                username="rolan",
                message="ok",
            ),
            TimelineEvent(
                timestamp="2026-06-25T12:30:00Z",
                username="rolan",
                message="Database connection issue was fixed after rollback",
            ),
        ],
    )

    prompt = build_ai_summary_prompt(data)

    assert "Database outage" in prompt
    assert "rolan" in prompt
    assert "rolan acknowledged and accepted the assignment" in prompt
    assert "Database connection issue was fixed after rollback" in prompt