import type { Readable } from "svelte/store";
import { PathPart } from "../url";

function createTeamId(): Readable<string> {
  const TeamPath = PathPart(0);
  TeamPath.subscribe((param) => {
    if (!param) {
      fetch("/team", { method: "POST" })
        .then((r) => {
          const createdId = r.headers
            .get("Location")
            .match(/\/team\/([^/]+)/)[1];
          TeamPath.set(createdId);
        })
        .catch((e) => console.error("error creating new team", e));
    }
  });

  return {
    subscribe: TeamPath.subscribe,
  };
}

export const TeamId = createTeamId();

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

export interface Team {
  id: string;
  name: string;
}

export interface Encounter {
  id: number;
  name: string;
  events?: Event[];
}

export interface Event {
  id: number;
  label: string;
  color: string;
  instances: EventInstance[];
}

export interface EventInstance {
  id: number;
  offset_sec: number;
}

export interface Assignment {
  id: number;
  member: number;
  spell_id?: number;
}

export interface RosterMember {
  encounter_id: number;
  member_id: number;
  spec?: number;
}
