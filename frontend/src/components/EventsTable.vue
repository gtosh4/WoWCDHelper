<template>
<v-card outlined tile>
  <v-data-table
    :headers="headers"
    :items="items"
    :items-per-page="-1"
    hide-default-header
    hide-default-footer
  >
    <template #item="{ item }">
      <Event :eventId="item.id" />
    </template>
  </v-data-table>
</v-card>
</template>
<script>
import Event from './Event'

export default {
  data: () => ({
    headers: [
        { text: 'Time',        value: 'time',    align: 'right' },
        { text: 'Label',       value: 'label',   align: 'left'  },
        { text: 'Assignments', value: 'assigns', align: 'left'  },
    ],
  }),

  props: {
  },

  computed: {
    items() {
      return [...this.$store.getters['events/orderedEvents'], {}]
    },
  },

  methods: {
    addItem() {
      const event = {}
      this.$store.commit('events/set', event)
      this.$nextTick(() => this.$el.querySelector(`#event-${event.id} input`).focus())
    },
  },

  mounted() {
  },

  components: {
    Event,
  },
};
</script>
<style>
</style>
