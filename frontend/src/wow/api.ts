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

export const Classes = fetch("/wow/classes")
  .then((r) => r.json() as Promise<PlayerClass[]>)
  .then((r) => new Map(r.map((cls) => [cls.id, cls])));

export const ClassList = Classes.then((cs) =>
  [...cs.values()].sort(SortClassByName)
);

const specs = new Map<number, Promise<Specialization>>();

export function Spec(id?: number): Promise<Specialization> {
  if (!id) return Promise.reject("'id' undefined");
  if (specs.has(id)) {
    return specs.get(id);
  }
  const p = fetch(`/wow/spec/${id}`).then(
    (r) => r.json() as Promise<Specialization>
  );

  specs.set(id, p);

  return p;
}
