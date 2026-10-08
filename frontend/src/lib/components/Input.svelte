<script>
  export let label = "";
  export let id = "";
  export let value = "";
  export let helper = "";
  export let error = "";
  export let type = "text";
  export let autocomplete = undefined;
  export let maxlength;
  export let required = false;
  export let placeholder = "";
</script>

<label class="field" for={id}>
  <span class="field-label">{label}{#if required}<span class="required">*</span>{/if}</span>
  <input class:error={Boolean(error)} {id} {type} bind:value {maxlength} {required} {placeholder} {autocomplete} aria-invalid={Boolean(error)} aria-describedby={error ? `${id}-error` : helper ? `${id}-helper` : undefined} />
  {#if error}<span class="field-error" id={`${id}-error`} role="alert">{error}</span>{:else if helper}<span class="field-helper" id={`${id}-helper`}>{helper}</span>{/if}
</label>

<style>
  .field { display:grid; gap:var(--space-2); color:var(--text); }
  .field-label { font-size:var(--font-14); font-weight:650; }
  .required { margin-left:var(--space-1); color:var(--warning); }
  input { width:100%; min-height:var(--input-height); padding:0 var(--space-3); border:1px solid var(--border); border-radius:var(--radius-control); color:var(--text); background:var(--surface-2); transition:border-color 180ms ease,box-shadow 180ms ease; }
  input:hover { border-color:color-mix(in srgb,var(--accent) 35%,var(--border)); }
  input:focus-visible { border-color:var(--accent); outline:2px solid var(--accent); outline-offset:2px; }
  input.error { border-color:var(--danger); }
  .field-helper { color:var(--muted); font-size:var(--font-12); }
  .field-error { color:var(--danger); font-size:var(--font-12); }
</style>
