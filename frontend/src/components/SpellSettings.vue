<template>
<v-card>
  <v-card-title>
    {{ player.name }}: {{ spell.name }}
  </v-card-title>
</v-card>
</template>
<script>
import { assignProps, player } from '../store/utils'

export default {
  data: () => ({
  }),

  props: {
    assignId: {
      required: true,
    },
  },

  computed: {
    ...assignProps(['name', 'className', 'specName', 'spell', 'playerId']),
    ...player(),
  },

  methods: {
  },

  mounted() {
  },

  components: {
  },
};
</script>
<style>
</style>
