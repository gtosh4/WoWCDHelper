import { Readable, writable } from "svelte/store";
import { Spec } from "../wow/api";
import { CurrentTeam, SortMemberIds } from "./api";

export function RoleMembers(roleName: string): Readable<number[]> {
  const w = writable<number[]>([]);

  const subs = [];

  CurrentTeam.subscribe((t) => {
    while (subs.length > 0) {
      const unsub = subs.pop();
      unsub();
    }
    t.forEach((m) => {
      if (!m.config.primarySpec) return;

      const spec = Spec(m.config.primarySpec);

      const sub = spec.subscribe((s) => {
        if (!s) return;

        w.update((members) => {
          const set = new Set(members);
          if (s.role.name == roleName) {
            set.add(m.id);
          } else {
            set.delete(m.id);
          }
          return [...set].sort(SortMemberIds);
        });
      });

      subs.push(sub);
    });
  });

  return { subscribe: w.subscribe };
}
