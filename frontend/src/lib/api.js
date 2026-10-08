import { API_BASE_URL } from "./config.js";
import { clearSessionState } from "./stores.js";

export class ApiError extends Error {
  constructor(status, message) { super(message); this.status = status; }
}

function expireSession() {
  clearSessionState();
  window.dispatchEvent(new Event("session:expired"));
}

async function fetchResponse(path, options = {}) {
  try {
    return await fetch(`${API_BASE_URL}${path}`, { credentials: "include", ...options });
  } catch (error) {
    throw new ApiError(0, "Can't reach the server.");
  }
}

async function readError(response) {
  const text = await response.text();
  if (response.status === 401) expireSession();
  throw new ApiError(response.status, text || `Request failed (${response.status}).`);
}

async function authPost(path, payload) {
  const response = await fetchResponse(path, { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(payload) });
  return { status: response.status, text: await response.text() };
}

async function jsonGet(path, publicEndpoint = false) {
  const response = await fetchResponse(path);
  if (!response.ok) {
    if (response.status === 401 && !publicEndpoint) return readError(response);
    return readError(response);
  }
  try { return await response.json(); }
  catch (error) { throw new ApiError(response.status, "The server returned invalid JSON."); }
}

export const api = {
  checkSession: () => fetchResponse("/account"),
  login: (data) => authPost("/login", data),
  register: (data) => authPost("/register", data),
  logout: async () => {
    const response = await fetchResponse("/logout", { method: "POST" });
    return { status: response.status, text: await response.text() };
  },
  listVersions: () => jsonGet("/versions", true),
  listServers: () => jsonGet("/servers"),
  createServer: async (data) => {
    const response = await fetchResponse("/servers/create", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(data) });
    if (!response.ok) return readError(response);
    return (await response.text()).trim();
  }
};
