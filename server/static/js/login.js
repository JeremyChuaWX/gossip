"use strict";

const loginForm = document.getElementById("login-form");
loginForm.onsubmit = async (event) => {
    event.preventDefault();
    const formData = new FormData(loginForm);
    const username = formData.get("username");
    const password = formData.get("password");
    try {
        await login(username, password);
    } catch {
        alert("Error logging in");
        return;
    }
    window.location.replace("/home");
};

/**
 * @param {string} username
 * @param {string} password
 */
async function login(username, password) {
    await fetch("/login", {
        method: "POST",
        headers: {
            "content-type": "application/json",
        },
        body: JSON.stringify({
            username: username,
            password: password,
        }),
    });
}
