import os
import time

import pytest

from client import ApiClient


@pytest.fixture(scope="session")
def base_url() -> str:
    """Points at the gateway. Override with API_BASE_URL for local runs, e.g.
    API_BASE_URL=http://localhost:8080 pytest
    In CI this is set to http://gateway:8080 (see gitlab-ci.yml snippet).
    """
    return os.environ.get("API_BASE_URL", "http://localhost:8080")


@pytest.fixture(scope="session")
def client(base_url) -> ApiClient:
    return ApiClient(base_url)


@pytest.fixture(scope="session")
def run_id() -> str:
    # Keeps test data unique across repeated CI runs against the same DB.
    return str(int(time.time() * 1000))


@pytest.fixture(scope="session")
def user_credentials(run_id):
    return {
        "name": "CI Test User",
        "email": f"ci-{run_id}@example.com",
        "password": "ci-password-123",
    }


@pytest.fixture(scope="session")
def registered_user(client, user_credentials):
    """Registers one throwaway user for the whole test session."""
    client.request("POST", "/register", expected_status=201, json_body=user_credentials)
    return user_credentials


@pytest.fixture(scope="session")
def auth_token(client, registered_user):
    resp = client.request(
        "POST",
        "/login",
        expected_status=200,
        json_body={
            "email": registered_user["email"],
            "password": registered_user["password"],
        },
    )
    data = resp.json()
    assert isinstance(data.get("token"), str) and data["token"], "login did not return a usable token"
    return data["token"]


@pytest.fixture
def meeting_payload():
    return {
        "title": "CI Test Meeting",
        "description": "Created by automated tests",
        "meeting_time": "2030-07-08T18:00:00Z",
        "type": "Спорт",
        "address": "CI Test address",
        "latitude": 55.7522,
        "longitude": 48.7446,
        "max_people": 5,
    }


@pytest.fixture
def request_payload():
    return {
        "title": "CI Test Request",
        "description": "Created by automated tests",
        "requester_address": "CI Test address",
        "type": "Помощь",
        "deadline": "2030-07-09T18:00:00Z",
    }


@pytest.fixture
def created_meeting(client, auth_token, meeting_payload):
    """A meeting created fresh for a single test."""
    resp = client.request("POST", "/meetings", expected_status=201, json_body=meeting_payload, token=auth_token)
    return resp.json()


@pytest.fixture
def created_request(client, auth_token, request_payload):
    """A request created fresh for a single test."""
    resp = client.request("POST", "/requests", expected_status=201, json_body=request_payload, token=auth_token)
    return resp.json()
