"use strict";

import { registerLogoutButton } from "./functions.js";

registerLogoutButton();

const createRoomForm = document.getElementById("create-room-form");
createRoomForm.onsubmit = async (event) => {
    event.preventDefault();
    const formData = new FormData(createRoomForm);
    const roomName = formData.get("room-name");
    try {
        await createRoom(roomName);
    } catch {
        alert("Error creating room");
        return;
    }
    window.location.replace("/home");
};

/**
 * @param {string} roomName
 */
async function createRoom(roomName) {
    await fetch("/api/rooms/create", {
        method: "POST",
        headers: {
            "content-type": "application/json",
        },
        body: JSON.stringify({
            roomName: roomName,
        }),
    });
}
