<script>
  import { onDestroy } from "svelte";
  import { ArrowRight, Box, ShieldCheck } from "lucide-svelte";
  import { fade } from "svelte/transition";
  import Button from "../lib/components/Button.svelte";
  import Captcha from "../lib/components/Captcha.svelte";
  import Input from "../lib/components/Input.svelte";
  import { api } from "../lib/api.js";
  import { pushToast } from "../lib/stores.js";

  export let mode = "login";
  export let warning = "";
  export let navigate;
  let username = "";
  let password = "";
  let captcha;
  let captchaToken = "";
  let captchaReady = false;
  let loading = false;
  let error = "";
  let success = "";
  let destroyed = false;
  $: registering = mode === "register";
  $: canSubmit = username.trim().length > 0 && password.length > 0 && password.length <= 72 && captchaReady && captchaToken.length > 0 && !loading;
  onDestroy(() => { destroyed = true; });

  async function submit(event) {
    event.preventDefault();
    error = "";
    success = "";
    const token = captcha?.getToken() || "";
    if (!username.trim() || !password || password.length > 72 || !token) {
      error = !token ? "Complete the security check to continue." : "Enter a username and a password up to 72 characters.";
      captcha?.reset();
      return;
    }
    loading = true;
    try {
      const credentials = { username: username.trim(), password, captchaToken: token };
      if (registering) {
        const registerResult = await api.register(credentials);
        if (registerResult.status !== 200) { error = registerResult.text || `Registration failed (${registerResult.status}).`; return; }
        captcha.reset();
        const freshToken = await captcha.freshToken();
        const loginResult = await api.login({ username: username.trim(), password, captchaToken: freshToken });
        if (loginResult.status !== 200) {
          success = "Account created. Log in with your new account.";
          mode = "login";
          password = "";
          pushToast("Account created", "success");
          return;
        }
      } else {
        const loginResult = await api.login(credentials);
        if (loginResult.status !== 200) { error = loginResult.text || `Login failed (${loginResult.status}).`; return; }
      }
      navigate("#/dashboard");
    } catch (requestError) {
      error = requestError.message || "Can't reach the server. Try again.";
    } finally {
      captcha?.reset();
      loading = false;
    }
  }
</script>

<main class="auth-page" in:fade={{ duration:180 }}>
  <section class="auth-card" aria-labelledby="auth-title">
    <a class="brand" href="#/dashboard" aria-label="Shmolph home"><span class="brand-mark"><Box size={18} /></span><span class="wordmark">SHMOLPH</span></a>
    <div class="auth-heading"><span class="eyebrow">{registering ? "NEW WORLD AWAITS" : "SERVER ACCESS"}</span><h1 id="auth-title">{registering ? "Create account" : "Welcome back"}</h1><p>{registering ? "One account. Your own place to play." : "Sign in to your worlds."}</p></div>
    {#if warning}<div class="notice" role="status">{warning}</div>{/if}
    <form on:submit={submit}>
      <Input id="auth-username" label="Username" bind:value={username} required autocomplete="username" placeholder="Your username" />
      <Input id="auth-password" label="Password" type="password" bind:value={password} required maxlength={72} autocomplete={registering ? "new-password" : "current-password"} placeholder="Your password" helper="Up to 72 characters" />
      <div class="captcha-box"><Captcha bind:this={captcha} on:token={(event) => { captchaToken = event.detail; captchaReady = true; }} on:error={(event) => error = event.detail} /></div>
      {#if error}<div class="form-message error" role="alert" aria-live="polite">{error}</div>{/if}
      {#if success}<div class="form-message success" role="status" aria-live="polite">{success}</div>{/if}
      <Button type="submit" disabled={!canSubmit} loading={loading}>{registering ? "Create account" : "Log in"}<ArrowRight size={16} /></Button>
    </form>
    <div class="auth-foot"><ShieldCheck size={16} /><span>Session protected by HttpOnly cookies.</span></div>
    <button class="auth-switch" type="button" on:click={() => { error = ""; success = ""; mode = registering ? "login" : "register"; navigate(registering ? "#/login" : "#/register"); }}>{registering ? "Already have an account? Log in" : "New to Shmolph? Create an account"}</button>
  </section>
</main>

<style>
  .auth-page { display:grid; min-height:100vh; place-items:center; padding:var(--space-5); background:transparent; }
  .auth-card { width:min(100%,440px); padding:var(--space-6); border:1px solid var(--border); border-radius:var(--radius-card); background:var(--surface-1); box-shadow:var(--shadow-lg); }
  .brand { display:inline-flex; gap:var(--space-3); align-items:center; margin-bottom:var(--space-7); color:var(--text); text-decoration:none; }
  .brand-mark { display:grid; width:40px; height:40px; place-items:center; border:1px solid var(--accent); border-radius:var(--radius-control); color:var(--text-on-accent); background:var(--accent); box-shadow:var(--accent-glow); }
  .wordmark { font-family:var(--pixel-font); font-size:12px; letter-spacing:1px; }
  .auth-heading { margin-bottom:var(--space-5); }
  .eyebrow { color:var(--accent); font-family:var(--pixel-font); font-size:9px; letter-spacing:1px; }
  h1 { margin:var(--space-3) 0 var(--space-2); color:var(--text); font-size:var(--font-32); letter-spacing:-.035em; }
  .auth-heading p { margin:0; color:var(--muted); font-size:var(--font-14); }
  form { display:grid; gap:var(--space-4); }
  .captcha-box { display:grid; min-height:var(--captcha-height); place-items:center; overflow:hidden; }
  .form-message { padding:var(--space-3); border:1px solid var(--border); border-radius:var(--radius-control); font-size:var(--font-14); }
  .form-message.error { border-color:var(--danger); color:var(--danger); background:var(--danger-soft); }
  .form-message.success { border-color:var(--accent); color:var(--accent); background:var(--accent-soft); }
  .notice { margin-bottom:var(--space-4); padding:var(--space-3); border:1px solid var(--warning); border-radius:var(--radius-control); color:var(--warning); background:var(--warning-soft); font-size:var(--font-14); }
  .auth-foot { display:flex; gap:var(--space-2); align-items:center; margin-top:var(--space-4); padding-top:var(--space-4); border-top:1px solid var(--border); color:var(--muted); font-size:var(--font-12); }
  :global(.auth-foot svg) { color:var(--accent); }
  .auth-switch { width:100%; min-height:var(--target); margin-top:var(--space-3); border:0; border-radius:var(--radius-control); color:var(--accent); background:transparent; font-size:var(--font-14); font-weight:650; }
  .auth-switch:hover { background:var(--surface-2); }
  @media(max-width:420px) { .auth-card { padding:var(--space-4); } }
</style>
