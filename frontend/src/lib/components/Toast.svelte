<script>
  import { fly } from "svelte/transition";
  export let items = [];
</script>

<div class="toasts" aria-live="polite">
  {#each items as toast (toast.id)}<div class={`toast toast-${toast.type}`} in:fly={{ y: 8, duration: 170 }}>{toast.message}</div>{/each}
</div>

<style>
  .toasts { position:fixed; z-index:20; right:var(--space-5); bottom:var(--space-5); display:grid; gap:var(--space-2); width:min(var(--toast-width),calc(100vw - var(--space-7))); }
  .toast { padding:var(--space-3) var(--space-4); border:1px solid var(--border); border-radius:var(--radius-control); color:var(--text); background:var(--surface-1); box-shadow:var(--shadow-lg); font-size:var(--font-14); }
  .toast-error { border-color: var(--danger); }
  .toast-success { border-color: var(--accent); }
</style>
