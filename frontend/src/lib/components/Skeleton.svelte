<script>
  export let lines = 3;
</script>

<div class="skeleton-list" aria-label="Loading">
  {#each Array(lines) as _}<div class="skeleton"></div>{/each}
</div>

<style>
  .skeleton-list { display:grid; gap:var(--space-3); }
  .skeleton { height:var(--skeleton-height); border:1px solid var(--border); border-radius:var(--radius-card); background:linear-gradient(90deg,var(--surface-2),var(--surface-3),var(--surface-2)); background-size:200% 100%; animation:shimmer 1.4s ease-in-out infinite; }
  @keyframes shimmer { from { background-position: 100% 0; } to { background-position: -100% 0; } }
</style>
