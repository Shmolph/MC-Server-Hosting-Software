<script>
  import { onMount, onDestroy, tick } from "svelte";
  import { ChevronDown, Compass, Crown, Eye, Gem, HeartPulse, Leaf, Maximize2, Minimize2, Search, Send, Shield, Skull, Sparkles, Trash2, UserRound, Utensils, X } from "lucide-svelte";
  import Badge from "../lib/components/Badge.svelte";
  import Button from "../lib/components/Button.svelte";
  import CopyButton from "../lib/components/CopyButton.svelte";
  import Modal from "../lib/components/Modal.svelte";
  import { api } from "../lib/api.js";
  import { pushToast } from "../lib/stores.js";

  export let server;
  export let onRefreshStatus = async () => {};

  const gamemodes = [
    { id: "survival", label: "Survival", icon: Leaf, color: "survival" },
    { id: "creative", label: "Creative", icon: Sparkles, color: "creative" },
    { id: "adventure", label: "Adventure", icon: Compass, color: "adventure" },
    { id: "spectator", label: "Spectator", icon: Eye, color: "spectator" }
  ];
  const popularItems = ["diamond", "netherite_ingot", "diamond_sword", "netherite_sword", "bow", "arrow", "golden_apple", "enchanted_golden_apple", "cooked_beef", "bread", "torch", "oak_log", "cobblestone", "tnt", "ender_pearl", "elytra", "shield", "totem_of_undying", "diamond_pickaxe", "diamond_axe", "diamond_shovel", "diamond_helmet", "diamond_chestplate", "diamond_leggings", "diamond_boots", "water_bucket", "lava_bucket", "firework_rocket", "oak_planks", "iron_ingot", "gold_ingot", "emerald", "redstone", "obsidian", "crafting_table", "furnace", "chest", "ender_chest", "beacon", "command_block"];
  const reasons = ["Griefing", "Cheating", "Spam", "Be respectful", "AFK", "Cooling off, come back later"];
  const ansiPattern = /\u001b(?:\[[0-?]*[ -/]*[@-~]|\][^\u0007]*(?:\u0007|\u001b\\))/g;
  const validPlayerName = (name) => /^[A-Za-z0-9_]{1,16}$/.test(String(name || ""));

  let activeTab = "players";
  let playersData = { max: 0, online: [], banned: [] };
  let playersOffline = false;
  let playersBusy = false;
  let playerSearch = "";
  let playerBusy = {};
  let serverActionBusy = "";
  let broadcastText = "";
  let senderName = "Server";
  let managePlayer = null;
  let managePanelElement;
  let manageReturnFocus;
  let modalKind = "";
  let modalPlayer = null;
  let reason = "";
  let confirmMessage = "";
  let itemSearch = "";
  let selectedItem = "diamond";
  let customItem = "";
  let itemCount = 1;
  let teleportTarget = "";
  let privateMessage = "";
  let consoleText = "";
  let consoleNotSetup = false;
  let consoleLines = [];
  let sessionLines = [];
  let consoleInput = "";
  let consoleBusy = false;
  let commandHistory = [];
  let historyIndex = 0;
  let consoleExpanded = false;
  let jumpVisible = false;
  let consoleOutput;
  let consoleInputElement;
  let playersTimer;
  let logsTimer;
  let playersRequestInFlight = false;
  let playersRefreshPending = false;
  let logsRequestInFlight = false;
  let gamemodeCache = {};
  let mounted = false;
  let currentPollKey = "";
  let previousServerStatus = "";
  let consoleWasAtBottom = true;
  let userScrolledDuringFetch = false;

  $: onlinePlayers = Array.isArray(playersData.online) ? playersData.online : [];
  $: bannedPlayers = Array.isArray(playersData.banned) ? playersData.banned : [];
  $: filteredPlayers = onlinePlayers.filter((player) => String(player.name).toLowerCase().includes(playerSearch.trim().toLowerCase()));
  $: filteredItems = popularItems.filter((item) => readableItem(item).toLowerCase().includes(itemSearch.trim().toLowerCase()) || item.includes(itemSearch.trim().toLowerCase()));
  $: selectedCustomItem = normalizeCustomItem(customItem);
  $: itemIsValid = customItem.trim() ? Boolean(selectedCustomItem) : Boolean(selectedItem);
  $: pendingAction = modalKind === "kick" || modalKind === "ban" ? modalKind : modalKind === "confirm-op" ? "op" : modalKind === "confirm-clear" ? "clear" : modalKind === "confirm-kill" ? "kill" : modalKind === "give" ? "give" : "";
  $: if (mounted) configurePolling(activeTab, server?.status);
  $: if (mounted && server?.status !== previousServerStatus) {
    if (server?.status === "starting") activeTab = "console";
    else if (server?.status === "online" && previousServerStatus === "starting") activeTab = "players";
    previousServerStatus = server?.status || "";
  }

  onMount(() => {
    mounted = true;
    document.addEventListener("visibilitychange", onVisibilityChange);
    document.addEventListener("keydown", onDocumentKeydown);
    return () => {
      mounted = false;
      clearInterval(playersTimer);
      clearInterval(logsTimer);
      document.removeEventListener("visibilitychange", onVisibilityChange);
      document.removeEventListener("keydown", onDocumentKeydown);
    };
  });

  onDestroy(() => {
    clearInterval(playersTimer);
    clearInterval(logsTimer);
  });

  function configurePolling(tab, status) {
    const key = `${tab}:${status}`;
    if (key === currentPollKey) return;
    currentPollKey = key;
    clearInterval(playersTimer);
    clearInterval(logsTimer);
    if (tab === "players" && status === "online") {
      fetchPlayers();
      playersTimer = setInterval(() => { if (document.visibilityState === "visible") fetchPlayers(); }, 5000);
    }
    if (tab === "console" && (status === "online" || status === "starting")) {
      fetchLogs();
      logsTimer = setInterval(() => { if (document.visibilityState === "visible") fetchLogs(); }, 1500);
    }
  }

  function onVisibilityChange() {
    if (document.visibilityState !== "visible") return;
    if (activeTab === "players" && server.status === "online") fetchPlayers();
    if (activeTab === "console" && (server.status === "online" || server.status === "starting")) fetchLogs();
  }

  function onDocumentKeydown(event) {
    if (event.key === "Escape" && consoleExpanded) consoleExpanded = false;
    else if (event.key === "Escape" && managePlayer && !modalKind) closeManage();
  }

  function selectTab(nextTab) {
    if (nextTab === "players" && server.status !== "online") return;
    if (nextTab === "console" && server.status !== "online" && server.status !== "starting") return;
    activeTab = nextTab;
  }

  function onTabKeydown(event) {
    if (!['ArrowLeft', 'ArrowRight', 'Home', 'End'].includes(event.key)) return;
    event.preventDefault();
    const available = ["players", "console"].filter((tab) => tab === "players" ? server.status === "online" : server.status === "online" || server.status === "starting");
    if (!available.length) return;
    let index = available.indexOf(activeTab);
    if (event.key === "Home") index = 0;
    else if (event.key === "End") index = available.length - 1;
    else index = (index + (event.key === "ArrowRight" ? 1 : -1) + available.length) % available.length;
    activeTab = available[index];
    document.getElementById(`server-tab-${activeTab}`)?.focus();
  }

  async function fetchPlayers() {
    if (playersRequestInFlight || activeTab !== "players" || document.visibilityState !== "visible" || server.status !== "online") return;
    playersRequestInFlight = true;
    try {
      const listed = await api.sendServerCommand(server.id, "list");
      if (listed.status === 401) return;
      if (listed.status === 409) { playersOffline = true; return; }
      if (listed.status !== 200) return;
      const listText = String(listed.text || "").trim();
      const maxMatch = listText.match(/max of\s+(\d+)/i);
      const namesText = listText.includes(":") ? listText.slice(listText.lastIndexOf(":") + 1).trim() : "";
      const names = namesText ? namesText.split(",").map((name) => name.trim()).filter(Boolean) : [];
      const nextOnline = [];
      for (const name of names) {
        const mode = validPlayerName(name) ? gamemodeCache[name] || "unknown" : "unknown";
        nextOnline.push({ name, gamemode: mode, op: false });
      }
      for (const player of nextOnline) {
        if (!validPlayerName(player.name)) continue;
        const modeResult = await api.sendServerCommand(server.id, `data get entity ${player.name} playerGameType`);
        if (modeResult.status === 401) return;
        if (modeResult.status === 409) { playersOffline = true; return; }
        if (modeResult.status !== 200) continue;
        const token = String(modeResult.text || "").trim().split(/\s+/).at(-1);
        const modes = { "0": "survival", "1": "creative", "2": "adventure", "3": "spectator" };
        if (Object.prototype.hasOwnProperty.call(modes, token)) gamemodeCache[player.name] = player.gamemode = modes[token];
      }
      const access = await api.getServerAccess(server.id);
      const ops = Array.isArray(access.ops) ? access.ops : [];
      const banned = Array.isArray(access.banned) ? access.banned : [];
      playersData = {
        max: Number(maxMatch?.[1]) || 0,
        online: nextOnline.map((player) => ({ ...player, op: ops.some((op) => String(op).toLowerCase() === String(player.name).toLowerCase()) })),
        banned
      };
      playersOffline = false;
      if (managePlayer) managePlayer = playersData.online.find((player) => player.name === managePlayer.name) || managePlayer;
    } catch (error) {
      if (error.status === 409) playersOffline = true;
    } finally {
      playersRequestInFlight = false;
      if (playersRefreshPending) {
        playersRefreshPending = false;
        if (activeTab === "players" && document.visibilityState === "visible") setTimeout(fetchPlayers, 0);
      }
    }
  }

  async function fetchLogs() {
    if (logsRequestInFlight || !["online", "starting"].includes(server.status)) return;
    logsRequestInFlight = true;
    consoleWasAtBottom = isConsoleAtBottom();
    userScrolledDuringFetch = false;
    try {
      consoleText = await api.getServerLogs(server.id);
      consoleNotSetup = false;
      await updateLogDisplay();
    } catch (error) {
      if (error.status === 409) {
        consoleNotSetup = true;
        consoleText = "";
        consoleLines = [];
      }
    } finally {
      logsRequestInFlight = false;
    }
  }

  async function updateLogDisplay() {
    const plainLines = consoleText.replace(ansiPattern, "").split(/\r?\n/).filter((line) => !line.includes("Thread RCON Client"));
    let inheritedLevel = "normal";
    consoleLines = plainLines.map((text, index) => {
      const upper = text.toUpperCase();
      const continuation = /^\s*at\s|^Caused by\b/.test(text);
      const initLine = text.includes("[init]");
      if (!continuation && !initLine) {
        if (/ERROR|FATAL|SEVERE|EXCEPTION/.test(upper)) inheritedLevel = "error";
        else if (/WARN/.test(upper)) inheritedLevel = "warning";
        else inheritedLevel = "normal";
      }
      const prefix = text.match(/^(\[\d{2}:\d{2}:\d{2}(?:\s+[^\]]+)?\]:?\s*)/);
      return { key: `${index}:${text}`, text, prefix: prefix?.[0] || "", rest: prefix ? text.slice(prefix[0].length) : text, level: initLine ? "init" : inheritedLevel };
    });
    await tick();
    if (consoleWasAtBottom && !userScrolledDuringFetch && consoleOutput) consoleOutput.scrollTop = consoleOutput.scrollHeight;
    else updateJumpVisibility();
  }

  function isConsoleAtBottom() {
    if (!consoleOutput) return true;
    return consoleOutput.scrollHeight - consoleOutput.scrollTop - consoleOutput.clientHeight < 36;
  }

  function updateJumpVisibility() {
    jumpVisible = !isConsoleAtBottom();
  }

  function handleConsoleScroll() {
    if (logsRequestInFlight) userScrolledDuringFetch = true;
    updateJumpVisibility();
  }

  async function runPlayerCommand(player, command, successMessage, busyLabel = "command") {
    if (!validPlayerName(player?.name)) {
      pushToast("Unsupported name", "error");
      return { status: 400 };
    }
    const key = player.name;
    if (playerBusy[key]) return { status: 0 };
    playerBusy = { ...playerBusy, [key]: busyLabel };
    try {
      appendSessionCommand(command);
      const result = await api.sendServerCommand(server.id, command);
      if (result.status === 200) {
        const reply = String(result.text || "").trim();
        if (reply) appendSessionReply(reply);
        const chatMessage = busyLabel === "msg";
        const showReply = reply && (!chatMessage || isChatReplyError(reply));
        pushToast(showReply ? `Server says: ${truncate(reply, 120)}` : successMessage, showReply ? "info" : "success");
        schedulePlayerRefresh();
      } else {
        showCommandError(result.status);
      }
      return result;
    } catch (error) {
      pushToast("Failed to send command", "error");
      return { status: 0 };
    } finally {
      const next = { ...playerBusy };
      delete next[key];
      playerBusy = next;
    }
  }

  async function runServerAction(command, label, chatMessage = false) {
    if (serverActionBusy) return;
    serverActionBusy = label;
    try {
      appendSessionCommand(command);
      const result = await api.sendServerCommand(server.id, command);
      if (result.status === 200) {
        const reply = String(result.text || "").trim();
        if (reply) appendSessionReply(reply);
        const showReply = reply && (!chatMessage || isChatReplyError(reply));
        pushToast(showReply ? `Server says: ${truncate(reply, 120)}` : label, showReply ? "info" : "success");
        schedulePlayerRefresh();
      } else showCommandError(result.status);
    } catch (error) {
      pushToast("Failed to send command", "error");
    } finally {
      serverActionBusy = "";
    }
  }

  function showCommandError(status) {
    if (status === 400) pushToast("That command was rejected", "error");
    else if (status === 409) { pushToast("Server isn't online", "error"); onRefreshStatus(); }
    else if (status === 404) { pushToast("Server not found", "error"); onRefreshStatus(); }
    else if (status !== 401) pushToast("Failed to send command", "error");
  }

  function schedulePlayerRefresh() {
    setTimeout(() => {
      if (playersRequestInFlight) playersRefreshPending = true;
      else fetchPlayers();
    }, 800);
  }

  function appendSession(command, reply) {
    appendSessionCommand(command);
    if (reply) appendSessionReply(reply);
  }

  function appendSessionCommand(command) {
    appendSessionRow({ key: `${Date.now()}-${sessionLines.length}-command`, text: `> ${command}`, kind: "command" });
  }

  function appendSessionReply(reply) {
    appendSessionRow({ key: `${Date.now()}-${sessionLines.length}-reply`, text: reply, kind: "reply" });
  }

  function appendSessionRow(line) {
    const wasAtBottom = isConsoleAtBottom();
    sessionLines = [...sessionLines, line];
    tick().then(() => {
      if (activeTab === "console" && consoleOutput && wasAtBottom) consoleOutput.scrollTop = consoleOutput.scrollHeight;
      else updateJumpVisibility();
    });
  }

  function truncate(text, max) {
    return text.length > max ? `${text.slice(0, max - 1)}…` : text;
  }

  function readableItem(item) {
    return item.split("_").map((part) => part.charAt(0).toUpperCase() + part.slice(1)).join(" ");
  }

  function normalizeCustomItem(value) {
    const trimmed = value.trim().toLowerCase();
    if (!trimmed || !/^[a-z0-9_]+(:[a-z0-9_/.]+)?$/.test(trimmed)) return "";
    return trimmed.includes(":") ? trimmed : `minecraft:${trimmed}`;
  }

  function cleanReason(event) {
    reason = event.currentTarget.value.replace(/[\r\n]+/g, " ").slice(0, 120);
  }

  function cleanChatMessage(event, field) {
    const cleaned = event.currentTarget.value.replace(/[\r\n]/g, "").slice(0, 150);
    if (field === "broadcast") broadcastText = cleaned;
    else privateMessage = cleaned;
  }

  function cleanSenderName(event) {
    senderName = event.currentTarget.value.replace(/[^A-Za-z0-9_ ]/g, "").slice(0, 16);
  }

  function chatSender() {
    return senderName.trim() || "Server";
  }

  function tellrawCommand(target, prefix, color, style, message) {
    const cleanMessage = message.replace(/[\r\n]/g, "").trim();
    if (!cleanMessage || cleanMessage.length > 150) return "";
    const component = [
      { text: prefix, color, [style]: true },
      { text: cleanMessage, color: "white" }
    ];
    return `tellraw ${target} ${JSON.stringify(component)}`;
  }

  function sendBroadcast() {
    const command = tellrawCommand("@a", `[${chatSender()}] `, "gold", "bold", broadcastText);
    if (!command || command.length >= 256) {
      pushToast("Message is too long to send", "error");
      return;
    }
    runServerAction(command, "Message sent", true);
    broadcastText = "";
  }

  function isChatReplyError(reply) {
    return /unknown|invalid|expected|incorrect argument/i.test(reply);
  }

  function commandReason(base, value) {
    const trimmed = value.replace(/[\r\n]+/g, " ").trim();
    return `${base}${trimmed ? ` ${trimmed}` : ""}`;
  }

  function requestPlayerAction(player, action) {
    if (!validPlayerName(player?.name) || playerBusy[player?.name]) return;
    modalPlayer = player;
    reason = "";
    if (action === "kick" || action === "ban") { modalKind = action; return; }
    if (action === "op") { confirmMessage = "This gives them full control of the server."; modalKind = "confirm-op"; return; }
    if (action === "clear" || action === "kill") { modalKind = `confirm-${action}`; return; }
    if (action === "give") { itemSearch = ""; selectedItem = "diamond"; customItem = ""; itemCount = 1; modalKind = "give"; return; }
    performPlayerAction(player, action);
  }

  async function performPlayerAction(player, action, argument = "") {
    if (!validPlayerName(player?.name)) return;
    const name = player.name;
    const commandBuilders = {
      op: () => `op ${name}`,
      deop: () => `deop ${name}`,
      gamemode: () => gamemodes.some((item) => item.id === argument) ? `gamemode ${argument} ${name}` : "",
      kick: () => commandReason(`kick ${name}`, argument),
      ban: () => commandReason(`ban ${name}`, argument),
      unban: () => validPlayerName(argument) ? `pardon ${argument}` : "",
      heal: () => `effect give ${name} minecraft:instant_health 1 10`,
      feed: () => `effect give ${name} minecraft:saturation 1 10`,
      clear: () => `clear ${name}`,
      kill: () => `kill ${name}`,
      give: () => {
        const item = customItem.trim() ? selectedCustomItem : `minecraft:${selectedItem}`;
        const count = Number(itemCount);
        return itemIsValid && Number.isInteger(count) && count >= 1 && count <= 64 ? `give ${name} ${item} ${count}` : "";
      },
      teleport: () => validPlayerName(argument) && argument !== name ? `tp ${name} ${argument}` : "",
      msg: () => validPlayerName(name) ? tellrawCommand(name, `[${chatSender()} -> You] `, "light_purple", "italic", privateMessage) : ""
    };
    const command = commandBuilders[action]?.();
    if (!command || command.length >= 256 || command.includes("\n") || command.startsWith("-")) {
      pushToast("That command was rejected", "error");
      return;
    }
    const success = {
      op: `${name} is now OP`, deop: `${name} is no longer OP`,
      gamemode: `${name}'s gamemode changed to ${argument}`,
      kick: `${name} was kicked`, ban: `${name} was banned`, unban: `${argument} was unbanned`,
      heal: `Healed ${name}`, feed: `Fed ${name}`, clear: `Cleared ${name}'s inventory`, kill: `${name} was killed`,
      give: `Gave ${Number(itemCount)} x ${readableItem(customItem.trim() ? selectedCustomItem.split(":").at(-1) : selectedItem)} to ${name}`,
      teleport: `Teleported ${name} to ${argument}`, msg: "Message sent"
    };
    await runPlayerCommand(player, command, success[action] || "Command sent", action === "gamemode" ? `gamemode:${argument}` : action);
    if (action === "msg") privateMessage = "";
    if (action === "kick" || action === "ban") modalKind = "";
    if (action === "unban") setTimeout(fetchPlayers, 800);
    if (managePlayer?.name === name && ["kick", "ban"].includes(action)) closeManage();
  }

  async function confirmModalAction() {
    if (!modalPlayer) return;
    const player = modalPlayer;
    const kind = modalKind;
    if (playerBusy[player.name]) return;
    if (kind === "kick") await performPlayerAction(player, "kick", reason);
    else if (kind === "ban") await performPlayerAction(player, "ban", reason);
    else if (kind === "confirm-op") await performPlayerAction(player, "op");
    else if (kind === "confirm-clear") await performPlayerAction(player, "clear");
    else if (kind === "confirm-kill") await performPlayerAction(player, "kill");
    else if (kind === "give") await performPlayerAction(player, "give");
    if (["confirm-op", "confirm-clear", "confirm-kill", "give"].includes(kind)) modalKind = "";
  }

  function playerColor(name) {
    let hash = 0;
    for (const char of String(name)) hash = (hash * 31 + char.charCodeAt(0)) | 0;
    return `hsl(${Math.abs(hash) % 360} 54% 42%)`;
  }

  function startManage(player) {
    if (!validPlayerName(player.name)) return;
    manageReturnFocus = document.activeElement;
    managePlayer = player;
    teleportTarget = onlinePlayers.find((item) => item.name !== player.name && validPlayerName(item.name))?.name || "";
    privateMessage = "";
    tick().then(() => managePanelElement?.querySelector("button:not(:disabled)")?.focus());
  }

  function closeManage() {
    const returnFocus = manageReturnFocus;
    manageReturnFocus = null;
    managePlayer = null;
    tick().then(() => returnFocus?.focus?.());
  }

  function trapManageFocus(event) {
    if (event.key !== "Tab" || !managePanelElement) return;
    const focusable = [...managePanelElement.querySelectorAll("button:not(:disabled),input:not(:disabled),select:not(:disabled),textarea:not(:disabled)")];
    if (!focusable.length) return;
    const first = focusable[0];
    const last = focusable[focusable.length - 1];
    if (event.shiftKey && document.activeElement === first) { event.preventDefault(); last.focus(); }
    else if (!event.shiftKey && document.activeElement === last) { event.preventDefault(); first.focus(); }
  }

  async function sendConsoleCommand(command = consoleInput) {
    const trimmed = command.trim();
    if (!trimmed || consoleBusy || server.status !== "online") return;
    if (trimmed.length > 256 || /[\r\n]/.test(trimmed) || trimmed.startsWith("-")) {
      sessionLines = [...sessionLines, { key: `${Date.now()}-invalid`, text: "Rejected: invalid command", kind: "error" }];
      return;
    }
    consoleBusy = true;
    consoleInput = "";
    commandHistory = [...commandHistory.filter((item) => item !== trimmed), trimmed];
    historyIndex = commandHistory.length;
    appendSession(trimmed, "");
    try {
      const result = await api.sendServerCommand(server.id, trimmed);
      if (result.status === 200) {
        const reply = String(result.text || "").trim();
        if (reply) {
          sessionLines = [...sessionLines, { key: `${Date.now()}-reply`, text: reply, kind: "reply" }];
          pushToast(`Server says: ${truncate(reply, 120)}`, "info");
        }
      } else {
        const errors = { 400: "Rejected: invalid command", 409: "Server is not online.", 404: "Server not found.", 500: "Failed to send command." };
        const message = errors[result.status] || "Failed to send command.";
        sessionLines = [...sessionLines, { key: `${Date.now()}-${result.status}`, text: message, kind: "error" }];
      }
    } catch (error) {
      sessionLines = [...sessionLines, { key: `${Date.now()}-failed`, text: "Failed to send command.", kind: "error" }];
    } finally {
      consoleBusy = false;
      await tick();
      consoleInputElement?.focus();
      if (consoleOutput && isConsoleAtBottom()) consoleOutput.scrollTop = consoleOutput.scrollHeight;
    }
  }

  function historyKeydown(event) {
    if (event.key !== "ArrowUp" && event.key !== "ArrowDown") return;
    event.preventDefault();
    if (!commandHistory.length) return;
    historyIndex = event.key === "ArrowUp" ? Math.max(0, historyIndex - 1) : Math.min(commandHistory.length, historyIndex + 1);
    consoleInput = commandHistory[historyIndex] || "";
  }

  function openItemDialog() {
    itemSearch = "";
    selectedItem = "diamond";
    customItem = "";
    itemCount = 1;
    modalKind = "give";
  }

  function selectMode(player, mode) {
    if (String(player.gamemode).toLowerCase() === mode) return;
    performPlayerAction(player, "gamemode", mode);
  }

  function toggleOp(player) {
    if (player.op) performPlayerAction(player, "deop");
    else requestPlayerAction(player, "op");
  }

  function statusTabHint(tab) {
    if (tab === "players") return "Players are available when the server is online.";
    return "Console logs are available while the server is starting or online.";
  }
