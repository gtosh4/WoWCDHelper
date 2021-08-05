import {
  derived,
  get,
  Readable,
  Updater,
  Writable,
  writable,
} from "svelte/store";
import { HashPath, HashPathPart } from "../url";

export interface Member {
  id: number;
  name: string;
  classId: number;
  config: MemberConfig;
}

export interface MemberConfig {
  specs: number[];
  primary_spec: number;
}

export interface Team {
  id: number;
}

function createTeam() {
  let teamId: string | null = null;
  const team = writable<Map<number, Member>>(new Map());

  const update = (r: Member[]) =>
    team.update((t) => {
      t.clear();
      r.forEach((m) => t.set(m.id, m));
      return t;
    });

  const reload = () => {
    if (teamId) {
      fetch(`/team/${teamId}`)
        .then((r) => r.json())
        .then(update)
        .catch((e) => console.error(`error loading team ${teamId}`, e));
    }
  };

  const TeamPath = HashPathPart(0);
  TeamPath.subscribe((param) => {
    if (!param) {
      fetch("/team/new", { method: "POST" })
        .then((r) => {
          const createdId = r.headers.get("Location").replace(/\/team\//, "");
          teamId = createdId;
          TeamPath.set(createdId);
        })
        .catch((e) => console.error("error creating new team", e));
    } else if (param != teamId) {
      teamId = param;
      reload();
    }
  });

  return {
    subscribe: team.subscribe,
    set: team.set,
    update: team.update,
    reload,
    teamID() {
      return teamId;
    },
  };
}

export const CurrentTeam = createTeam();

export function SortMemberIds(a: number, b: number) {
  const team = get(CurrentTeam);
  const mA = team.get(a);
  const mB = team.get(b);
  if (mA.classId != mB.classId) {
    return mA.classId - mB.classId;
  } else if (mA.config.primary_spec != mB.config.primary_spec) {
    return mA.config.primary_spec - mB.config.primary_spec;
  } else if (mA.name < mB.name) return -1;
  else if (mA.name > mB.name) return 1;
  else return mA.id - mB.id;
}

export function TeamMember(id: number) {
  const subscribe = derived(CurrentTeam, (t) => t.get(id)).subscribe;

  const set = (m: Member) => {
    if (m.id != id) return;

    fetch(`/team/${CurrentTeam.teamID()}/${id}`, {
      method: "PUT",
      body: JSON.stringify(m),
    })
      .then(() => CurrentTeam.update((t) => t.set(id, m)))
      .catch((e) =>
        console.error("error updating member", { member: m, err: e })
      );
  };
  const update = (f: Updater<Member>) => {
    CurrentTeam.update((t) => {
      const oldM = t.get(id);
      const newM = f(oldM);
      if (newM.id != id) return;

      fetch(`/team/${CurrentTeam.teamID()}/${id}`, {
        method: "PUT",
        body: JSON.stringify(newM),
      })
        .then(() => t.set(id, newM))
        .catch((e) =>
          console.error("error updating member", { oldM, newM, err: e })
        );

      return t;
    });
  };

  const remove = () => {
    fetch(`/team/${CurrentTeam.teamID()}/${id}`, { method: "DELETE" }).then(
      () =>
        CurrentTeam.update((t) => {
          t.delete(id);
          return t;
        })
    );
  };

  return {
    subscribe,
    set,
    update,
    remove,
  };
}
