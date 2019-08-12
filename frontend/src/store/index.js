import Vue from 'vue'
import Vuex from 'vuex'
import assigns from './modules/assigns'
import events from './modules/events'
import VuexPersist from 'vuex-persist'

Vue.use(Vuex)

import moment from 'moment'
import { spec } from '../components/wow_info';

const vuexPersist = new VuexPersist({
  key: 'wow-cd-helper',
  storage: localStorage,

  saveState(key, state, storage) {
    storage.setItem(key,JSON.stringify({
      name: state.name,
      events: state.events.events,
      assigns: state.assigns.assigns,
    }))
  },

  restoreState(key, storage) {
    const raw = JSON.parse(storage.getItem(key)) || {}
    Object.values(raw.events || {}).forEach(v => v.time = moment.duration(v.time))
    Object.values(raw.assigns).forEach(a => {
      if (a.spell) {
        const player = raw.assigns[a.playerId]
        if (player.specName) {
          spec(player.className, player.specName).spells.forEach(specSpell => {
            if (specSpell.id == a.spell.id) {
              a.spell = specSpell
            }
          })
        }
      }
    })
    return {
      name: raw.name,
      events: {events: raw.events},
      assigns: {assigns: raw.assigns},
    }
  },
})


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
  },
  modules: {
    assigns,
    events,
  },
  plugins: [vuexPersist.plugin],
  strict: debug,
})
