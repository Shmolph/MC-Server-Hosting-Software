# Shmolph.cloud

Shmolph.cloud is a low-cost Minecraft server hosting solution for people who want to play together without dealing with a complicated setup or an expensive hosting plan.

The goal is to make hosting feel approachable: choose a server, bring your friends, and spend your time building and playing instead of wrestling with infrastructure. It is intended for small friend groups, communities, and anyone who wants a simple home for a Minecraft world.

## Frontend

The frontend currently contains a standalone landing page in `index.html`, a signup page in `signup.html`, a login page in `login.html`, and a signed-in dashboard in `account.html`. The pages introduce the service and provide the initial account flows.

The pages use a single visual theme built around deep leaf green, soft cream, mint, and warm gold. They have responsive layouts for desktop and mobile screens, light entrance animations, and do not require a framework or build step.

The signup page loads Cloudflare Turnstile with the test site key `1x00000000000000000000AA`, then sends a JSON `username`, `password`, `captchaKey`, and `captchaToken` payload to `http://localhost:3030/register`. The page reports successful registration for HTTP `200`, a taken username or rejected signup for HTTP `409`, and an internal server problem for HTTP `500`. The backend must be running locally, and it must validate the Turnstile token before production use; the frontend alone cannot provide CAPTCHA security.

The signup form requires passwords to be 8 to 72 characters long and include at least one uppercase letter, one number, and one symbol.

The login page sends credentials to `http://localhost:3030/login` with `credentials: "include"`, so the browser can store the server's HttpOnly `session` cookie. A successful login redirects to `account.html`. The frontend never reads or stores that cookie, and it shows a generic message for HTTP `401` responses.

After a successful signup, the frontend requests a fresh Turnstile token and logs the new account in automatically before redirecting to `account.html`. This second CAPTCHA step is necessary because Turnstile tokens are single-use.

The signed-in dashboard currently provides the welcome screen, an honest empty server workspace, and a short setup checklist. Server creation and server management are not connected yet, so the dashboard does not display invented server status or usage data.

Before showing `account.html`, the page calls `GET http://localhost:3030/account` with `credentials: "include"`. A `200` reveals the dashboard, a `401` redirects to `login.html`, and a network or unexpected server error shows a session-check error without revealing the account content. The session cookie remains HttpOnly and is never read or stored by JavaScript.

## Run locally

Open `index.html` directly in a browser, or serve the `frontend` folder with any static web server.

For example:

```bash
cd frontend
python3 -m http.server 5000
```

Then visit `http://localhost:5000`.
