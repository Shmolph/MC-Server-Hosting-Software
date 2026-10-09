<script>
  import { onMount } from "svelte";
  import { fade, slide } from "svelte/transition";
  import Navbar from "./lib/components/Navbar.svelte";
  import Toast from "./lib/components/Toast.svelte";
  import Dashboard from "./routes/Dashboard.svelte";
  import ServerPage from "./routes/ServerPage.svelte";
  import AuthPage from "./routes/AuthPage.svelte";
  import CreateServerModal from "./routes/CreateServerModal.svelte";
  import NotFound from "./routes/NotFound.svelte";
  import { api } from "./lib/api.js";
  import { clearSessionState, session, toasts, pushToast } from "./lib/stores.js";
  import { initializeTheme } from "./lib/theme.js";

  let route = "#/dashboard";
  let sessionState = "checking";
  let warning = "";
  let servers = [];
  let serversLoading = true;
  let serversError = "";
  let createOpen = false;
  let busyActions = {};

  const currentRoute = () => (window.location.hash || "#/dashboard").replace(/^#\/?/, "");
  const path = (value) => value.split("/")[0] || "dashboard";
  $: activePage = path(route);
  $: routeParts = route.replace(/^#\/?/, "").split("/");
  $: serverId = routeParts[1] || "";
  $: authMode = activePage === "register" ? "register" : "login";
  $: isAuthPage = activePage === "login" || activePage === "register";

  onMount(() => {
    initializeTheme();
    const onRoute = () => { route = currentRoute(); checkSession(); };
    const onExpired = () => {
      clearSessionState();
      servers = [];
      sessionState = "anonymous";
      route = "#/login";
      window.location.hash = route;
    };
    window.addEventListener("hashchange", onRoute);
    window.addEventListener("session:expired", onExpired);
    route = currentRoute();
    checkSession();
    return () => {
      window.removeEventListener("hashchange", onRoute);
      window.removeEventListener("session:expired", onExpired);
    };
  });

  async function checkSession() {
    sessionState = "checking";
    warning = "";
    try {
      const response = await api.checkSession();
      if (response.status === 200) {
        sessionState = "authenticated";
        session.set({ status: "authenticated", username: null });
        if (isAuthPage) { route = "#/dashboard"; window.location.hash = route; return; }
        await refreshServers();
      } else if (response.status === 401) {
        clearSessionState();
        sessionState = "anonymous";
        servers = [];
        if (!isAuthPage) { route = "#/login"; window.location.hash = route; }
      } else {
        sessionState = isAuthPage ? "anonymous" : "error";
        warning = "Can't reach the server. Check your connection and try again.";
      }
    } catch (error) {
      sessionState = isAuthPage ? "anonymous" : "error";
      warning = "Can't reach the server. Check your connection and try again.";
    }
  }

  async function refreshServers({ quiet = false } = {}) {
    if (!quiet) { serversLoading = true; serversError = ""; }
    try {
      servers = await api.listServers();
      serversError = "";
      return servers;
    } catch (error) {
      if (!quiet) serversError = error.message;
      return null;
    } finally {
      if (!quiet) serversLoading = false;
    }
  }

  async function performServerAction(server, action) {
    const id = Number(server.id);
    if (busyActions[id]) return { status: 0, servers };
    busyActions = { ...busyActions, [id]: action };
    try {
      const status = action === "setup" ? await api.setupServer(id) : action === "delete" ? await api.deleteServer(id) : await api.powerServer(id, action);
      if (status === 200) {
        const messages = { setup: "Server set up", start: "Server started", stop: "Server stopped", restart: "Server restarting", delete: "Server deleted" };
        pushToast(messages[action], "success");
      } else if (status === 404) {
        pushToast("Server not found", "error");
      } else if (status === 409 && action === "delete") {
        pushToast("Stop the server first", "error");
      } else if (status !== 401) {
        pushToast("Something went wrong. Try again.", "error");
      }
      const updatedServers = status === 401 ? servers : await refreshServers({ quiet: true });
      if (action === "delete" && status === 200) navigate("#/dashboard");
      return { status, servers: updatedServers };
    } catch (error) {
      pushToast("Cannot reach the backend", "error");
      return { status: 0, servers };
    } finally {
      const next = { ...busyActions };
      delete next[id];
      busyActions = next;
    }
  }

  function refreshQuietly() { return refreshServers({ quiet: true }); }

  async function logout() {
    try {
      const response = await api.logout();
      if (response.status === 200 || response.status === 401) {
        clearSessionState();
        servers = [];
        sessionState = "anonymous";
        route = "#/login";
        window.location.hash = route;
      } else {
        pushToast(response.text || "Logout failed. Try again.", "error");
      }
    } catch (error) {
      pushToast(error.message, "error");
    }
  }

  function navigate(pathname) { route = pathname; window.location.hash = pathname; }
  async function onCreated() { createOpen = false; await refreshServers(); }
</script>

{#if sessionState === "checking"}
  <div class="startup"><span class="startup-mark">S</span><span>Checking session</span></div>
{:else if isAuthPage && sessionState !== "authenticated"}
  <AuthPage mode={authMode} {warning} {navigate} />
{:else if sessionState === "error"}
  <AuthPage mode="login" {warning} {navigate} />
{:else}
  <div class="app-shell" in:fade={{ duration: 160 }}>
    <Navbar onLogout={logout} />
    <div class="route-wrap" in:slide={{ duration: 160, y: 5 }}>
      {#if activePage === "dashboard"}
        <Dashboard {servers} {busyActions} {performServerAction} loading={serversLoading} error={serversError} onRetry={refreshServers} onCreate={() => createOpen = true} />
      {:else if activePage === "servers" && serverId}
        <ServerPage {serverId} {servers} {busyActions} onAction={performServerAction} onRefresh={refreshQuietly} {navigate} />
      {:else}
        <NotFound navigate={navigate} />
      {/if}
    </div>
    {#if createOpen}
      <CreateServerModal onClose={() => createOpen = false} onCreated={onCreated} />
    {/if}
  </div>
{/if}
<Toast items={$toasts} />

<style>
  .startup { display:grid; min-height:100vh; align-content:center; justify-items:center; gap:var(--space-3); color:var(--muted); background:var(--bg); }
  .startup-mark { display:grid; width:44px; height:44px; place-items:center; border-radius:12px; color:var(--text-on-accent); background:var(--accent); box-shadow:var(--accent-glow); font-family:var(--pixel-font); }
  .app-shell { min-height:100vh; }
  .route-wrap { min-height:calc(100vh - var(--navbar-height)); }
</style>
