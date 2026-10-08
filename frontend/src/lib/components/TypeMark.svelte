<script>
  import { Box, Layers, Sparkles, Hammer } from "lucide-svelte";
  export let type = "Vanilla";
  const icons = { Vanilla: Box, Forge: Hammer, Fabric: Layers, Paper: Sparkles };
  $: TypeIcon = icons[type] || Box;
</script>

<span class={`type-mark type-${String(type).toLowerCase()}`}><svelte:component this={TypeIcon} size={15} /><span>{type}</span></span>
<style>
  .type-mark { display:inline-flex; min-height:var(--size-badge); gap:var(--space-2); align-items:center; padding:0 var(--space-3); border:1px solid var(--border); border-radius:var(--radius-pill); color:var(--muted); background:var(--surface-2); font-size:var(--font-12); font-weight:650; }
  .type-paper { color:var(--type-paper); background:var(--type-paper-soft); }
  .type-fabric { color:var(--type-fabric); background:var(--type-fabric-soft); }
  .type-forge { color:var(--type-forge); background:var(--type-forge-soft); }
  .type-vanilla { color:var(--type-vanilla); background:var(--type-vanilla-soft); }
</style>
