<script>
  import { ArrowLeft, Box } from "lucide-svelte";
  import Button from "../lib/components/Button.svelte";
  export let navigate;
</script>

<main class="not-found" aria-labelledby="not-found-title">
  <span class="mark"><Box size={22} /></span>
  <p class="eyebrow">404 · LOST CHUNK</p>
  <h1 id="not-found-title">This page isn't here.</h1>
  <Button on:click={() => navigate("#/dashboard")}><ArrowLeft size={16} />Back to dashboard</Button>
</main>

<style>
  .not-found { display:grid; min-height:calc(100vh - var(--navbar-height)); align-content:center; justify-items:start; gap:var(--space-4); width:min(calc(100% - var(--space-7)),var(--content-width)); margin:0 auto; }
  .mark { display:grid; width:48px; height:48px; place-items:center; border:1px solid var(--accent); border-radius:var(--radius-card); color:var(--text-on-accent); background:var(--accent); box-shadow:var(--accent-glow); }
  .eyebrow { margin:0; color:var(--accent); font-family:var(--pixel-font); font-size:10px; }
  h1 { margin:0 0 var(--space-2); color:var(--text); font-size:var(--font-32); }
</style>
