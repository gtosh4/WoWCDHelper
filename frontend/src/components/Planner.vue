<template>
<v-flex py-0 grow id="cd-planner"><v-card outlined tile>
  <v-card-title class="pa-1"><v-toolbar dense flat>
    <v-toolbar-title>
      <v-text-field
        single-line
        hide-details
        placeholder="Assignments"
        id="plan-name"
        v-model="name"
      />
    </v-toolbar-title>
    <v-toolbar-items>
      <v-tooltip top>
        <template #activator="{ on }">
          <v-btn v-on="on" tile x-small icon @click="clearAll">
            <v-icon>mdi-backspace</v-icon>
          </v-btn>
        </template>
        <span>Clear All Assignments</span>
      </v-tooltip>
      <v-tooltip top>
        <template #activator="{ on }">
          <v-btn v-on="on" tile x-small icon @click="deleteAll">
            <v-icon>mdi-delete</v-icon>
          </v-btn>
        </template>
        <span>Delete All Events</span>
      </v-tooltip>
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
        <Event :eventId="item.id" />
      </template>

      <template #footer>
        <v-footer>
          <v-spacer />

          <v-dialog persistent v-model="showImport" @keydown.esc.stop="showImport = false">
            <template #activator="{ on: dialog }">
              <v-tooltip top>
                <template #activator="{ on: tooltip }">
                  <v-btn v-on="{...tooltip, ...dialog}" tile x-small icon>
                    <v-icon>mdi-import</v-icon>
                    </v-btn>
                </template>
                <span>Import</span>
              </v-tooltip>
            </template>

            <Import @close="showImport = false" />
          </v-dialog>

          <v-dialog persistent v-model="showExport" @keydown.esc.stop="showExport = false">
            <template #activator="{ on: dialog }">
              <v-tooltip top>
                <template #activator="{ on: tooltip }">
                  <v-btn v-on="{...tooltip, ...dialog}" tile x-small icon>
                    <v-icon>mdi-export</v-icon>
                  </v-btn>
                </template>
                <span>Export</span>
              </v-tooltip>
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
    ],
    assignees: {},
    showExport: false,
    showImport: false,
  }),

  props: {
  },

  computed: {
    items() {
      return [...this.$store.getters['events/orderedEvents'], {}]
    },

    name: {
      get() {
        return this.$store.state.name
      },

      set(v) {
        this.$store.commit("setName", v)
      }
    }
  },

  methods: {
    addItem() {
      const event = {}
      this.$store.commit('events/set', event)
      this.$nextTick(() => this.$el.querySelector(`#event-${event.id} input`).focus())
    },

    clearAll() {
      this.$store.commit('events/clearAll')
    },

    deleteAll() {
      this.$store.commit('events/import', {})
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
#cd-planner tr:nth-child(even) {
  background-color: "grey";
}
</style>
