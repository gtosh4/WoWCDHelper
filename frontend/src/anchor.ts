import { derived, writable } from "svelte/store";

export const CreateAnchor = () => {
  const c = writable(new Set<string>());

  const addClass = (name: string) => {
    c.update((names: Set<string>) => {
      names.add(name);
      return names;
    });
  };

  const removeClass = (name: string) => {
    c.update((names) => {
      names.delete(name);
      return names;
    });
  };

  return {
    subscribe: derived(c, (names) => [...names].join(" ")).subscribe,
    addClass,
    removeClass,
  };
};
