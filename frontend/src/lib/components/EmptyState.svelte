<script>
  export let title = "Nothing here yet";
</script>

<div class="empty"><strong>{title}</strong><slot /></div>
<style>
  .empty { display:grid; min-height:calc(var(--space-8) * 3); place-items:center; gap:var(--space-3); padding:var(--space-6); border:1px dashed var(--border); border-radius:var(--radius-card); color:var(--muted); text-align:center; }
  strong { color: var(--text); }
</style>
