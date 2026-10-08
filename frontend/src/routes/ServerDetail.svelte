<script>
  import { onMount, onDestroy } from "svelte";
  import { ArrowLeft, Ban, Copy, Download, FileText, Play, RotateCcw, Send, Square, Trash2, UserX } from "lucide-svelte";
  import Badge from "../lib/components/Badge.svelte";
  import Button from "../lib/components/Button.svelte";
  import EmptyState from "../lib/components/EmptyState.svelte";
  import Tabs from "../lib/components/Tabs.svelte";
  import { api } from "../lib/api.js";
  import { pushToast } from "../lib/stores.js";

  export let serverId;
  export let navigate;
  let server;
  let tab = "Console";
  let loading = true;
  let error = "";
  let stats = null;
  let logs = [];
  let players = [];
  let files = [];
  let backups = [];
  let command = "";
  let consoleElement;
  let autoScroll = true;
  let jumpVisible = false;
  let working = "";
  let statsTimer;
  const tabs = ["Console", "Players", "Files", "Backups", "Settings"];

  onMount(async () => { await load(); statsTimer = setInterval(loadStats, 4000); });
  onDestroy(() => clearInterval(statsTimer));
  async function load() { loading = true; try { const servers = await api.listServers(); server = servers.find((item) => String(item.id) === String(serverId)); if (!server) throw new Error("This feature isn't connected yet."); await loadStats(); } catch (requestError) { error = requestError.message; } finally { loading = false; } }
  async function loadStats() { if (!server) return; try { stats = await api.getStats(server.id); } catch (requestError) { stats = null; } }
  async function action(name) { working = name; try { await api[`${name}Server`](server.id); pushToast(`${name} request sent`, "success"); await load(); } catch (requestError) { pushToast(requestError.message, "error"); } finally { working = ""; } }
  async function loadTab(nextTab) { tab = nextTab; if (nextTab === "Players") players = await api.listPlayers(server.id).catch(() => []); if (nextTab === "Files") files = await api.listFiles(server.id).catch(() => []); if (nextTab === "Backups") backups = await api.listBackups(server.id).catch(() => []); }
  async function sendCommand() { if (!command.trim()) return; working = "command"; try { const response = await api.sendCommand(server.id, command.trim()); logs = [...logs, `> ${command.trim()}`, response]; command = ""; await tick(); scrollLatest(); } catch (requestError) { pushToast(requestError.message, "error"); } finally { working = ""; } }
  function scrollLatest() { if (consoleElement) consoleElement.scrollTop = consoleElement.scrollHeight; jumpVisible = false; }
  function handleScroll() { if (!consoleElement) return; jumpVisible = consoleElement.scrollHeight - consoleElement.scrollTop - consoleElement.clientHeight > 30; autoScroll = !jumpVisible; }
  async function copyConsole() { await navigator.clipboard?.writeText(logs.join("\n")); pushToast("Console copied", "success"); }
  let killOpen = false;
  function kill() { killOpen = true; }
  function confirmKill() { killOpen = false; action("kill"); }
  import { tick } from "svelte";
</script>

