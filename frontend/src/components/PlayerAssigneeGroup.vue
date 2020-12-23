<template>
  <v-card outlined tile class="assignee-group">
    <v-card-title class="assignee-title">
      <v-btn
        tile
        text
        left
        width="100%"
        @click="expanded = !expanded"
      >
        <span>assignees</span>
        <v-icon :class="expandedClass">
          $vuetify.icons.expand
        </v-icon>
      </v-btn>
    </v-card-title>

    <v-list :style="{display: expanded ? '' : 'none'}">
      <v-list-item v-for="(player, index) in sortedPlayers" :key="index" class="player-container">
        <PlayerAssignee :assign-id="player.id" @config="openSettings" />
      </v-list-item>
    </v-list>

    <v-dialog v-model="showSettings">
      <AssigneeSettings
        v-if="showSettings"
        :assign-id="settingsAssignId"
        @close="settingsAssignId = null"
      />
    </v-dialog>
  </v-card>
</template>
<script>
import PlayerAssignee from './PlayerAssignee.vue'
import AssigneeSettings from './AssigneeSettings.vue'

import { mapGetters } from 'vuex'

export default {
  components: {
    PlayerAssignee,
    AssigneeSettings,
  },

  data: () => ({
    expanded: true,
    settingsAssignId: null,
  }),

  computed: {
    ...mapGetters('assigns', ['players', 'spells']),

    expandedClass() {
      const c = ['v-data-table__expand-icon']
      if (this.expanded) {
        c.push('v-data-table__expand-icon--active')
      }
      return c
    },

    showSettings: {
      get() {
        // return this.settingsAssignId != null
        return false
      },

      set(v) {
        if (!v) {
          this.settingsAssignId = null
        }
      },
    },

    sortedPlayers() {
      return [...this.players].sort((a, b) => {
        let c = a.className.localeCompare(b.className)
        if (c != 0) return c
        if (a.specName && b.specName) {
          c = a.specName.localeCompare(b.specName)
          if (c != 0) return c
        }
        c = a.id - b.id
        return c
      })
    },
  },

  methods: {
    openSettings(assignId) {
      this.settingsAssignId = assignId
    },
  },
}
</script>
<style lang="scss">
.assignee-group {
  .assignee-title.v-card__title {
    padding: 0;
  }
  
  .v-list-item {
    padding-left: 0;
    padding-right: 0;
  }

  .player-container.v-list-item {
    min-height: 20px;
    padding-left: 4px;
    margin-bottom: 4px;

    &:not(:first-child) {
      padding-top: 4px;
    }
  }
}
</style>
