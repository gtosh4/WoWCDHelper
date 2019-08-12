function nextid(state) {
  let next = 1
  while (state.events[next] !== undefined) {
    next++
  }
  return next
}

const state = {
  events: {},
}

const getters = {
  orderedEvents(state) {
    return [...Object.values(state.events)].sort((a, b) => {
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
      return a.id > b.id
    })
  },
  nextid,
}

const actions = {

}

const mutations = {
  import(state, events) {
    state.events = {...events}
  },

  set(state, event) {
    if (!event) return

    if (event.id === undefined) {
      do {
        event.id = nextid(state)
      } while (state.events[event.id])
    }
    if (event.assignments === undefined) {
      event.assignments = []
    }
    state.events = {...state.events, [event.id]: event}
  },

  delete(state, id) {
    const n = {...state.events}
    delete n[id]
    state.events = n
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
      [event.id]: {
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

  moveAssignment(state, {from, to}) {
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
        [event.id]: {
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
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
}
