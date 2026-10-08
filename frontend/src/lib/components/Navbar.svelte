<script>
  import { Box, LogOut, Moon, Sun } from "lucide-svelte";
  import { pushToast } from "../stores.js";
  import { theme, cycleTheme } from "../theme.js";
  export let onLogout;
  export let username = "";
  let busy = false;
  async function logout() { busy = true; try { await onLogout(); } catch (error) { pushToast(error.message || "Logout failed. Try again.", "error"); } finally { busy = false; } }
</script>

<header class="navbar"><a class="brand" href="#/dashboard" aria-label="Shmolph home"><span class="mark"><Box size={18} /></span><span class="wordmark">SHMOLPH</span></a><nav aria-label="Primary"><a href="#/dashboard">Dashboard</a></nav><div class="user-area"><button class="theme-button" type="button" aria-label={`Theme: ${$theme}. Change theme`} on:click={() => cycleTheme($theme)}>{#if $theme === "dark"}<Moon size={16} />{:else}<Sun size={16} />{/if}<span>{`${$theme[0].toUpperCase()}${$theme.slice(1)}`}</span></button><span class="avatar">{username ? username[0].toUpperCase() : "S"}</span><span class="username">{username || "My account"}</span><button type="button" disabled={busy} on:click={logout}><LogOut size={16} />{busy ? "Logging out" : "Log out"}</button></div></header>

<style>
  .navbar { position:sticky; z-index:8; top:0; display:flex; min-height:68px; gap:32px; align-items:center; padding:0 max(24px,calc((100vw - 1200px)/2)); border-bottom:1px solid var(--border); background:color-mix(in srgb,var(--surface) 88%,transparent); backdrop-filter:blur(18px); }
  .brand { display:flex; gap:10px; align-items:center; color:var(--text); text-decoration:none; }
  .mark { display:grid; width:32px; height:32px; place-items:center; border-radius:8px; color:var(--text-on-accent); background:var(--accent); box-shadow:var(--accent-glow); }
  .wordmark { font-family:var(--pixel-font); font-size:12px; letter-spacing:1px; }
  nav { display:flex; gap:20px; flex:1; }
  nav a { color:var(--muted); font-size:14px; font-weight:650; text-decoration:none; }
  nav a:hover { color:var(--text); }
  .user-area { display:flex; gap:10px; align-items:center; }
  .avatar { display:grid; width:32px; height:32px; place-items:center; border:1px solid var(--border); border-radius:50%; color:var(--accent); background:var(--surface-2); font-weight:700; }
  .username { color:var(--muted); font-size:13px; }
  button { display:flex; min-height:40px; gap:8px; align-items:center; padding:0 11px; border:1px solid var(--border); border-radius:8px; color:var(--muted); background:var(--surface-2); font-weight:650; }
  button:hover:not(:disabled) { border-color:var(--accent); color:var(--text); }
  button:disabled { opacity:.6; }
  .theme-button { display:flex; gap:6px; align-items:center; }
  @media(max-width:600px) { .navbar { min-height:60px; gap:12px; padding:0 16px; } nav,.username { display:none; } .user-area { margin-left:auto; } .wordmark { font-size:10px; } .theme-button span { display:none; } }
</style>
