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
      events: {events: raw.events},
      assigns: {assigns: raw.assigns},
    }
  },
})


const debug = process.env.NODE_ENV !== 'production'

export default new Vuex.Store({
  mutations: {
    deleteAssign(state, id) {
      const es = {...state.events.events}
      es.forEach(e => e.assignments.filter(a => a != id))
      state.events.events = es

      const n = {...state.assigns.assigns}
      const playerId = state.assigns.assigns[id].playerId
      delete n[id]
      if (playerId && id == playerId) {
        Object.values(n).forEach(a => {
          if (a.playerId == playerId) {
            delete n[a.id]
          }
        })
      }
      state.assigns.assigns = n
    },
  },
  modules: {
    assigns,
    events,
  },
  plugins: [vuexPersist.plugin],
  strict: debug,
})
