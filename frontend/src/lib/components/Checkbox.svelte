<script>
  export let checked = false;
  export let id = "";
  export let disabled = false;
</script>

<label class="check" for={id}><input {id} type="checkbox" bind:checked {disabled} /><span class="box" aria-hidden="true"></span><span class="label"><slot /></span></label>

<style>
  .check { display:flex; min-height:var(--target); gap:var(--space-3); align-items:center; color:var(--text); cursor:pointer; }
  input { position:absolute; width:1px; height:1px; opacity:0; }
  .box { display:grid; width:var(--space-5); height:var(--space-5); flex:0 0 var(--space-5); place-items:center; border:1px solid var(--border); border-radius:var(--space-1); background:var(--surface-2); }
  input:checked + .box { border-color:var(--accent); background:var(--accent); }
  input:checked + .box::after { width:6px; height:10px; content:""; border:solid var(--text-on-accent); border-width:0 2px 2px 0; transform:rotate(45deg) translate(-1px,-1px); }
  input:focus-visible + .box { outline:2px solid var(--accent); outline-offset:2px; }
  input:disabled ~ .label { color:var(--dim); }
  .label { font-size:var(--font-14); }
</style>
