<template>
  <v-col id="cd-planner" py-0>
    <v-card outlined tile>
      <v-card-title class="pa-1">
        <v-toolbar dense flat>
          <v-toolbar-title>
            <v-text-field
              id="plan-name"
              v-model="name"
              single-line
              hide-details
              placeholder="Assignments"
            />
          </v-toolbar-title>
          <v-toolbar-items>
            <v-tooltip top>
              <template #activator="{ on }">
                <v-btn
                  tile
                  x-small
                  icon
                  v-on="on"
                  @click="clearAll"
                >
                  <v-icon>mdi-backspace</v-icon>
                </v-btn>
              </template>
              <span>Clear All Assignments</span>
            </v-tooltip>
            <v-tooltip top>
              <template #activator="{ on }">
                <v-btn
                  tile
                  x-small
                  icon
                  v-on="on"
                  @click="deleteAll"
                >
                  <v-icon>mdi-delete</v-icon>
                </v-btn>
              </template>
              <span>Delete All Events</span>
            </v-tooltip>

          </v-toolbar-items>
          <v-spacer />
          <v-toolbar-items>
            <v-tabs v-model="tab">
              <v-tab href="#planner">
                Table
              </v-tab>
              <v-tab href="#timeline">
                Timeline
              </v-tab>
            </v-tabs>
          </v-toolbar-items>
        </v-toolbar>
      </v-card-title>

      <v-tabs-items v-model="tab">

        <v-tab-item value="planner">
          <EventsTable />
        </v-tab-item>

        <v-tab-item value="timeline">
          <v-lazy>
            <LogTimeline />
          </v-lazy>
        </v-tab-item>

      </v-tabs-items>

      <v-card-actions>
        <v-dialog
          v-model="showImport"
          persistent
          max-width="800"
          @keydown.esc.stop="showImport = false"
        >
          <template #activator="{ on: dialog }">
            <v-tooltip top>
              <template #activator="{ on: tooltip }">
                <v-btn
                  tile
                  x-small
                  class="mr-1"
                  v-on="{...tooltip, ...dialog}"
                >
                  <v-icon>mdi-import</v-icon>Import
                </v-btn>
              </template>
              <span>Import</span>
            </v-tooltip>
          </template>

          <Import @close="showImport = false" />
        </v-dialog>

        <v-dialog
          v-model="showExport"
          persistent
          max-width="800"
          @keydown.esc.stop="showExport = false"
        >
          <template #activator="{ on: dialog }">
            <v-tooltip top>
              <template #activator="{ on: tooltip }">
                <v-btn tile x-small v-on="{...tooltip, ...dialog}">
                  <v-icon>mdi-export</v-icon>Export
                </v-btn>
              </template>
              <span>Export</span>
            </v-tooltip>
          </template>

          <Export @close="showExport = false" />
        </v-dialog>
      </v-card-actions>
    </v-card>
  </v-col>
</template>
<script>
import EventsTable from './EventsTable'
import LogTimeline from './LogTimeline'
import Import from './Import'
import Export from './Export'

export default {
  components: {
    EventsTable,
    LogTimeline,
    Import,
    Export,
  },

  data: () => ({
    showExport: false,
    showImport: false,
  }),

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
}
</script>
<style>
#cd-planner {
  padding: 0 4px 0 0;
}
</style>
