function nextid(state) {
  let next = 1
  while (state.assigns[next] !== undefined) {
    next++
  }
  return next
}

const state = {
  assigns: {},
}

const getters = {
  players(state) {
    return Object.values(state.assigns).filter(a => a.spell === undefined)
  },

  spells(state) {
    return Object.values(state.assigns).filter(a => a.spell !== undefined)
  },

  nextid,
}

const actions = {

}

const mutations = {
  import(state, assigns) {
    state.assigns = {...assigns}
  },

  set(state, assign) {
    if (!assign) return

    if (assign.id === undefined) {
      do {
        assign.id = nextid(state)
      } while (state.assigns[assign.id])
    }
    if (assign.playerId === undefined) {
      assign.playerId = assign.id
    }
    state.assigns = {...state.assigns, [assign.id]: assign}
    return assign
  },

  delete(state, id) {
    const n = {...state.assigns}
    const assign = state.assigns[id]
    delete n[id]
    if (assign != null && assign.id == assign.playerId) {
      Object.values(n).forEach(a => {
        if (a.playerId == id) {
          delete n[a.id]
        }
      })
    }
    state.assigns = n
  },
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
}
