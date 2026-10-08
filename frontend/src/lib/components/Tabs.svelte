<script>
  export let tabs = [];
  export let active = "";
  export let onSelect = () => {};
</script>

<div class="tabs" role="tablist">
  {#each tabs as tab}<button class:active={active === tab} type="button" role="tab" aria-selected={active === tab} on:click={() => onSelect(tab)}>{tab}</button>{/each}
</div>

<style>
  .tabs { display:flex; gap:var(--space-1); overflow-x:auto; border-bottom:1px solid var(--border); }
  button { min-height:var(--target); padding:0 var(--space-3); border:0; border-bottom:2px solid transparent; color:var(--muted); background:transparent; font-size:var(--font-14); font-weight:750; white-space:nowrap; cursor:pointer; }
  button.active, button:hover { color: var(--text); border-color: var(--accent); }
</style>
