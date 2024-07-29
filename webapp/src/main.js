"use strict";

const SERVER_URL = "127.0.0.1:3000";

let sessionId = "";

let loginForm = document.getElementById("login-form");
loginForm.onsubmit = (event) => {
    event.preventDefault();
    const formData = new FormData(loginForm);
    const username = formData.get("username");
    const password = formData.get("password");
    login(username, password);
};

/**
 * @param {string} username
 * @param {string} password
 */
async function login(username, password) {
    let res = await fetch(`${SERVER_URL}/login`, {
        method: "POST",
        headers: {
            "content-type": "application/json",
        },
        body: {
            username: username,
            password: password,
        },
    });
    const data = await res.json();
    sessionId = data.session.id;
}
