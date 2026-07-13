"""Endpoint-level contract checks: one endpoint, one behaviour, per test.

These don't chain user journeys together (see test_workflows.py for that) —
each test only cares whether a given request returns the right status code
and a JSON body of the expected shape.
"""
import pytest


# --------------------------------------------------------------------------
# /register
# --------------------------------------------------------------------------

class TestRegister:
    def test_empty_body_is_rejected(self, client):
        resp = client.request("POST", "/register", expected_status=400, raw_body="")
        assert isinstance(client.json_or_none(resp), (dict, type(None)))

    def test_malformed_json_is_rejected(self, client):
        client.request("POST", "/register", expected_status=400, raw_body="{bad json")

    def test_valid_registration_returns_message(self, client, run_id):
        email = f"contract-{run_id}@example.com"
        resp = client.request(
            "POST",
            "/register",
            expected_status=201,
            json_body={"name": "Contract User", "email": email, "password": "pw-123456"},
        )
        data = resp.json()
        assert isinstance(data, dict)
        assert isinstance(data.get("message"), str)

    def test_duplicate_registration_is_rejected(self, client, registered_user):
        client.request("POST", "/register", expected_status=400, json_body=registered_user)


# --------------------------------------------------------------------------
# /login
# --------------------------------------------------------------------------

class TestLogin:
    def test_empty_body_is_rejected(self, client):
        client.request("POST", "/login", expected_status=400, raw_body="")

    def test_malformed_json_is_rejected(self, client):
        client.request("POST", "/login", expected_status=400, raw_body="{bad json")

    def test_wrong_password_is_unauthorized(self, client, registered_user):
        client.request(
            "POST",
            "/login",
            expected_status=401,
            json_body={"email": registered_user["email"], "password": "definitely-wrong"},
        )

    def test_valid_login_returns_bearer_token(self, client, registered_user):
        resp = client.request(
            "POST",
            "/login",
            expected_status=200,
            json_body={"email": registered_user["email"], "password": registered_user["password"]},
        )
        data = resp.json()
        assert isinstance(data.get("token"), str) and data["token"]
        assert data.get("type") == "Bearer"
        assert isinstance(data.get("expiresIn"), int)

    def test_me_requires_and_accepts_token(self, client, auth_token):
        client.request("GET", "/me", expected_status=401)  # no token
        resp = client.request("GET", "/me", expected_status=200, token=auth_token)
        assert isinstance(resp.json(), dict)


# --------------------------------------------------------------------------
# /meetings
# --------------------------------------------------------------------------

class TestMeetingsContract:
    def test_get_meetings_returns_array(self, client):
        resp = client.request("GET", "/meetings", expected_status=200)
        data = resp.json()
        assert isinstance(data.get("meetings"), (list, type(None)))

    def test_post_requires_auth(self, client, meeting_payload):
        client.request("POST", "/meetings", expected_status=401, json_body=meeting_payload)

    def test_post_rejects_empty_body(self, client, auth_token):
        client.request("POST", "/meetings", expected_status=400, json_body={}, token=auth_token)

    def test_post_rejects_invalid_time_format(self, client, auth_token, meeting_payload):
        bad = dict(meeting_payload, meeting_time="not-rfc3339")
        client.request("POST", "/meetings", expected_status=400, json_body=bad, token=auth_token)

    def test_post_creates_meeting_with_expected_shape(self, client, created_meeting, meeting_payload):
        assert isinstance(created_meeting.get("id"), int) and created_meeting["id"] > 0
        assert created_meeting.get("type") == meeting_payload["type"]

    def test_get_meeting_by_id_returns_details(self, client, created_meeting, meeting_payload):
        resp = client.request("GET", f"/meetings/{created_meeting['id']}", expected_status=200)
        data = resp.json()
        assert data.get("title") == meeting_payload["title"]


# --------------------------------------------------------------------------
# /requests
# --------------------------------------------------------------------------

