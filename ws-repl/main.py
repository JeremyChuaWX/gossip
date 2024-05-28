import argparse
import json
import readline
import threading
import time
from datetime import datetime, timezone

import requests
import websocket

SESSION_ID_HEADER = "x-session-id"
USER_ID_HEADER = "x-user-id"
WS_URL = "ws://127.0.0.1:3000/users/connect"
LOGIN_URL = "http://127.0.0.1:3000/login"
DETAILS_URL = "http://127.0.0.1:3000/users"
ROOM_ID = "97bde354-70f6-4f38-a078-8c43c881af93"
PASSWORD = "123"


def login(username: str) -> str:
    response = requests.post(
        LOGIN_URL, json={"username": username, "password": PASSWORD}
    )
    body = response.json()
    return body["data"]["sessionId"]


def get_user_id(session_id: str):
    response = requests.get(DETAILS_URL, headers={SESSION_ID_HEADER: session_id})
    body = response.json()
    return body["data"]["user"]["id"]


def websocket_repl(uri: str, session_id: str, user_id: str, args: any):
    websocket.enableTrace(args.debug)
    if args.mode == "display":
        handlers = {
            "on_message": on_message,
            "on_error": on_error,
            "on_close": on_close,
        }
    elif args.mode == "input":
        handlers = {
            "on_open": on_open,
        }
    else:
        raise Exception("invalid mode")
    ws = websocket.WebSocketApp(
        uri, header={SESSION_ID_HEADER: session_id, USER_ID_HEADER: user_id}, **handlers
    )
    ws.run_forever()


def on_message(ws, message):
    print("message:", message)


def on_error(ws, error):
    print("error:", error)


def on_close(ws, code, message):
    print("connection closed")


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
                    "userId": ws.header[USER_ID_HEADER],
                    "body": message,
                    "timestamp": datetime.now(timezone.utc).astimezone().isoformat(),
                }
            )
        )
    time.sleep(1)
    ws.close()


def parse_args():
    parser = argparse.ArgumentParser()
    # fmt: off
    parser.add_argument("--mode", choices=("display", "input"), type=str)
    parser.add_argument("--debug", default=False, type=bool)
    # fmt: on
    return parser.parse_args()


def main():
    args = parse_args()
    username = input("username: ")
    session_id = login(username)
    user_id = get_user_id(session_id)
    websocket_repl(WS_URL, session_id, user_id, args)


if __name__ == "__main__":
    main()
