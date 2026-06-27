from app.schemas import TimelineEvent
from app.timeline_preprocessor import (
    build_timeline_context,
    is_ack_for_previous_assignment,
    preprocess_timeline,
)


def test_regular_noise_message_is_removed():
    timeline = [
        TimelineEvent(
            timestamp="2026-06-25T12:00:00Z",
            username="aslan",
            message="ок",
        ),
        TimelineEvent(
            timestamp="2026-06-25T12:01:00Z",
            username="ivan",
            message="Проверил логи, ошибка была в подключении к базе",
        ),
    ]

    cleaned = preprocess_timeline(timeline)

    assert len(cleaned) == 1
    assert cleaned[0].message == "Проверил логи, ошибка была в подключении к базе"


def test_ack_after_assignment_is_kept():
    timeline = [
        TimelineEvent(
            timestamp="2026-06-25T12:00:00Z",
            username="teamlead",
            message="Ролан, возьми задачу на себя",
        ),
        TimelineEvent(
            timestamp="2026-06-25T12:01:00Z",
            username="rolan",
            message="ок",
        ),
    ]

    cleaned = preprocess_timeline(timeline)

    assert len(cleaned) == 2
    assert cleaned[1].message == "rolan acknowledged and accepted the assignment."


def test_ack_after_non_assignment_is_removed():
    timeline = [
        TimelineEvent(
            timestamp="2026-06-25T12:00:00Z",
            username="teamlead",
            message="Инцидент начался после деплоя",
        ),
        TimelineEvent(
            timestamp="2026-06-25T12:01:00Z",
            username="rolan",
            message="ок",
        ),
    ]

    cleaned = preprocess_timeline(timeline)

    assert len(cleaned) == 1
    assert cleaned[0].message == "Инцидент начался после деплоя"


def test_assignment_ack_detection():
    previous_event = TimelineEvent(
        timestamp="2026-06-25T12:00:00Z",
        username="teamlead",
        message="rolan, take this task",
    )
    current_event = TimelineEvent(
        timestamp="2026-06-25T12:01:00Z",
        username="rolan",
        message="ok",
    )

    assert is_ack_for_previous_assignment(current_event, previous_event) is True


def test_build_timeline_context_contains_acknowledgement():
    timeline = [
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
            username="teamlead",
            message="Incident closed after rollback",
        ),
    ]

    context = build_timeline_context(timeline)

    assert "rolan, take this task" in context
    assert "rolan acknowledged and accepted the assignment" in context
    assert "Incident closed after rollback" in context


def test_duplicate_events_are_removed():
    timeline = [
        TimelineEvent(
            timestamp="2026-06-25T12:00:00Z",
            username="aslan",
            message="Database connection failed after deployment",
        ),
        TimelineEvent(
            timestamp="2026-06-25T12:00:00Z",
            username="aslan",
            message="Database connection failed after deployment",
        ),
    ]

    cleaned = preprocess_timeline(timeline)

    assert len(cleaned) == 1