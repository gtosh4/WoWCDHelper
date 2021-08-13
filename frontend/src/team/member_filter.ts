import { Spec } from "../wow/api";
import type { Member, RosterMember } from "./team_api";

export function RoleMembers<T extends RosterMember | Member>(
  members: T[],
  roleType: string
): Promise<T[]> {
  return Promise.all(
    members.map((m) => {
      let specId: number;
      if ("config" in m) {
        specId = m.config.primarySpec;
      } else if ("spec" in m) {
        specId = m.spec;
      }
      return Spec(specId).then((s) => ({ m, s }));
    })
  )
    .then((entries) =>
      entries.filter((e) => e.s.role.type == roleType).map((e) => e.m)
    )
    .then((ms) => ms || []);
}

export function SortMembers(a: Member, b: Member) {
  if (a.classId != b.classId) {
    return a.classId - b.classId;
  } else if (a.config.primarySpec != b.config.primarySpec) {
    return a.config.primarySpec - b.config.primarySpec;
  } else if (a.name < b.name) return -1;
  else if (a.name > b.name) return 1;
  else return a.id - b.id;
}
