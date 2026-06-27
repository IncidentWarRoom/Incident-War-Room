import re
from typing import Iterable

from .schemas import TimelineEvent


NOISE_MESSAGES = {
    "",
    "+",
    "++",
    "-",
    ".",
    "ok",
    "okay",
    "yes",
    "yeah",
    "да",
    "нет",
    "ок",
    "окей",
    "ага",
    "угу",
    "понял",
    "поняла",
    "принял",
    "приняла",
    "спс",
    "спасибо",
    "ща",
    "сейчас",
}

ACK_MESSAGES = {
    "ok",
    "okay",
    "yes",
    "yeah",
    "да",
    "ок",
    "окей",
    "ага",
    "угу",
    "понял",
    "поняла",
    "принял",
    "приняла",
    "беру",
    "взял",
    "взяла",
    "взял в работу",
    "взяла в работу",
    "беру в работу",
}

ASSIGNMENT_KEYWORDS = {
    "assign",
    "assigned",
    "take",
    "owner",
    "responsible",
    "task",
    "назнач",
    "взять",
    "возьми",
    "берешь",
    "бери",
    "ответственный",
    "задач",
    "работу",
    "на тебя",
    "на него",
    "на нее",
}

CLOSE_KEYWORDS = {
    "close",
    "closed",
    "resolved",
    "fixed",
    "done",
    "закры",
    "решено",
    "решили",
    "исправ",
    "почин",
    "готово",
}

SYSTEM_KEYWORDS = {
    "incident created",
    "incident closed",
    "alert",
    "severity",
    "status changed",
    "инцидент создан",
    "инцидент закрыт",
    "алерт",
    "статус",
    "severity",
}


def normalize_message(message: str | None) -> str:
    if not message:
        return ""

    message = message.strip()
    message = re.sub(r"\s+", " ", message)

    return message


def normalize_username(username: str | None) -> str:
    if not username:
        return ""

    username = username.strip().lower()
    username = username.removeprefix("@")

    return username


def message_contains_any(message: str, keywords: set[str]) -> bool:
    normalized = message.lower()

    for keyword in keywords:
        if keyword in normalized:
            return True

    return False


def is_ack_message(message: str | None) -> bool:
    normalized = normalize_message(message).lower()

    return normalized in ACK_MESSAGES


def is_noise_message(message: str | None) -> bool:
    normalized = normalize_message(message).lower()

    if normalized in NOISE_MESSAGES:
        return True

    if re.fullmatch(r"[+\-.]+", normalized):
        return True

    return False


def mentions_user(message: str | None, username: str | None) -> bool:
    normalized_message = normalize_message(message).lower()
    normalized_username = normalize_username(username)

    if not normalized_message or not normalized_username:
        return False

    possible_mentions = {
        normalized_username,
        f"@{normalized_username}",
    }

    for mention in possible_mentions:
        if mention in normalized_message:
            return True

    return False


def is_assignment_message(event: TimelineEvent) -> bool:
    message = normalize_message(event.message)

    if not message:
        return False

    return message_contains_any(message, ASSIGNMENT_KEYWORDS)


def is_close_message(event: TimelineEvent) -> bool:
    message = normalize_message(event.message)

    if not message:
        return False

    return message_contains_any(message, CLOSE_KEYWORDS)


def is_system_message(event: TimelineEvent) -> bool:
    message = normalize_message(event.message)

    if not message:
        return False

    return message_contains_any(message, SYSTEM_KEYWORDS)


def is_ack_for_previous_assignment(
    current_event: TimelineEvent,
    previous_event: TimelineEvent | None,
) -> bool:
    if previous_event is None:
        return False

    if not is_ack_message(current_event.message):
        return False

    if not is_assignment_message(previous_event):
        return False

    if mentions_user(previous_event.message, current_event.username):
        return True

    if normalize_username(previous_event.username) != normalize_username(current_event.username):
        return True

    return False


def copy_event_with_message(event: TimelineEvent, message: str) -> TimelineEvent:
    return TimelineEvent(
        timestamp=event.timestamp,
        username=event.username,
        message=message,
    )


def remove_duplicate_events(events: Iterable[TimelineEvent]) -> list[TimelineEvent]:
    unique_events = []
    seen = set()

    for event in events:
        key = (
            event.timestamp,
            normalize_username(event.username),
            normalize_message(event.message).lower(),
        )

        if key in seen:
            continue

        seen.add(key)
        unique_events.append(event)

    return unique_events


def preprocess_timeline(timeline: list[TimelineEvent]) -> list[TimelineEvent]:
    if not timeline:
        return []

    cleaned_events = []
    previous_original_event = None

    for event in timeline:
        message = normalize_message(event.message)

        if not message:
            previous_original_event = event
            continue

        if is_ack_for_previous_assignment(event, previous_original_event):
            cleaned_events.append(
                copy_event_with_message(
                    event,
                    f"{event.username} acknowledged and accepted the assignment.",
                )
            )
            previous_original_event = event
            continue

        if is_noise_message(message):
            previous_original_event = event
            continue

        if is_assignment_message(event):
            cleaned_events.append(copy_event_with_message(event, message))
            previous_original_event = event
            continue

        if is_close_message(event):
            cleaned_events.append(copy_event_with_message(event, message))
            previous_original_event = event
            continue

        if is_system_message(event):
            cleaned_events.append(copy_event_with_message(event, message))
            previous_original_event = event
            continue

        cleaned_events.append(copy_event_with_message(event, message))
        previous_original_event = event

    return remove_duplicate_events(cleaned_events)


def build_timeline_context(timeline: list[TimelineEvent]) -> str:
    cleaned_timeline = preprocess_timeline(timeline)

    if not cleaned_timeline:
        return "Timeline is empty or contains no useful events."

    lines = []

    for event in cleaned_timeline:
        lines.append(f"- [{event.timestamp}] {event.username}: {event.message}")

    return "\n".join(lines)