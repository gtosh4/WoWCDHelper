import Vue from 'vue'
import Vuex from 'vuex'
import assigns from './modules/assigns'
import dragassign from './modules/dragassign'
import events from './modules/events'
import { updateURLPlugin } from './plugins/url'

Vue.use(Vuex)

/* global process */
const debug = process.env.NODE_ENV !== 'production'

export interface State {
  name: string
  logURL?: string
}

export default new Vuex.Store<State>({
  state: {
    name: "",
  },

  mutations: {
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
