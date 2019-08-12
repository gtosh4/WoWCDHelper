<template>
<v-chip v-if="draggedAssign" label :color="color" class="insert-assign" />
</template>
<script>
import moment from 'moment'
import { eventProps, assignProps } from '../store/utils'

export default {
  data: () => ({
  }),

  props: {
    eventId: {
      required: true,
    },
    draggedAssign: {
      type: Object,
    },
  },

  computed: {
    ...eventProps(['time']),

    assignId() {
      return this.draggedAssign ? this.draggedAssign.assignId : undefined
    },
    ...assignProps(['spell']),

    cooldown() {
      return this.spell ? moment.duration(this.spell.cd, 'seconds') : undefined
    },

    reorder() {
      return this.draggedAssign && this.draggedAssign.sourceId ? this.draggedAssign.sourceId == this.eventId : false
    },

    color() {
      if (!this.reorder &&
        (this.spell && this.spell.cd) && 
        (this.prevTime && (this.time.asSeconds() - this.prevTime.asSeconds() < this.spell.cd))
      ) {
        return "red lighten-4"
      }
      return "green lighten-4"
    },

    prevTime() {
      if (this.draggedAssign == null || !this.time || !this.spell) return null
      const events = this.$store.getters["events/orderedEvents"].filter(
        e => e.assignments.indexOf(this.assignId) >= 0 && e.time && e.time.asSeconds() < this.time.asSeconds()
      )
      if (events.length == 0) return null
      return events[events.length-1].time
    },
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
.insert-assign {
  width: 25px;
  opacity: 0.8;
}
</style>
