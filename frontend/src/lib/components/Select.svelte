<script>
  export let label = "";
  export let id = "";
  export let value = "";
  export let options = [];
  export let disabled = false;
  export let error = "";
</script>

<label class="field" for={id}><span>{label}</span><select {id} bind:value {disabled} aria-invalid={Boolean(error)}>{#each options as option}<option value={typeof option === "string" ? option : option.value}>{typeof option === "string" ? option : option.label}</option>{/each}<slot /></select>{#if error}<span class="error" role="alert">{error}</span>{/if}</label>

<style>
  .field { display:grid; gap:var(--space-2); color:var(--text); font-size:var(--font-14); font-weight:650; }
  select { width:100%; min-height:var(--input-height); padding:0 var(--space-3); border:1px solid var(--border); border-radius:var(--radius-control); color:var(--text); background:var(--surface-2); }
  select:focus-visible { border-color:var(--accent); outline:2px solid var(--accent); outline-offset:2px; }
  select:disabled { opacity:.55; }
  .error { color:var(--danger); font-size:var(--font-12); }
</style>
