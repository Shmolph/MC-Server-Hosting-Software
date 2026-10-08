const dashboardApi = window.dashboardApi;
const accountContent = document.querySelector("#account-content");
const accountError = document.querySelector("#account-error");
const logoutButton = document.querySelector("#logout-button");
const menuToggle = document.querySelector("#menu-toggle");
const sidebar = document.querySelector(".sidebar");
const globalMessage = document.querySelector("#global-message");
const serverListView = document.querySelector("#server-list-view");
const serverDetailView = document.querySelector("#server-detail-view");
const dashboardView = document.querySelector("#dashboard-view");
const profileView = document.querySelector("#profile-view");
const settingsView = document.querySelector("#settings-view");
const serverList = document.querySelector("#server-list");
const dashboardServerList = document.querySelector("#dashboard-server-list");
const dashboardStats = document.querySelector("#dashboard-stats");
const serverSearch = document.querySelector("#server-search");
const statusFilter = document.querySelector("#status-filter");
const createDialog = document.querySelector("#create-dialog");
const createForm = document.querySelector("#create-form");
const createError = document.querySelector("#create-error");
const createRam = document.querySelector("#create-ram");
const createRamValue = document.querySelector("#create-ram-value");
const createName = document.querySelector("#create-name");
const createType = document.querySelector("#create-type");
const createStepLabels = document.querySelectorAll("[data-step-label]");
let createStep = 1;
const toast = document.querySelector("#toast");
const confirmDialog = document.querySelector("#confirm-dialog");
const confirmMessage = document.querySelector("#confirm-message");
const confirmAccept = document.querySelector("#confirm-accept");
const confirmCancel = document.querySelector("#confirm-cancel");
const themeToggle = document.querySelector("#theme-toggle");
let servers = [];
let selectedServerId = null;
let activeTab = "console";
let statsTimer = null;
const consoleLogs = new Map();

function setGlobalError(message = "") { globalMessage.textContent = message; }
function showToast(message) { toast.textContent = message; toast.hidden = false; setTimeout(() => { toast.hidden = true; }, 2400); }
function statusClass(status) { return status === "Online" ? "status-online" : status === "Offline" ? "status-offline" : status === "Error" ? "status-error" : "status-transitioning"; }
function dotClass(status) { return status === "Online" ? "online" : status === "Offline" ? "" : "transitioning"; }
function selectedServer() { return servers.find((server) => server.id === selectedServerId); }
function filteredServers() { const query = serverSearch.value.trim().toLowerCase(); const status = statusFilter.value; return servers.filter((server) => (!query || server.name.toLowerCase().includes(query)) && (status === "All" || server.status === status)); }
function setTheme(theme) { document.documentElement.dataset.theme = theme; localStorage.setItem("shmolph-theme", theme); themeToggle.textContent = `Theme: ${theme.charAt(0).toUpperCase()}${theme.slice(1)}`; document.querySelectorAll("[data-theme-choice]").forEach((button) => button.classList.toggle("active", button.dataset.themeChoice === theme)); }
function cycleTheme() { const themes = ["system", "light", "dark"]; setTheme(themes[(themes.indexOf(document.documentElement.dataset.theme) + 1) % themes.length]); }
function setPreference(name, value) { localStorage.setItem(`shmolph-${name}`, value); if (name === "compact") document.documentElement.classList.toggle("compact", value === "true"); if (name === "font-size") document.documentElement.dataset.consoleSize = value; }
function setView(viewName) { [dashboardView, serverListView, serverDetailView, profileView, settingsView].forEach((view) => { view.hidden = view.id !== `${viewName}-view`; }); document.querySelectorAll("[data-view]").forEach((button) => button.classList.toggle("active", button.dataset.view === viewName)); sidebar.classList.remove("open"); if (viewName === "dashboard") renderDashboard(); }
function confirmAction(message) { return new Promise((resolve) => { confirmMessage.textContent = message; confirmDialog.showModal(); const finish = (result) => { confirmDialog.close(); confirmAccept.removeEventListener("click", accept); confirmCancel.removeEventListener("click", cancel); resolve(result); }; const accept = () => finish(true); const cancel = () => finish(false); confirmAccept.addEventListener("click", accept); confirmCancel.addEventListener("click", cancel); }); }
function openCreateServer() { createStep = 1; createError.textContent = ""; createForm.reset(); createType.value = "Vanilla"; createRam.value = 4; createRamValue.textContent = "4 GB"; updateCreateStep(); createDialog.showModal(); }
function updateCreateStep() { document.querySelectorAll("[data-step]").forEach((step) => { step.hidden = Number(step.dataset.step) !== createStep; }); createStepLabels.forEach((label) => label.classList.toggle("active", Number(label.dataset.stepLabel) <= createStep)); document.querySelector("#create-back").hidden = createStep === 1; document.querySelector("#create-next").hidden = createStep === 3; document.querySelector("#submit-create").hidden = createStep !== 3; if (createStep === 3) { document.querySelector("#review-name").textContent = createName.value.trim(); document.querySelector("#review-type").textContent = createType.value; document.querySelector("#review-version").textContent = document.querySelector("#create-version").value; document.querySelector("#review-ram").textContent = `${createRam.value} GB`; } }
function validateCreateName() { const valid = createName.value.trim().length >= 3 && createName.value.trim().length <= 32; document.querySelector("#create-name-hint").textContent = valid ? "Name looks good" : "3 to 32 characters"; return valid; }

