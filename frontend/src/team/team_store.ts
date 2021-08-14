import { writable, derived, readable } from "svelte/store";
import type { Readable, Writable } from "svelte/store";
import {
  Encounter,
  Member,
  RosterMember,
  Team,
  TeamId,
  Event,
} from "./team_api";
import { Classes, Spec } from "../wow/api";
import { apiResource, LoadingState, resourceWritable } from "../store_helpers";
import { encounterEvent } from "../assignments/assignment_api";

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

      init = false;
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
  private teamId: string;
  encounterId: number;

  encounter: resourceWritable<Encounter>;
  _encounter: Writable<Encounter>;

  members: Readable<RosterMember[]>;
  memberAPI: apiResource<RosterMember[]>;

  events: resourceWritable<Event[]>;
  private _events: Writable<Event[]>;

  private _eventStores: Map<number, encounterEvent>;

  constructor(s: store, encounterId: number) {
    this.teamId = s.teamId;
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

    this._events = writable(undefined);
    this.events = new resourceWritable(
      `/team/${s.teamId}/envounter/${this.encounterId}/events`,
      undefined,
      this._events
    );

    this._eventStores = new Map();
  }

  public event(eventId: number): encounterEvent {
    if (this._eventStores.has(eventId)) {
      return this._eventStores.get(eventId);
    }
    const s = new encounterEvent(this.teamId, this.encounterId, eventId, (v) =>
      this._events.update((evts) => {
        const idx = evts ? evts.findIndex((e) => e.id == v.id) : -1;
        if (idx >= 0) {
          evts[idx] = v;
        }
        return evts;
      })
    );
    this._eventStores.set(eventId, s);
    return s;
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
