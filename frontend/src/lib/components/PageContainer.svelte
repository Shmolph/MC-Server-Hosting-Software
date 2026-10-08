<script>
  export let narrow = false;
</script>

<main class:narrow class="page-container"><slot /></main>
<style>
  .page-container { width:min(calc(100% - var(--space-7)),var(--content-width)); margin:0 auto; padding:var(--space-7) 0 var(--space-8); }
  .page-container.narrow { max-width:var(--content-narrow); }
  @media (max-width:600px) { .page-container { width:min(calc(100% - var(--space-6)),var(--content-width)); padding:var(--space-5) 0 var(--space-7); } }
</style>
