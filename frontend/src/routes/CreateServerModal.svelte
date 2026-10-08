<script>
  import { onMount } from "svelte";
  import { fade, slide } from "svelte/transition";
  import { Check, Hammer, Layers, Sparkles, X } from "lucide-svelte";
  import Button from "../lib/components/Button.svelte";
  import Checkbox from "../lib/components/Checkbox.svelte";
  import Input from "../lib/components/Input.svelte";
  import Modal from "../lib/components/Modal.svelte";
  import Select from "../lib/components/Select.svelte";
  import Slider from "../lib/components/Slider.svelte";
  import Stepper from "../lib/components/Stepper.svelte";
  import Skeleton from "../lib/components/Skeleton.svelte";
  import { api } from "../lib/api.js";
  import { SERVER_OPTIONS } from "../lib/config.js";
  import { pushToast } from "../lib/stores.js";

  export let onClose = () => {};
  export let onCreated = () => {};
  let step = 1;
  let name = "";
  let type = "Paper";
  let version = "";
  let ram = 2048;
  let eula = false;
  let versions = [];
  let versionsLoading = true;
  let versionsError = "";
  let error = "";
  let loading = false;
  const typeIcons = { Vanilla: Sparkles, Forge: Hammer, Fabric: Layers, Paper: Layers };
  $: validName = name.trim().length >= 3 && name.trim().length <= 32;
  $: canSubmit = validName && Boolean(type) && Boolean(version) && Number.isInteger(Number(ram)) && Number(ram) >= 1024 && Number(ram) <= 4096 && eula && !loading;

  onMount(loadVersions);
  async function loadVersions() {
    versionsLoading = true;
    versionsError = "";
    try {
      const items = await api.listVersions();
      if (!Array.isArray(items) || items.some((item) => typeof item !== "string")) throw new Error("The version list returned an invalid response.");
      versions = items;
      version = items[0] || "";
      if (!items.length) versionsError = "No Minecraft releases are available.";
    } catch (requestError) { versionsError = requestError.message; }
    finally { versionsLoading = false; }
  }
  function nameError() {
    const length = name.trim().length;
    if (length < 3) return "name too short (min 3 characters)";
    if (length > 32) return "name too long (max 32 characters)";
    return "";
  }
  function next() {
    error = "";
    if (step === 1 && !validName) { error = nameError(); return; }
    if (step === 2 && (!version || versionsError)) { error = versionsError || "Choose a Minecraft version."; return; }
    step = Math.min(3, step + 1);
  }
  function back() { error = ""; step = Math.max(1, step - 1); }
  async function create() {
    error = "";
    if (!validName) { step = 1; error = nameError(); return; }
    if (!canSubmit) { error = "Choose a version, memory size, and accept the EULA."; return; }
    loading = true;
    try {
      const id = await api.createServer({ name: name.trim(), type, version, ram_mb: Number(ram), eula: true });
      pushToast(`Server ${id} created`, "success");
      onClose();
      await onCreated();
    } catch (requestError) { error = requestError.message; }
    finally { loading = false; }
  }
</script>

