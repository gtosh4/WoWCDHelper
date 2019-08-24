<template>
<v-card>
  <v-card-title>
    {{ player.name }}
  </v-card-title>
</v-card>
</template>
<script>
import { assignProps } from '../store/utils'

export default {
  data: () => ({
  }),

  props: {
    assignId: {
      required: true,
    },
  },

  computed: {
    ...assignProps(['name', 'className', 'specName']),
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
