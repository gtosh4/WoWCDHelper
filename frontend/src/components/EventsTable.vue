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
        <Event :event-id="item.id" @config="configEvent=item.id" />
      </template>
    </v-data-table>
    <v-dialog v-model="showConfig" persistent max-width="300px">
      <EventConfig :event-id="configEvent" @close="configEvent = null" />
    </v-dialog>
  </v-card>
</template>
<script>
import Event from './Event'
import EventConfig from './EventConfig'

export default {

  components: {
    Event,
    EventConfig,
  },

  props: {
  },
  data: () => ({
    headers: [
      { text: 'Time',        value: 'time',    align: 'right' },
      { text: 'Label',       value: 'label',   align: 'left'  },
      { text: 'Assignments', value: 'assigns', align: 'left'  },
    ],
    configEvent: null,
  }),

  computed: {
    items() {
      return [
        ...this.$store.getters['events/ordered'],
        {}, // Empty event for the "new" row
      ]
    },

    showConfig() {
      return this.configEvent != null
    }
  },

  mounted() {
  },

  methods: {
    addItem() {
      const event = {}
      this.$store.commit('events/set', event)
      this.$nextTick(() => this.$el.querySelector(`#event-${event.id} input`).focus())
    },
  },
}
</script>
<style>
</style>
