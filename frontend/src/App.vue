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
<script>
import Planner from './components/Planner';
import Cooldowns from './components/Cooldowns'

import {loadFromURL} from './store/plugins/url'

export default {
  name: 'App',
  data: () => ({
  }),
  mounted() {
  },
  watch: {
    '$route':{
      handler(to) {
        loadFromURL(this.$store, to)
      },
      immediate: true,
    },
  },
  methods: {
  },
  components: {
    Planner,
    Cooldowns,
  },
};
</script>
