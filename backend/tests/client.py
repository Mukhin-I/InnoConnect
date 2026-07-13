"""Small HTTP helper used across the test suite.

Wraps `requests` so tests can assert on status code in one line and get
back parsed JSON (or None for empty bodies) without repeating boilerplate.
"""
from __future__ import annotations

import requests


class ApiError(AssertionError):
    """Raised when a response doesn't match what the test expected."""


class ApiClient:
    def __init__(self, base_url: str):
        self.base_url = base_url.rstrip("/")
        self.session = requests.Session()

    def request(
        self,
        method: str,
        path: str,
        expected_status: int | None = None,
        json_body=None,
        raw_body: str | None = None,
        token: str | None = None,
    ) -> requests.Response:
        headers = {}
        if token:
            headers["Authorization"] = f"Bearer {token}"

        kwargs = {"headers": headers, "timeout": 15}
        if raw_body is not None:
            # Used for deliberately malformed/empty payloads.
            kwargs["data"] = raw_body
            headers["Content-Type"] = "application/json"
        elif json_body is not None:
            kwargs["json"] = json_body

        resp = self.session.request(method, f"{self.base_url}{path}", **kwargs)

        if expected_status is not None and resp.status_code != expected_status:
            raise ApiError(
                f"{method} {path}: expected HTTP {expected_status}, "
                f"got {resp.status_code}: {resp.text[:300]!r}"
            )
        return resp

    def json_or_none(self, resp: requests.Response):
        if resp.text == "":
            return None
        try:
            return resp.json()
        except ValueError as err:
            raise ApiError(f"response is not valid JSON: {err}: {resp.text[:300]!r}") from err
