/** @type string | undefined */
let currentRoom = undefined;

/** @type string | undefined */
let currentUsername = undefined;

/** @type WebSocket | undefined */
let ws = undefined;

/** @type string[] */
let rooms = [];

const HTTP_URL = "http://127.0.0.1:3000/room/";
const WS_URL = "ws://127.0.0.1:3000/room/";

const availableRoomsList = document.getElementById("available-rooms");
const newRoomForm = document.getElementById("new-room");
const joinRoomForm = document.getElementById("join-room");
const currentStatusHeader = document.getElementById("current-status");
const messagesTextArea = document.getElementById("messages");
const sendMessageForm = document.getElementById("send-message");

function bootstrap() {
    newRoomForm.onsubmit = newRoom;
    joinRoomForm.onsubmit = joinRoom;
    sendMessageForm.onsubmit = sendMessage;
}

async function newRoom(e) {
    e.preventDefault();

    const newRoomData = new FormData(newRoomForm);
    const name = newRoomData.get("room-name");
    if (name === null || name === "") {
        console.log("empty room name");
        return;
    }
    const newRoomURL = new URL(HTTP_URL);
    newRoomURL.searchParams.append("name", name);
    await fetch(newRoomURL, { method: "POST" });
    const data = await fetch(HTTP_URL).then((res) => res.json());
    rooms = data.rooms;
    render();
}

function joinRoom(e) {
    e.preventDefault();

    const joinRoomData = new FormData(joinRoomForm);
    currentRoom = joinRoomData.get("room-name");
    currentUsername = joinRoomData.get("username");
    const joinRoomURL = new URL(currentRoom, WS_URL);
    joinRoomURL.searchParams.append("username", currentUsername);
    initWS(joinRoomURL);
}

function sendMessage(e) {
    e.preventDefault();

    const sendMessageData = new FormData(sendMessageForm);
    const message = sendMessageData.get("message");
    ws.send(JSON.stringify({ message }));
}

function render() {
    // render li for available rooms
    const roomsLi = rooms.map((room) => {
        const elem = document.createElement("li");
        elem.appendChild(document.createTextNode(room));
        return elem;
    });
    availableRoomsList.replaceChildren(...roomsLi);

    // render header for current status: room (user)
}

/**
 * @param {URL} url
 */
function initWS(url) {
    ws = new WebSocket(url);
    ws.onopen = console.log;
    ws.onmessage = console.log;
    ws.onerror = console.log;
}

bootstrap();
