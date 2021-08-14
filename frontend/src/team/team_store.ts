import { writable, derived, get, readable } from "svelte/store";
import type {
  Readable,
  Subscriber,
  Updater,
  Unsubscriber,
  Writable,
} from "svelte/store";
import { Encounter, Member, RosterMember, Team, TeamId } from "./team_api";
import equal from "fast-deep-equal/es6";
import { Classes, Spec } from "../wow/api";

class apiResource<T> {
  url: string;
  listener: (v: T) => void;

  constructor(url: string, listener?: (v: T) => void) {
    this.url = url;
    if (listener) {
      this.listener = listener;
    } else {
      this.listener = () => {};
    }
  }

  get(): Promise<T> {
    const p = fetch(this.url).then((r) => r.json() as Promise<T>);
    p.then((v) => this.listener(v));
    return p;
  }

  put(v: T): Promise<T> {
    return fetch(this.url, {
      method: "PUT",
      body: JSON.stringify(v),
      headers: { "Content-Type": "application/json" },
    }).then((r) => {
      const contentType = r.headers.get("content-type");
      if (contentType && contentType.indexOf("application/json") >= 0) {
        const p = r.json() as Promise<T>;
        p.then((v) => this.listener(v));
        return p;
      } else {
        return this.get();
      }
    });
  }

  remove(): Promise<void> {
    return fetch(this.url, { method: "DELETE" }).then(() =>
      this.listener(undefined)
    );
  }
}

export enum LoadingState {
  Uninitialized,
  Loading,
  Loaded,
}

class resourceWritable<T> extends apiResource<T> implements Writable<T> {
  public state: LoadingState;
  private _w: Writable<T>;

  constructor(url: string, listener?: (v: T) => void, w?: Writable<T>) {
    super(url, listener);
    if (w) {
      this._w = w;
    } else {
      this._w = writable(undefined, () => {
        this.reload();
      });
    }
    this.state = LoadingState.Uninitialized;
  }

  get(): Promise<T> {
    if (this.state == LoadingState.Uninitialized) {
      this.state = LoadingState.Loading;
    }
    return super.get().then((v) => {
      this._w.set(v);
      this.state = LoadingState.Loaded;
      return v;
    });
  }

  put(v: T): Promise<T> {
    this._w.set(v);
    return super.put(v).then((v2) => {
      if (!equal(get(this._w), v2)) {
        this._w.set(v2);
      }
      return v2;
    });
  }

  remove(): Promise<void> {
    this._w.set(undefined);
    return super.remove();
  }

  reload(): Promise<void> {
    return this.get().then(() => {});
  }

  subscribe(
    run: Subscriber<T>,
    invalidate?: (value?: T) => void
  ): Unsubscriber {
    return this._w.subscribe(run, invalidate);
  }

  set(value: T): void {
    this.put(value);
  }

  update(updater: Updater<T>): void {
    const oldV = get(this._w);
    const newV = updater(oldV);
    if (!equal(oldV, newV)) {
      this.set(newV);
    }
  }
}

class cell {
  memberId: number;
  encounterId: number;

  rosterMember: resourceWritable<RosterMember>;
  _rosterMember: Writable<RosterMember>;

  constructor(s: store, r: row, c: column) {
    this.memberId = r.memberId;
    this.encounterId = c.encounterId;

    this._rosterMember = writable(undefined);
    this.rosterMember = new resourceWritable(
      `/team/${s.teamId}/encounter/${this.encounterId}/roster/${this.memberId}`,
      undefined,
      this._rosterMember
    );
  }
}

function rosterMembersDerived<T>(
  r: Readable<T>,
  f: (t: T) => cell[]
): Readable<RosterMember[]> {
  return readable(undefined, (set) => {
    const unsubs = [];

    const unsubAll = () => {
      while (unsubs.length > 0) {
        const uns = unsubs.pop();
        uns();
      }
    };

    r.subscribe((es) => {
      unsubAll();

      if (es == undefined) {
        set(undefined);
        return;
      }

      // on first subscribe, only send once
      // instead of every time we subscribe to a cell
      let init = true;

      const cells = f(es);
      const rms = Array(cells.length).fill(undefined);
      cells.forEach((c, idx) => {
        c._rosterMember.subscribe((rm) => {
          rms[idx] = rm;
          if (!init) {
            set(rms);
          }
        });
      });
      set(rms);
    });

    return unsubAll;
  });
}

