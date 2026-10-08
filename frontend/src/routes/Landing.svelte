<script>
  import { ArrowRight, LogIn, Server, ShieldCheck } from "lucide-svelte";
  import Button from "../lib/components/Button.svelte";
  export let navigate;
</script>

<div class="landing"><header class="landing-top"><a class="brand" href="#/home"><span class="brand-mark">S</span><span>shmolph.cloud</span></a><nav><button type="button" on:click={() => navigate("#/login")}><LogIn size={15} />Log in</button><Button on:click={() => navigate("#/signup")}>Sign up <ArrowRight size={15} /></Button></nav></header><main><section class="landing-intro"><p class="kicker">Minecraft hosting</p><h1>Simple control for your server.</h1><p>Affordable hosting with a clear path from account to world.</p><div><Button on:click={() => navigate("#/signup")}>Create account <ArrowRight size={15} /></Button><button class="quiet-action" type="button" on:click={() => navigate("#/login")}>Log in</button></div></section><section class="landing-points"><div><Server size={18} /><strong>Server control</strong><span>Start, stop and manage your world.</span></div><div><ShieldCheck size={18} /><strong>Clear by default</strong><span>No hidden setup maze.</span></div></section></main></div>

<style>
  .landing { min-height:100vh; padding:0 5vw; background:var(--bg); }
  .landing-top { display:flex; min-height:72px; align-items:center; justify-content:space-between; border-bottom:1px solid var(--border); }
  .landing-top nav { display:flex; gap:8px; align-items:center; }
  .landing-top nav > button { display:inline-flex; gap:6px; align-items:center; min-height:42px; padding:0 12px; border:0; border-radius:8px; color:var(--muted); background:transparent; font-weight:750; }
  .landing-top nav > button:hover { color:var(--text); background:var(--surface-2); }
  .landing-intro { max-width:650px; padding:clamp(80px,16vh,160px) 0 90px; }
  h1 { max-width:620px; margin:0 0 18px; font-size:clamp(42px,7vw,82px); letter-spacing:-.06em; line-height:.95; }
  .landing-intro > p:not(.kicker) { max-width:460px; color:var(--muted); font-size:17px; }
  .landing-intro > div { display:flex; gap:9px; align-items:center; margin-top:28px; }
  .quiet-action { min-height:42px; padding:0 12px; border:0; color:var(--muted); background:transparent; font-weight:750; }
  .landing-points { display:grid; grid-template-columns:repeat(2,minmax(0,260px)); gap:10px; padding-bottom:48px; }
  .landing-points div { display:grid; gap:7px; padding:18px; border:1px solid var(--border); border-radius:10px; background:var(--surface); }
  :global(.landing-points svg) { color:var(--accent); }
  .landing-points span { color:var(--muted); font-size:12px; }
  @media (max-width:560px) { .landing-top nav > button:first-child { display:none; } .landing-points { grid-template-columns:1fr; } }
</style>
