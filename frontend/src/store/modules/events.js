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

  addAssignment(state, {id, assignId}) {
    const event = state.events[id]
    if (!event) return
    state.events = {
      ...state.events,
      [event.id]: {
        ...event,
        assignments: [...event.assignments, assignId],
      },
    }
  },

  removeAssignment(state, {id, assignId}) {
    const event = state.events[id]
    if (!event) return
    state.events = {
      ...state.events,
      [id]: {
        ...event,
        assignments: event.assignments.filter(a => a != assignId),
      },
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
