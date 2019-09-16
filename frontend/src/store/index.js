import Vue from 'vue'
import Vuex from 'vuex'
import assigns from './modules/assigns'
import events from './modules/events'
import dragassign from './modules/dragassign'

Vue.use(Vuex)

import { updateURLPlugin } from './plugins/url'

const debug = process.env.NODE_ENV !== 'production'

function clearAssign(state, id) {
  const playerId = state.assigns.assigns[id].playerId
  const deleteplayer = playerId && id == playerId

  const es = {...state.events.events}
  Object.values(es).forEach(e => e.assignments = e.assignments.filter(a => {
    const isDeletedId = a == id
    const isDeletedPlayer = deleteplayer && state.assigns.assigns[a].playerId == playerId
    return !isDeletedId && !isDeletedPlayer
  }))
  state.events.events = es
}

export default new Vuex.Store({
  mutations: {
    clearAssign,

    deleteAssign(state, id) {
      const playerId = state.assigns.assigns[id].playerId
      const deleteplayer = playerId && id == playerId

      clearAssign(state, id)
      
      const n = {...state.assigns.assigns}
      delete n[id]
      if (deleteplayer) {
        Object.values(n).forEach(a => {
          if (a.playerId == playerId) {
            delete n[a.id]
          }
        })
      }
      state.assigns.assigns = n
    },

    setName(state, name) {
      state.name = name
    },

    setLogURL(state, logURL) {
      state.logURL = logURL
    },
  },
  modules: {
    assigns,
    events,
    dragassign,
  },
  plugins: [
    updateURLPlugin,
  ],
  strict: debug,
})
