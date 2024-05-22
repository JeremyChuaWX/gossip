import asyncio
import json
import readline

import requests
import websockets

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


async def websocket_repl(uri: str, session_id: str):
    async with websockets.connect(
        uri, extra_headers={SESSION_ID_HEADER: session_id}
    ) as ws:
        print(f"connected to {uri}")
        while True:
            body = input("message ('q' to quit): ")
            if body.lower() == "q":
                break
            message = {
                "roomId": ROOM_ID,
                "userId": USER_ID,
                "body": body,
            }
            await ws.send(json.dumps(message))
            print(f"sent: {body}")
            response = await ws.recv()
            print(f"received: {json.loads(response)}")


def main():
    username = input("username: ")
    session_id = login(username)
    print(session_id)
    asyncio.get_event_loop().run_until_complete(websocket_repl(WS_URL, session_id))


if __name__ == "__main__":
    main()
