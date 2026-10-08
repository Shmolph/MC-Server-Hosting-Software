<script>
  import { Check, Copy } from "lucide-svelte";
  import { pushToast } from "../stores.js";
  export let value = "";
  export let label = "Copy address";
  let copied = false;
  async function copy() {
    try {
      await navigator.clipboard.writeText(value);
      copied = true;
      pushToast("Address copied", "success");
      setTimeout(() => copied = false, 1600);
    } catch (error) { pushToast("Clipboard access failed. Select and copy the address.", "error"); }
  }
</script>

<button class="copy-button" type="button" aria-label={copied ? "Copied" : label} title={copied ? "Copied!" : label} on:click={copy}>{#if copied}<Check size={16} />Copied!{:else}<Copy size={16} />Copy{/if}</button>

<style>
  .copy-button { display:inline-flex; min-height:var(--button-height); gap:var(--space-2); align-items:center; justify-content:center; padding:0 var(--space-3); border:1px solid var(--border); border-radius:var(--radius-control); color:var(--muted); background:var(--surface-2); font-size:var(--font-12); font-weight:700; transition:160ms ease; }
  .copy-button:hover { border-color:var(--accent); color:var(--text); box-shadow:var(--accent-glow); }
</style>
