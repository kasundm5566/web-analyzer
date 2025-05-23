async function analyzeUrl() {
    const urlInput = document.getElementById("urlInput");
    const resultBox = document.getElementById("result");
    const url = urlInput.value.trim();

    resultBox.textContent = "🔍 Analyzing...";
    resultBox.style.color = "black";

    if (!url) {
        resultBox.textContent = "❌ Please enter a URL.";
        resultBox.style.color = "red";
        return;
    }

    try {
        const response = await fetch("http://localhost:8080/analyze-url", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({url})
        });

        // Handle non-200 responses
        if (!response.ok) {
            const errorText = await response.text();
            resultBox.textContent = `❌ Error ${response.status}: ${errorText}`;
            resultBox.style.color = "red";
            return;
        }

        const data = await response.json();
        displayResult(data);
    } catch (err) {
        resultBox.textContent = `❌ Network or server error: ${err.message}`;
        resultBox.style.color = "red";
    }
}

function displayResult(data) {
    const resultBox = document.getElementById("result");

    resultBox.textContent = "";

    resultBox.innerHTML = `
        <p><strong>HTML Version:</strong> ${data.htmlVersion}</p>
        <p><strong>Page Title:</strong> ${data.pageTitle}</p>
        <p><strong>Headings Count:</strong> ${data.headingsCount}</p>
        <p><strong>Internal Links Count:</strong> ${data.internalLinksCount}</p>
        <p><strong>External Links Count:</strong> ${data.externalLinksCount}</p>
        <p><strong>Inaccessible Links Count:</strong> ${data.inaccessibleLinksCount}</p>
        <p><strong>Inaccessible Links:</strong></p>
        <ul>
            ${data.inaccessibleLinks.map(link => `<li>${link}</li>`).join("")}
        </ul>
        <p><strong>Contains Login Form:</strong> ${data.containsLoginForm ? "Yes" : "No"}</p>
    `;
}

function login() {
    const usernameInput = document.getElementById("username");
    const passwordInput = document.getElementById("password");

    const usernameError = document.getElementById("usernameError");
    const passwordError = document.getElementById("passwordError");

    const username = usernameInput.value.trim();
    const password = passwordInput.value.trim();

    if (!username) {
        usernameError.style.visibility = "visible";
        usernameError.textContent = "❌ Please enter username.";
        usernameError.style.color = "red";
        return;
    } else {
        usernameError.style.visibility = "hidden";
    }
    if (!password) {
        passwordError.style.visibility = "visible";
        passwordError.textContent = "❌ Please enter password.";
        passwordError.style.color = "red";
        return;
    } else {
        passwordError.style.visibility = "hidden";
    }
}