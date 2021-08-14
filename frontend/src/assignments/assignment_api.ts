import { apiResource } from "../store_helpers";
import type { Assignment, Event } from "../team/team_api";

export class encounterEvent extends apiResource<Event> {
  private teamId: string;
  private encounterId: number;
  private _assignments: Map<number, apiResource<Assignment[]>>;

  constructor(
    teamId: string,
    encounterId: number,
    eventId: number,
    listener?: (v: Event) => void
  ) {
    super(
      `/team/${teamId}/envounter/${encounterId}/event/${eventId}`,
      listener
    );

    this.teamId = teamId;
    this.encounterId = encounterId;

    this._assignments = new Map();
  }

  public assignments(evtInstId: number): apiResource<Assignment[]> {
    if (this._assignments.has(evtInstId)) {
      return this._assignments.get(evtInstId);
    }
    const a = new apiResource<Assignment[]>(
      `/team/${this.teamId}/encounter/${this.encounterId}/assignments/${evtInstId}`
    );
    this._assignments.set(evtInstId, a);
    return a;
  }
}
