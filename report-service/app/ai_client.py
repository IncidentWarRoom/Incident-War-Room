import logging
import os
from typing import Any

import requests

from .ai_title import build_ai_title_prompt
from .ai_summary import build_ai_summary_prompt
from .schemas import ReportRequest


logger = logging.getLogger(__name__)


DEEPSEEK_API_URL = "https://api.deepseek.com/chat/completions"
DEFAULT_MODEL = "deepseek-v4-flash"
DEFAULT_TIMEOUT_SECONDS = 30


def get_fallback_summary(data: ReportRequest) -> str:
    incident = data.incident

    return (
        "AI summary is currently unavailable. "
        f"Incident '{incident.title}' has status '{incident.status or 'Unknown'}' "
        f"and severity '{incident.severity or 'Unknown'}'. "
        "Please check the incident timeline for details."
    )


def extract_summary_from_response(response_data: dict[str, Any]) -> str:
    try:
        summary = response_data["choices"][0]["message"]["content"]
    except (KeyError, IndexError, TypeError):
        raise ValueError("Invalid DeepSeek response format")

    if not summary or not summary.strip():
        raise ValueError("DeepSeek returned an empty summary")

    return summary.strip()


def generate_ai_summary(data: ReportRequest) -> str:
    api_key = os.getenv("DEEPSEEK_API_KEY")

    if not api_key:
        logger.warning("DEEPSEEK_API_KEY is not set. Using fallback summary.")
        return get_fallback_summary(data)

    prompt = build_ai_summary_prompt(data)

    payload = {
        "model": os.getenv("DEEPSEEK_MODEL", DEFAULT_MODEL),
        "messages": [
            {
                "role": "system",
                "content": "You are an assistant that writes clear incident report summaries.",
            },
            {
                "role": "user",
                "content": prompt,
            },
        ],
        "temperature": 0.2,
        "max_tokens": 1200,
        "stream": False,
    }

    headers = {
        "Authorization": f"Bearer {api_key}",
        "Content-Type": "application/json",
    }

    try:
        response = requests.post(
            DEEPSEEK_API_URL,
            headers=headers,
            json=payload,
            timeout=DEFAULT_TIMEOUT_SECONDS,
        )

        response.raise_for_status()

        response_data = response.json()

        return extract_summary_from_response(response_data)

    except requests.RequestException as error:
        logger.exception("DeepSeek API request failed: %s", error)
        return get_fallback_summary(data)

    except ValueError as error:
        logger.exception("Failed to parse DeepSeek response: %s", error)
        return get_fallback_summary(data)


def get_fallback_title(data: ReportRequest) -> str:
    timeline_messages = " ".join(
        event.message.lower()
        for event in data.timeline
    )

    if "s3_bucket_name" in timeline_messages or "bucket" in timeline_messages:
        return "Missing S3 bucket configuration"

    if "5xx" in timeline_messages or "500" in timeline_messages:
        return "Report service 5xx errors"

    if "pdf" in timeline_messages and "report" in timeline_messages:
        return "Report PDF generation failure"

    if data.incident.severity:
        return f"{data.incident.severity.capitalize()} incident resolved"

    return f"Incident {data.incident.id} report"


def clean_ai_title(title: str) -> str:
    cleaned_title = title.strip()
    cleaned_title = cleaned_title.strip('"')
    cleaned_title = cleaned_title.strip("'")
    cleaned_title = cleaned_title.rstrip(".")

    return cleaned_title.strip()


def generate_ai_title(data: ReportRequest) -> str:
    api_key = os.getenv("DEEPSEEK_API_KEY")

    if not api_key:
        logger.warning("DEEPSEEK_API_KEY is not set. Using fallback AI title.")
        return get_fallback_title(data)

    prompt = build_ai_title_prompt(data)
    original_title = data.incident.title.strip().lower()

    payload = {
        "model": os.getenv("DEEPSEEK_MODEL", DEFAULT_MODEL),
        "messages": [
            {
                "role": "system",
                "content": (
                    "You generate new short technical incident titles. "
                    "Never copy placeholder titles such as problem_4 or Issue in production."
                ),
            },
            {
                "role": "user",
                "content": prompt,
            },
        ],
        "temperature": 0.4,
        "max_tokens": 60,
        "stream": False,
    }

    headers = {
        "Authorization": f"Bearer {api_key}",
        "Content-Type": "application/json",
    }

    try:
        response = requests.post(
            DEEPSEEK_API_URL,
            headers=headers,
            json=payload,
            timeout=DEFAULT_TIMEOUT_SECONDS,
        )

        response.raise_for_status()

        response_data = response.json()
        title = extract_summary_from_response(response_data)
        title = clean_ai_title(title)

        logger.info("DeepSeek raw AI title: %s", title)

        if not title:
            return get_fallback_title(data)

        if title.strip().lower() == original_title:
            logger.warning(
                "DeepSeek copied original title '%s'. Retrying with stricter prompt.",
                data.incident.title,
            )

            retry_payload = {
                **payload,
                "messages": [
                    {
                        "role": "system",
                        "content": (
                            "You generate new short technical incident titles. "
                            "You must not copy the original title. "
                            "Use the timeline and root cause."
                        ),
                    },
                    {
                        "role": "user",
                        "content": (
                            prompt
                            + "\n\nThe previous answer copied the original title. "
                            + "Generate a different technical title based on the timeline."
                        ),
                    },
                ],
                "temperature": 0.6,
            }

            retry_response = requests.post(
                DEEPSEEK_API_URL,
                headers=headers,
                json=retry_payload,
                timeout=DEFAULT_TIMEOUT_SECONDS,
            )

            retry_response.raise_for_status()

            retry_data = retry_response.json()
            retry_title = extract_summary_from_response(retry_data)
            retry_title = clean_ai_title(retry_title)

            logger.info("DeepSeek retry AI title: %s", retry_title)

            if retry_title and retry_title.strip().lower() != original_title:
                return retry_title

            return get_fallback_title(data)

        return title

    except requests.RequestException as error:
        logger.exception("DeepSeek API request for AI title failed: %s", error)
        return get_fallback_title(data)

    except ValueError as error:
        logger.exception("Failed to parse DeepSeek AI title response: %s", error)
        return get_fallback_title(data)