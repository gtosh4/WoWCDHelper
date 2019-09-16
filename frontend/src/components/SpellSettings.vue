<template>
<v-card>
  <v-card-title>
    {{ player.name }}: {{ spell.name }}
  </v-card-title>
</v-card>
</template>
<script>
import { assignProps, player, spell } from '../store/utils'

export default {
  data: () => ({
  }),

  props: {
    assignId: {
      required: true,
    },
  },

  computed: {
    ...assignProps(['name', 'className', 'specName', 'playerId']),
    ...spell(),
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