async function refreshServers() {
  servers = await dashboardApi.listServers();
  renderServerList();
  renderDashboard();
}

function renderServerList() {
  if (!servers.length) {
    serverList.innerHTML = '<div class="empty-state"><strong>No servers yet</strong><button class="primary-button" type="button" data-empty-create>Create server</button></div>';
    return;
  }

  serverList.innerHTML = filteredServers().map((server) => `<article class="server-row" data-server-id="${server.id}">
    <button class="server-name" type="button" data-open-server="${server.id}" aria-label="Open ${server.name}"><span class="server-dot ${dotClass(server.status)}"></span><span class="server-title"><strong>${server.name}</strong><span>${server.address}</span></span></button>
    <span class="status ${statusClass(server.status)}">${server.status}</span>
    <div class="data-block"><strong>${server.players}/${server.maxPlayers}</strong><span>Players</span></div>
    <div class="data-block"><strong>${server.version}</strong><span>${server.type}</span></div>
    <div class="data-block"><strong>${server.ramUsed} / ${server.ramAllocated} GB</strong><span>RAM</span></div>
    <div class="row-actions"><button class="small-button" type="button" data-server-action="start" data-server-id="${server.id}" ${server.status !== "Offline" ? "disabled" : ""}>Start</button><button class="small-button" type="button" data-server-action="stop" data-server-id="${server.id}" ${server.status !== "Online" ? "disabled" : ""}>Stop</button><button class="small-button" type="button" data-server-action="restart" data-server-id="${server.id}" ${server.status !== "Online" ? "disabled" : ""}>Restart</button></div>
  </article>`).join("");
}

function renderDashboard() {
  const onlineServers = servers.filter((server) => server.status === "Online").length;
  const players = servers.reduce((total, server) => total + Number(server.players || 0), 0);
  const memory = servers.reduce((total, server) => total + Number(server.ramUsed || 0), 0);
  dashboardStats.innerHTML = `<div class="stat"><span class="stat-label">Servers</span><strong>${servers.length}</strong></div><div class="stat"><span class="stat-label">Online now</span><strong>${onlineServers}</strong></div><div class="stat"><span class="stat-label">Players now</span><strong>${players}</strong></div><div class="stat"><span class="stat-label">Memory used</span><strong>${memory.toFixed(1)}<small> GB</small></strong></div>`;
  const preview = servers.slice(0, 5);
  dashboardServerList.innerHTML = preview.length ? preview.map((server) => `<article class="server-row dashboard-server-row"><button class="server-name" type="button" data-open-server="${server.id}"><span class="server-dot ${dotClass(server.status)}"></span><span class="server-title"><strong>${server.name}</strong><span>${server.version} · ${server.type}</span></span></button><span class="status ${statusClass(server.status)}">${server.status}</span><div class="data-block"><strong>${server.players}/${server.maxPlayers}</strong><span>Players</span></div><div class="row-actions"><button class="small-button" type="button" data-server-action="start" data-server-id="${server.id}" ${server.status !== "Offline" ? "disabled" : ""}>Start</button><button class="small-button" type="button" data-server-action="stop" data-server-id="${server.id}" ${server.status !== "Online" ? "disabled" : ""}>Stop</button></div></article>`).join("") : '<div class="empty-state"><strong>No servers yet</strong><button class="primary-button" type="button" data-empty-create>Create server</button></div>';
}

