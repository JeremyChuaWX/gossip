"use strict";

const joinRoomForm = document.getElementById("join-room-form");
joinRoomForm.onsubmit = async (event) => {
    event.preventDefault();
    const formData = new FormData(joinRoomForm);
    const roomId = formData.get("room-id");
    try {
        await joinRoom(roomId);
    } catch {
        alert("Error creating room");
        return;
    }
    window.location.replace("/home");
};

/**
 * @param {string} roomId
 */
async function joinRoom(roomId) {
    await fetch("/rooms/join", {
        method: "POST",
        headers: {
            "content-type": "application/json",
        },
        body: JSON.stringify({
            roomId: roomId,
        }),
    });
}
