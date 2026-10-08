import { writable } from "svelte/store";

export const session = writable({ status: "checking", username: null });
export const toasts = writable([]);
export const serverCache = writable(null);
let toastSequence = 0;

export function pushToast(message, type = "info") {
  const id = `${Date.now()}-${toastSequence++}`;
  toasts.update((items) => [...items, { id, message, type }]);
  setTimeout(() => toasts.update((items) => items.filter((item) => item.id !== id)), 3200);
}

export function clearSessionState() {
  session.set({ status: "anonymous", username: null });
  serverCache.set(null);
}
