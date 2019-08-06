<template>
  <v-app>
    <v-content>
      <Planner />
    </v-content>
  </v-app>
</template>
<script>
import Planner from './components/Planner';

export default {
  name: 'App',
  components: {
    Planner,
  },
  data: () => ({
  }),
};
</script>
