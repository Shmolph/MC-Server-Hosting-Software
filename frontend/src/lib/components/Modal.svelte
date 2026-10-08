<script>
  import { afterUpdate, createEventDispatcher, onDestroy } from "svelte";
  export let open = false;
  export let labelledBy = "dialog-title";
  const dispatch = createEventDispatcher();
  let dialog;
  let returnFocus;
  let wasOpen = false;
  afterUpdate(() => {
    if (open && !wasOpen && dialog && !dialog.open) {
      returnFocus = document.activeElement;
      dialog.showModal();
      requestAnimationFrame(() => dialog.querySelector("input,button,select,a,[tabindex='0']")?.focus());
    } else if (!open && wasOpen && dialog?.open) dialog.close();
    wasOpen = open;
  });
  function requestClose() { dispatch("close"); }
  function onCancel(event) { event.preventDefault(); requestClose(); }
  function onClose() { returnFocus?.focus?.(); dispatch("closed"); }
  function onBackdrop(event) { if (event.target === dialog) requestClose(); }
  onDestroy(() => { if (dialog?.open) dialog.close(); });
</script>

<dialog bind:this={dialog} class="modal" aria-labelledby={labelledBy} on:cancel={onCancel} on:close={onClose} on:click={onBackdrop}>
  <slot />
</dialog>

<style>
  .modal { width:min(calc(100vw - var(--space-8)),560px); max-height:calc(100vh - var(--space-8)); overflow:auto; padding:0; border:1px solid var(--border); border-radius:var(--radius-card); color:var(--text); background:var(--surface); box-shadow:var(--shadow-lg); }
  .modal::backdrop { background:var(--scrim); backdrop-filter:blur(4px); }
  .modal[open] { animation:dialog-enter 180ms ease-out; }
  @keyframes dialog-enter { from { opacity:0; transform:translateY(8px) scale(.99); } to { opacity:1; transform:translateY(0) scale(1); } }
</style>
