"use strict";

const signupForm = document.getElementById("signup-form");
signupForm.onsubmit = async (event) => {
    event.preventDefault();
    const formData = new FormData(signupForm);
    const username = formData.get("username");
    const password = formData.get("password");
    try {
        await signup(username, password);
    } catch {
        alert("Error signing up");
        return;
    }
    const urlParams = new URLSearchParams(window.location.search);
    const prev = urlParams.get("prev");
    if (prev) {
        window.location.replace(`/login?prev=${prev}`);
    } else {
        window.location.replace("/login");
    }
};

/**
 * @param {string} username
 * @param {string} password
 */
async function signup(username, password) {
    await fetch("/api/signup", {
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
