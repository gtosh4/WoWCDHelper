import moment from 'moment'
import {Module} from 'vuex'
import { State } from '..'
import {toColor} from '../../components/colour_utils'
import Vue from 'vue'

function nextEventId(state: EventState) {
  let next = 1
  while (state.events[next] !== undefined) {
    next++
  }
  return next
}

export function sortEvents(events: {[id: string]: Event}) {
  return Object.values(events).sort((a, b) => {
    if (a.time != null && typeof(a.time.asSeconds) == 'function' && b.time != null && typeof(b.time.asSeconds) == 'function') {
      const t = a.time.asSeconds() - b.time.asSeconds()
      if (t != 0) return t
    }
    if (a.time != null) {
      return -1
    }
    if (b.time != null) {
      return 1
    }
    return a.id! - b.id!
  })
}

export interface Event {
  id?: number
  assignments: number[]
  colour?: {r: number, g: number, b: number, a?: number}
  time?: moment.Duration
  label: string
}

export interface AssignMove {
  from: {id: number, index: number}
  to: {id: number, index: number}
}

export interface EventState {
  events: {[id: string]: Event}
}

export default {
  namespaced: true,
  
  state: {
    events: {},  
  },

  getters: {
    ordered(state) {
      return sortEvents(state.events)
    },
  
    allEventColours(state) {
      return [...Object.values(state.events).reduce((s, e) => {
        const c = toColor(e.colour)
        if (c) {
          s.add(c.hex())
        }
        return s
      }, new Set<string>())]
    },
  
    nextEventId,
  },

  mutations: {
    import(state, events: {[id: string]: Event}) {
      state.events = {...events}
    },
  
    set(state, event?: Event) {
      if (!event) return
  
      if (event.id === undefined) {
        event.id = nextEventId(state)
      }
      if (event.assignments === undefined) {
        event.assignments = []
      }
      state.events = {...state.events, [event.id]: event}
    },
  
    delete(state, id) {
      Vue.delete(state.events, id)
    },
  
    clearAll(state) {
      const es = {...state.events}
      Object.values(es).forEach(e => e.assignments = [])
      state.events = es
    },

    clearAssignment(state, assignId) {
      const es = {...state.events}
      Object.values(es).forEach(e => e.assignments = [...e.assignments.filter(a => a != assignId)])
      state.events = es
    },
  
    addAssignment(state, {id, assignId, index}) {
      const event = state.events[id]
      if (!event) return
      if (index === undefined) {
        index = event.assignments.length
      }
      const assignments = [...event.assignments]
      assignments.splice(index, 0, assignId)
  
      state.events = {
        ...state.events,
        [event.id!]: {
          ...event,
          assignments,
        },
      }
    },
  
    removeAssignment(state, {id, index}) {
      const event = state.events[id]
      if (!event) return
  
      const assignments = [...event.assignments]
      assignments.splice(index, 1)
  
      state.events = {
        ...state.events,
        [id]: {
          ...event,
          assignments,
        },
      }
    },
  
    moveAssignment(state, move: AssignMove) {
      const from = move.from, to = move.to

      if (from.id == to.id) {
        const event = state.events[from.id]
        const assignments = event.assignments
        const assignId = assignments.splice(from.index, 1)[0]
        if (from.index < to.index) {
          to.index--
        }
        assignments.splice(to.index, 0, assignId)
  
        state.events = {
          ...state.events,
          [event.id!]: {
            ...event,
            assignments,
          }
        }
      } else {
        const fromE = state.events[from.id],
          toE = state.events[to.id]
        const fromAs = [...fromE.assignments],
          toAs = [...toE.assignments]
        const assignId = fromAs.splice(from.index, 1)[0]
        toAs.splice(to.index, 0, assignId)
  
        state.events = {
          ...state.events,
          [from.id]: {
            ...fromE,
            assignments: fromAs,
          },
          [to.id]: {
            ...toE,
            assignments: toAs,
          }
        }
      }
    },
  },
} as Module<EventState, State>
