"use strict";

let signupForm = document.getElementById("signup-form");
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
    window.location.replace("/login");
};

/**
 * @param {string} username
 * @param {string} password
 */
async function signup(username, password) {
    await fetch("/signup", {
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
