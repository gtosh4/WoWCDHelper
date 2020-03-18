<template>
<v-flex py-0 id="cd-planner"><v-card outlined tile>
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
    <v-spacer />
    <v-toolbar-items>
      <v-tabs v-model="tab">
        <v-tab href="#planner">Table</v-tab>
        <v-tab href="#timeline">Timeline</v-tab>
      </v-tabs>
    </v-toolbar-items>
  </v-toolbar></v-card-title>

  <v-tabs-items v-model="tab">

    <v-tab-item value="planner">
      <EventsTable />
    </v-tab-item>

    <v-tab-item value="timeline">
      <LogTimeline />
    </v-tab-item>

  </v-tabs-items>

  <v-card-actions>
    <v-dialog persistent v-model="showImport" @keydown.esc.stop="showImport = false" max-width="800">
      <template #activator="{ on: dialog }">
        <v-tooltip top>
          <template #activator="{ on: tooltip }">
            <v-btn v-on="{...tooltip, ...dialog}" tile x-small class="mr-1">
              <v-icon>mdi-import</v-icon>Import
            </v-btn>
          </template>
          <span>Import</span>
        </v-tooltip>
      </template>

      <Import @close="showImport = false" />
    </v-dialog>

    <v-dialog persistent v-model="showExport" @keydown.esc.stop="showExport = false" max-width="800">
      <template #activator="{ on: dialog }">
        <v-tooltip top>
          <template #activator="{ on: tooltip }">
            <v-btn v-on="{...tooltip, ...dialog}" tile x-small>
              <v-icon>mdi-export</v-icon>Export
            </v-btn>
          </template>
          <span>Export</span>
        </v-tooltip>
      </template>

      <Export @close="showExport = false" />
    </v-dialog>
  </v-card-actions>
</v-card></v-flex>
</template>
<script>
import EventsTable from './EventsTable'
import LogTimeline from './LogTimeline'
import Import from './Import'
import Export from './Export'

export default {
  data: () => ({
    showExport: false,
    showImport: false,
  }),

  props: {
  },

  computed: {
    name: {
      get() {
        return this.$store.state.name
      },

      set(v) {
        this.$store.commit("setName", v)
      }
    },

    tab: {
      get() {
        var t = this.$route.query.tab
        return t || 'planner'
      },

      set(v) {
        if (!v) return

        const r = {...this.$route}
        r.query = {...r.query, tab: v}
        this.$router.push(r)
      },
    },
  },

  methods: {
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
    EventsTable,
    LogTimeline,
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
