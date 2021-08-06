import { Readable, writable } from "svelte/store";
import { Spec } from "../wow/api";
import { CurrentTeam, SortMembers } from "./api";

export function RoleMembers(roleName: string): Readable<Promise<number[]>> {
  const w = writable<Promise<number[]>>(Promise.resolve([]));

  CurrentTeam.subscribe((p) =>
    w.set(
      p.then((t) =>
        Promise.all(
          [...t.values()].map((m) => {
            const specId = m?.config?.primarySpec;
            return Spec(specId).then((s) => ({ m, s }));
          })
        ).then((entries) => {
          console.log("updating", { roleName, entries });
          return entries
            .filter((e) => e.s.role.name == roleName)
            .map((e) => e.m)
            .sort(SortMembers)
            .map((m) => m.id);
        })
      )
    )
  );

  return { subscribe: w.subscribe };
}
