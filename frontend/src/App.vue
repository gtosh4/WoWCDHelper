<template>
  <v-app>
    <v-content>
      <v-container fluid>
        <v-row no-gutters>
          <Planner />
          <Cooldowns />
        </v-row>
      </v-container>
    </v-content>
  </v-app>
</template>

<script lang="ts">
import Vue from 'vue'
import Planner from './components/Planner.vue'
import Cooldowns from './components/Cooldowns.vue'

import {loadFromURL} from './store/plugins/url'

export default Vue.extend({
  name: 'App',
  data: () => ({
  }),

  watch: {
    '$route':{
      handler(to) {
        loadFromURL(this.$store, to)
      },
      immediate: true,
    },
  },

  components: {
    Planner,
    Cooldowns,
  },
});
</script>

<style lang="scss">
#app {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
