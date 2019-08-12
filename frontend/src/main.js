import Vue from 'vue'
import Sortable from 'sortablejs'
import VueResource from 'vue-resource'

import App from './App.vue'
import vuetify from './plugins/vuetify';
import store from './store'

Vue.use(VueResource)

Vue.directive('sortable', {
  inserted: function (el, binding) {
    new Sortable(el, binding.value || {})
  }
});

Vue.config.productionTip = false

new Vue({
  render: h => h(App),
  vuetify,
  store,
}).$mount('#app')