async function runServerAction(button, action, serverId) {
  button.disabled = true;
  const originalText = button.textContent;
  button.textContent = "...";
  setGlobalError();
  try {
    await dashboardApi[`${action}Server`](serverId);
    await refreshServers();
    if (selectedServerId) renderDetail();
  } catch (error) {
    setGlobalError(error.message);
  } finally {
    button.disabled = false;
    button.textContent = originalText;
  }
}

function showServerDetail(serverId) {
  selectedServerId = serverId;
  activeTab = "console";
  serverListView.hidden = true;
  serverDetailView.hidden = false;
  renderDetail();
}

function showServerList() {
  clearInterval(statsTimer);
  selectedServerId = null;
  serverDetailView.hidden = true;
  serverListView.hidden = false;
  refreshServers().catch((error) => setGlobalError(error.message));
}

function renderDetail() {
  const server = selectedServer();
  if (!server) return;
  document.querySelector("#detail-name").textContent = server.name;
  const detailStatus = document.querySelector("#detail-status");
  detailStatus.textContent = server.status;
  detailStatus.className = `status ${statusClass(server.status)}`;
  document.querySelectorAll("[data-detail-action]").forEach((button) => {
    const action = button.dataset.detailAction;
    button.disabled = action === "start" ? server.status !== "Offline" : action === "stop" || action === "restart" ? server.status !== "Online" : false;
  });
  document.querySelectorAll("[data-tab]").forEach((tab) => tab.classList.toggle("active", tab.dataset.tab === activeTab));
  renderTab();
  refreshStats();
  clearInterval(statsTimer);
  statsTimer = setInterval(refreshStats, 4000);
}

async function refreshStats() {
  if (!selectedServerId) return;
  try {
    const stats = await dashboardApi.getStats(selectedServerId);
    document.querySelector("#stats-strip").innerHTML = `<div class="stat"><span class="stat-label">CPU</span><strong>${stats.cpu}%</strong></div><div class="stat"><span class="stat-label">RAM</span><strong>${stats.ramUsed}<small> / ${stats.ramAllocated} GB</small></strong></div><div class="stat"><span class="stat-label">Players</span><strong>${stats.players}<small> / ${stats.maxPlayers}</small></strong></div>`;
  } catch (error) { setGlobalError(error.message); }
}

function renderTab() {
  const panel = document.querySelector("#tab-panel");
  if (activeTab === "console") renderConsole(panel);
  if (activeTab === "players") renderPlayers(panel);
  if (activeTab === "files") renderFiles(panel);
  if (activeTab === "backups") renderBackups(panel);
  if (activeTab === "settings") renderSettings(panel);
}

function renderConsole(panel) {
  const server = selectedServer();
  const logs = consoleLogs.get(server.id) || [];
  consoleLogs.set(server.id, logs);
  panel.innerHTML = `<pre id="console-log" class="console">${logs.join("\n")}</pre><form id="console-form" class="console-form"><input id="command-input" class="input" autocomplete="off" placeholder="Enter command"><button class="primary-button" type="submit">Send</button><button id="clear-console" class="secondary-button" type="button">Clear</button></form>`;
  const consoleLog = panel.querySelector("#console-log");
  consoleLog.scrollTop = consoleLog.scrollHeight;
  panel.querySelector("#console-form").addEventListener("submit", async (event) => {
    event.preventDefault();
    const commandInput = panel.querySelector("#command-input");
    const command = commandInput.value.trim();
    if (!command) return;
    event.submitter.disabled = true;
    try { logs.push(`> ${command}`, await dashboardApi.sendCommand(server.id, command)); renderConsole(panel); } catch (error) { setGlobalError(error.message); event.submitter.disabled = false; }
  });
  panel.querySelector("#clear-console").addEventListener("click", () => { consoleLogs.set(server.id, []); renderConsole(panel); });
}

