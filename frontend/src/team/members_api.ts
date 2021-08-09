import { derived, get, Readable, Updater, writable } from "svelte/store";
import { RosterMember, TeamId } from "./team_api";
import type { Member } from "./team_api";
import { Encounters } from "./encounters_api";

interface memberEncountersStore extends Readable<Promise<RosterMember[]>> {
  set(rm: RosterMember): Promise<any>;
  remove(): Promise<any>;
}

interface memberStore extends Readable<Promise<Member>> {
  update(f: Updater<Member>);
  set(m: Member);
  remove();
  encounters: memberEncountersStore;
}

function createMember(id: number): memberStore {
  const subscribe = derived(Members, (p) => p.then((t) => t.get(id))).subscribe;

  const set = (m: Member) => {
    if (m.id != id) return;

    return fetch(`/team/${get(TeamId)}/member/${id}`, {
      method: "PUT",
      body: JSON.stringify(m),
      headers: { "Content-Type": "application/json" },
    }).then(() => Members.reload());
  };

  const update = (f: Updater<Member>) =>
    get({ subscribe }).then((oldM) => {
      const newM = f(oldM);
      if (newM != undefined) {
        return set(newM);
      } else {
        return Promise.resolve();
      }
    });

  const remove = () =>
    fetch(`/team/${get(TeamId)}/member/${id}`, {
      method: "DELETE",
    }).then(() => Members.reload());

  const encounterW = writable<Promise<RosterMember[]>>(Promise.resolve([]));

  const reloadEncounters = () => {
    const p = fetch(`/team/${get(TeamId)}/member/${id}/encounters`).then(
      (r) => r.json() as Promise<RosterMember[]>
    );
    encounterW.set(p);
    get(Encounters).then((es) =>
      es.forEach((e) => Encounters.encounter(e.id).roster.reload())
    );
    return p;
  };

  subscribe(reloadEncounters);

  const setEncounters = (rm: RosterMember): Promise<any> =>
    fetch(`/team/${get(TeamId)}/member/${id}/encounters`, {
      method: "POST",
      body: JSON.stringify(rm),
      headers: { "Content-Type": "application/json" },
    }).then(reloadEncounters);

  const removeEncounters = (): Promise<any> =>
    fetch(`/team/${get(TeamId)}/member/${id}/encounters`, {
      method: "DELETE",
    }).then(reloadEncounters);

  return {
    subscribe,
    set,
    update,
    remove,
    encounters: {
      subscribe: encounterW.subscribe,
      set: setEncounters,
      remove: removeEncounters,
    },
  };
}

function createMembers() {
  const w = writable<Promise<Map<number, Member>>>(Promise.resolve(new Map()));
  let loadedTeamId: string | null = null;

  const reload = () => {
    const teamId = get(TeamId);
    if (teamId) {
      w.set(
        fetch(`/team/${teamId}/members`)
          .then((r) => r.json() as Promise<Member[]>)
          .then((r) => {
            loadedTeamId = teamId;
            return new Map(r.map((m) => [m.id, m]));
          })
      );
    }
  };

  TeamId.subscribe((teamId) => {
    if (teamId != loadedTeamId) {
      reload();
    }
  });

  const addMember = (m: Member) => {
    fetch(`/team/${get(TeamId)}/member`, {
      method: "POST",
      body: JSON.stringify(m),
      headers: { "Content-Type": "application/json" },
    })
      .then(reload)
      .catch((e) => console.error(`error adding member`, { m, e }));
  };

  const memberStores = new Map<number, memberStore>();

  const member = (id: number): memberStore => {
    if (memberStores.has(id)) {
      return memberStores.get(id);
    } else {
      const store = createMember(id);
      memberStores.set(id, store);
      return store;
    }
  };

  return {
    subscribe: w.subscribe,
    reload,
    member,
    addMember,
  };
}

export const Members = createMembers();

export function SortMembers(a: Member, b: Member) {
  if (a.classId != b.classId) {
    return a.classId - b.classId;
  } else if (a.config.primarySpec != b.config.primarySpec) {
    return a.config.primarySpec - b.config.primarySpec;
  } else if (a.name < b.name) return -1;
  else if (a.name > b.name) return 1;
  else return a.id - b.id;
}
