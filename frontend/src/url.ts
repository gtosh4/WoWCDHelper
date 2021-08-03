import { derived, Updater, writable } from "svelte/store";

const isBrowser = typeof window !== "undefined";

const href = writable(isBrowser ? window.location.href : "https://example.com");

const URL = isBrowser ? window.URL : require("url").URL;

if (isBrowser) {
  const originalPushState = history.pushState;
  const originalReplaceState = history.replaceState;

  const updateHref = () => href.set(window.location.href);

  history.pushState = function () {
    originalPushState.apply(this, arguments);
    updateHref();
  };

  history.replaceState = function () {
    originalReplaceState.apply(this, arguments);
    updateHref();
  };

  window.addEventListener("popstate", updateHref);
  window.addEventListener("hashchange", updateHref);
}

export default {
  subscribe: derived(href, ($href): URL => new URL($href)).subscribe,
  update: (f: Updater<URL>) => href.update((h) => f(new URL(h)).toString()),
  set: (urlHref) => href.set(urlHref),
};
