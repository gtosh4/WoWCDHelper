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
  props: {
    assignId: {
      type: String,
      required: true,
    },
  },

  computed: {
    ...assignProps(['name', 'className', 'specName']),
  },
}
</script>
<style>
</style>
