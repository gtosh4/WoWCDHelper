import Vue from 'vue'
import VueRouter from 'vue-router'

Vue.use(VueRouter)

const routes = [
  {
    path: '/:data?',
  },
]

export default new VueRouter({
  routes
})