async function renderPlayers(panel) {
  panel.innerHTML = '<div class="empty-state">Loading players...</div>';
  try {
    const players = await dashboardApi.listPlayers(selectedServerId);
    panel.innerHTML = players.length ? `<table class="table"><thead><tr><th>Player</th><th>Joined</th><th>Role</th><th></th></tr></thead><tbody>${players.map((player) => `<tr><td><strong>${player.name}</strong></td><td>${player.joined}</td><td>${player.op ? "Op" : "Player"}</td><td><div class="table-actions"><button class="small-button" type="button" data-player-action="op">Op</button><button class="small-button" type="button" data-player-action="kick">Kick</button><button class="small-button" type="button" data-player-action="ban">Ban</button></div></td></tr>`).join("")}</tbody></table>` : '<div class="empty-state">No players online</div>';
    panel.querySelectorAll("[data-player-action]").forEach((button) => button.addEventListener("click", () => showToast(`${button.dataset.playerAction} queued`)));
  } catch (error) { setGlobalError(error.message); }
}

async function renderFiles(panel) {
  panel.innerHTML = '<div class="empty-state">Loading files...</div>';
  try {
    const files = await dashboardApi.listFiles(selectedServerId);
    panel.innerHTML = `<div class="list">${files.map((file) => `<div class="list-row"><strong>${file.name}</strong><span>${file.size}</span><span>${file.modified}</span><button class="small-button" type="button" data-download-file>Download</button></div>`).join("")}</div>`;
    panel.querySelectorAll("[data-download-file]").forEach((button) => button.addEventListener("click", () => showToast("Download is not connected")));
  } catch (error) { setGlobalError(error.message); }
}

async function renderBackups(panel) {
  panel.innerHTML = '<div class="empty-state">Loading backups...</div>';
  try {
    const backups = await dashboardApi.listBackups(selectedServerId);
    panel.innerHTML = `<div class="backup-toolbar"><button id="create-backup" class="primary-button" type="button">+ Create backup</button></div>${backups.length ? `<div class="list">${backups.map((backup) => `<div class="list-row"><strong>${backup.name}</strong><span>${backup.size}</span><span>${backup.created}</span><div class="table-actions"><button class="small-button" type="button" data-restore="${backup.id}">Restore</button><button class="small-button" type="button" data-delete-backup="${backup.id}">Delete</button></div></div>`).join("")}</div>` : '<div class="empty-state">No backups yet</div>'}`;
    panel.querySelector("#create-backup").addEventListener("click", async (event) => { event.target.disabled = true; try { await dashboardApi.createBackup(selectedServerId); renderBackups(panel); } catch (error) { setGlobalError(error.message); event.target.disabled = false; } });
    panel.querySelectorAll("[data-restore]").forEach((button) => button.addEventListener("click", async () => { button.disabled = true; try { await dashboardApi.restoreBackup(selectedServerId, button.dataset.restore); showToast("Backup restored"); } catch (error) { setGlobalError(error.message); button.disabled = false; } }));
    panel.querySelectorAll("[data-delete-backup]").forEach((button) => button.addEventListener("click", async () => { button.disabled = true; try { await dashboardApi.deleteBackup(selectedServerId, button.dataset.deleteBackup); renderBackups(panel); } catch (error) { setGlobalError(error.message); button.disabled = false; } }));
  } catch (error) { setGlobalError(error.message); }
}

