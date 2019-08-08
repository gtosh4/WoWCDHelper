import Vue from 'vue'
import Vuex from 'vuex'
import assigns from './modules/assigns'
import events from './modules/events'

Vue.use(Vuex)

const debug = process.env.NODE_ENV !== 'production'

export default new Vuex.Store({
  modules: {
    assigns,
    events
  },
  strict: debug,
})
