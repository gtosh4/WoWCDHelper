<template>
<v-flex py-0 grow id="cd-planner"><v-card outlined tile>
  <v-card-title><v-toolbar dense>
    <v-toolbar-title>
      <v-text-field
        single-line
        hide-details
        placeholder="Assignments"
        id="plan-name"
      />
    </v-toolbar-title>
    <v-spacer />
    <v-toolbar-items>
      <v-btn tile icon><v-icon>mdi-save</v-icon></v-btn>
    </v-toolbar-items>
  </v-toolbar></v-card-title>

  <v-card outlined tile>
    <v-data-table
      :headers="headers"
      :items="items"
      :items-per-page="-1"
      hide-default-header
      hide-default-footer
    >

      <template #item="{ item }">
        <Event :eventId="item.id" :id="`event-${item.id === undefined ? 'new': item.id}`" :class="item ? 'event-even' : 'event-odd' " />
      </template>

      <template #footer>
        <v-footer>
          <v-spacer />

          <v-dialog persistent v-model="showImport" @keydown.esc.stop="showImport = false">
            <template #activator="{ on }">
              <v-btn v-on="on" tile x-small icon><v-icon>mdi-import</v-icon></v-btn>
            </template>

            <Import @close="showImport = false" />
          </v-dialog>

          <v-dialog persistent v-model="showExport" @keydown.esc.stop="showExport = false">
            <template #activator="{ on }">
              <v-btn v-on="on" tile x-small icon><v-icon>mdi-export</v-icon></v-btn>
            </template>

            <Export @close="showExport = false" />
          </v-dialog>
        </v-footer>
      </template>

    </v-data-table>
  </v-card>
</v-card></v-flex>
</template>
<script>
import Event from './Event'
import Import from './Import'
import Export from './Export'

export default {
  data: () => ({
    headers: [
        { text: 'Time',        value: 'time',    align: 'right' },
        { text: 'Label',       value: 'label',   align: 'left'  },
        { text: 'Assignments', value: 'assigns', align: 'left',  width: '100%' },
        { text: '',            value: 'clear',   align: 'right'},
    ],
    assignees: {},
    showExport: false,
    showImport: false,
  }),

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
    Import,
    Export,
  },
};
</script>
<style>
#cd-planner {
  padding: 0 4px 0 0;
}
</style>
