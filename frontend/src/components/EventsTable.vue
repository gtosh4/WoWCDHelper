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
      <Event :eventId="item.id" @config="configEvent=item.id" />
    </template>
  </v-data-table>
  <v-dialog persistent v-model="showConfig" max-width="300px">
    <EventConfig :eventId="configEvent" @close="configEvent = null" />
  </v-dialog>
</v-card>
</template>
<script>
import Event from './Event'
import EventConfig from './EventConfig'

export default {
  data: () => ({
    headers: [
        { text: 'Time',        value: 'time',    align: 'right' },
        { text: 'Label',       value: 'label',   align: 'left'  },
        { text: 'Assignments', value: 'assigns', align: 'left'  },
    ],
    configEvent: null,
  }),

  props: {
  },

  computed: {
    items() {
      return [
        ...this.$store.getters['events/orderedEvents'],
        {}, // Empty event for the "new" row
      ]
    },

    showConfig() {
      return this.configEvent != null
    }
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
    EventConfig,
  },
};
</script>
<style>
</style>
