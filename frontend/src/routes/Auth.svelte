<script>
  import { onDestroy, onMount } from "svelte";
  import { ArrowRight, ShieldCheck } from "lucide-svelte";
  import Button from "../lib/components/Button.svelte";
  import { api } from "../lib/api.js";
  import { pushToast } from "../lib/stores.js";

  export let mode = "login";
  export let navigate;
  export let sessionWarning = "";
  const turnstileSiteKey = "1x00000000000000000000AA";
  let username = "";
  let password = "";
  let widgetId = null;
  let token = "";
  let loading = false;
  let error = "";
  let turnstileReady = false;
  let responseWasCreated = false;
  let freshTokenResolve = null;
  $: isSignup = mode === "signup";

  onMount(async () => {
    try {
      await loadTurnstile();
      widgetId = window.turnstile.render("#turnstile-widget", { sitekey: turnstileSiteKey, callback: (value) => { token = value; freshTokenResolve?.(value); freshTokenResolve = null; } });
      turnstileReady = true;
    } catch (loadError) { error = "Could not load the CAPTCHA. Refresh and try again."; }
  });

  function loadTurnstile() {
    if (window.turnstile) return Promise.resolve();
    return new Promise((resolve, reject) => { const script = document.createElement("script"); script.src = "https://challenges.cloudflare.com/turnstile/v0/api.js"; script.async = true; script.defer = true; script.onload = resolve; script.onerror = reject; document.head.appendChild(script); });
  }
  function resetCaptcha() { if (widgetId !== null) window.turnstile.reset(widgetId); token = ""; }
  function removeCaptcha() { if (widgetId !== null && window.turnstile) window.turnstile.remove(widgetId); widgetId = null; token = ""; }
  function responseMessage(status, text) {
    if (!isSignup && status === 401) return "Invalid username or password.";
    const fallback = { 400: "Check the username and password.", 403: "CAPTCHA failed. Please try again.", 409: "That username is already taken.", 500: "The server could not complete the request." };
    return text || fallback[status] || "Request failed.";
  }
  async function submit() {
    error = "";
    if (!username.trim() || !password || (isSignup && (password.length < 8 || password.length > 72))) { error = isSignup ? "Enter a valid username and an 8 to 72 character password." : "Enter your username and password."; return; }
    if (!token) { error = "Complete the CAPTCHA check."; return; }
    loading = true;
    responseWasCreated = false;
    try {
      const response = isSignup ? await api.register({ username: username.trim(), password, captchaToken: token }) : await api.login({ username: username.trim(), password, captchaToken: token });
      if (response.status === 200) {
        responseWasCreated = isSignup;
        if (isSignup) {
          resetCaptcha();
          error = "";
          const loginToken = await new Promise((resolve) => { freshTokenResolve = resolve; resetCaptcha(); });
          const loginResponse = await api.login({ username: username.trim(), password, captchaToken: loginToken });
          if (loginResponse.status !== 200) throw new Error(loginResponse.status === 401 ? "Account created, but automatic login failed." : await loginResponse.text());
        }
        navigate("#/dashboard");
        return;
      }
      const text = await response.text();
      error = responseMessage(response.status, text);
    } catch (requestError) { error = isSignup && responseWasCreated ? "Account created. Please log in." : requestError.message; }
    finally { resetCaptcha(); loading = false; }
  }

  onDestroy(removeCaptcha);
</script>

<div class="auth-page"><div class="auth-card"><a class="brand" href="#/dashboard"><span class="brand-mark">S</span><span>shmolph.cloud</span></a><div class="auth-heading"><p class="kicker">{isSignup ? "Create account" : "Welcome back"}</p><h1>{isSignup ? "Start hosting simply." : "Sign in."}</h1><span class="auth-hint">{isSignup ? "One account. One clear control panel." : "Your servers, in one place."}</span></div>{#if sessionWarning}<div class="notice">{sessionWarning}</div>{/if}<form on:submit|preventDefault={submit}><label>Username<input bind:value={username} autocomplete="username" required placeholder="Username"></label><label>Password<input bind:value={password} type="password" maxlength="72" autocomplete={isSignup ? "new-password" : "current-password"} required placeholder="Password"></label>{#if isSignup}<div class="password-hint">8 to 72 characters · capital · number · symbol</div>{/if}<div id="turnstile-widget"></div><Button type={isSignup ? "submit" : "button"} loading={loading}>{isSignup ? "Create account" : "Log in"}<ArrowRight size={15} /></Button>{#if error}<div class="form-error" role="alert">{error}</div>{/if}</form><div class="auth-foot"><ShieldCheck size={15} /><span>Session cookies stay HttpOnly.</span></div><button class="switch-auth" type="button" on:click={() => navigate(isSignup ? "#/login" : "#/signup")}>{isSignup ? "Already have an account? Log in" : "Need an account? Sign up"}</button></div></div>

<style>
  .auth-page { display:grid; min-height:100vh; place-items:center; padding:24px; background:var(--bg); }
  .auth-card { width:min(100%,430px); padding:32px; border:1px solid var(--border); border-radius:12px; background:var(--surface); box-shadow:var(--shadow-strong); }
  .auth-card .brand { margin-bottom:52px; }
  .auth-heading { margin-bottom:26px; }
  .auth-heading h1 { margin:0 0 8px; font-size:32px; letter-spacing:-.03em; }
  .auth-hint { color:var(--muted); }
  form { display:grid; gap:15px; }
  label { display:grid; gap:7px; color:var(--muted); font-size:12px; font-weight:750; }
  input { min-height:44px; padding:0 12px; border:1px solid var(--border); border-radius:8px; color:var(--text); background:var(--surface-2); }
  input:focus { border-color:var(--accent); outline:0; }
  #turnstile-widget { min-height:66px; }
  .password-hint,.auth-foot { display:flex; gap:7px; align-items:center; color:var(--muted); font-size:11px; }
  .notice,.form-error { padding:10px 12px; border-radius:8px; color:var(--warning); background:var(--warning-soft); font-size:12px; }
  .form-error { color:var(--danger); background:var(--danger-soft); }
  .auth-foot { margin-top:22px; padding-top:16px; border-top:1px solid var(--border); }
  .switch-auth { width:100%; margin-top:18px; border:0; color:var(--accent); background:transparent; font-size:12px; font-weight:750; }
</style>
