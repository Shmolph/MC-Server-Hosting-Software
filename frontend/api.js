const API_BASE_URL = "http://localhost:3030";

class ApiError extends Error {
  constructor(status, message) {
    super(message);
    this.status = status;
  }
}

async function request(path, options = {}) {
  let response;
  try {
    response = await fetch(`${API_BASE_URL}${path}`, { credentials: "include", ...options });
  } catch (error) {
    throw new ApiError(0, "This feature isn't connected yet.");
  }
  const responseText = await response.text();
  if (!response.ok) throw new ApiError(response.status, response.status === 404 || response.status === 501 ? "This feature isn't connected yet." : responseText || "This feature isn't connected yet.");
  if (!responseText) return null;
  try { return JSON.parse(responseText); } catch (error) { throw new ApiError(response.status, "This feature isn't connected yet."); }
}

async function statusRequest(path, options = {}) {
  try { return await fetch(`${API_BASE_URL}${path}`, { credentials: "include", ...options }); } catch (error) { throw new ApiError(0, "This feature isn't connected yet."); }
}

function unavailable() { throw new ApiError(501, "This feature isn't connected yet."); }

window.dashboardApi = {
  API_BASE_URL,
  register: (data) => statusRequest("/register", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(data) }),
  login: (data) => statusRequest("/login", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(data) }),
  checkSession: () => statusRequest("/account"),
  logout: () => statusRequest("/logout", { method: "POST" }),
  listServers: () => request("/servers"),
  createServer: (data) => request("/servers", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify(data) }),
  startServer: (id) => request(`/servers/${encodeURIComponent(id)}/start`, { method: "POST" }),
  stopServer: (id) => request(`/servers/${encodeURIComponent(id)}/stop`, { method: "POST" }),
  restartServer: (id) => request(`/servers/${encodeURIComponent(id)}/restart`, { method: "POST" }),
  killServer: (id) => request(`/servers/${encodeURIComponent(id)}/kill`, { method: "POST" }),
  deleteServer: (id) => request(`/servers/${encodeURIComponent(id)}`, { method: "DELETE" }),
  sendCommand: (id, command) => request(`/servers/${encodeURIComponent(id)}/console`, { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ command }) }),
  getStats: (id) => request(`/servers/${encodeURIComponent(id)}/stats`),
  listPlayers: (id) => request(`/servers/${encodeURIComponent(id)}/players`),
  listFiles: (id) => request(`/servers/${encodeURIComponent(id)}/files`),
  listBackups: (id) => request(`/servers/${encodeURIComponent(id)}/backups`),
  renameServer: (id, name) => request(`/servers/${encodeURIComponent(id)}`, { method: "PATCH", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ name }) }),
  updateServerSettings: (id, settings) => request(`/servers/${encodeURIComponent(id)}`, { method: "PATCH", headers: { "Content-Type": "application/json" }, body: JSON.stringify(settings) }),
  createBackup: (id) => request(`/servers/${encodeURIComponent(id)}/backups`, { method: "POST" }),
  restoreBackup: (id, backupId) => request(`/servers/${encodeURIComponent(id)}/backups/${encodeURIComponent(backupId)}/restore`, { method: "POST" }),
  deleteBackup: (id, backupId) => request(`/servers/${encodeURIComponent(id)}/backups/${encodeURIComponent(backupId)}`, { method: "DELETE" }),
  unavailable
};
