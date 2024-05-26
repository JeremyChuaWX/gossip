import json
import readline
import threading
import time
from datetime import datetime, timezone

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
    ws.run_forever()


def on_message(ws, message):
    print("message:", message)


def on_error(ws, error):
    print("error:", error)


def on_close(ws, code, message):
    print("connectionclosed")


def on_open(ws):
    print("connection opened")
    threading.Thread(target=run, args=(ws,)).start()


def run(ws: websocket.WebSocketApp):
    while True:
        message = input("message:")
        if message.lower() == "q":
            break
        ws.send(
            json.dumps(
                {
                    "roomId": ROOM_ID,
                    "userId": USER_ID,
                    "body": message,
                    "timestamp": datetime.now(timezone.utc).astimezone().isoformat(),
                }
            )
        )
    time.sleep(1)
    ws.close()


def main():
    username = input("username: ")
    session_id = login(username)
    websocket_repl(WS_URL, session_id)


if __name__ == "__main__":
    main()
