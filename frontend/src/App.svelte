<script>
  import { onMount } from "svelte";
  import { Gauge, Menu, Moon, Server, Settings as SettingsIcon, Sun, UserCircle } from "lucide-svelte";
  import Dashboard from "./routes/Dashboard.svelte";
  import Landing from "./routes/Landing.svelte";
  import Servers from "./routes/Servers.svelte";
  import ServerDetail from "./routes/ServerDetail.svelte";
  import CreateServer from "./routes/CreateServer.svelte";
  import Profile from "./routes/Profile.svelte";
  import Settings from "./routes/Settings.svelte";
  import Auth from "./routes/Auth.svelte";
  import Toast from "./lib/components/Toast.svelte";
  import { api } from "./lib/api.js";
  import { session, toasts } from "./lib/stores.js";
  import { theme, cycleTheme } from "./lib/theme.js";

  let route = "#/home";
  let menuOpen = false;
  let sessionWarning = "";
  let sessionStatus = "checking";
  let unsubscribeTheme;

  onMount(() => {
    const syncRoute = () => { const nextRoute = window.location.hash || "#/home"; route = nextRoute; menuOpen = false; checkSession(nextRoute); };
    window.addEventListener("hashchange", syncRoute);
    syncRoute();
    unsubscribeTheme = theme.subscribe(() => {});
    return () => { window.removeEventListener("hashchange", syncRoute); unsubscribeTheme?.(); };
  });

  async function checkSession(routeValue = route) {
    const routeName = routeValue.replace(/^#\/?/, "").split("/")[0] || "home";
    if (routeName === "home") { sessionStatus = "ready"; return; }
    sessionStatus = "checking";
    session.set({ status: "checking", username: null });
    try {
      const response = await api.checkSession();
      if (response.status === 200) {
        sessionStatus = "authenticated";
        session.set({ status: "authenticated", username: null });
        if (routeName === "login" || routeName === "signup") window.location.hash = "#/dashboard";
      } else if (response.status === 401) {
        sessionStatus = "anonymous";
        sessionWarning = routeName === "signup" ? "No active session. Create an account to continue." : "No active session. Log in to continue.";
        session.set({ status: "anonymous", username: null });
        if (routeName === "dashboard" || routeName === "servers" || routeName === "profile" || routeName === "settings" || routeName === "create") window.location.hash = "#/login";
      } else { sessionStatus = "error"; sessionWarning = "Can't reach the server. Try again when the backend is running."; }
    } catch (error) { sessionStatus = "error"; sessionWarning = "Can't reach the server. Try again when the backend is running."; if (routeName === "dashboard" || routeName === "servers") session.set({ status: "error", username: null }); }
  }

  function navigate(next) { window.location.hash = next; }
  async function logout() { try { const response = await api.logout(); if (response.status === 200 || response.status === 401) navigate("#/login"); else sessionWarning = "The server could not end your session."; } catch (error) { sessionWarning = "Could not reach the server. You are still signed in."; } }
  $: routeParts = route.replace(/^#\/?/, "").split("/");
  $: current = routeParts[0] || "home";
  $: serverId = routeParts[1];
  $: authRoute = current === "login" || current === "signup";
  $: protectedRoute = !authRoute && current !== "home";
</script>

{#if sessionStatus === "checking"}<div class="loading-screen">Checking session…</div>{:else if current === "home"}<Landing {navigate}/>{:else if authRoute && sessionStatus !== "authenticated"}<Auth mode={current} {navigate} sessionWarning={sessionWarning} />{:else if sessionStatus === "error" && protectedRoute}<Auth mode="login" {navigate} sessionWarning={sessionWarning} />{:else}
  <div class="app"><header class="topbar"><a class="brand" href="#/dashboard"><span class="brand-mark">S</span><span>shmolph.cloud</span></a><div class="top-actions"><button class="utility-button menu-button" type="button" on:click={() => menuOpen = !menuOpen}><Menu size={17} />Menu</button><button class="utility-button" type="button" on:click={() => cycleTheme($theme)}>{#if $theme === "light"}<Sun size={16} />Light{:else if $theme === "dark"}<Moon size={16} />Dark{:else}<Sun size={16} />System{/if}</button><button class="utility-button" type="button" on:click={() => navigate("#/profile")}><UserCircle size={16} />Profile</button><button class="utility-button" type="button" on:click={() => navigate("#/settings")}><SettingsIcon size={16} />Settings</button><button class="utility-button" type="button" on:click={logout}>Log out</button></div></header><div class:open={menuOpen} class="layout"><aside class="sidebar"><p class="sidebar-title">Control panel</p><nav class="nav-list"><a class:active={current === "dashboard"} class="nav-item" href="#/dashboard"><Gauge size={17} />Dashboard</a><a class:active={current === "servers" || current === "create"} class="nav-item" href="#/servers"><Server size={17} />Servers</a><button class="nav-item" type="button" disabled><Server size={17} />Backups <span class="soon">Soon</span></button><button class="nav-item" type="button" disabled><span>$</span>Billing <span class="soon">Soon</span></button><a class:active={current === "settings"} class="nav-item" href="#/settings"><SettingsIcon size={17} />Settings</a></nav></aside><main class="content">{#if current === "dashboard"}<Dashboard {navigate}/>{:else if current === "servers" && serverId}<ServerDetail {serverId} {navigate}/>{:else if current === "servers"}<Servers {navigate}/>{:else if current === "create"}<CreateServer {navigate}/>{:else if current === "profile"}<Profile onLogout={logout}/>{:else if current === "settings"}<Settings />{:else}<Dashboard {navigate}/>{/if}</main></div></div>
{/if}

<Toast items={$toasts} />

<style>
  .loading-screen { display:grid; min-height:100vh; place-items:center; color:var(--muted); background:var(--bg); }
  .menu-button { display:none; }
  @media (max-width:800px) { .menu-button { display:inline-flex; } .layout { position:relative; } .layout.open .sidebar { transform:translateX(0); } .sidebar { position:fixed; z-index:5; top:64px; bottom:0; left:0; width:240px; transform:translateX(-105%); transition:transform 180ms ease; box-shadow:var(--shadow-strong); } }
</style>