class TestRequestsContract:
    def test_get_requests_returns_array(self, client):
        resp = client.request("GET", "/requests", expected_status=200)
        data = resp.json()
        assert isinstance(data.get("requests"), (list, type(None)))

    def test_post_requires_auth(self, client, request_payload):
        client.request("POST", "/requests", expected_status=401, json_body=request_payload)

    def test_post_rejects_empty_body(self, client, auth_token):
        client.request("POST", "/requests", expected_status=400, json_body={}, token=auth_token)

    def test_post_rejects_invalid_deadline_format(self, client, auth_token, request_payload):
        bad = dict(request_payload, deadline="not-rfc3339")
        client.request("POST", "/requests", expected_status=400, json_body=bad, token=auth_token)

    def test_post_creates_request_with_expected_shape(self, client, created_request, request_payload):
        assert isinstance(created_request.get("id"), int) and created_request["id"] > 0
        assert created_request.get("title") == request_payload["title"]

    def test_get_request_by_id_returns_details(self, client, created_request):
        resp = client.request("GET", f"/requests/{created_request['id']}", expected_status=200)
        data = resp.json()
        # Note: unlike /meetings/{id} (field "id"), this endpoint returns "request_id".
        assert data.get("request_id") == created_request["id"]


# --------------------------------------------------------------------------
# Chats (meeting chat, request chat, generic /chats/{id} + messages)
# --------------------------------------------------------------------------

class TestChatsContract:
    def test_meeting_chat_bad_id_is_rejected(self, client, auth_token):
        client.request("GET", "/meetings/not-a-number/chat", expected_status=400, token=auth_token)

    def test_meeting_chat_requires_auth(self, client, created_meeting):
        client.request("GET", f"/meetings/{created_meeting['id']}/chat", expected_status=401)

    def test_meeting_chat_returns_chat_object(self, client, auth_token, created_meeting):
        resp = client.request(
            "GET", f"/meetings/{created_meeting['id']}/chat", expected_status=200, token=auth_token
        )
        data = resp.json()
        assert isinstance(data.get("chat_id"), int) and data["chat_id"] > 0
        assert data.get("type") in (2, "MEETING")

    def test_request_chat_bad_id_is_rejected(self, client, auth_token):
        client.request("POST", "/requests/not-a-number/chat", expected_status=400, token=auth_token)

    def test_request_chat_requires_auth(self, client, created_request):
        client.request("POST", f"/requests/{created_request['id']}/chat", expected_status=401)

    def test_request_chat_returns_chat_object(self, client, auth_token, created_request):
        resp = client.request(
            "POST", f"/requests/{created_request['id']}/chat", expected_status=200, token=auth_token
        )
        data = resp.json()
        assert isinstance(data.get("chat_id"), int) and data["chat_id"] > 0
        assert data.get("type") in (1, "REQUEST")

    def test_get_chats_list_requires_auth_and_returns_array(self, client, auth_token):
        client.request("GET", "/chats", expected_status=401)
        resp = client.request("GET", "/chats", expected_status=200, token=auth_token)
        assert isinstance(resp.json().get("chats"), (list, type(None)))

    def test_chat_by_id_bad_id_is_rejected(self, client, auth_token):
        client.request("GET", "/chats/not-a-number", expected_status=400, token=auth_token)

    def test_post_message_rejects_empty_and_malformed_body(self, client, auth_token, created_meeting):
        chat = client.request(
            "GET", f"/meetings/{created_meeting['id']}/chat", expected_status=200, token=auth_token
        ).json()
        chat_id = chat["chat_id"]
        client.request("POST", f"/chats/{chat_id}/messages", expected_status=400, raw_body="", token=auth_token)
        client.request(
            "POST", f"/chats/{chat_id}/messages", expected_status=400, raw_body="{bad json", token=auth_token
        )

    def test_post_message_returns_message_shape(self, client, auth_token, created_meeting):
        chat = client.request(
            "GET", f"/meetings/{created_meeting['id']}/chat", expected_status=200, token=auth_token
        ).json()
        chat_id = chat["chat_id"]
        resp = client.request(
            "POST",
            f"/chats/{chat_id}/messages",
            expected_status=200,
            json_body={"text": "contract test message"},
            token=auth_token,
        )
        data = resp.json()
        assert isinstance(data.get("id"), int)
        assert data.get("text") == "contract test message"
        assert isinstance(data.get("sender"), dict)
        assert isinstance(data.get("sent_at"), str) and data["sent_at"]

    def test_get_messages_returns_array(self, client, auth_token, created_meeting):
        chat = client.request(
            "GET", f"/meetings/{created_meeting['id']}/chat", expected_status=200, token=auth_token
        ).json()
        chat_id = chat["chat_id"]
        resp = client.request("GET", f"/chats/{chat_id}/messages", expected_status=200, token=auth_token)
        assert isinstance(resp.json().get("messages"), (list, type(None)))
