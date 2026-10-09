<script>
  import { onMount } from "svelte";
  import { ArrowLeft, RefreshCw } from "lucide-svelte";
  import Badge from "../lib/components/Badge.svelte";
  import Button from "../lib/components/Button.svelte";
  import Card from "../lib/components/Card.svelte";
  import CopyButton from "../lib/components/CopyButton.svelte";
  import PageContainer from "../lib/components/PageContainer.svelte";
  import Skeleton from "../lib/components/Skeleton.svelte";
  import ServerTools from "./ServerTools.svelte";

  export let serverId;
  export let servers = [];
  export let busyActions = {};
  export let onAction = async () => ({ status: 0 });
  export let onRefresh = async () => null;
  export let navigate = () => {};

  let server = null;
  let loading = true;
  let error = "";
  let confirming = "";
  let actionNotFound = false;
  let pollTimer;
  let pollingReady = false;
  let requestInFlight = false;
  $: working = server ? busyActions[server.id] : "";
  $: if (!actionNotFound && !loading && Array.isArray(servers)) {
    const latest = servers.find((item) => String(item.id) === String(serverId));
    if (latest) setServer(latest);
    else if (!servers.length) setServer(null);
  }

  onMount(() => {
    pollingReady = true;
    schedulePolling(server?.status);
    refreshServer();
    const onVisibilityChange = () => {
      if (document.visibilityState === "visible") refreshServer(true);
    };
    document.addEventListener("visibilitychange", onVisibilityChange);
    return () => {
      clearInterval(pollTimer);
      document.removeEventListener("visibilitychange", onVisibilityChange);
    };
  });

  function schedulePolling(status) {
    clearInterval(pollTimer);
    pollTimer = setInterval(() => {
      if (document.visibilityState === "visible") refreshServer(true);
    }, status === "starting" ? 2500 : 5000);
  }

  function setServer(nextServer) {
    const previousStatus = server?.status;
    server = nextServer;
    if (pollingReady && previousStatus !== nextServer?.status) schedulePolling(nextServer?.status);
  }

  async function refreshServer(silent = false) {
    if (requestInFlight) return;
    requestInFlight = true;
    if (!silent) { loading = true; error = ""; }
    try {
      const latestServers = await onRefresh();
      if (!Array.isArray(latestServers)) throw new Error("Can't reach the server. Try again.");
      if (!actionNotFound) setServer(latestServers.find((item) => String(item.id) === String(serverId)) || null);
    } catch (requestError) {
      if (!silent) error = requestError.message;
    } finally {
      requestInFlight = false;
      if (!silent) loading = false;
    }
  }

  async function executeAction(action) {
    confirming = "";
    if (!server || working) return;
    const result = await onAction(server, action);
    if (result.status === 404) {
      actionNotFound = true;
      setServer(null);
    } else if (Array.isArray(result.servers)) {
      setServer(result.servers.find((item) => String(item.id) === String(serverId)) || server);
    }
  }

  function requestAction(action) {
    if (action === "stop" || action === "restart" || action === "delete") confirming = action;
    else executeAction(action);
  }
</script>

