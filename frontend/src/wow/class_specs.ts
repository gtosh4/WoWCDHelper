import { Classes, Spec, Specialization } from "../wow/api";

export function ClassSpecs(classId: number): Promise<Specialization[]> {
  return Classes.then((classes) => {
    const c = classes.get(classId);
    if (!c) return [];

    return Promise.all(
      c.specializations.map((specId) => {
        return Spec(specId.id);
      })
    );
  });
}
