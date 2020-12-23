<template>
  <v-app>
    <v-main>
      <v-container fluid>
        <v-row no-gutters>
          <Planner />
          <Cooldowns />
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<script lang="ts">
import Vue from 'vue'
import Planner from './components/Planner.vue'
import Cooldowns from './components/Cooldowns.vue'

import {loadFromURL} from './store/plugins/url'

export default Vue.extend({
  name: 'App',

  components: {
    Planner,
    Cooldowns,
  },
  
  watch: {
    '$route':{
      handler(to) {
        loadFromURL(this.$store, to)
      },
      immediate: true,
    },
  },
})
</script>

<style lang="scss">
#app {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
