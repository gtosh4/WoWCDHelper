const state = {
  info: null,
}

const getters = {
}

const actions = {
}

const mutations = {
  start(state, info) {
    state.info = info
  },

  stop(state) {
    state.info = null
  },
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations,
}
