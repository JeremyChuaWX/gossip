"use strict";

const urlParams = new URLSearchParams(window.location.search);
const prev = urlParams.get("prev");

const signupLink = document.getElementById("signup-link");
if (prev) {
    signupLink.href = `/signup?prev=${prev}`;
}

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
    if (prev) {
        window.location.replace(prev);
    } else {
        window.location.replace("/home");
    }
};

/**
 * @param {string} username
 * @param {string} password
 */
async function login(username, password) {
    await fetch("/api/login", {
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
