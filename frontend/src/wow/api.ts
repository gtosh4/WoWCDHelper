import { derived, Readable } from "svelte/store";
import { writable } from "svelte/store";

export interface PlayerClass {
  id: number;
  name: string;
  power_type: PowerType;
  specializations: { id: number }[];
}

export interface PowerType {
  name: string;
  id: number;
}

export interface Specialization {
  name: string;
  id: number;
  playable_class: { name: string; id: number };
  role: { type: string; name: string };
}

export function SortClassByName(a, b) {
  if (a.name < b.name) return -1;
  else if (a.name > b.name) return 1;
  else return 0;
}

function createClasses() {
  const w = writable(new Map<number, PlayerClass>());
  const r = { subscribe: w.subscribe };

  fetch("/wow/classes")
    .then((r) => r.json() as Promise<PlayerClass[]>)
    .then((r) =>
      w.update((m) => {
        m.clear();
        r.forEach((cls) => m.set(cls.id, cls));
        return m;
      })
    );

  return r;
}

export const Classes = createClasses();

export const ClassList = derived(Classes, (p) =>
  [...p.values()].sort(SortClassByName)
);

const specs = new Map<number, Readable<Specialization | undefined>>();

export function Spec(id: number): Readable<Specialization | undefined> {
  if (specs.has(id)) {
    return specs.get(id);
  }
  const w = writable<Specialization | undefined>(undefined);
  const r = { subscribe: w.subscribe };
  specs.set(id, r);

  fetch(`/wow/spec/${id}`)
    .then((r) => r.json())
    .then((r) => w.set(r))
    .catch((e) => console.error(`error fetching spec ${id}`, e));

  return r;
}
