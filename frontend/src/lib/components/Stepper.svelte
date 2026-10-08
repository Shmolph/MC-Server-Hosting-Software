<script>
  export let steps = [];
  export let current = 1;
</script>

<div class="stepper">{#each steps as step, index}<span class:active={index + 1 <= current}>{index + 1} {step}</span>{/each}</div>
<style>
  .stepper { display:grid; grid-template-columns:repeat(3,1fr); gap:var(--space-2); }
  span { min-height:var(--target); padding:var(--space-2) var(--space-1); border-bottom:2px solid var(--border); color:var(--muted); font-size:var(--font-12); font-weight:750; text-align:center; }
  span.active { border-color: var(--accent); color: var(--accent); }
</style>
