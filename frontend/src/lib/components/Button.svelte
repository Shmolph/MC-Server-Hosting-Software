<script>
  export let variant = "primary";
  export let loading = false;
  export let disabled = false;
  export let type = "button";
</script>

<button class={`button button-${variant}`} {type} disabled={disabled || loading} on:click>
  {#if loading}<span class="button-spinner" aria-hidden="true"></span>{/if}
  <slot />
</button>

<style>
  .button { display:inline-flex; min-height:var(--target); gap:var(--space-2); align-items:center; justify-content:center; padding:0 var(--space-4); border:1px solid var(--border); border-radius:var(--radius-control); color:var(--text); background:var(--surface-2); font-size:var(--font-14); font-weight:750; cursor:pointer; transition:transform 180ms ease-out,border-color 180ms ease-out,background 180ms ease-out,box-shadow 180ms ease-out; }
  .button:hover:not(:disabled) { transform:translateY(-1px); border-color:var(--accent); }
  .button:active:not(:disabled) { transform:translateY(0); }
  .button:disabled { cursor:not-allowed; opacity:.55; }
  .button-primary { border-color:var(--accent); color:var(--text-on-accent); background:var(--accent); box-shadow:var(--accent-glow); }
  .button-primary:hover:not(:disabled) { box-shadow:var(--accent-glow),var(--shadow-sm); }
  .button-danger { border-color:var(--danger); color:var(--text-on-danger); background:var(--danger); }
  .button-ghost { background: transparent; }
  .button-spinner { display:inline-block; width:16px; height:16px; border:2px solid currentColor; border-right-color:transparent; border-radius:50%; animation:spin 700ms linear infinite; }
  @keyframes spin { to { transform: rotate(360deg); } }
</style>
