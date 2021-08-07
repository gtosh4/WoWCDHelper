import {
  derived,
  get,
  Readable,
  Updater,
  Writable,
  writable,
} from "svelte/store";
import { HashPathPart } from "../url";

export interface Member {
  id: number;
  team: string;
  name: string;
  classId: number;
  config: MemberConfig;
}

export interface MemberConfig {
  specs: number[];
  primarySpec: number;
}

function createTeam() {
  let teamId: string | null = null;
  const w = writable<Promise<Map<number, Member>>>(Promise.resolve(new Map()));

  const reload = () => {
    if (teamId) {
      w.set(
        fetch(`/team/${teamId}`)
          .then((r) => r.json() as Promise<Member[]>)
          .then((r) => new Map(r.map((m) => [m.id, m])))
      );
    }
  };

  const TeamPath = HashPathPart(0);
  TeamPath.subscribe((param) => {
    if (!param) {
      fetch("/team/new", { method: "POST" })
        .then((r) => {
          const createdId = r.headers
            .get("Location")
            .match(/\/team\/([^/]+)/)[1];
          teamId = createdId;
          TeamPath.set(createdId);
        })
        .catch((e) => console.error("error creating new team", e));
    } else if (param != teamId) {
      teamId = param;
      reload();
    }
  });

  const addMember = (m: Member) => {
    fetch(`/team/${teamId}/member`, {
      method: "POST",
      body: JSON.stringify(m),
      headers: { "Content-Type": "application/json" },
    })
      .then(reload)
      .catch((e) => console.error(`error adding member`, { m, e }));
  };

  return {
    subscribe: w.subscribe,
    reload,
    addMember,
    teamID() {
      return teamId;
    },
  };
}

export const CurrentTeam = createTeam();

export function SortMembers(a: Member, b: Member) {
  if (a.classId != b.classId) {
    return a.classId - b.classId;
  } else if (a.config.primarySpec != b.config.primarySpec) {
    return a.config.primarySpec - b.config.primarySpec;
  } else if (a.name < b.name) return -1;
  else if (a.name > b.name) return 1;
  else return a.id - b.id;
}

interface teamMember extends Readable<Promise<Member>> {
  update(f: Updater<Member>);
  set(m: Member);
  remove();
}

export function TeamMember(id: number): teamMember {
  const subscribe = derived(CurrentTeam, (p) =>
    p.then((t) => t.get(id))
  ).subscribe;

  const set = async (m: Member) => {
    if (m.id != id) return;

    return fetch(`/team/${CurrentTeam.teamID()}/${id}`, {
      method: "PUT",
      body: JSON.stringify(m),
      headers: { "Content-Type": "application/json" },
    }).then(() => CurrentTeam.reload());
  };

  const update = async (f: Updater<Member>) =>
    get(CurrentTeam).then((t) => {
      const oldM = t.get(id);
      const newM = f(oldM);
      if (newM.id != id) return;

      return fetch(`/team/${CurrentTeam.teamID()}/${id}`, {
        method: "PUT",
        body: JSON.stringify(newM),
        headers: { "Content-Type": "application/json" },
      }).then(() => CurrentTeam.reload());
    });

  const remove = async () =>
    fetch(`/team/${CurrentTeam.teamID()}/${id}`, { method: "DELETE" }).then(
      () => CurrentTeam.reload()
    );

  return {
    subscribe,
    set,
    update,
    remove,
  };
}
