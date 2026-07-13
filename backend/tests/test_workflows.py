"""End-to-end journeys, each test tells one full story rather than
checking a single endpoint in isolation. If one of these fails partway,
the failing step's name in the traceback tells you exactly where the
journey broke.
"""


def test_register_then_login(client, run_id):
    email = f"workflow-{run_id}@example.com"
    password = "workflow-password-123"

    register_resp = client.request(
        "POST",
        "/register",
        expected_status=201,
        json_body={"name": "Workflow User", "email": email, "password": password},
    )
    assert isinstance(register_resp.json().get("message"), str)

    login_resp = client.request(
        "POST", "/login", expected_status=200, json_body={"email": email, "password": password}
    )
    token = login_resp.json().get("token")
    assert isinstance(token, str) and token

    me_resp = client.request("GET", "/me", expected_status=200, token=token)
    assert me_resp.json().get("email") == email


def test_create_and_apply_to_meeting(client, auth_token, meeting_payload):
    created = client.request(
        "POST", "/meetings", expected_status=201, json_body=meeting_payload, token=auth_token
    ).json()
    meeting_id = created["id"]

    # Anyone can view meeting details without a token.
    details = client.request("GET", f"/meetings/{meeting_id}", expected_status=200).json()
    assert details["title"] == meeting_payload["title"]

    # The creator applies to their own meeting.
    apply_resp = client.request("POST", f"/meetings/{meeting_id}", expected_status=200, token=auth_token)
    assert apply_resp.json() == 200  # endpoint echoes the status code as its body

    # It now shows up in the public listing.
    listing = client.request("GET", "/meetings", expected_status=200).json()
    assert any(m.get("id") == meeting_id for m in listing["meetings"])


def test_create_and_apply_to_request(client, auth_token, request_payload):
    created = client.request(
        "POST", "/requests", expected_status=201, json_body=request_payload, token=auth_token
    ).json()
    request_id = created["id"]

    details = client.request("GET", f"/requests/{request_id}", expected_status=200).json()
    assert details["request_id"] == request_id

    # "Applying" to a request means opening a chat on it.
    chat = client.request(
        "POST", f"/requests/{request_id}/chat", expected_status=200, token=auth_token
    ).json()
    assert chat["type"] in (1, "REQUEST")

    listing = client.request("GET", "/requests", expected_status=200).json()
    assert any(r.get("id") == request_id for r in listing["requests"])


def test_write_and_read_meeting_chat_message(client, auth_token, created_meeting):
    chat = client.request(
        "GET", f"/meetings/{created_meeting['id']}/chat", expected_status=200, token=auth_token
    ).json()
    chat_id = chat["chat_id"]

    message_text = "hello from the meeting chat workflow test"
    sent = client.request(
        "POST",
        f"/chats/{chat_id}/messages",
        expected_status=200,
        json_body={"text": message_text},
        token=auth_token,
    ).json()
    assert sent["text"] == message_text

    history = client.request(
        "GET", f"/chats/{chat_id}/messages", expected_status=200, token=auth_token
    ).json()
    assert any(m.get("text") == message_text for m in history["messages"])


def test_write_and_read_request_chat_message(client, auth_token, created_request):
    chat = client.request(
        "POST", f"/requests/{created_request['id']}/chat", expected_status=200, token=auth_token
    ).json()
    chat_id = chat["chat_id"]

    message_text = "hello from the request chat workflow test"
    sent = client.request(
        "POST",
        f"/chats/{chat_id}/messages",
        expected_status=200,
        json_body={"text": message_text},
        token=auth_token,
    ).json()
    assert sent["text"] == message_text

    history = client.request(
        "GET", f"/chats/{chat_id}/messages", expected_status=200, token=auth_token
    ).json()
    assert any(m.get("text") == message_text for m in history["messages"])

    # Both this and the meeting chat should also be reachable via the
    # generic chat list/detail endpoints, wired the same way.
    all_chats = client.request("GET", "/chats", expected_status=200, token=auth_token).json()
    assert any(c.get("chat_id") == chat_id for c in all_chats["chats"])

    detail = client.request("GET", f"/chats/{chat_id}", expected_status=200, token=auth_token).json()
    assert detail["chat_id"] == chat_id
