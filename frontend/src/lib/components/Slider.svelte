<script>
  export let label = "";
  export let id = "";
  export let value = 1024;
  export let min = 1024;
  export let max = 4096;
  export let step = 1024;
  $: display = `${(Number(value) / 1024).toLocaleString(undefined, { maximumFractionDigits: 1 })} GB`;
</script>

<div class="slider"><div class="slider-head"><label for={id}>{label}</label><output for={id}>{display}</output></div><input {id} type="range" bind:value {min} {max} {step} /><div class="range-labels"><span>{min} MB</span><span>{max} MB</span></div></div>

<style>
  .slider { display:grid; gap:var(--space-3); }
  .slider-head,.range-labels { display:flex; align-items:center; justify-content:space-between; }
  label { color:var(--text); font-size:var(--font-14); font-weight:650; }
  output { color:var(--accent); font-family:var(--mono-font); font-size:var(--font-14); font-weight:700; }
  input { width:100%; accent-color:var(--accent); }
  .range-labels { color:var(--muted); font-size:var(--font-12); }
</style>
