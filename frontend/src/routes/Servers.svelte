<script>
  import { onMount } from "svelte";
  import { Search } from "lucide-svelte";
  import Badge from "../lib/components/Badge.svelte";
  import Button from "../lib/components/Button.svelte";
  import EmptyState from "../lib/components/EmptyState.svelte";
  import Skeleton from "../lib/components/Skeleton.svelte";
  import { api } from "../lib/api.js";

  export let navigate;
  let servers = [];
  let query = "";
  let status = "All";
  let loading = true;
  let error = "";
  $: filtered = servers.filter((server) => (!query || server.name.toLowerCase().includes(query.toLowerCase())) && (status === "All" || String(server.status).toLowerCase() === status));
  onMount(loadServers);
  async function loadServers() { loading = true; error = ""; try { servers = await api.listServers(); } catch (requestError) { error = requestError.message; } finally { loading = false; } }
</script>

<div class="page-head"><div><p class="kicker">Servers</p><h1>All servers</h1></div><Button on:click={() => navigate("#/create")}>+ Create server</Button></div>
<div class="toolbar"><label class="search"><Search size={16} /><input bind:value={query} placeholder="Search servers" aria-label="Search servers" /></label><select class="select" bind:value={status} aria-label="Filter by status"><option>All</option><option value="running">Running</option><option value="stopped">Stopped</option></select></div>
{#if loading}<Skeleton lines={4} />{:else if error}<div class="feature-error">{error}</div>{:else if !filtered.length}<EmptyState title="No servers yet, create your first one"><Button on:click={() => navigate("#/create")}>+ Create server</Button></EmptyState>{:else}<div class="server-list">{#each filtered as server}<div class="server-row"><button class="server-name" type="button" on:click={() => navigate(`#/servers/${server.id}`)}><span class:online={String(server.status).toLowerCase() === "running"} class="server-dot"></span><span class="server-title"><strong>{server.name}</strong><span>{server.type} · {server.version}</span></span></button><Badge status={server.status} /><div class="data-block"><strong>{(Number(server.ram_mb) / 1024).toLocaleString(undefined, { maximumFractionDigits: 2 })} GB</strong><span>Memory (RAM)</span></div><div class="data-block"><strong>{server.port}</strong><span>Port</span></div></div>{/each}</div>{/if}

<style>
  .toolbar { display: flex; gap: 8px; margin-bottom: 12px; }
  .search { display: flex; flex: 1; max-width: 360px; gap: 8px; align-items: center; padding: 0 11px; border: 1px solid var(--border); border-radius: 8px; color: var(--muted); background: var(--surface); }
  .search input { width: 100%; min-height: 40px; border: 0; outline: 0; color: var(--text); background: transparent; }
  .select { max-width: 140px; }
</style>