function renderSettings(panel) {
  const server = selectedServer();
  panel.innerHTML = `<form id="settings-form" class="settings-grid"><div class="setting setting-wide"><label for="settings-name">Server name</label><input id="settings-name" class="input" value="${server.name}" minlength="3" maxlength="32" required></div><div class="setting"><label for="settings-type">Type</label><select id="settings-type" class="select"><option ${server.type === "Vanilla" ? "selected" : ""}>Vanilla</option><option ${server.type === "Paper" ? "selected" : ""}>Paper</option><option ${server.type === "Fabric" ? "selected" : ""}>Fabric</option><option ${server.type === "Forge" ? "selected" : ""}>Forge</option></select></div><div class="setting"><label for="settings-version">Minecraft version</label><select id="settings-version" class="select"><option ${server.version === "1.21.4" ? "selected" : ""}>1.21.4</option><option ${server.version === "1.21.1" ? "selected" : ""}>1.21.1</option><option ${server.version === "1.20.6" ? "selected" : ""}>1.20.6</option><option ${server.version === "1.20.1" ? "selected" : ""}>1.20.1</option></select></div><div class="setting setting-wide"><label for="settings-ram">RAM</label><div class="range-row"><input id="settings-ram" type="range" min="1" max="16" value="${server.ramAllocated}"><span id="settings-ram-value" class="range-value">${server.ramAllocated} GB</span></div></div><div class="settings-actions"><span></span><button class="primary-button" type="submit">Save changes</button></div><div class="danger-zone"><strong>Delete server</strong><span>Type ${server.name} to confirm.</span><input id="delete-confirmation" class="input" placeholder="Server name"><button id="delete-server" class="danger-button danger-outline" type="button" disabled>Delete permanently</button></div></form>`;
  const settingsRam = panel.querySelector("#settings-ram");
  settingsRam.addEventListener("input", () => { panel.querySelector("#settings-ram-value").textContent = `${settingsRam.value} GB`; });
  panel.querySelector("#settings-form").addEventListener("submit", async (event) => { event.preventDefault(); event.submitter.disabled = true; try { await dashboardApi.renameServer(server.id, panel.querySelector("#settings-name").value.trim()); await dashboardApi.updateServerSettings(server.id, { type: panel.querySelector("#settings-type").value, version: panel.querySelector("#settings-version").value, ram: settingsRam.value }); await refreshServers(); renderDetail(); showToast("Settings saved"); } catch (error) { setGlobalError(error.message); event.submitter.disabled = false; } });
  const deleteConfirmation = panel.querySelector("#delete-confirmation");
  const deleteButton = panel.querySelector("#delete-server");
  deleteConfirmation.addEventListener("input", () => { deleteButton.disabled = deleteConfirmation.value !== server.name; });
  deleteButton.addEventListener("click", async () => { deleteButton.disabled = true; try { await dashboardApi.deleteServer(server.id); showServerList(); showToast("Server deleted"); } catch (error) { setGlobalError(error.message); deleteButton.disabled = false; } });
}

async function handleCreateServer(event) {
  event.preventDefault();
  const submit = document.querySelector("#submit-create");
  const name = document.querySelector("#create-name").value.trim();
  createError.textContent = "";
  if (name.length < 3 || name.length > 32) { createError.textContent = "Server name must be 3 to 32 characters."; return; }
  submit.disabled = true;
  try { await dashboardApi.createServer({ name, type: document.querySelector("#create-type").value, version: document.querySelector("#create-version").value, ram: createRam.value }); createDialog.close(); createForm.reset(); createRam.value = 4; createRamValue.textContent = "4 GB"; await refreshServers(); showToast("Server created"); } catch (error) { createError.textContent = error.message; } finally { submit.disabled = false; }
}

