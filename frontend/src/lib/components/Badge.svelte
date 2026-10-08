<script>
  export let status = "offline";
  $: normalized = String(status).toLowerCase();
  $: tone = normalized === "running" || normalized === "online" ? "online" : normalized === "starting" || normalized === "stopping" ? "transitioning" : normalized === "error" ? "error" : "offline";
  $: label = normalized === "running" ? "Running" : normalized === "stopped" ? "Stopped" : status;
</script>

<span class={`badge badge-${tone}`}><span class="dot"></span>{label}</span>

<style>
  .badge { display:inline-flex; gap:var(--space-2); align-items:center; min-height:var(--size-badge); padding:0 var(--space-3); border-radius:var(--radius-pill); font-size:var(--font-12); font-weight:750; }
  .dot { width:var(--space-2); height:var(--space-2); border-radius:var(--radius-pill); background:currentColor; }
  .badge-online { color: var(--online); background: var(--online-soft); }
  .badge-offline { color: var(--muted); background: var(--muted-soft); }
  .badge-transitioning { color: var(--warning); background: var(--warning-soft); }
  .badge-error { color: var(--danger); background: var(--danger-soft); }
  .badge-transitioning .dot,.badge-online .dot { animation: pulse 1.8s ease-in-out infinite; }
  @keyframes pulse { 50% { opacity: .35; } }
</style>
