<script>
  export let checked = false;
  export let label = "";
</script>

<label class="toggle"><span>{label}</span><input type="checkbox" bind:checked /><span class="track"></span></label>

<style>
  .toggle { display:flex; min-height:var(--target); align-items:center; justify-content:space-between; gap:var(--space-3); color:var(--muted); cursor:pointer; }
  input { position: absolute; opacity: 0; }
  .track { width:calc(var(--space-7) + var(--space-2)); height:var(--space-5); padding:var(--space-1); border-radius:var(--radius-pill); background:var(--surface-3); transition:160ms ease; }
  .track::after { display:block; width:var(--space-4); height:var(--space-4); content:""; border-radius:var(--radius-pill); background:var(--muted); transition:160ms ease; }
  input:checked + .track { background: var(--accent); }
  input:checked + .track::after { transform:translateX(var(--space-4)); background:var(--text-on-accent); }
</style>
