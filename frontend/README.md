# Shmolph.cloud

Shmolph.cloud is a low-cost Minecraft server hosting solution for people who want to play together without dealing with a complicated setup or an expensive hosting plan.

The goal is to make hosting feel approachable: choose a server, bring your friends, and spend your time building and playing instead of wrestling with infrastructure. It is intended for small friend groups, communities, and anyone who wants a simple home for a Minecraft world.

## Frontend

The frontend is a Svelte 5 + Vite app with hash routes for `#/dashboard`, `#/login`, and `#/register`. Hash routing keeps app navigation refresh-safe with the Go file server.

All backend calls live in `src/lib/api.js` and use `credentials: "include"`. Server features call their expected real endpoints and show `This feature isn't connected yet.` for unavailable APIs. No fake servers, players, logs, stats, or backups are rendered.

The theme supports Light, Dark, and System. The preference is stored in `localStorage` as UI state only; the HttpOnly session cookie is never read or stored by JavaScript. The API URL, public server hostname, Turnstile key, and server type options are centralized in `src/lib/config.js`. The public hostname is still `FILL_IN_PUBLIC_HOST` and must be replaced with the address players use to connect. Cloudflare Turnstile currently uses the test site key `1x00000000000000000000AA` on the auth routes.

To add a page, create a Svelte component in `src/routes/`, add its hash route to `src/Root.svelte`, and compose it from shared components in `src/lib/components/` and the global tokens in `src/global.css`.

## Run locally

Install dependencies:

```bash
cd frontend
npm install
```

Run Vite on the backend-approved origin:

```bash
npm run dev -- --port 5000 --strictPort
```

Open `http://localhost:5000`. Stop the Go static file server first; Vite and the Go file server cannot listen on port `5000` at the same time.

Build the production files:

```bash
npm run build
```

The default output folder is `frontend/dist`, which contains the built `index.html`. Set `VITE_OUTPUT_DIR` to change it, for example `VITE_OUTPUT_DIR=../backend/public npm run build`. Point the Go static file server at the resulting output folder.
