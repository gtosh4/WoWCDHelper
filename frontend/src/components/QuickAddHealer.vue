<template>
<v-card outlined tile class="add-players-quick">
  <v-card-title>
    <v-btn tile text left
      @click="expanded = !expanded"
      width="100%"
    >
      <span>healers</span>
      <v-icon :class="expandedClass">$vuetify.icons.expand</v-icon>
    </v-btn>
  </v-card-title>

  <v-list :style="{display: expanded ? '' : 'none'}">
    <v-list-item v-for="(healer, i) in healers" :key="i" class="player-select">
      <WowIcon :className="healer.className" />
      <WowIcon :className="healer.className" :specName="healer.specName" />
      <span :style="{color: classColour(healer.className)}" class="mx-1">{{healer.className}} - {{ healer.specName }}</span>
      <div class="add-player-action">
        <span>{{ counts[healer.className][healer.specName] }}</span>
        <v-btn tile x-small icon @click="$emit('addPlayer', healer.className, healer.specName)">
          <v-icon>mdi-account-plus</v-icon>
        </v-btn>
      </div>
    </v-list-item>
  </v-list>
</v-card>
</template>
<script>
import WowIcon from './WowIcon'

import { mapGetters } from 'vuex'
import { classColour, healers } from './wow_info'

export default {
  data: () => ({
    healers,
    expanded: true,
  }),

  props: {
  },

  computed: {
    ...mapGetters('assigns', ['players']),

    expandedClass() {
      const c = ['v-data-table__expand-icon']
      if (this.expanded) {
        c.push('v-data-table__expand-icon--active')
      }
      return c
    },

    counts() {
      const cnt = healers.reduce((c, {className, specName}) => {
        if (!c[className]) c[className] = {}
        c[className][specName] = 0
        return c
      }, {})
      return this.players.reduce((cnt, {className, specName}) => {
        const classCnt = cnt[className]
        if (classCnt != null) {
          if (classCnt[specName] != null) {
            classCnt[specName]++
          }
        }
        return cnt
      }, cnt)
    },
  },

  methods: {
    classColour,

    specCount(className, specName) {
      return ((this.players[className] || {})[specName] || []).length
    },
  },

  components: {
    WowIcon,
  },
}
</script>
<style>
.add-players-quick > .v-chip {
  width: 100%;
}

.add-players-quick .v-card__title {
  padding: 0;
}

.add-players-quick .v-list {
  padding-top: 0;
}

.add-players-quick .v-list-item {
  min-height: 20px;
  padding-left: 4px;
}


.add-players-quick .add-player-action {
  margin-left: auto;
}

.add-players-quick .add-player-action span {
  margin-right: 8px;
}
</style>