class row {
  memberId: number;

  member: resourceWritable<Member>;
  _member: Writable<Member>;

  encounters: Readable<RosterMember[]>;
  encounterAPI: apiResource<RosterMember[]>;

  constructor(s: store, memberId: number) {
    this.memberId = memberId;

    this._member = writable(undefined);
    this.member = new resourceWritable(
      `/team/${s.teamId}/member/${this.memberId}`,
      (m) => {
        // preload the class specs
        Classes.then((cs) => {
          const c = cs.get(m.classId);
          c.specializations.forEach((spec) => Spec(spec.id));
        });
      },
      this._member
    );

    this.encounters = rosterMembersDerived(s.Encounters, (es) =>
      es.map((e) => s.cell(this.memberId, e.id))
    );
    this.encounterAPI = new apiResource(
      `/team/${s.teamId}/member/${this.memberId}/encounters`,
      (rms) => {
        if (rms) {
          rms.forEach((m) => {
            if (m) {
              const cell = s.cell(m.member_id, m.encounter_id);
              cell._rosterMember.set(m);
              cell.rosterMember.state = LoadingState.Loaded;
            }
          });
        }
      }
    );
  }
}

class column {
  encounterId: number;

  encounter: resourceWritable<Encounter>;
  _encounter: Writable<Encounter>;

  members: Readable<RosterMember[]>;
  memberAPI: apiResource<RosterMember[]>;

  constructor(s: store, encounterId: number) {
    this.encounterId = encounterId;

    this._encounter = writable(undefined);
    this.encounter = new resourceWritable(
      `/team/${s.teamId}/encounter/${this.encounterId}`,
      undefined,
      this._encounter
    );

    this.members = rosterMembersDerived(s.Members, (ms) =>
      ms.map((m) => s.cell(m.id, this.encounterId))
    );
    this.memberAPI = new apiResource(
      `/team/${s.teamId}/encounter/${this.encounterId}/encounters`,
      (rms) => {
        if (rms) {
          rms.forEach((m) => {
            if (m) {
              const cell = s.cell(m.member_id, m.encounter_id);
              cell._rosterMember.set(m);
              cell.rosterMember.state = LoadingState.Loaded;
            }
          });
        }
      }
    );
  }
}

class store {
  teamId: string;

  private cells: Map<string, cell>;
  private rows: Map<number, row>;
  private columns: Map<number, column>;

  Team: resourceWritable<Team>;
  Members: resourceWritable<Member[]>;
  Encounters: resourceWritable<Encounter[]>;

  constructor(teamId: string) {
    this.teamId = teamId;

    this.cells = new Map();
    this.rows = new Map();
    this.columns = new Map();

    this.Team = new resourceWritable(`/team/${teamId}`);

    this.Members = new resourceWritable(`/team/${teamId}/members`);
    this.Members.subscribe((ms) => {
      if (!ms) return;
      ms.forEach((m) => {
        const row = this.row(m.id);
        row._member.set(m);
      });
    });

    this.Encounters = new resourceWritable(`/team/${teamId}/encounters`);
    this.Encounters.subscribe((es) => {
      if (!es) return;
      es.forEach((e) => {
        const col = this.column(e.id);
        col._encounter.set(e);
      });
    });
  }

  row(memberId: number): row {
    if (this.rows.has(memberId)) {
      return this.rows.get(memberId);
    }
    const r = new row(this, memberId);
    this.rows.set(memberId, r);
    return r;
  }

  column(encounterId: number): column {
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

  cell(memberId: number, encounterId: number): cell {
    const key = this.cellKey(memberId, encounterId);
    if (this.cells.has(key)) {
      return this.cells.get(key);
    }
    const c = new cell(this, this.row(memberId), this.column(encounterId));
    this.cells.set(key, c);
    return c;
  }

  newMember(m: Member): Promise<row> {
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

  newEncounter(e: Encounter): Promise<column> {
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
