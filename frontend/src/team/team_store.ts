import { writable } from "svelte/store";
import type { Readable, Updater } from "svelte/store";
import type { Member, RosterMember } from "./team_api";

interface AsyncRWStore<T> extends Readable<Promise<T>> {
  update(f: Updater<T>): Promise<T>;
}

interface cellKey {
  memberId: number;
  encounterId: number;
}

interface cell extends row, column {
  rosterMember: AsyncRWStore<RosterMember>;
}

interface row {
  teamId: string;
  memberId: number;
  memberInfo: AsyncRWStore<Member>;
  memberEncounters: AsyncRWStore<RosterMember[]>;
}

interface column {
  teamId: string;
  encounterId: number;
  encounterMembers: AsyncRWStore<RosterMember[]>;
}

class TeamStore {
  teamId: string;
  cells: Map<cellKey, cell>;

  constructor(teamId: Readable<string>) {
    this.cells = new Map();
    this.teamId = teamId;

    teamId.subscribe(() => {});
  }
}