function initializeUi() {
  setTheme(document.documentElement.dataset.theme || "system");
  menuToggle.addEventListener("click", () => sidebar.classList.toggle("open"));
  setPreference("compact", localStorage.getItem("shmolph-compact") || "false");
  document.querySelector("[data-pref=compact]").checked = localStorage.getItem("shmolph-compact") === "true";
  document.querySelector("[data-pref=compact]").addEventListener("change", (event) => setPreference("compact", String(event.target.checked)));
  document.querySelectorAll("[data-font-size]").forEach((button) => button.addEventListener("click", () => { setPreference("font-size", button.dataset.fontSize); document.querySelectorAll("[data-font-size]").forEach((item) => item.classList.toggle("active", item === button)); }));
  themeToggle.addEventListener("click", cycleTheme);
  document.querySelectorAll("[data-theme-choice]").forEach((button) => button.addEventListener("click", () => setTheme(button.dataset.themeChoice)));
  document.querySelectorAll("[data-view]").forEach((button) => button.addEventListener("click", () => setView(button.dataset.view)));
  serverSearch.addEventListener("input", renderServerList);
  statusFilter.addEventListener("change", renderServerList);
  document.querySelector("[data-profile-logout]").addEventListener("click", () => logoutButton.click());
  document.querySelector("#create-server-button").addEventListener("click", openCreateServer);
  document.querySelector("#dashboard-create-button").addEventListener("click", openCreateServer);
  createName.addEventListener("input", validateCreateName);
  document.querySelectorAll("[data-type]").forEach((button) => button.addEventListener("click", () => { createType.value = button.dataset.type; document.querySelectorAll("[data-type]").forEach((item) => item.classList.toggle("active", item === button)); }));
  document.querySelectorAll("[data-ram]").forEach((button) => button.addEventListener("click", () => { createRam.value = button.dataset.ram; createRamValue.textContent = `${createRam.value} GB`; document.querySelectorAll("[data-ram]").forEach((item) => item.classList.toggle("active", item === button)); }));
  document.querySelector("#create-next").addEventListener("click", () => { if (createStep === 1 && !validateCreateName()) { createError.textContent = "Enter a server name from 3 to 32 characters."; return; } createError.textContent = ""; createStep += 1; updateCreateStep(); });
  document.querySelector("#create-back").addEventListener("click", () => { createError.textContent = ""; createStep -= 1; updateCreateStep(); });
  document.querySelector("#cancel-create").addEventListener("click", () => createDialog.close());
  document.querySelector("#close-create").addEventListener("click", () => createDialog.close());
  createRam.addEventListener("input", () => { createRamValue.textContent = `${createRam.value} GB`; });
  createForm.addEventListener("submit", handleCreateServer);
  document.querySelector("#back-to-servers").addEventListener("click", showServerList);
  document.querySelectorAll("[data-tab]").forEach((tab) => tab.addEventListener("click", () => { activeTab = tab.dataset.tab; renderTab(); document.querySelectorAll("[data-tab]").forEach((item) => item.classList.toggle("active", item === tab)); }));
  serverList.addEventListener("click", (event) => { const openButton = event.target.closest("[data-open-server]"); const actionButton = event.target.closest("[data-server-action]"); const emptyCreate = event.target.closest("[data-empty-create]"); if (openButton) showServerDetail(openButton.dataset.openServer); if (actionButton) runServerAction(actionButton, actionButton.dataset.serverAction, actionButton.dataset.serverId); if (emptyCreate) openCreateServer(); });
  dashboardServerList.addEventListener("click", (event) => { const openButton = event.target.closest("[data-open-server]"); const actionButton = event.target.closest("[data-server-action]"); const emptyCreate = event.target.closest("[data-empty-create]"); if (openButton) showServerDetail(openButton.dataset.openServer); if (actionButton) runServerAction(actionButton, actionButton.dataset.serverAction, actionButton.dataset.serverId); if (emptyCreate) openCreateServer(); });
  serverDetailView.addEventListener("click", async (event) => { const actionButton = event.target.closest("[data-detail-action]"); if (!actionButton) return; const action = actionButton.dataset.detailAction; if (action === "kill" && !(await confirmAction("Kill this server?"))) return; runServerAction(actionButton, action, selectedServerId); });
  logoutButton.addEventListener("click", async () => { logoutButton.disabled = true; try { const response = await dashboardApi.logout(); if (response.status === 200 || response.status === 401) { window.location.replace("login.html"); return; } setGlobalError("The server could not end your session."); } catch (error) { setGlobalError("Could not reach the server. You are still signed in."); } finally { logoutButton.disabled = false; } });
}

async function checkSessionAndStart() {
  let response;
  try {
    response = await dashboardApi.checkSession();
  } catch (error) { accountError.hidden = false; return; }
  if (response.status === 401) { window.location.replace("login.html"); return; }
  if (response.status !== 200) { accountError.hidden = false; return; }
  accountContent.hidden = false;
  initializeUi();
  try { await refreshServers(); } catch (error) { setGlobalError(error.message); }
}

checkSessionAndStart();
