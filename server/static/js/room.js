"use strict";

/**
 * @typedef {Object} Message
 * @property {string} roomId
 * @property {string} userId
 * @property {string} username
 * @property {string} body
 * @property {string} timestamp
 */

import { registerLogoutButton } from "./functions.js";

registerLogoutButton();

const roomId = document.URL.split("/").pop();

const leaveRoomButton = document.getElementById("leave-room-button");
leaveRoomButton.onclick = async (event) => {
    event.preventDefault();
    await leaveRoom();
};

const messageBox = document.getElementById("message-box");
messageBox.onsubmit = (event) => {
    event.preventDefault();
    const formData = new FormData(messageBox);
    const body = formData.get("body");
    sendMessage(body);
    messageBox.reset();
};

const messages = document.getElementById("messages");

const messageTemplate = document.getElementById("message-template");

const ws = new WebSocket(wsURL());
ws.onopen = (event) => {
    console.log("onopen", event);
};
ws.onmessage = (event) => {
    /** @type Message */
    const message = JSON.parse(event.data);
    /** @type HTMLElement */
    const messageElement = messageTemplate.content.cloneNode(true);
    // prettier-ignore
    {
    messageElement.querySelector("#message-template-username").textContent = message.username;
    messageElement.querySelector("#message-template-body").textContent = message.body;
    messageElement.querySelector("#message-template-timestamp").textContent = new Date(message.timestamp).toLocaleString();
    }
    messages.appendChild(messageElement);
};
ws.onerror = (event) => {
    console.log("onerror", event);
};
ws.onclose = (event) => {
    console.log("onclose", event);
};

async function leaveRoom() {
    if (!roomId) {
        console.error("invalid roomId", roomId);
        return;
    }
    await fetch("/api/rooms/leave", {
        method: "POST",
        headers: {
            "content-type": "application/json",
        },
        body: JSON.stringify({
            roomId: roomId,
        }),
    });
    ws.close();
    window.location.replace("/home");
}

/**
 * @param {string} body
 */
function sendMessage(body) {
    if (!roomId) {
        console.error("invalid roomId", roomId);
        return;
    }
    if (body.length === 0) {
        console.error("empty body", body);
        return;
    }
    /** @type Message */
    const message = {
        roomId: roomId,
        body: body,
        timestamp: new Date().toISOString(),
    };
    ws.send(JSON.stringify(message));
}

function wsURL() {
    let scheme = "ws";
    if (document.location.protocol === "https:") {
        scheme += "s";
    }
    return `${scheme}://${document.location.hostname}:${document.location.port}/api/connect`;
}
