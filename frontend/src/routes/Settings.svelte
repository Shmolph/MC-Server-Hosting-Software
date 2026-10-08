<script>
  import { onMount } from "svelte";
  import Toggle from "../lib/components/Toggle.svelte";
  import { theme, applyTheme } from "../lib/theme.js";
  let compact = false;
  let timestamps = true;
  let autoscroll = true;
  let fontSize = "medium";
  onMount(() => { compact = localStorage.getItem("shmolph-compact") === "true"; timestamps = localStorage.getItem("shmolph-timestamps") !== "false"; autoscroll = localStorage.getItem("shmolph-autoscroll") !== "false"; fontSize = localStorage.getItem("shmolph-console-size") || "medium"; });
  function save(key, value) { localStorage.setItem(`shmolph-${key}`, String(value)); }
</script>

<div class="page-head"><div><p class="kicker">Settings</p><h1>Preferences</h1></div></div>
<section class="settings-card"><div class="setting-section"><h2>Appearance</h2><div class="theme-options">{#each ["light", "dark", "system"] as choice}<button class:active={$theme === choice} class="theme-card" type="button" on:click={() => applyTheme(choice)}>{choice[0].toUpperCase() + choice.slice(1)}</button>{/each}</div><Toggle label="Compact spacing" bind:checked={compact} on:change={() => save("compact", compact)} /></div><div class="setting-section"><h2>Console</h2><div class="font-size"><span>Font size</span>{#each ["small", "medium", "large"] as choice}<button class:active={fontSize === choice} type="button" on:click={() => { fontSize = choice; save("console-size", choice); }}>{choice[0].toUpperCase()}</button>{/each}</div><Toggle label="Show timestamps" bind:checked={timestamps} on:change={() => save("timestamps", timestamps)} /><Toggle label="Auto-scroll" bind:checked={autoscroll} on:change={() => save("autoscroll", autoscroll)} /></div><div class="setting-section"><h2>Notifications <small>Soon</small></h2><Toggle label="Server alerts" checked={false} /></div></section>

<style>
  .settings-card { max-width:720px; padding:20px; border:1px solid var(--border); border-radius:12px; background:var(--surface); box-shadow:var(--shadow); }
  .setting-section { padding:18px 0; border-bottom:1px solid var(--border); }
  .setting-section:last-child { border-bottom:0; }
  .setting-section h2 { margin:0 0 12px; font-size:14px; }
  small { margin-left:6px; color:var(--warning); font-size:10px; }
  .theme-options { display:grid; grid-template-columns:repeat(3,1fr); gap:8px; margin-bottom:10px; }
  .theme-card,.font-size button { min-height:48px; border:1px solid var(--border); border-radius:8px; color:var(--muted); background:var(--surface-2); font-weight:750; }
  .theme-card.active,.theme-card:hover,.font-size button.active,.font-size button:hover { border-color:var(--accent); color:var(--text); }
  .font-size { display:flex; gap:7px; align-items:center; justify-content:space-between; margin-bottom:8px; color:var(--muted); }
  .font-size button { min-width:40px; min-height:34px; }
  .font-size span { margin-right:auto; }
  @media (max-width:480px) { .theme-options { grid-template-columns:1fr; } }
</style>
