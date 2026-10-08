import { writable } from "svelte/store";

const storedTheme = localStorage.getItem("shmolph-theme") || "dark";
export const theme = writable(storedTheme);

export function applyTheme(value) {
  document.documentElement.dataset.theme = value;
  localStorage.setItem("shmolph-theme", value);
  theme.set(value);
}

export function initializeTheme() {
  applyTheme(localStorage.getItem("shmolph-theme") || "dark");
}

export function cycleTheme(value) {
  const choices = ["system", "light", "dark"];
  applyTheme(choices[(choices.indexOf(value) + 1) % choices.length]);
}