{#if loading}<h1>Loading server</h1>{:else if error}<div class="feature-error">{error}</div>{:else if server}<button class="back-link" type="button" on:click={() => navigate("#/servers")}><ArrowLeft size={15} /> All servers</button><div class="detail-head"><div class="detail-title"><div><p class="kicker">Server</p><h1>{server.name}</h1></div><Badge status={server.status} /></div><div class="detail-actions"><Button variant="secondary" loading={working === "start"} disabled={server.status !== "Offline"} on:click={() => action("start")}><Play size={15} />Start</Button><Button variant="secondary" loading={working === "stop"} disabled={server.status !== "Online"} on:click={() => action("stop")}><Square size={15} />Stop</Button><Button variant="secondary" loading={working === "restart"} disabled={server.status !== "Online"} on:click={() => action("restart")}><RotateCcw size={15} />Restart</Button><Button variant="danger" loading={working === "kill"} on:click={kill}>Kill</Button></div></div><div class="server-stats">{#each [["CPU", stats ? `${stats.cpu}%` : "—"], ["Memory", stats ? `${stats.ramUsed} / ${stats.ramAllocated} GB` : "—"], ["Players", stats ? `${stats.players} / ${stats.maxPlayers}` : "—"]] as stat}<div class="stat-card"><span class="stat-label">{stat[0]}</span><strong>{stat[1]}</strong></div>{/each}</div><Tabs {tabs} active={tab} onSelect={loadTab} />
{#if tab === "Console"}<div class="console-wrap"><pre bind:this={consoleElement} class="console" on:scroll={handleScroll}>{logs.join("\n")}</pre>{#if jumpVisible}<button class="jump" type="button" on:click={scrollLatest}>Jump to latest</button>{/if}<form class="console-form" on:submit|preventDefault={sendCommand}><input class="input" bind:value={command} placeholder="Enter command" aria-label="Console command"><Button variant="primary" loading={working === "command"}><Send size={15} />Send</Button><Button variant="secondary" type="button" on:click={() => logs = []}>Clear</Button><Button variant="secondary" type="button" on:click={copyConsole}><Copy size={15} />Copy</Button></form></div>{:else if tab === "Players"}<div class="list-panel">{#if !players.length}<EmptyState title="No players online." />{:else}{#each players as player}<div class="list-row"><strong>{player.name}</strong><span>{player.op ? "Op" : "Player"}</span><div><Button variant="secondary"><UserX size={14} />Kick</Button><Button variant="danger"><Ban size={14} />Ban</Button></div></div>{/each}{/if}</div>{:else if tab === "Files"}<div class="list-panel">{#if !files.length}<EmptyState title="This feature isn't connected yet." />{:else}{#each files as file}<div class="list-row"><strong><FileText size={14} />{file.name}</strong><span>{file.size} · {file.modified}</span><Button variant="secondary"><Download size={14} />Download</Button></div>{/each}{/if}</div>{:else if tab === "Backups"}<div class="list-panel"><Button variant="primary">Create backup</Button>{#if !backups.length}<EmptyState title="No backups yet." />{:else}{#each backups as backup}<div class="list-row"><strong>{backup.name}</strong><span>{backup.size} · {backup.created}</span><div><Button variant="secondary">Restore</Button><Button variant="danger"><Trash2 size={14} />Delete</Button></div></div>{/each}{/if}</div>{:else}<div class="settings-card"><label>Server name<input class="input" value={server.name}></label><label>Memory (RAM)<input class="input" type="range" min="1" max="16" value={server.ramAllocated}></label><label>Version<select class="select"><option>{server.version}</option></select></label><Button variant="primary">Save changes</Button><div class="danger-zone"><strong>Delete server</strong><span>Type the server name to confirm.</span><input class="input" placeholder={server.name}><Button variant="danger" disabled>Delete server</Button></div></div>{/if}{/if}

{#if killOpen}<div class="confirm-scrim" role="presentation" on:click={(event) => event.target === event.currentTarget && (killOpen = false)}><div class="confirm-card" role="dialog" aria-modal="true" aria-labelledby="kill-title" tabindex="-1"><h2 id="kill-title">Kill server?</h2><p>This will stop the server immediately.</p><div><Button variant="secondary" on:click={() => killOpen = false}>Cancel</Button><Button variant="danger" on:click={confirmKill}>Kill server</Button></div></div></div>{/if}

<style>
  .detail-title { display:flex; gap:12px; align-items:center; }
  .detail-actions { display:flex; flex-wrap:wrap; gap:8px; justify-content:flex-end; }
  .console-wrap { position:relative; }
  .console { overflow:auto; height:390px; margin:16px 0 8px; padding:16px; border:1px solid var(--border); border-radius:10px; color:var(--console-text); background:var(--console-bg); font:12px/1.7 ui-monospace, SFMono-Regular, Consolas, monospace; white-space:pre-wrap; }
  .console-form { display:flex; gap:8px; }
  .jump { position:absolute; right:16px; bottom:68px; min-height:34px; padding:0 11px; border:1px solid var(--accent); border-radius:999px; color:var(--text-on-accent); background:var(--accent); font-size:11px; font-weight:800; }
  .list-panel { display:grid; gap:8px; padding-top:16px; }
  .list-row { display:flex; min-height:54px; align-items:center; justify-content:space-between; gap:12px; padding:10px 12px; border:1px solid var(--border); border-radius:8px; background:var(--surface); }
  .list-row strong { display:flex; gap:7px; align-items:center; }
  .list-row span { color:var(--muted); font-size:12px; }
  .settings-card { display:grid; gap:16px; max-width:640px; padding:20px; border:1px solid var(--border); border-radius:10px; background:var(--surface); }
  .settings-card > label { display:grid; gap:7px; color:var(--muted); font-size:12px; font-weight:750; }
  .stat-label { display:block; margin-bottom:8px; color:var(--muted); font-size:11px; font-weight:750; }
  .confirm-scrim { position:fixed; z-index:10; inset:0; display:grid; place-items:center; padding:20px; background:var(--scrim); }
  .confirm-card { width:min(100%,420px); padding:22px; border:1px solid var(--border); border-radius:12px; background:var(--surface); box-shadow:var(--shadow-strong); }
  .confirm-card p { color:var(--muted); }
  .confirm-card > div { display:flex; justify-content:flex-end; gap:8px; }
  .danger-zone { display:grid; gap:8px; padding:16px; border:1px solid var(--danger); border-radius:10px; background:var(--danger-soft); }
  .danger-zone span { color:var(--muted); font-size:12px; }
</style>
