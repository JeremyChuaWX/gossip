import readline

import rel
import requests
import websocket

SESSION_ID_HEADER = "x-session-id"
WS_URL = "ws://127.0.0.1:3000/users/connect"
LOGIN_URL = "http://127.0.0.1:3000/login"
USER_ID = "2150c411-4385-4e23-b919-e760181d8321"
ROOM_ID = "7852da78-f7b0-4278-b9c8-5d64766f65cd"
PASSWORD = "123"


def login(username: str) -> str:
    response = requests.post(
        LOGIN_URL, json={"username": username, "password": PASSWORD}
    )
    body = response.json()
    return body["data"]["sessionId"]


def websocket_repl(uri: str, session_id: str, debug=False):
    websocket.enableTrace(debug)
    ws = websocket.WebSocketApp(
        uri,
        header={SESSION_ID_HEADER: session_id},
        on_open=on_open,
        on_message=on_message,
        on_error=on_error,
        on_close=on_close,
    )
    ws.run_forever(dispatcher=rel, reconnect=5)
    rel.signal(2, rel.abort)
    rel.dispatch()


def on_message(ws, message):
    print(message)


def on_error(ws, error):
    print(error)


def on_close(ws, code, message):
    print("closed")


def on_open(ws):
    print("opened")


def main():
    username = input("username: ")
    session_id = login(username)
    websocket_repl(WS_URL, session_id)


if __name__ == "__main__":
    main()
