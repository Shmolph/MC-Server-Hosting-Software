<script>
  import { ArrowLeft, Check } from "lucide-svelte";
  import Button from "../lib/components/Button.svelte";
  import Modal from "../lib/components/Modal.svelte";
  import Stepper from "../lib/components/Stepper.svelte";
  import { api } from "../lib/api.js";
  import { pushToast } from "../lib/stores.js";

  export let navigate;
  let step = 1;
  let name = "";
  let type = "Paper";
  let version = "1.21.4";
  let ram = 2048;
  let eula = false;
  let error = "";
  let loading = false;
  const SERVER_OPTIONS = { types: ["Vanilla", "Paper", "Fabric", "Forge"], versions: ["1.21.4", "1.21.3", "1.20.6"] };
  function validName() { return name.trim().length >= 3 && name.trim().length <= 32; }
  function nameValidationMessage() {
    const length = name.trim().length;
    if (length < 3) return "name too short (min 3 characters)";
    if (length > 32) return "name too long (max 32 characters)";
    return "";
  }
  function next() { error = ""; if (step === 1 && !validName()) { error = nameValidationMessage(); return; } step = Math.min(3, step + 1); }
  async function create() {
    error = "";
    if (!validName()) { step = 1; error = nameValidationMessage(); return; }
    if (!type || !version || !Number.isInteger(Number(ram)) || Number(ram) < 1024 || Number(ram) > 4096 || !eula) { error = "Choose valid server settings and agree to the EULA."; return; }
    loading = true;
    try {
      await api.createServer({ name: name.trim(), type, version, ram_mb: Number(ram), eula: true });
      await api.listServers();
      pushToast("Server created", "success");
      navigate("#/servers");
    } catch (requestError) { error = requestError.message; }
    finally { loading = false; }
  }
</script>

<Modal open={true} on:close={() => navigate("#/servers")}>
  <div class="modal-content"><header><div><p class="kicker">New server</p><h1>Create server</h1></div><button class="close" type="button" on:click={() => navigate("#/servers")}>×</button></header><Stepper steps={["Details", "Resources", "Review"]} current={step} />
  {#if step === 1}<div class="step-body"><label>Server name<input class="input" bind:value={name} minlength="3" maxlength="32" placeholder="My Minecraft server"><span>{name.trim().length < 3 ? "Enter at least 3 characters" : `${name.trim().length}/32`}</span></label><div class="type-grid">{#each SERVER_OPTIONS.types as serverType}<button class:active={type === serverType} type="button" class="type-card" on:click={() => type = serverType}><strong>{serverType}</strong>{#if type === serverType}<Check size={15} />{/if}</button>{/each}</div></div>{:else if step === 2}<div class="step-body"><label>Minecraft version<select class="select" bind:value={version}>{#each SERVER_OPTIONS.versions as item}<option>{item}</option>{/each}</select></label><label>Memory (RAM)<input type="range" min="1024" max="4096" step="512" bind:value={ram}><strong class="range-value">{ram} MB ({(Number(ram) / 1024).toLocaleString(undefined, { maximumFractionDigits: 1 })} GB)</strong></label><label class="eula"><input type="checkbox" bind:checked={eula}>I agree to the Minecraft <a href="https://aka.ms/MinecraftEULA" target="_blank" rel="noreferrer">EULA</a></label></div>{:else}<div class="review"><div><span>Name</span><strong>{name.trim()}</strong></div><div><span>Type</span><strong>{type}</strong></div><div><span>Version</span><strong>{version}</strong></div><div><span>Memory (RAM)</span><strong>{ram} MB</strong></div><div><span>EULA</span><strong>{eula ? "Agreed" : "Not agreed"}</strong></div></div>{/if}
  {#if error}<p class="form-error">{error}</p>{/if}<footer><Button variant="ghost" on:click={() => navigate("#/servers")} >Cancel</Button>{#if step > 1}<Button variant="secondary" disabled={loading} on:click={() => step -= 1}>Back</Button>{/if}{#if step < 3}<Button on:click={next}>Next</Button>{:else}<Button loading={loading} disabled={!eula} on:click={create}>Create server</Button>{/if}</footer></div>
</Modal>

<style>
  .modal-content { padding:20px; }
  header { display:flex; justify-content:space-between; align-items:start; margin-bottom:16px; }
  h1 { margin:0; font-size:22px; }
  .close { min-width:40px; min-height:40px; border:0; color:var(--muted); background:transparent; font-size:24px; }
  .step-body { display:grid; gap:18px; padding:20px 0; }
  label { display:grid; gap:8px; color:var(--muted); font-size:12px; font-weight:750; }
  .eula { display:flex; flex-wrap:wrap; align-items:center; gap:7px; }
  .eula input { width:18px; height:18px; accent-color:var(--accent); }
  .eula a { color:var(--accent); }
  label span { font-size:11px; font-weight:500; }
  .type-grid { display:grid; grid-template-columns:repeat(2,1fr); gap:8px; }
  .type-card { position:relative; display:grid; gap:7px; min-height:92px; padding:14px; border:1px solid var(--border); border-radius:10px; color:var(--muted); background:var(--surface-2); text-align:left; }
  .type-card.active,.type-card:hover { border-color:var(--accent); color:var(--text); }
  :global(.type-card svg) { position:absolute; right:12px; top:12px; color:var(--accent); }
  .type-card strong { color:var(--text); }
  input[type=range] { width:100%; accent-color:var(--accent); }
  .range-value { color:var(--accent); }
  .review { display:grid; border:1px solid var(--border); border-radius:8px; overflow:hidden; }
  .review div { display:flex; justify-content:space-between; padding:13px; border-bottom:1px solid var(--border); }
  .review div:last-child { border-bottom:0; }
  .review span { color:var(--muted); }
  footer { display:flex; justify-content:flex-end; gap:8px; padding-top:16px; border-top:1px solid var(--border); }
  .form-error { margin:0; color:var(--danger); font-size:12px; }
</style>
