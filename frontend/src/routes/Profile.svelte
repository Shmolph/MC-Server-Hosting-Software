<script>
  import Button from "../lib/components/Button.svelte";
  export let onLogout;
</script>

<div class="page-head"><div><p class="kicker">Profile</p><h1>Account</h1></div></div>
<section class="profile-card"><div class="profile-row"><span class="avatar">S</span><div><strong>Signed in</strong><span>Username unavailable</span></div></div><section><h2>Change password <small>Soon</small></h2><Button variant="secondary" disabled>Change password</Button></section><section><h2>Active session</h2><Button variant="secondary" on:click={onLogout}>Log out</Button></section><section class="danger-section"><h2>Danger zone</h2><Button variant="danger" disabled>Delete account · Soon</Button></section></section>