<Modal open={true} labelledBy="create-server-title" on:close={onClose}>
  <div class="modal-content">
    <header class="modal-header"><div><span class="eyebrow">NEW WORLD</span><h2 id="create-server-title">Create Server</h2></div><button class="icon-close" type="button" aria-label="Close dialog" on:click={onClose}><X size={18} /></button></header>
    <Stepper steps={["Details", "Resources", "Review"]} current={step} />
    {#if step === 1}
      <section class="form-step" in:slide={{ x: 12, duration: 180 }}>
        <Input id="server-name" label="Server name" bind:value={name} maxlength={32} required helper={`${name.trim().length}/32 characters`} error={name ? nameError() : ""} placeholder="My Minecraft world" />
        <fieldset><legend>Choose a server type</legend><div class="type-grid">{#each SERVER_OPTIONS.types as option}<button class:selected={type === option} class="type-choice" type="button" aria-pressed={type === option} on:click={() => type = option}>{#if typeIcons[option]}<svelte:component this={typeIcons[option]} size={20} />{/if}<strong>{option}</strong>{#if type === option}<Check class="selected-mark" size={16} />{/if}</button>{/each}</div></fieldset>
      </section>
    {:else if step === 2}
      <section class="form-step" in:slide={{ x: 12, duration: 180 }}>
        {#if versionsLoading}<Skeleton lines={1} />{:else if versionsError}<div class="version-error" role="alert"><span>{versionsError}</span><Button variant="secondary" on:click={loadVersions}>Retry versions</Button></div>{:else}<Select id="minecraft-version" label="Minecraft version" options={versions} bind:value={version} />{/if}
        <Slider id="server-ram" label="Memory (RAM)" min={1024} max={4096} step={1024} bind:value={ram} />
        <div class="eula-box"><Checkbox id="minecraft-eula" bind:checked={eula}>I agree to the <a href="https://www.minecraft.net/en-us/eula" target="_blank" rel="noreferrer">Minecraft EULA</a></Checkbox></div>
      </section>
    {:else}
      <section class="form-step" in:fade={{ duration: 160 }}>
        <div class="review-list">{#each [["Name", name.trim()], ["Type", type], ["Version", version], ["Memory (RAM)", `${(Number(ram) / 1024).toLocaleString(undefined, { maximumFractionDigits: 0 })} GB`], ["EULA", eula ? "Accepted" : "Required"]] as row}<div><span>{row[0]}</span><strong>{row[1]}</strong></div>{/each}</div>
      </section>
    {/if}
    {#if error}<p class="form-error" role="alert">{error}</p>{/if}
    <footer class="modal-actions"><Button variant="ghost" disabled={loading} on:click={onClose}>Cancel</Button>{#if step > 1}<Button variant="secondary" disabled={loading} on:click={back}>Back</Button>{/if}{#if step < 3}<Button disabled={step === 2 && (versionsLoading || Boolean(versionsError) || !version)} on:click={next}>Continue</Button>{:else}<Button loading={loading} disabled={!canSubmit} on:click={create}>Create Server</Button>{/if}</footer>
  </div>
</Modal>

<style>
  .modal-content { padding:var(--space-5); }
  .modal-header { display:flex; align-items:flex-start; justify-content:space-between; gap:var(--space-4); margin-bottom:var(--space-4); }
  .eyebrow { color:var(--accent); font-family:var(--pixel-font); font-size:9px; letter-spacing:1px; }
  h2 { margin:var(--space-2) 0 0; color:var(--text); font-size:var(--font-24); }
  .icon-close { display:grid; width:var(--target); height:var(--target); place-items:center; border:1px solid var(--border); border-radius:var(--radius-control); color:var(--muted); background:var(--surface-2); }
  .icon-close:hover { color:var(--text); border-color:var(--accent); }
  .form-step { display:grid; gap:var(--space-5); min-height:var(--create-step-height); padding:var(--space-5) 0; }
  fieldset { display:grid; gap:var(--space-3); margin:0; padding:0; border:0; }
  legend { margin-bottom:var(--space-3); color:var(--text); font-size:var(--font-14); font-weight:650; }
  .type-grid { display:grid; grid-template-columns:repeat(2,minmax(0,1fr)); gap:var(--space-3); }
  .type-choice { position:relative; display:flex; min-height:80px; gap:var(--space-3); align-items:center; padding:var(--space-4); border:1px solid var(--border); border-radius:var(--radius-card); color:var(--muted); background:var(--surface-2); text-align:left; transition:180ms ease-out; }
  .type-choice strong { color:var(--text); font-size:var(--font-14); }
  .type-choice:hover,.type-choice.selected { border-color:var(--accent); box-shadow:var(--accent-glow); }
  :global(.selected-mark) { margin-left:auto; color:var(--accent); }
  .version-error { display:flex; min-height:var(--target); gap:var(--space-3); align-items:center; justify-content:space-between; color:var(--danger); font-size:var(--font-14); }
  .eula-box { padding:var(--space-3); border:1px solid var(--border); border-radius:var(--radius-control); background:var(--surface-2); }
  .eula-box a { color:var(--accent); text-decoration:underline; text-underline-offset:3px; }
  .review-list { overflow:hidden; border:1px solid var(--border); border-radius:var(--radius-control); }
  .review-list div { display:flex; min-height:var(--target); align-items:center; justify-content:space-between; gap:var(--space-3); padding:0 var(--space-4); border-bottom:1px solid var(--border); }
  .review-list div:last-child { border-bottom:0; }
  .review-list span { color:var(--muted); font-size:var(--font-14); }
  .review-list strong { color:var(--text); font-size:var(--font-14); }
  .form-error { margin:0; padding:var(--space-3); border:1px solid var(--danger); border-radius:var(--radius-control); color:var(--danger); background:var(--danger-soft); font-size:var(--font-14); }
  .modal-actions { display:flex; flex-wrap:wrap; justify-content:flex-end; gap:var(--space-2); padding-top:var(--space-4); border-top:1px solid var(--border); }
  @media(max-width:420px) { .modal-content { padding:var(--space-4); } .type-choice { padding:var(--space-3); } }
</style>