<PageContainer>
  {#if loading}
    <div class="detail-skeleton" aria-label="Loading server"><Skeleton lines={4} /></div>
  {:else if error}
    <Card variant="error" className="detail-error"><h1>Server couldn’t load</h1><p>{error}</p><Button variant="secondary" on:click={() => refreshServer()}><RefreshCw size={16} />Retry</Button></Card>
  {:else if !server}
    <section class="not-found"><p class="eyebrow">SERVER</p><h1>Server not found</h1><a href="#/dashboard"><ArrowLeft size={16} />Back to dashboard</a></section>
  {:else}
    <a class="back-link" href="#/dashboard" on:click|preventDefault={() => navigate("#/dashboard")}><ArrowLeft size={16} />Dashboard</a>
    <header class="detail-heading"><div><p class="eyebrow">SERVER</p><h1>{server.name}</h1></div><div class="status-group"><Badge status={server.status} />{#if server.status === "starting"}<p class="starting-help">Server is loading. Players can join once it shows Online.</p>{/if}</div></header>

    <Card className="info-panel">
      <h2>Server information</h2>
      <dl class="server-info">
        <div><dt>Type</dt><dd>{server.type}</dd></div>
        <div><dt>Version</dt><dd>{server.version}</dd></div>
        <div><dt>Memory (RAM)</dt><dd>{Number(server.ram_mb).toLocaleString()} MB</dd></div>
        <div><dt>Port</dt><dd>{server.port}</dd></div>
        <div class="connect-address"><dt>Connect address</dt><dd><code>localhost:{server.port}</code><CopyButton value={`localhost:${server.port}`} /></dd></div>
      </dl>
    </Card>

    <section class="controls">
      <h2>Power controls</h2>
      {#if server.status === "not setup"}
        <div class="control-actions"><Button variant="primary" loading={working === "setup"} disabled={Boolean(working)} on:click={() => requestAction("setup")}>{working === "setup" ? "Setting up..." : "Set up"}</Button><Button variant="danger" loading={working === "delete"} disabled={Boolean(working)} on:click={() => requestAction("delete")}>{working === "delete" ? "Deleting..." : "Delete"}</Button></div>
      {:else if server.status === "offline"}
        <div class="control-actions"><Button variant="primary" loading={working === "start"} disabled={Boolean(working)} on:click={() => requestAction("start")}>{working === "start" ? "Starting..." : "Start"}</Button><Button variant="danger" loading={working === "delete"} disabled={Boolean(working)} on:click={() => requestAction("delete")}>{working === "delete" ? "Deleting..." : "Delete"}</Button></div>
      {:else if server.status === "starting" || server.status === "online"}
        <div class="control-actions"><Button variant="secondary" loading={working === "restart"} disabled={Boolean(working)} on:click={() => requestAction("restart")}>{working === "restart" ? "Restarting..." : "Restart"}</Button><Button variant="danger" loading={working === "stop"} disabled={Boolean(working)} on:click={() => requestAction("stop")}>{working === "stop" ? "Stopping..." : "Stop"}</Button><Button variant="danger" disabled={true}>Delete</Button></div>
        <p class="disabled-hint">Stop the server first</p>
      {/if}
      {#if confirming}
        <div class="confirm-row" role="group" aria-label="Confirm power action"><p>{confirming === "delete" ? "This permanently deletes the server and its world. This cannot be undone." : "Players will be disconnected."}</p><Button variant="ghost" disabled={Boolean(working)} on:click={() => confirming = ""}>Cancel</Button><Button variant={confirming === "restart" ? "secondary" : "danger"} disabled={Boolean(working)} on:click={() => executeAction(confirming)}>{confirming === "delete" ? "Delete server" : `Confirm ${confirming}`}</Button></div>
      {/if}
    </section>
    <ServerTools {server} onRefreshStatus={() => refreshServer(true)} />
  {/if}
</PageContainer>

<style>
  .detail-skeleton { max-width:780px; margin:var(--space-6) auto; }
  :global(.detail-error) { display:grid; justify-items:start; gap:var(--space-3); padding:var(--space-5); }
  :global(.detail-error h1),:global(.detail-error p) { margin:0; }
  :global(.detail-error p) { color:var(--muted); }
  .back-link,.not-found a { display:inline-flex; gap:var(--space-2); align-items:center; color:var(--muted); font-size:var(--font-14); text-decoration:none; }
  .back-link:hover,.not-found a:hover { color:var(--accent); }
  .detail-heading { display:flex; gap:var(--space-4); align-items:center; justify-content:space-between; margin:var(--space-5) 0 var(--space-6); }
  .status-group { display:grid; justify-items:end; gap:var(--space-2); }
  .starting-help { max-width:310px; margin:0; color:var(--muted); font-size:var(--font-12); line-height:1.5; text-align:right; }
  .eyebrow { margin:0 0 var(--space-2); color:var(--accent); font-family:var(--pixel-font); font-size:10px; }
  h1 { margin:0; color:var(--text); font-size:var(--font-32); }
  h2 { margin:0; color:var(--text); font-size:var(--font-18); }
  :global(.info-panel) { display:grid; gap:var(--space-5); padding:var(--space-5); }
  .server-info { display:grid; grid-template-columns:repeat(2,minmax(0,1fr)); gap:var(--space-4); margin:0; }
  .server-info > div { display:grid; gap:var(--space-2); min-width:0; }
  dt { color:var(--muted); font-size:var(--font-12); }
  dd { display:flex; min-width:0; gap:var(--space-2); align-items:center; margin:0; color:var(--text); font-size:var(--font-14); }
  .connect-address { grid-column:1 / -1; }
  code { font-family:var(--mono-font); }
  .controls { display:grid; gap:var(--space-3); margin-top:var(--space-6); }
  .control-actions { display:flex; flex-wrap:wrap; gap:var(--space-2); }
  .disabled-hint { margin:calc(-1 * var(--space-2)) 0 0; color:var(--muted); font-size:var(--font-12); }
  .confirm-row { display:flex; flex-wrap:wrap; gap:var(--space-2); align-items:center; margin-top:var(--space-2); }
  .confirm-row p { flex:1 1 100%; margin:0; color:var(--muted); font-size:var(--font-13); }
  .not-found { display:grid; justify-items:start; gap:var(--space-3); max-width:540px; margin:var(--space-8) auto; }
  @media(max-width:520px) { h1 { font-size:var(--font-24); } .server-info { grid-template-columns:1fr; } .connect-address { grid-column:auto; } .connect-address dd { flex-wrap:wrap; } .control-actions { display:grid; grid-template-columns:1fr; } .starting-help { max-width:220px; } }
</style>
