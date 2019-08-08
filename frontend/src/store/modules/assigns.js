const state = {
  assigns: {},
  nextid: 1,
}

const getters = {
  players(state) {
    return Object.values(state.assigns).filter(a => a.ability === undefined)
  },

  abilities(state) {
    return Object.values(state.assigns).filter(a => a.ability !== undefined)
  },
}

const actions = {

}

const mutations = {
  set(state, assign) {
    if (!assign) return

    if (!assign.id) {
      do {
        assign.id = state.nextid++
      } while (state.assigns[assign.id])
    }
    state.assigns = {...state.assigns, [assign.id]: assign}
  },

  delete(state, id) {
    const n = {...state.assigns}
    delete n[id]
    Object.values(n).forEach(a => {
      if (a.parentId == id) {
        delete n[a.id]
      }
    })
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
