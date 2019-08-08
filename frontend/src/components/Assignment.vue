<template>
<v-chip
  v-if="!ability"
  label
  close
  class="assignment player"
  :color="classColour"
  @click:close="remove"
>
  <v-icon class="handle">mdi-drag</v-icon>
  <span>{{ name }}</span>
</v-chip>
<v-chip label class="assignment ability" v-else>
</v-chip>
</template>
<script>
import {classes, classIcon, specIcon, spec} from './wow_info'
import { mapState } from 'vuex';

export default {
  data: () => ({
  }),

  props: {
    eventId: {
      required: true,
    },
    assignId: {
      required: true,
    },
  },

  computed: {
    ...mapState('assigns', {
      name(state) {
        return state.assigns[this.assignId].name
      },
      className(state) {
        return state.assigns[this.assignId].className
      },
      specName(state) {
        return state.assigns[this.assignId].specName
      },
      ability(state) {
        return state.assigns[this.assignId].ability
      },
    }),
    classColour() {
      const c = classes[this.className].colour
      return `rgba(${c.r}, ${c.g}, ${c.b}, 0.5)`
    }
  },

  methods: {
    classIcon,
    specIcon,
    spec,

    remove() {
      this.$store.commit('events/removeAssignment', {id: this.eventId, assignId: this.assignId})
    },
  },
};
</script>
<style>
.assignment {
  margin-left: 4px;
  margin-right: 4px;
}
</style>
