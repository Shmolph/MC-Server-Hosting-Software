<script>
  import { onMount, onDestroy, createEventDispatcher } from "svelte";
  import { TURNSTILE_SITE_KEY } from "../config.js";
  const dispatch = createEventDispatcher();
  let container;
  let widgetId = null;
  let token = "";
  let ready = false;
  let resolveFresh = null;
  let rejectedFresh = null;
  let scriptPromise;

  function loadScript() {
    if (window.turnstile) return Promise.resolve();
    if (scriptPromise) return scriptPromise;
    scriptPromise = new Promise((resolve, reject) => {
      const existing = document.querySelector('script[data-turnstile="true"]');
      const script = existing || document.createElement("script");
      script.src = "https://challenges.cloudflare.com/turnstile/v0/api.js?render=explicit";
      script.async = true;
      script.defer = true;
      script.dataset.turnstile = "true";
      script.addEventListener("load", resolve, { once: true });
      script.addEventListener("error", reject, { once: true });
      if (!existing) document.head.appendChild(script);
    });
    return scriptPromise;
  }

  onMount(async () => {
    try {
      await loadScript();
      if (!container) return;
      widgetId = window.turnstile.render(container, {
        sitekey: TURNSTILE_SITE_KEY,
        theme: document.documentElement.dataset.theme === "light" ? "light" : "dark",
        callback(value) {
          token = value;
          dispatch("token", value);
          resolveFresh?.(value);
          resolveFresh = null;
        },
        "expired-callback"() { token = ""; dispatch("token", ""); },
        "error-callback"() { token = ""; dispatch("token", ""); }
      });
      ready = true;
    } catch (error) {
      dispatch("error", "CAPTCHA could not load. Refresh and try again.");
    }
  });

  export function getToken() {
    return widgetId === null ? "" : window.turnstile.getResponse(widgetId);
  }

  export function reset() {
    token = "";
    dispatch("token", "");
    if (widgetId !== null) window.turnstile.reset(widgetId);
  }

  export function freshToken() {
    if (widgetId === null) return Promise.reject(new Error("CAPTCHA is not ready."));
    return new Promise((resolve, reject) => {
      resolveFresh = resolve;
      rejectedFresh = reject;
      reset();
    });
  }

  onDestroy(() => {
    if (widgetId !== null && window.turnstile) window.turnstile.remove(widgetId);
    if (rejectedFresh) rejectedFresh(new Error("CAPTCHA was closed."));
  });
</script>

<div class="captcha"><div class="widget" bind:this={container}></div>{#if !ready}<span class="captcha-status">Loading security check…</span>{/if}</div>

<style>
  .captcha { display:grid; min-height:76px; justify-items:center; gap:8px; }
  .widget { display:grid; justify-items:center; max-width:100%; }
  .captcha-status { color:var(--muted); font-size:12px; }
</style>
