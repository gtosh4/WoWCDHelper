import { derived, get, Readable, Updater, writable } from "svelte/store";
import { Member, TeamId } from "./team_api";
import type { Encounter, RosterMember } from "./team_api";

interface rosterMemberStore
  extends Readable<Promise<RosterMember | undefined>> {
  set(m: RosterMember | Member): Promise<any>;
  remove(): Promise<any>;
}

interface rosterStore extends Readable<Promise<RosterMember[]>> {
  set(rs: RosterMember[]): Promise<RosterMember[]>;
  reload(): Promise<any>;
  member(id: number): rosterMemberStore;
}

interface encounterStore extends Readable<Promise<Encounter>> {
  update(f: Updater<Encounter>): Promise<Encounter>;
  remove(): Promise<any>;
  roster: rosterStore;
}

function createRoster(encId: number): rosterStore {
  const w = writable<Promise<RosterMember[]>>(Promise.resolve([]));

  const reload = () => {
    const p = fetch(`/team/${get(TeamId)}/encounter/${encId}/roster`).then(
      (r) => r.json() as Promise<RosterMember[]>
    );
    w.set(p);
    return p;
  };

  Encounters.subscribe(reload);

  const set = (rs: RosterMember[]) => {
    return fetch(`/team/${get(TeamId)}/encounter/${encId}/roster`, {
      method: "PUT",
      body: JSON.stringify(rs),
      headers: { "Content-Type": "application/json" },
    }).then(() => reload());
  };

  const memberStores = new Map<number, rosterMemberStore>();

  const member = (id: number): rosterMemberStore => {
    if (memberStores.has(id)) {
      return memberStores.get(id);
    } else {
      const subscribe = derived({ subscribe: w.subscribe }, (p) =>
        p.then((roster) => roster.find((m) => m.member_id == id))
      ).subscribe;

      const set = (m: RosterMember | Member): Promise<any> => {
        let rm: RosterMember;
        if ("member_id" in m) {
          rm = m;
        } else {
          rm = {
            encounter_id: encId,
            member_id: m.id,
            spec: m.config.primarySpec,
          };
        }
        const p = fetch(
          `/team/${get(TeamId)}/encounter/${encId}/roster/${rm.member_id}`,
          {
            method: "PUT",
            body: JSON.stringify(rm),
            headers: { "Content-Type": "application/json" },
          }
        );
        get(Encounters).then((es) =>
          es.forEach((e) => Encounters.encounter(e.id).roster.reload())
        );
        return p.then(() => reload());
      };

      const remove = () => {
        const p = fetch(
          `/team/${get(TeamId)}/encounter/${encId}/roster/${id}`,
          {
            method: "DELETE",
          }
        );
        get(Encounters).then((es) =>
          es.forEach((e) => Encounters.encounter(e.id).roster.reload())
        );

        return p.then(() => reload());
      };

      const store = {
        subscribe,
        set,
        remove,
      };
      memberStores.set(id, store);
      return store;
    }
  };

  return {
    subscribe: w.subscribe,
    set,
    reload,
    member,
  };
}

function createEncounters() {
  const w = writable<Promise<Encounter[]>>(Promise.resolve([]));
  let loadedTeamId: string | null = null;

  const reload = () => {
    const teamId = get(TeamId);
    if (teamId) {
      const p = fetch(`/team/${teamId}/encounters`)
        .then((r) => r.json() as Promise<Encounter[]>)
        .then((r) => {
          loadedTeamId = teamId;
          return r;
        });
      w.set(p);
      return p;
    }
    return Promise.resolve<Encounter[]>([]);
  };

  TeamId.subscribe((teamId) => {
    if (teamId != loadedTeamId) {
      reload();
    }
  });

  const addEncounter = (e: Encounter) => {
    fetch(`/team/${get(TeamId)}/encounter`, {
      method: "POST",
      body: JSON.stringify(e),
      headers: { "Content-Type": "application/json" },
    }).then(reload);
  };

  const rosterStores = new Map<number, rosterStore>();

  const roster = (id: number) => {
    if (rosterStores.has(id)) {
      return rosterStores.get(id);
    } else {
      const store = createRoster(id);
      rosterStores.set(id, store);
      return store;
    }
  };

  const encounter = (id: number): encounterStore => {
    const subscribe = derived(w, (p) =>
      p.then((t) => t.find((e) => e.id == id))
    ).subscribe;

    const update = (f: Updater<Encounter>) =>
      get({ subscribe }).then((oldE) => {
        const newE = f(oldE);
        if (newE == undefined) {
          return Promise.resolve(oldE);
        }
        return fetch(`/team/${get(TeamId)}/encounter/${id}`, {
          method: "PUT",
          body: JSON.stringify(newE),
          headers: { "Content-Type": "application/json" },
        })
          .then(() => reload())
          .then((p) => p.find((e) => e.id == id));
      });

    const remove = () =>
      fetch(`/team/${get(TeamId)}/encounter/${id}`, {
        method: "DELETE",
      }).then(() => reload());

    return {
      subscribe,
      update,
      remove,
      roster: roster(id),
    };
  };

  return {
    subscribe: w.subscribe,
    reload,
    encounter,
    addEncounter,
  };
}

export const Encounters = createEncounters();
