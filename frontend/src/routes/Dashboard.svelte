<script>
    import { Plus, RefreshCw, Server, Sparkles } from "lucide-svelte";
    import { fade } from "svelte/transition";
    import Badge from "../lib/components/Badge.svelte";
    import Button from "../lib/components/Button.svelte";
    import Card from "../lib/components/Card.svelte";
    import EmptyState from "../lib/components/EmptyState.svelte";
    import Skeleton from "../lib/components/Skeleton.svelte";
    import CopyButton from "../lib/components/CopyButton.svelte";
    import TypeMark from "../lib/components/TypeMark.svelte";
    import PageContainer from "../lib/components/PageContainer.svelte";
    import { PUBLIC_SERVER_HOST } from "../lib/config.js";

    export let servers = [];
    export let loading = true;
    export let error = "";
    export let onRetry = () => {};
    export let onCreate = () => {};
    $: publicHostReady = PUBLIC_SERVER_HOST && !PUBLIC_SERVER_HOST.startsWith("FILL_IN");
  </script>

  <PageContainer>
    <header class="page-heading">
      <div><span class="eyebrow"><Sparkles size={14} />SERVER CONTROL</span><h1>Your Servers</h1><p>All your worlds, in one place.</p></div>
      <div class="heading-actions"><span class="count">{servers.length}<small>{servers.length === 1 ? "server" : "servers"}</small></span><Button on:click={onCreate}><Plus size={17} />Create Server</Button></div>
    </header>
    {#if loading}
      <div class="server-grid" aria-label="Loading servers">{#each [1, 2, 3] as item (item)}<Skeleton lines={4} />{/each}</div>
    {:else if error}
      <Card variant="error" className="error-card"><div class="error-content"><span class="error-icon">!</span><div><h2>Servers couldn’t load</h2><p>{error}</p></div><Button variant="secondary" on:click={onRetry}><RefreshCw size={16} />Retry</Button></div></Card>
    {:else if !servers.length}
      <EmptyState title="No servers yet"><svg class="empty-art" viewBox="0 0 132 100" role="img" aria-label="Empty server blocks"><path d="M15 36 43 22l28 14v32L43 82 15 68V36Z"/><path d="m43 22 28 14-28 15-28-15 28-14Z"/><path d="M43 51v31"/><path d="m72 43 21-11 23 12v27L93 83 72 71V43Z"/><path d="m93 32 23 12-23 13-21-14 21-11Z"/><path d="M93 57v26"/></svg><span>A fresh world is one step away.</span><Button on:click={onCreate}><Plus size={18} />Create your first server</Button></EmptyState>
    {:else}
      <div class="server-grid">{#each servers as server (server.id)}<div class="server-entry" in:fade={{ duration: 180 }}><Card variant="server" className="server-card">
        <div class="card-top"><div class="title-group"><span class="server-symbol"><Server size={18} /></span><h2>{server.name}</h2></div><Badge status={server.status} /></div>
        <TypeMark type={server.type} />
        <div class="server-facts"><div><span>VERSION</span><strong>{server.version}</strong></div><div><span>MEMORY</span><strong>{(Number(server.ram_mb) / 1024).toLocaleString(undefined, { maximumFractionDigits: 2 })} GB</strong></div></div>
        <div class="address-line"><div class="address-copy"><span>CONNECT ADDRESS</span>{#if publicHostReady}<code>{PUBLIC_SERVER_HOST}:{server.port}</code>{:else}<span class="host-notice">Set public host to enable address</span>{/if}</div>{#if publicHostReady}<CopyButton value={`${PUBLIC_SERVER_HOST}:${server.port}`} />{/if}</div>
        <footer><span>Server controls</span><span class="future-dot" aria-hidden="true"></span></footer>
      </Card></div>{/each}</div>
    {/if}
  </PageContainer>

  <style>
    .page-heading { display:flex; align-items:flex-end; justify-content:space-between; gap:var(--space-5); margin-bottom:var(--space-7); }
    .eyebrow { display:flex; gap:var(--space-2); align-items:center; margin-bottom:var(--space-3); color:var(--accent); font-family:var(--pixel-font); font-size:10px; letter-spacing:1px; }
    h1 { margin:0; color:var(--text); font-size:var(--font-32); font-weight:720; letter-spacing:-.035em; }
    .page-heading p { margin:var(--space-2) 0 0; color:var(--muted); font-size:var(--font-14); }
    .heading-actions { display:flex; gap:var(--space-4); align-items:center; }
    .count { display:flex; gap:var(--space-2); align-items:baseline; color:var(--text); font-size:var(--font-24); font-variant-numeric:tabular-nums; font-weight:700; }
    .count small { color:var(--muted); font-size:var(--font-12); font-weight:550; }
    .server-grid { display:grid; grid-template-columns:repeat(2,minmax(0,1fr)); gap:var(--space-4); }
    .server-entry { min-width:0; }
    :global(.server-card) { display:grid; gap:var(--space-4); min-width:0; padding:var(--space-5); }
    .card-top { display:flex; gap:var(--space-3); align-items:center; justify-content:space-between; }
    .title-group { display:flex; min-width:0; gap:var(--space-3); align-items:center; }
    .server-symbol { display:grid; width:40px; height:40px; flex:0 0 40px; place-items:center; border:1px solid var(--border); border-radius:var(--radius-control); color:var(--accent); background:var(--surface-2); }
    h2 { overflow:hidden; margin:0; color:var(--text); font-size:var(--font-20); font-weight:700; text-overflow:ellipsis; white-space:nowrap; }
    .server-facts { display:grid; grid-template-columns:1fr 1fr; gap:var(--space-3); }
    .server-facts div,.address-copy { display:grid; gap:var(--space-2); }
    .server-facts span,.address-copy > span { color:var(--dim); font-family:var(--pixel-font); font-size:9px; letter-spacing:.6px; }
    .server-facts strong { color:var(--text); font-size:var(--font-14); font-variant-numeric:tabular-nums; }
    .address-line { display:flex; min-width:0; gap:var(--space-2); align-items:center; justify-content:space-between; padding:var(--space-3); border:1px solid var(--border); border-radius:var(--radius-control); background:var(--surface-2); }
    .address-copy { min-width:0; }
    code { overflow:hidden; color:var(--text); font-family:var(--mono-font); font-size:var(--font-12); text-overflow:ellipsis; white-space:nowrap; }
    :global(.server-card footer) { display:flex; min-height:var(--space-6); gap:var(--space-2); align-items:center; padding-top:var(--space-3); border-top:1px solid var(--border); color:var(--dim); font-size:var(--font-12); }
    .future-dot { width:6px; height:6px; border-radius:var(--radius-pill); background:var(--border); }
    :global(.server-grid > .skeleton) { min-height:260px; }
    :global(.error-card) { padding:var(--space-5); }
    .error-content { display:flex; gap:var(--space-4); align-items:center; }
    .error-content h2 { margin-bottom:var(--space-1); }
    .error-content p { margin:0; color:var(--muted); font-size:var(--font-14); }
    .error-icon { display:grid; width:40px; height:40px; flex:0 0 40px; place-items:center; border-radius:var(--radius-pill); color:var(--danger); background:var(--danger-soft); font-size:var(--font-20); font-weight:800; }
    .empty-art { width:132px; height:100px; color:var(--accent); fill:var(--accent-soft); stroke:currentColor; stroke-width:2; stroke-linejoin:round; }
    @media(max-width:760px) { .page-heading { align-items:flex-start; flex-direction:column; } .server-grid { grid-template-columns:1fr; } .heading-actions { width:100%; justify-content:space-between; } }
    @media(max-width:400px) { h1 { font-size:var(--font-24); } :global(.server-card) { padding:var(--space-4); } }
  </style>