</script>

<section class="server-tools">
  <div class="tools-tabs" role="tablist" aria-label="Server tools" tabindex="-1" on:keydown={onTabKeydown}>
    <button id="server-tab-players" class:active={activeTab === "players"} type="button" role="tab" aria-selected={activeTab === "players"} aria-controls="server-panel-players" disabled={server.status !== "online"} title={server.status === "online" ? "Players" : statusTabHint("players")} tabindex={activeTab === "players" ? 0 : -1} on:click={() => selectTab("players")}>Players ({onlinePlayers.length}/{playersData.max || 0})</button>
    <button id="server-tab-console" class:active={activeTab === "console"} type="button" role="tab" aria-selected={activeTab === "console"} aria-controls="server-panel-console" disabled={server.status !== "online" && server.status !== "starting"} title={server.status === "online" || server.status === "starting" ? "Console" : statusTabHint("console")} tabindex={activeTab === "console" ? 0 : -1} on:click={() => selectTab("console")}>Console</button>
  </div>
  {#if server.status !== "online"}<p class="tab-hint">{statusTabHint("players")}</p>{/if}

  {#if activeTab === "players"}
    <div id="server-panel-players" class="tab-panel" role="tabpanel" aria-labelledby="server-tab-players" tabindex="0">
      {#if server.status !== "online"}
        <div class="offline-empty"><UserRound size={22} /><p>Players are unavailable while the server is {server.status}.</p></div>
      {:else}
        <div class="server-quick-actions">
          <form class="broadcast-form" on:submit|preventDefault={sendBroadcast}>
            <label for="broadcast-message">Broadcast message</label><div class="broadcast-row"><div><input id="broadcast-message" bind:value={broadcastText} on:input={(event) => cleanChatMessage(event, "broadcast")} maxlength="150" placeholder="Message all players" aria-label="Broadcast message" /><span class="message-counter">{broadcastText.length}/150</span></div><div><label for="sender-name">Sender name</label><input id="sender-name" bind:value={senderName} on:input={cleanSenderName} maxlength="16" aria-label="Sender name" /></div><Button type="submit" variant="secondary" loading={serverActionBusy === "Message sent"} disabled={!broadcastText.trim() || Boolean(serverActionBusy)}>Broadcast</Button></div>
          </form>
          <div class="world-actions" aria-label="Server quick actions">
            <Button variant="secondary" loading={serverActionBusy === "Set day"} disabled={Boolean(serverActionBusy)} on:click={() => runServerAction("time set day", "Time set to day")}>Day</Button>
            <Button variant="secondary" loading={serverActionBusy === "Set night"} disabled={Boolean(serverActionBusy)} on:click={() => runServerAction("time set night", "Time set to night")}>Night</Button>
            <Button variant="secondary" loading={serverActionBusy === "Clear weather"} disabled={Boolean(serverActionBusy)} on:click={() => runServerAction("weather clear", "Weather cleared")}>Clear weather</Button>
            <Button variant="secondary" loading={serverActionBusy === "Set rain"} disabled={Boolean(serverActionBusy)} on:click={() => runServerAction("weather rain", "Rain started")}>Rain</Button>
            <Button variant="secondary" loading={serverActionBusy === "Save world"} disabled={Boolean(serverActionBusy)} on:click={() => runServerAction("save-all", "World saved")}>Save world</Button>
          </div>
        </div>
        <div class="players-heading"><div><h2>Online players</h2><span>{onlinePlayers.length} of {playersData.max || 0} slots</span></div>{#if onlinePlayers.length > 8}<label class="player-search"><Search size={16} /><input bind:value={playerSearch} placeholder="Search players" aria-label="Search players" /></label>{/if}</div>
        {#if playersOffline}<div class="offline-empty"><p>Server isn't online.</p></div>
        {:else if !onlinePlayers.length}<div class="offline-empty"><UserRound size={22} /><p>Nobody is online right now.</p><div class="connect-line"><code>localhost:{server.port}</code><CopyButton value={`localhost:${server.port}`} /></div></div>
        {:else if !filteredPlayers.length}<div class="offline-empty"><p>No players match that search.</p></div>
        {:else}<div class="player-grid">{#each filteredPlayers as player (player.name)}
          {@const supported = validPlayerName(player.name)}
          {@const busy = Boolean(playerBusy[player.name])}
          {@const mode = String(player.gamemode || "unknown").toLowerCase()}
          <article class="player-card">
            <header class="player-header"><span class="avatar" style={`--avatar-color:${playerColor(player.name)}`}>{String(player.name || "?").slice(0, 1).toUpperCase()}</span><div class="player-title"><strong>{player.name}</strong><span class="player-chips">{#if player.op}<span class="op-badge"><Crown size={12} />OP</span>{/if}{#if gamemodes.some((item) => item.id === mode)}<span class={`mode-chip mode-${mode}`}>{#each gamemodes.filter((item) => item.id === mode) as item}<svelte:component this={item.icon} size={13} />{item.label}{/each}</span>{:else}<span class="mode-chip mode-unknown">? Unknown</span>{/if}</span></div><Button variant="ghost" disabled={!supported || busy} title={supported ? "Manage player" : "Unsupported name"} on:click={() => startManage(player)}>Manage</Button></header>
            <div class="mode-switch" role="group" aria-label={`Gamemode for ${player.name}`} title={supported ? "Change gamemode" : "Unsupported name"}>{#each gamemodes as item}<button type="button" class:chosen={mode === item.id} class={`mode-button mode-${item.color}`} aria-label={item.label} aria-pressed={mode === item.id} aria-busy={playerBusy[player.name] === `gamemode:${item.id}`} disabled={!supported || busy} on:click={() => selectMode(player, item.id)}>{#if playerBusy[player.name] === `gamemode:${item.id}`}<span class="tiny-spinner" aria-hidden="true"></span>{:else}<svelte:component this={item.icon} size={15} />{/if}</button>{/each}</div>
            <div class="player-actions"><Button variant={player.op ? "secondary" : "ghost"} loading={playerBusy[player.name] === (player.op ? "deop" : "op")} disabled={!supported || busy} title={supported ? (player.op ? "Remove OP" : "Give OP") : "Unsupported name"} on:click={() => toggleOp(player)}>{player.op ? "Remove OP" : "Give OP"}</Button><Button variant="secondary" loading={playerBusy[player.name] === "kick"} disabled={!supported || busy} title={supported ? "Kick player" : "Unsupported name"} on:click={() => requestPlayerAction(player, "kick")}>Kick</Button><Button variant="danger" loading={playerBusy[player.name] === "ban"} disabled={!supported || busy} title={supported ? "Ban player" : "Unsupported name"} on:click={() => requestPlayerAction(player, "ban")}>Ban</Button></div>
          </article>
        {/each}</div>{/if}
        <details class="banned-section"><summary><span>Banned players</span><span>{bannedPlayers.length}<ChevronDown size={16} /></span></summary>{#if bannedPlayers.length}<ul>{#each bannedPlayers as banned (banned)}<li><span>{banned}</span><Button variant="secondary" disabled={!validPlayerName(banned) || Boolean(playerBusy[banned])} title={validPlayerName(banned) ? "Unban player" : "Unsupported name"} loading={playerBusy[banned] === "unban"} on:click={() => runPlayerCommand({ name: banned }, `pardon ${banned}`, `${banned} was unbanned`, "unban")}>Unban</Button></li>{/each}</ul>{:else}<p class="no-banned">No banned players.</p>{/if}</details>
      {/if}
    </div>
  {:else}
    <div id="server-panel-console" class:expanded={consoleExpanded} class="console-panel" role="tabpanel" aria-labelledby="server-tab-console" tabindex="0">
      <header class="console-titlebar"><div class="terminal-mark"><span class:online={server.status === "online"} class:starting={server.status === "starting"}></span><strong>Console</strong><span class="terminal-server-name">{server.name}</span></div><button class="expand-button" type="button" aria-label={consoleExpanded ? "Collapse console" : "Expand console"} title={consoleExpanded ? "Collapse console" : "Expand console"} on:click={() => consoleExpanded = !consoleExpanded}>{#if consoleExpanded}<Minimize2 size={17} />{:else}<Maximize2 size={17} />{/if}</button></header>
      {#if server.status !== "online" && server.status !== "starting"}<div class="console-empty">Console is available while the server is starting or online.</div>{:else if consoleNotSetup}<div class="console-empty">Set up the server to see logs.</div>{:else}<div class="console-output" bind:this={consoleOutput} on:scroll={handleConsoleScroll} aria-label="Server log output">{#each consoleLines as line (line.key)}<div class={`log-line level-${line.level}`}><span class="log-prefix">{line.prefix}</span>{line.rest}</div>{/each}{#each sessionLines as line (line.key)}<div class={`log-line session-${line.kind}`}>{line.text}</div>{/each}</div>{/if}
      {#if jumpVisible}<button class="jump-latest" type="button" on:click={() => { consoleOutput.scrollTop = consoleOutput.scrollHeight; jumpVisible = false; }}>Jump to latest</button>{/if}
      <form class="console-command" on:submit|preventDefault={() => sendConsoleCommand()}><span class="prompt-mark" aria-hidden="true">&gt;</span><input bind:this={consoleInputElement} bind:value={consoleInput} on:keydown={historyKeydown} maxlength="256" autocomplete="off" spellcheck="false" placeholder={server.status === "online" ? "Enter a command" : "Server must be online to send commands"} disabled={server.status !== "online" || consoleBusy} aria-label="Console command" /><Button type="submit" variant="primary" loading={consoleBusy} disabled={server.status !== "online" || consoleBusy || !consoleInput.trim()}><Send size={15} />Send</Button></form>
    </div>
  {/if}
</section>

{#if managePlayer}
  <div class="manage-scrim" role="presentation" on:click={(event) => event.target === event.currentTarget && closeManage()}>
    <div class="manage-panel" bind:this={managePanelElement} role="dialog" aria-modal="true" aria-labelledby="manage-title" tabindex="-1" on:keydown={trapManageFocus}><header><div><span class="eyebrow">PLAYER</span><h2 id="manage-title">{managePlayer.name}</h2></div><button class="icon-button" type="button" aria-label="Close manage panel" on:click={closeManage}><X size={18} /></button></header>
      <section class="manage-group"><h3>Permissions</h3><Button variant="secondary" loading={playerBusy[managePlayer.name] === (managePlayer.op ? "deop" : "op")} disabled={Boolean(playerBusy[managePlayer.name])} on:click={() => managePlayer.op ? performPlayerAction(managePlayer, "deop") : requestPlayerAction(managePlayer, "op")}>{managePlayer.op ? "Remove OP" : "Give OP"}</Button></section>
      <section class="manage-group"><h3>Gamemode</h3><div class="manage-modes">{#each gamemodes as item}<Button variant={String(managePlayer.gamemode).toLowerCase() === item.id ? "primary" : "secondary"} loading={playerBusy[managePlayer.name] === `gamemode:${item.id}`} disabled={Boolean(playerBusy[managePlayer.name])} on:click={() => selectMode(managePlayer, item.id)}><svelte:component this={item.icon} size={15} />{item.label}</Button>{/each}</div></section>
      <section class="manage-group"><h3>Items</h3><Button variant="secondary" disabled={Boolean(playerBusy[managePlayer.name])} on:click={openItemDialog}>Give item</Button></section>
      <section class="manage-group"><h3>Quick actions</h3><div class="manage-actions"><Button variant="secondary" loading={playerBusy[managePlayer.name] === "heal"} disabled={Boolean(playerBusy[managePlayer.name])} on:click={() => performPlayerAction(managePlayer, "heal")}><HeartPulse size={15} />Heal</Button><Button variant="secondary" loading={playerBusy[managePlayer.name] === "feed"} disabled={Boolean(playerBusy[managePlayer.name])} on:click={() => performPlayerAction(managePlayer, "feed")}><Utensils size={15} />Feed</Button><Button variant="secondary" loading={playerBusy[managePlayer.name] === "clear"} disabled={Boolean(playerBusy[managePlayer.name])} on:click={() => requestPlayerAction(managePlayer, "clear")}><Trash2 size={15} />Clear inventory</Button><Button variant="danger" loading={playerBusy[managePlayer.name] === "kill"} disabled={Boolean(playerBusy[managePlayer.name])} on:click={() => requestPlayerAction(managePlayer, "kill")}><Skull size={15} />Kill</Button></div>
        <label>Teleport to online player<select bind:value={teleportTarget}><option value="">Choose player</option>{#each onlinePlayers.filter((item) => item.name !== managePlayer.name && validPlayerName(item.name)) as other (other.name)}<option value={other.name}>{other.name}</option>{/each}</select></label><Button variant="secondary" disabled={!teleportTarget || Boolean(playerBusy[managePlayer.name])} on:click={() => performPlayerAction(managePlayer, "teleport", teleportTarget)}>Teleport</Button>
        <form class="private-message" on:submit|preventDefault={() => performPlayerAction(managePlayer, "msg")}><label for="private-message">Send private message</label><div><input id="private-message" bind:value={privateMessage} on:input={(event) => cleanChatMessage(event, "private")} maxlength="150" placeholder="Message" /><span class="message-counter">{privateMessage.length}/150</span><Button type="submit" variant="secondary" disabled={!privateMessage.trim() || Boolean(playerBusy[managePlayer.name])}>Send</Button></div></form>
      </section>
      <section class="manage-group"><h3>Moderation</h3><div class="manage-actions"><Button variant="secondary" disabled={Boolean(playerBusy[managePlayer.name])} on:click={() => requestPlayerAction(managePlayer, "kick")}>Kick</Button><Button variant="danger" disabled={Boolean(playerBusy[managePlayer.name])} on:click={() => requestPlayerAction(managePlayer, "ban")}>Ban</Button></div></section>
    </div>
  </div>
{/if}

<Modal open={Boolean(modalKind)} className="player-action-modal" labelledBy="player-action-title" on:close={() => modalKind = ""}>
  <div class="action-dialog">
    <header><div><span class="eyebrow">{modalKind === "give" ? "ITEMS" : "PLAYER ACTION"}</span><h2 id="player-action-title">{modalKind === "kick" ? `Kick ${modalPlayer?.name}` : modalKind === "ban" ? `Ban ${modalPlayer?.name}` : modalKind === "give" ? `Give item to ${modalPlayer?.name}` : modalKind === "confirm-op" ? `Give OP to ${modalPlayer?.name}?` : modalKind === "confirm-clear" ? `Clear ${modalPlayer?.name}'s inventory?` : `Kill ${modalPlayer?.name}?`}</h2></div><button class="icon-button" type="button" aria-label="Close dialog" on:click={() => modalKind = ""}><X size={18} /></button></header>
    {#if modalKind === "kick" || modalKind === "ban"}
      {#if modalKind === "ban"}<p class="dialog-warning">They can't rejoin until unbanned.</p>{/if}
      <label for="moderation-reason">Reason (optional)</label><textarea id="moderation-reason" rows="2" maxlength="120" bind:value={reason} on:input={cleanReason} placeholder="Add a reason"></textarea><div class="reason-count">{reason.length}/120</div><div class="reason-chips">{#each reasons as preset}<button type="button" on:click={() => reason = preset}>{preset}</button>{/each}</div>
    {:else if modalKind === "give"}
      <label for="item-search">Search items</label><div class="search-field"><Search size={16} /><input id="item-search" bind:value={itemSearch} placeholder="Find an item" /></div><div class="item-list" role="listbox" aria-label="Popular items">{#each filteredItems as item (item)}<button type="button" class:selected={selectedItem === item && !customItem} role="option" aria-selected={selectedItem === item && !customItem} on:click={() => { selectedItem = item; customItem = ""; }}>{readableItem(item)}</button>{/each}</div><label for="custom-item">Custom item id</label><input id="custom-item" bind:value={customItem} placeholder="minecraft:item_id" /><label for="item-count">Count</label><div class="count-control"><input id="item-count" type="number" min="1" max="64" bind:value={itemCount} /><div>{#each [1, 16, 64] as count}<button type="button" class:chosen={Number(itemCount) === count} on:click={() => itemCount = count}>{count}</button>{/each}</div></div>
    {:else}<p class="confirm-copy">{modalKind === "confirm-op" ? confirmMessage : modalKind === "confirm-clear" ? "This clears the player's inventory." : "This eliminates the player in-game."}</p>{/if}
    <footer><Button variant="ghost" disabled={Boolean(modalPlayer && playerBusy[modalPlayer.name])} on:click={() => modalKind = ""}>Cancel</Button><Button variant="danger" loading={Boolean(modalPlayer && playerBusy[modalPlayer.name] === pendingAction)} disabled={Boolean(modalPlayer && playerBusy[modalPlayer.name]) || (modalKind === "give" && (!itemIsValid || !Number.isInteger(Number(itemCount)) || Number(itemCount) < 1 || Number(itemCount) > 64))} on:click={confirmModalAction}>{modalKind === "give" ? "Give" : modalKind === "confirm-op" ? "Give OP" : modalKind === "confirm-clear" ? "Clear inventory" : modalKind === "confirm-kill" ? "Kill" : modalKind === "kick" ? "Kick player" : "Ban player"}</Button></footer>
  </div>
</Modal>

<style>
  .server-tools { margin-top:var(--space-7); }
  .tools-tabs { display:flex; gap:var(--space-1); border-bottom:1px solid var(--border); }
  .tools-tabs button { min-height:var(--target); padding:0 var(--space-4); border:0; border-bottom:2px solid transparent; color:var(--muted); background:transparent; font-size:var(--font-14); font-weight:750; cursor:pointer; }
  .tools-tabs button.active { color:var(--text); border-bottom-color:var(--accent); }
  .tools-tabs button:disabled { cursor:not-allowed; opacity:.45; }
  .tools-tabs button:focus-visible,.mode-button:focus-visible,.icon-button:focus-visible,.expand-button:focus-visible { outline:2px solid var(--accent); outline-offset:2px; }
  .tab-hint { margin:var(--space-2) 0 0; color:var(--muted); font-size:var(--font-12); }
  .tab-panel { padding-top:var(--space-4); }
  .server-quick-actions { display:grid; gap:var(--space-3); margin-bottom:var(--space-5); padding:var(--space-4); border:1px solid var(--border); border-radius:var(--radius-card); background:var(--surface-1); }
  .broadcast-form,.private-message { display:grid; gap:var(--space-2); }
  .broadcast-form label,.private-message label,.manage-group label { color:var(--muted); font-size:var(--font-12); font-weight:700; }
  .broadcast-form > div,.private-message > div { display:flex; gap:var(--space-2); }
  .broadcast-row { flex-wrap:wrap; align-items:end; }
  .broadcast-row > div { display:grid; flex:1 1 180px; min-width:0; gap:var(--space-1); }
  .broadcast-row > div > label { color:var(--muted); font-size:var(--font-12); font-weight:700; }
  .message-counter { flex:0 0 auto; align-self:center; color:var(--muted); font-size:var(--font-11); white-space:nowrap; }
  input,textarea,select { min-width:0; min-height:var(--target); padding:0 var(--space-3); border:1px solid var(--border); border-radius:var(--radius-control); color:var(--text); background:var(--surface-2); font:inherit; }
  textarea { width:100%; min-height:76px; padding:var(--space-3); resize:vertical; }
  .broadcast-form input,.private-message input { flex:1; }
  .world-actions { display:flex; flex-wrap:wrap; gap:var(--space-2); }
  .players-heading { display:flex; flex-wrap:wrap; gap:var(--space-3); align-items:center; justify-content:space-between; margin:var(--space-4) 0; }
  .players-heading h2 { margin:0; color:var(--text); font-size:var(--font-18); }
  .players-heading span { color:var(--muted); font-size:var(--font-12); }
  .player-search,.search-field { display:flex; min-width:200px; gap:var(--space-2); align-items:center; padding:0 var(--space-3); border:1px solid var(--border); border-radius:var(--radius-control); color:var(--muted); background:var(--surface-2); }
  .player-search input,.search-field input { width:100%; min-height:var(--target); padding:0; border:0; outline:0; background:transparent; }
  .player-grid { display:grid; grid-template-columns:repeat(2,minmax(0,1fr)); gap:var(--space-3); }
  .player-card { display:grid; gap:var(--space-3); min-width:0; padding:var(--space-4); border:1px solid var(--border); border-radius:var(--radius-card); background:var(--surface-1); }
  .player-header { display:flex; min-width:0; gap:var(--space-3); align-items:center; }
  .avatar { display:grid; width:40px; height:40px; flex:0 0 40px; place-items:center; border-radius:var(--radius-control); color:white; background:var(--avatar-color); font-weight:800; }
  .player-title { display:grid; flex:1; min-width:0; gap:var(--space-1); }
  .player-title > strong { overflow:hidden; color:var(--text); text-overflow:ellipsis; white-space:nowrap; }
  .player-chips { display:flex; flex-wrap:wrap; gap:var(--space-1); }
  .op-badge,.mode-chip { display:inline-flex; min-height:22px; gap:4px; align-items:center; padding:0 7px; border-radius:var(--radius-pill); font-size:10px; font-weight:750; }
  .op-badge { color:#f3c85b; background:rgba(243,200,91,.16); }
  .mode-survival { color:#68cc8e; background:rgba(72,190,116,.14); }
  .mode-creative { color:#78b9ff; background:rgba(91,157,255,.16); }
  .mode-adventure { color:#e6a765; background:rgba(225,148,71,.15); }
  .mode-spectator { color:#c4a4ec; background:rgba(162,123,213,.16); }
  .mode-unknown { color:var(--muted); background:var(--muted-soft); }
  .mode-switch { display:grid; grid-template-columns:repeat(4,1fr); overflow:hidden; border:1px solid var(--border); border-radius:var(--radius-control); background:var(--surface-2); }
  .mode-button { display:grid; min-height:44px; place-items:center; border:0; border-right:1px solid var(--border); color:var(--muted); background:transparent; cursor:pointer; }
  .mode-button:last-child { border-right:0; }
  .mode-button.chosen { color:var(--text); background:var(--surface-3); box-shadow:inset 0 -2px currentColor; }
  .tiny-spinner { width:14px; height:14px; border:2px solid currentColor; border-right-color:transparent; border-radius:50%; animation:tiny-spin 700ms linear infinite; }
  @keyframes tiny-spin { to { transform:rotate(360deg); } }
  .mode-button.mode-survival { color:#68cc8e; }.mode-button.mode-creative { color:#78b9ff; }.mode-button.mode-adventure { color:#e6a765; }.mode-button.mode-spectator { color:#c4a4ec; }
  .mode-button:disabled,:global(.player-actions button:disabled) { cursor:not-allowed; opacity:.45; }
  .player-actions { display:grid; grid-template-columns:repeat(3,minmax(0,1fr)); gap:var(--space-2); }
  .offline-empty { display:grid; min-height:160px; justify-items:center; align-content:center; gap:var(--space-2); color:var(--muted); text-align:center; }
  .offline-empty p { margin:0; }
  .connect-line { display:flex; gap:var(--space-2); align-items:center; }
  .connect-line code { color:var(--text); font-family:var(--mono-font); }
  .banned-section { margin-top:var(--space-5); border-top:1px solid var(--border); border-bottom:1px solid var(--border); }
  .banned-section summary { display:flex; min-height:var(--target); align-items:center; justify-content:space-between; color:var(--text); font-weight:700; cursor:pointer; list-style:none; }
  .banned-section summary::-webkit-details-marker { display:none; }
  .banned-section summary > span:last-child { display:flex; gap:var(--space-2); align-items:center; color:var(--muted); }
  .banned-section[open] summary :global(svg) { transform:rotate(180deg); }
  .banned-section ul { display:grid; gap:var(--space-2); margin:0; padding:0 0 var(--space-3); list-style:none; }
  .banned-section li { display:flex; min-height:var(--target); gap:var(--space-3); align-items:center; justify-content:space-between; padding:0 var(--space-3); border:1px solid var(--border); border-radius:var(--radius-control); color:var(--text); background:var(--surface-1); }
  .no-banned { color:var(--muted); font-size:var(--font-13); }
  .console-panel { position:relative; overflow:hidden; margin-top:var(--space-4); border:1px solid #242b31; border-radius:var(--radius-card); color:#d5dde5; background:#0a0d10; box-shadow:0 12px 32px rgba(0,0,0,.25); }
  .console-panel.expanded { position:fixed; z-index:30; inset:12px; display:flex; flex-direction:column; margin:0; border-radius:10px; }
  .console-titlebar { display:flex; min-height:48px; gap:var(--space-3); align-items:center; justify-content:space-between; padding:0 var(--space-4); border-bottom:1px solid #242b31; background:#11161b; }
  .terminal-mark { display:flex; min-width:0; gap:var(--space-2); align-items:center; font-size:var(--font-13); }
  .terminal-mark strong { color:#eff5fa; }
  .terminal-server-name { overflow:hidden; color:#82909b; text-overflow:ellipsis; white-space:nowrap; }
  .terminal-mark > span:first-child { width:8px; height:8px; border-radius:50%; background:#58616b; }
  .terminal-mark > span.online { background:#35c77a; box-shadow:0 0 10px rgba(53,199,122,.55); }
  .terminal-mark > span.starting { background:#e4af46; box-shadow:0 0 10px rgba(228,175,70,.45); }
  .expand-button,.icon-button { display:grid; width:var(--target); height:var(--target); place-items:center; border:1px solid #303841; border-radius:var(--radius-control); color:#aeb9c2; background:#171d23; cursor:pointer; }
  .console-output { overflow:auto; height:380px; padding:var(--space-4); font:12px/1.65 var(--mono-font); white-space:pre-wrap; overflow-wrap:anywhere; }
  .expanded .console-output { flex:1; height:auto; }
  .log-line { color:#ccd5dd; }
  .log-line.level-error,.session-error { color:#ff7777; }
  .log-line.level-warning { color:#eec45d; }
  .log-line.level-init { color:#9aa5ae; }
  .log-prefix { color:#78838d; }
  .session-command { color:#58d7e9; }
  .session-reply { color:#dce4eb; }
  .console-empty { display:grid; min-height:180px; place-items:center; color:#b9c4ce; font:13px var(--mono-font); }
  .jump-latest { position:absolute; right:var(--space-4); bottom:76px; min-height:38px; padding:0 var(--space-3); border:1px solid #36aaba; border-radius:var(--radius-control); color:#e4fbff; background:#12616d; cursor:pointer; }
  .console-command { display:flex; min-height:58px; gap:var(--space-2); align-items:center; padding:var(--space-2) var(--space-3); border-top:1px solid #242b31; background:#11161b; }
  .prompt-mark { color:#55d6e5; font:700 15px var(--mono-font); }
  .console-command input { flex:1; min-height:42px; border:0; color:#e4ebf0; background:transparent; font:12px var(--mono-font); outline:0; }
  .console-command input::placeholder { color:#65717b; }
  .manage-scrim { position:fixed; z-index:25; inset:0; display:flex; justify-content:flex-end; background:var(--scrim); }
  .manage-panel { overflow:auto; width:min(100%,480px); height:100%; padding:var(--space-5); border-left:1px solid var(--border); background:var(--surface); box-shadow:var(--shadow-lg); }
  .manage-panel > header,.action-dialog > header { display:flex; gap:var(--space-3); align-items:center; justify-content:space-between; margin-bottom:var(--space-4); }
  .manage-panel > header h2,.action-dialog h2 { margin:var(--space-1) 0 0; color:var(--text); font-size:var(--font-20); overflow-wrap:anywhere; }
  .eyebrow { color:var(--accent); font-family:var(--pixel-font); font-size:9px; }
  .manage-group { display:grid; gap:var(--space-3); padding:var(--space-4) 0; border-top:1px solid var(--border); }
  .manage-group h3 { margin:0; color:var(--text); font-size:var(--font-14); }
  .manage-modes,.manage-actions { display:flex; flex-wrap:wrap; gap:var(--space-2); }
  .manage-group label { display:grid; gap:var(--space-2); }
  .action-dialog { display:grid; gap:var(--space-3); padding:var(--space-5); }
  :global(dialog.player-action-modal) { margin:auto; }
  .action-dialog > label { color:var(--muted); font-size:var(--font-13); font-weight:700; }
  .dialog-warning { margin:0; padding:var(--space-3); border-left:3px solid var(--warning); color:var(--warning); background:var(--warning-soft); font-size:var(--font-13); }
  .reason-count { justify-self:end; margin-top:calc(-1 * var(--space-2)); color:var(--muted); font-size:var(--font-11); }
  .reason-chips { display:flex; flex-wrap:wrap; gap:var(--space-2); }
  .reason-chips button,.count-control button { min-height:44px; padding:0 var(--space-3); border:1px solid var(--border); border-radius:var(--radius-pill); color:var(--muted); background:var(--surface-2); cursor:pointer; }
  .reason-chips button:hover,.count-control button.chosen { border-color:var(--accent); color:var(--text); }
  .action-dialog footer { display:flex; flex-wrap:wrap; justify-content:flex-end; gap:var(--space-2); margin-top:var(--space-2); }
  .dialog-warning + label { margin-top:var(--space-2); }
  .item-list { display:grid; grid-template-columns:repeat(2,minmax(0,1fr)); gap:var(--space-1); max-height:180px; overflow:auto; padding:var(--space-1); border:1px solid var(--border); border-radius:var(--radius-control); }
  .item-list button { min-height:44px; padding:0 var(--space-2); border:1px solid transparent; border-radius:var(--radius-control); color:var(--muted); background:transparent; text-align:left; cursor:pointer; }
  .item-list button.selected { border-color:var(--accent); color:var(--text); background:var(--accent-soft); }
  .count-control { display:flex; gap:var(--space-3); align-items:center; }
  .count-control input { width:90px; }
  .count-control > div { display:flex; gap:var(--space-1); }
  .confirm-copy { margin:0; color:var(--muted); line-height:1.5; }
  .tab-panel :global(button:focus-visible),.manage-panel :global(button:focus-visible),.action-dialog :global(button:focus-visible),.console-command :global(button:focus-visible) { outline:2px solid var(--accent); outline-offset:2px; }
  @media(max-width:720px) { .player-grid { grid-template-columns:1fr; } .console-output { height:320px; } .server-quick-actions { padding:var(--space-3); } }
  @media(max-width:480px) { .players-heading { align-items:stretch; flex-direction:column; } .player-search { min-width:0; } .broadcast-form > div,.private-message > div { align-items:stretch; } .broadcast-row > div { flex-basis:100%; } .world-actions { display:grid; grid-template-columns:repeat(2,minmax(0,1fr)); } .player-actions { gap:var(--space-1); } .player-actions :global(.button) { padding:0 var(--space-2); font-size:var(--font-12); } .manage-panel { width:100%; padding:var(--space-4); } .action-dialog { padding:var(--space-4); } .console-command { flex-wrap:wrap; } .console-command input { min-width:60%; } .console-panel.expanded { inset:0; border-radius:0; } :global(dialog.player-action-modal) { position:fixed; inset:auto 0 0; width:100%; max-width:none; max-height:90vh; margin:0; border-radius:var(--radius-card) var(--radius-card) 0 0; } }
</style>
