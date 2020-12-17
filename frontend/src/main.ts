import Vue from 'vue'
import VueResource from 'vue-resource'

import App from './App.vue'
import vuetify from './plugins/vuetify'
import router from './plugins/router'
import store from './store'

Vue.use(VueResource)

Vue.config.productionTip = false

new Vue({
  render: h => h(App),
  vuetify,
  store,
  router,
}).$mount('#app')
