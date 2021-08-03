import { writable } from "svelte/store";
import url from "../url";

export interface Member {
  id: number;
  name: string;
  className: string;
}

export interface Team {
  id: number;
}

function createTeam() {
  let teamId: string | null = null;
  const team = writable<Member[]>([]);

  let inCreate = false;
  url.subscribe((href) => {
    const param = href.searchParams.get("team");
    if (param && param != teamId) {
      teamId = param;
      if (!inCreate) {
        fetch(`/team/${teamId}`)
          .then((r) => r.json())
          .then((r) => {
            team.set(r);
          })
          .catch((e) => console.error(`error loading team ${teamId}`, e));
      }
    }
  });

  team.subscribe((ms: Member[]) => {
    let init: Promise<any> = Promise.resolve();
    if (!teamId) {
      init = fetch("/team/new", { method: "POST" })
        .then((r) => r.json())
        .then((r) => (teamId = r.id))
        .then(() => {
          inCreate = true;
          url.update((u) => {
            u.searchParams.set("team", teamId);
            return u;
          });
          inCreate = false;
        });
    }
    init
      .then(() => {
        fetch(`/team/${teamId}`, { method: "PUT", body: JSON.stringify(ms) })
          .then((r) => r.json())
          .then((r: Member[]) => {
            // mutate the member list so it doesn't trigger another fetch
            while (ms.length > 0) ms.pop();
            r.forEach((m) => ms.push(m));
          });
      })
      .catch((e) => console.error(`error saving team ${teamId}`, e));
  });

  return {
    subscribe: team.subscribe,
    set: team.set,
    update: team.update,

    updateMember: (m: Member) =>
      team.update((ms: Member[]) => {
        ms.filter((e) => e.id == m.id).map((e) => Object.assign(e, m));
        return ms;
      }),

    deleteMember: (id: number) =>
      team.update((ms: Member[]) => {
        return ms.filter((e) => e.id != id);
      }),

    addMember: (m: Member) =>
      team.update((ms: Member[]) => {
        ms.push(m);
        return ms;
      }),
  };
}

export const CurrentTeam = createTeam();
