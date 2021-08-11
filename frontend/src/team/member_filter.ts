import { Spec } from "../wow/api";
import type { Member, RosterMember } from "./team_api";

export function RoleMembers<T extends RosterMember | Member>(
  p: Promise<T[]>,
  roleType: string
): Promise<T[]> {
  return p
    .then((members) =>
      Promise.all(
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
    )
    .then((entries) =>
      entries.filter((e) => e.s.role.type == roleType).map((e) => e.m)
    );
}
