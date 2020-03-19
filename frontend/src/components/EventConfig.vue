<template>
  <v-card tile @keydown.esc.stop="close">
    <v-color-picker v-model="colourPicker" :swatches="swatches" show-swatches flat hide-mode-switch />

    <v-card-actions>
      <v-btn @click="clear">Clear</v-btn>
      <v-spacer />
      <v-btn @click="discard">Discard</v-btn>
      <v-btn @click="save">Save</v-btn>
    </v-card-actions>
  </v-card>
</template>
<script>
import { eventProps } from '../store/utils'

export default {
  data: () => ({
    colourPicker: null,
  }),

  props: {
    eventId: {
      required: true,
    },
  },

  computed: {
    ...eventProps(['colour']),

    swatches() {
      return this.$store.getters["events/allEventColours"].map(c => [c])
    },
  },

  watch: {
    colour: {
      handler() {
        this.colourPicker = this.colour
      },
      immediate: true,
    },
  },

  methods: {
    close() {
      this.$emit('close')
    },

    clear() {
      this.colour = null
      this.close()
    },

    discard() {
      this.close()
    },

    save() {
      this.colour = this.colourPicker
      this.close()
    },
  },
}
</script>
