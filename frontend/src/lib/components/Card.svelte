<script>
  export let variant = "default";
  export let className = "";
</script>

<article class={`card card-${variant} ${className}`}><slot /></article>
<style>
  .card { min-width:0; border:1px solid var(--border); border-radius:var(--radius-card); background:var(--surface-1); box-shadow:var(--shadow); }
  .card-server { padding:var(--space-5); transition:transform 180ms ease-out,border-color 180ms ease-out,box-shadow 180ms ease-out; }
  .card-server:hover { transform:translateY(-2px); border-color:color-mix(in srgb,var(--accent) 38%,var(--border)); box-shadow:var(--shadow-md),var(--accent-glow); }
  .card-error { padding:var(--space-5); }
</style>
