import { writable, derived, get } from "svelte/store";
import type {
  Readable,
  Subscriber,
  Updater,
  Unsubscriber,
  Writable,
} from "svelte/store";
import { Encounter, Member, RosterMember, Team, TeamId } from "./team_api";
import equal from "fast-deep-equal/es6";

interface AsyncRWStore<T> extends Readable<Promise<T>> {
  reload(): Promise<T>;
  set(v: T): void;
  update(f: Updater<T>): void;
  remove(): void;

  directUpdate(f: Updater<T>): void;
}

class apiResource<T> implements AsyncRWStore<T> {
  public url: string;
  private direct: Writable<Promise<T>>;

  constructor(url: string) {
    this.url = url;
    this.direct = writable(undefined);
  }

  public subscribe(
    run: Subscriber<Promise<T>>,
    invalidate?: (value?: Promise<T>) => void
  ): Unsubscriber {
    if (get(this.direct) == undefined) {
      this.reload();
    }
    return this.direct.subscribe(run, invalidate);
  }

  public reload(): Promise<T> {
    const p = fetch(this.url).then((r) => r.json() as Promise<T>);
    this.direct.set(p);
    return p;
  }

  public set(v: T): void {
    fetch(this.url, {
      method: "PUT",
      body: JSON.stringify(v),
      headers: { "Content-Type": "application/json" },
    }).then(this.reload);
  }

  public remove(): void {
    fetch(this.url, { method: "DELETE" });
  }

  public update(f: Updater<T>): void {
    get(this.direct).then((oldV) => {
      const newV = f(oldV);
      if (!equal(oldV, newV)) {
        return this.set(newV);
      }
    });
  }

  public directUpdate(f: Updater<T>): void {
    get(this.direct).then((oldV) => {
      const newV = f(oldV);
      if (!equal(oldV, newV)) {
        return this.direct.set(Promise.resolve(newV));
      }
    });
  }
}

class cell {
  public memberId: number;
  public memberInfo: AsyncRWStore<Member>;
  public encounterId: number;
  public encounterInfo: AsyncRWStore<Encounter>;
  public rosterMember: AsyncRWStore<RosterMember>;

  constructor(s: store, r: row, c: column) {
    this.memberId = r.memberId;
    this.memberInfo = r.memberInfo;
    this.encounterId = c.encounterId;
    this.encounterInfo = c.encounterInfo;

    this.rosterMember = new apiResource(
      `/team/${s.teamId}/encounters/${this.encounterId}/roster/${this.memberId}`
    );

    this.rosterMember.subscribe((p) =>
      p.then((rm) => {
        r.updateFromCell(rm);
        c.updateFromCell(rm);
      })
    );
  }
}

class row {
  public memberId: number;
  public memberInfo: AsyncRWStore<Member>;
  public memberEncounters: AsyncRWStore<RosterMember[]>;

  private inUpdate: boolean; // prevent cyclic updates cell -> row -> cell etc

  constructor(s: store, memberId: number) {
    this.memberId = memberId;

    this.memberInfo = new apiResource(`/team/${s.teamId}/member/${memberId}`);
    this.memberEncounters = new apiResource(
      `/team/${s.teamId}/member/${memberId}/encounters`
    );

    this.inUpdate = false;
    this.memberEncounters.subscribe((p) => {
      if (this.inUpdate) return;

      p.then((rms) => {
        rms.forEach((rm) => {
          s.cell(memberId, rm.encounter_id).rosterMember.directUpdate(() => rm);
        });
      });
    });
  }

  updateFromCell(rm: RosterMember) {
    this.inUpdate = true;
    this.memberEncounters.directUpdate((rms) => {
      const idx = rms.findIndex(
        (m) => m.encounter_id == rm.encounter_id && m.member_id == rm.member_id
      );
      if (idx >= 0) {
        rms[idx] = rm;
      } else {
        rms.push(rm);
      }
      return rms;
    });
    this.inUpdate = false;
  }
}

class column {
  public encounterId: number;
  public encounterInfo: AsyncRWStore<Encounter>;
  public encounterMembers: AsyncRWStore<RosterMember[]>;

  private inUpdate: boolean; // prevent cyclic updates cell -> column -> cell etc

  constructor(s: store, encounterId: number) {
    this.encounterId = encounterId;

    this.encounterInfo = new apiResource(`/team/${s.teamId}/encounter`);
    this.encounterMembers = new apiResource(
      `/team/${s.teamId}/encounter/${encounterId}/roster`
    );

    this.inUpdate = false;
    this.encounterMembers.subscribe((p) => {
      if (this.inUpdate) return;

      p.then((rms) => {
        rms.forEach((rm) => {
          s.cell(rm.member_id, encounterId).rosterMember.directUpdate(() => rm);
        });
      });
    });
  }

  updateFromCell(rm: RosterMember) {
    this.inUpdate = true;

    this.encounterMembers.directUpdate((rms) => {
      const idx = rms.findIndex(
        (m) => m.encounter_id == rm.encounter_id && m.member_id == rm.member_id
      );
      if (idx >= 0) {
        rms[idx] = rm;
      } else {
        rms.push(rm);
      }
      return rms;
    });
    this.inUpdate = false;
  }
}

class store {
  public teamId: string;

  private cells: Map<string, cell>;
  private rows: Map<number, row>;
  private columns: Map<number, column>;

  public Team: AsyncRWStore<Team>;
  public Members: AsyncRWStore<Member[]>;
  public Encounters: AsyncRWStore<Encounter[]>;

  constructor(teamId: string) {
    this.teamId = teamId;

    this.cells = new Map();
    this.rows = new Map();
    this.columns = new Map();

    this.Team = new apiResource(`/team/${teamId}`);
    this.Members = new apiResource(`/team/${teamId}/members`);
    this.Encounters = new apiResource(`/team/${teamId}/encounters`);
  }

  public row(memberId: number): row {
    if (this.rows.has(memberId)) {
      return this.rows.get(memberId);
    }
    const r = new row(this, memberId);
    this.rows.set(memberId, r);
    return r;
  }

  public column(encounterId: number): column {
    if (this.columns.has(encounterId)) {
      return this.columns.get(encounterId);
    }
    const c = new column(this, encounterId);
    this.columns.set(encounterId, c);
    return c;
  }

  cellKey(memberId: number, encounterId: number): string {
    return `${memberId}/${encounterId}`;
  }

  public cell(memberId: number, encounterId: number): cell {
    const key = this.cellKey(memberId, encounterId);
    if (this.cells.has(key)) {
      return this.cells.get(key);
    }
    const c = new cell(this, this.row(memberId), this.column(encounterId));
    this.cells.set(key, c);
    return c;
  }

  public newMember(m: Member): Promise<row> {
    return fetch(`/team/${this.teamId}/members`, {
      method: "POST",
      body: JSON.stringify(m),
      headers: { "Content-Type": "application/json" },
    }).then((r) => {
      const createdId = r.headers
        .get("Locations")
        .match(/\/team\/[^/]+\/member\/(\d+)/)[1];
      return this.row(+createdId);
    });
  }

  public newEncounter(e: Encounter): Promise<column> {
    return fetch(`/team/${this.teamId}/encounters`, {
      method: "POST",
      body: JSON.stringify(e),
      headers: { "Content-Type": "application/json" },
    }).then((r) => {
      const createdId = r.headers
        .get("Locations")
        .match(/\/team\/[^/]+\/encounter\/(\d+)/)[1];
      return this.column(+createdId);
    });
  }
}

export const TeamStore = derived(TeamId, (teamId) => new store(teamId));
