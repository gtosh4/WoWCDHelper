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

function hashSplit(url: URL): string[] {
  return url.hash.replace(/^#\//, "").split("/");
}

export const HashPath = {
  subscribe: derived(href, ($href): string[] => hashSplit(new URL($href)))
    .subscribe,

  update: (f: Updater<string[]>) => {
    const url = new URL(window.location.href);
    const hashPath = hashSplit(url);
    const nextHashPath = f(hashPath);
    url.hash = `#/${nextHashPath.join("/")}`;
    history.pushState(null, "", url.toString());
  },

  set: (hashPath) => {
    const url = new URL(window.location.href);
    url.hash = `#/${hashPath.join("/")}`;
    history.pushState(null, "", url.toString());
  },
};

export function HashPathPart(idx: number) {
  return {
    subscribe: derived(HashPath, (path) =>
      path && path.length > idx ? path[idx] : undefined
    ).subscribe,

    update: (f: Updater<string>) => {
      return HashPath.update((path) => {
        if (path.length > idx) {
          path[idx] = f(path[idx]);
        } else if (path.length == idx) {
          path.push(f(undefined));
        }
        return path;
      });
    },

    set: (value: string) => {
      HashPath.update((path) => {
        if (path.length > idx) {
          path[idx] = value;
        }
        if (path.length == idx) {
          path.push(value);
        }
        return path;
      });
    },
  };
}

export default {
  subscribe: derived(href, ($href): URL => new URL($href)).subscribe,
  update: (f: Updater<URL>) => href.update((h) => f(new URL(h)).toString()),
  set: (urlHref) => history.pushState(null, "", urlHref.toString()),
};
