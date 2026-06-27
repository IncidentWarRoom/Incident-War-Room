from dotenv import load_dotenv

load_dotenv()

from app.ai_client import generate_ai_summary
from app.schemas import IncidentInfo, Participant, ReportRequest, TimelineEvent


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

summary = generate_ai_summary(data)

print("\nAI SUMMARY:\n")
print(summary)