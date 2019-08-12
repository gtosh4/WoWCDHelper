import Vue from 'vue'
import Vuex from 'vuex'
import assigns from './modules/assigns'
import events from './modules/events'
import VuexPersist from 'vuex-persist'

Vue.use(Vuex)

import moment from 'moment'

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
    return {
      events: {events: raw.events},
      assigns: {assigns: raw.assigns},
    }
  },
})


const debug = process.env.NODE_ENV !== 'production'

export default new Vuex.Store({
  modules: {
    assigns,
    events,
  },
  plugins: [vuexPersist.plugin],
  strict: debug,
})
