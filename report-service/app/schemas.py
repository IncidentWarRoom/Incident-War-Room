from typing import List, Optional

from pydantic import BaseModel, Field


class IncidentInfo(BaseModel):
    id: str = Field(..., description="Incident unique identifier")
    title: str = Field(..., description="Incident title")
    createdAt: str = Field(..., description="Incident creation time")
    closedAt: Optional[str] = Field(None, description="Incident closing time")
    severity: Optional[str] = Field(None, description="Incident severity")
    status: Optional[str] = Field(None, description="Incident status")


class Participant(BaseModel):
    userId: int = Field(..., description="Telegram user ID")
    username: str = Field(..., description="Telegram username")


class TimelineEvent(BaseModel):
    timestamp: str = Field(..., description="Event time")
    username: str = Field(..., description="User who created the event")
    message: str = Field(..., description="Timeline event message")


class ReportRequest(BaseModel):
    incident: IncidentInfo
    participants: List[Participant]
    timeline: List[TimelineEvent]

class ReportResponse(BaseModel):
    reportUrl: str = Field(..., description="Public URL of the generated PDF report")