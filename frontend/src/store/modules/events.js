const state = {
  events: {},
  nextid: 1,
}

const getters = {

}

const actions = {

}

const mutations = {
  set(state, event) {
    if (!event) return
    event = {
      label: "",
      assignments: [],
      ...event,
    }

    if (!event.id) {
      do {
        event.id = state.nextid++
      } while (state.events[event.id])
    }
    state.events = {...state.events, [event.id]: event}
  },

  delete(state, id) {
    const n = {...state.assigns}
    delete(n, id)
    state.assigns = n
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
