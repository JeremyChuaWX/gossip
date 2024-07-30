"use strict";

const roomId = document.URL.split("/").pop();

const ws = new WebSocket("/connect");
ws.onopen = (event) => {
    console.log("onopen", event);
};
ws.onmessage = (event) => {
    console.log("onmessage event", event);
    console.log("event.data", JSON.parse(event.data));
};
ws.onerror = (event) => {
    console.log("onerror", event);
};
ws.onclose = (event) => {
    console.log("onclose", event);
};

const messageBox = document.getElementById("message-box");
messageBox.onsubmit = async (event) => {
    event.preventDefault();
    const formData = new FormData(messageBox);
    const body = formData.get("body");
    sendMessage(body);
};

/**
 * @param {string} body
 */
function sendMessage(body) {
    if (!roomId) {
        console.error("invalid roomId", roomId);
        return;
    }
    const message = {
        roomId: roomId,
        body: body,
        timestamp: new Date().toISOString(),
    };
    ws.send(JSON.stringify(message));
}
