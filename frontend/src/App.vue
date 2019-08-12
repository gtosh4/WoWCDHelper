<template>
  <v-app>
    <v-content>
      <v-container fluid>
        <v-layout>
          <Planner />
          <Cooldowns />
        </v-layout>
      </v-container>
    </v-content>
  </v-app>
</template>
<script>
import Planner from './components/Planner';
import Cooldowns from './components/Cooldowns'

export default {
  name: 'App',
  components: {
    Planner,
    Cooldowns,
  },
  data: () => ({
  }),
};
</script>
