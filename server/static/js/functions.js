export function registerLogoutButton() {
    const logoutButton = document.getElementById("logout-button");
    logoutButton.onclick = async (event) => {
        event.preventDefault();
        try {
            await logout();
        } catch {
            alert("Error logging out");
            return;
        }
        window.location.replace("/");
    };
}

async function logout() {
    await fetch("/api/logout", {
        method: "POST",
    });
}
