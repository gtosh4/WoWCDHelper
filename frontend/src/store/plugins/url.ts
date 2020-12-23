import moment from 'moment'
import { Route } from 'vue-router'
import { MutationPayload, Store } from 'vuex'
import { State } from '..'
import router from '../../plugins/router'
import { Assign, AssignState, PlayerAssign, SpellAssign } from '../modules/assigns'
import { Event, EventState } from '../modules/events'
import { DragAssignState } from '../modules/dragassign'

import jsonurl from 'json-url'
const codec = jsonurl('lzma')

function shortenEvent(event: Event) {
  return {
    id: event.id!,
    t: event.time?.asSeconds(),
    a: event.assignments,
    l: event.label,
    c: event.colour,
  }
}

function expandEvent(e: any): Event {
  return {
    id: e.id,
    time: e.t ? moment.duration(e.t, 'seconds') : undefined,
    assignments: e.a,
    label: e.l,
    colour: e.c,
  }
}

function shortenAssign(assign: Assign) {
  const a = {id: `${assign.id}`} as any
  if (assign.type == "player") {
    a.n = assign.name
    a.c = assign.className
    if (assign.specName) {
      a.s = assign.specName
    }
  } else if (assign.type == "spell") {
    // ONLY set `p` for non-player assigns
    a.p = `${assign.playerId}`
    a.s = {id: assign.spell.id, c: assign.spell.cfg}
  }
  return a
}

function expandAssign(a: any): Assign {
  if (a.p) {
    const assign = {
      id: a.id,
      playerId: a.p,
      type: "spell",
    } as SpellAssign
    if (a.s) {
      assign.spell = {
        id: a.s.id,
        cfg: a.s.c,
      }
    }
    return assign
  } else {
    const assign = {
      id: a.id,
      name: a.n,
      playerId: a.id,
      className: a.c,
      specName: a.s,
      type: "player",
    } as PlayerAssign
    return assign
  }
}

function shortenState(state: any): any {
  const s = {
    n: state.name,
    u: state.logURL,
  } as any
  if ('assigns' in state) {
    const assigns = (<AssignState> state.assigns).assigns
    s.a = Object.values(assigns).map(shortenAssign)
  }
  if ('events' in state) {
    const events = (<EventState> state.events).events
    s.e = Object.values(events).map(shortenEvent)
  }
  return s
}

function expandState(s: any): State {
  const state = {
    name: s.n,
    logURL: s.u,
    events: {
      events: {},
    } as EventState,
    assigns: {
      assigns: {},
    } as AssignState,
    dragassign: {} as DragAssignState,
  }

  if (s.a) {
    s.a.forEach((a: any) => {
      const assign = expandAssign(a)
      state.assigns.assigns[assign.id] = assign
    })
  }

  if (s.e) {
    s.e.forEach((e: any) => {
      const event = expandEvent(e)
      state.events.events[event.id!] = event
    })
  }
  return state
}

export function loadFromURL(store: Store<State>, route: Route) {
  const data = route.params.data
  if (!data) {
    return
  }

  codec.decompress(data).then(short => {
    const json = JSON.parse(short)
    const state = expandState(json)
    console.info("url", {json, data, state}) // eslint-disable-line
    store.replaceState(state)
  }).catch(err => console.error("Error decompressing", {data, err}))
}

export function updateURLPlugin(store: Store<State>) {
  const ignorePatterns = [
    /^dragassign/,
  ]
  store.subscribe((mutation: MutationPayload, state: State) => {
    const ignored = ignorePatterns.some(p => p.test(mutation.type))
    if (ignored) {
      return
    }
    const short = shortenState(state)
    const str = JSON.stringify(short)

    codec.compress(str).then(data => {
      console.info("url", {short, data, state}) // eslint-disable-line
      if (data == router.currentRoute.params.data) {
        return
      }
      const r = {
        path: `/${data}`,
        query: router.currentRoute.query,
      }
      router.push(r)
    }).catch(err => console.error("Error compressing", {state, short, err}))
  })
}
