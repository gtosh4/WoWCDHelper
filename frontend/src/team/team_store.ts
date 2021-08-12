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
  remove(): Promise<void>;

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
    const unsub = this.direct.subscribe(run, invalidate);
    if (get(this.direct) == undefined) {
      this.reload();
    }
    return unsub;
  }

  public reload(): Promise<T> {
    console.log("reloading", { url: this.url });
    const p = fetch(this.url).then((r) => r.json() as Promise<T>);
    this.direct.set(p);
    return p;
  }

  public set(v: T): void {
    fetch(this.url, {
      method: "PUT",
      body: JSON.stringify(v),
      headers: { "Content-Type": "application/json" },
    }).then((r) => {
      const contentType = r.headers.get("content-type");
      if (contentType && contentType.indexOf("application/json") >= 0) {
        this.direct.set(r.json());
      } else {
        this.reload();
      }
    });
  }

  public remove(): Promise<void> {
    return fetch(this.url, { method: "DELETE" }).then(() => {});
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
        this.direct.set(Promise.resolve(newV));
      }
    });
  }
}

class cell {
  private store: store;

  public memberId: number;
  public memberInfo: AsyncRWStore<Member>;
  public encounterId: number;
  public encounterInfo: AsyncRWStore<Encounter>;

  public rosterMember: Writable<Promise<RosterMember>>;

  constructor(s: store, r: row, c: column) {
    this.store = s;
    this.memberId = r.memberId;
    this.memberInfo = r.memberInfo;
    this.encounterId = c.encounterId;
    this.encounterInfo = c.encounterInfo;

    this.rosterMember = writable(Promise.resolve<RosterMember>(undefined));

    let firstRun = true;
    this.rosterMember.subscribe((p) => {
      if (!firstRun) {
        p.then((rm) => {
          console.trace("cell propagation", {
            memberId: this.memberId,
            encounterId: this.encounterId,
            rm,
          });
          r.updateFromCell(this.encounterId, rm);
          c.updateFromCell(this.memberId, rm);
        });
      } else {
        firstRun = false;
      }
    });
  }

  update(f: Updater<RosterMember>): Promise<RosterMember> {
    return get(this.rosterMember).then((oldV) => {
      const newV = f(oldV);
      if (!equal(oldV, newV)) {
        return fetch(
          `/team/${this.store.teamId}/encounter/${this.encounterId}/roster/${this.memberId}`,
          {
            method: "PUT",
            body: JSON.stringify(newV),
            headers: { "Content-Type": "application/json" },
          }
        ).then((r) => r.json() as Promise<RosterMember>);
      }
      return Promise.resolve(oldV);
    });
  }

  remove(): Promise<void> {
    return fetch(
      `/team/${this.store.teamId}/encounter/${this.encounterId}/roster/${this.memberId}`,
      { method: "DELETE" }
    ).then(() => this.rosterMember.set(Promise.resolve(undefined)));
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
      console.trace("updating row", {
        encounterId: this.memberId,
        inUpdate: this.inUpdate,
      });
      if (this.inUpdate) return;
      if (!p) return;

      p.then((rms) => {
        rms.forEach((rm) => {
          console.log("row cell update", { memberId: this.memberId, rm });
          s.cell(memberId, rm.encounter_id).rosterMember.set(
            Promise.resolve(rm)
          );
        });
      });
    });
  }

  remove(): Promise<void> {
    return this.memberInfo.remove();
  }

  updateFromCell(encounterId: number, rm?: RosterMember) {
    this.inUpdate = true;
    this.memberEncounters.directUpdate((rms) => {
      if (!rms) {
        return rm ? [rm] : rms;
      }
      const idx = rms.findIndex(
        (m) => m.encounter_id == encounterId && m.member_id == this.memberId
      );
      if (idx >= 0) {
        if (rm) {
          rms[idx] = rm;
        } else {
          rms.splice(idx, 1);
        }
      } else {
        if (rm) {
          rms.push(rm);
        }
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
      console.log("updating col", {
        encounterId: this.encounterId,
        inUpdate: this.inUpdate,
      });
      if (this.inUpdate) return;
      if (!p) return;

      p.then((rms) => {
        rms.forEach((rm) => {
          console.log("col cell update", { encounterId: this.encounterId, rm });
          s.cell(rm.member_id, encounterId).rosterMember.set(
            Promise.resolve(rm)
          );
        });
      });
    });
  }

  remove() {
    this.encounterInfo.remove();
  }

  updateFromCell(memberId: number, rm: RosterMember) {
    this.inUpdate = true;
    this.encounterMembers.directUpdate((rms) => {
      if (!rms) {
        return rm ? [rm] : rms;
      }
      const idx = rms.findIndex(
        (m) => m.encounter_id == this.encounterId && m.member_id == memberId
      );
      if (idx >= 0) {
        if (rm) {
          rms[idx] = rm;
        } else {
          rms.splice(idx, 1);
        }
      } else {
        if (rm) {
          rms.push(rm);
        }
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
    console.trace("new cell", { memberId, encounterId });
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
