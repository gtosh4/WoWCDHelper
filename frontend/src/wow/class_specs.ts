import { Readable, writable } from "svelte/store";
import { Classes, Spec, Specialization } from "../wow/api";

export function ClassSpecs(classId: number): Readable<Specialization[]> {
  const w = writable<Specialization[]>([]);

  const subs = [];

  Classes.subscribe((classes) => {
    if (!classes) return;
    const c = classes.get(classId);
    if (!c) return;

    while (subs.length > 0) {
      const unsub = subs.pop();
      unsub();
    }

    c.specializations.forEach((specId) => {
      const spec = Spec(specId.id);

      const sub = spec.subscribe((s) => {
        if (!s) return;

        w.update((specs) => {
          const set = new Set(specs);
          set.add(s);
          return [...set].sort((a, b) => a.id - b.id);
        });
      });

      subs.push(sub);
    });
  });

  return { subscribe: w.subscribe };
}
