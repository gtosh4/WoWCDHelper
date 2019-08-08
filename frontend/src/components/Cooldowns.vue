<template>
<v-container fluid>

  <v-layout>
    <v-flex xs12 class="pa-1">
      <v-card outlined tile>

        <v-list>
          <v-list-item v-for="(player, index) in sortedPlayers" :key="index" min-height="20px">
            <Assignee :assignId="player.id" />
          </v-list-item>
        </v-list>

      </v-card>
    </v-flex>
  </v-layout>

  <v-layout wrap>
    <v-flex xs12 class="pa-1">
      <v-card outlined tile>
        <v-chip
          label
          color="transparent"
          @click="expanded = {...expanded, healers: !expanded['healers']}"
          style="width: 100%"
        >
          <v-layout align-center fill-height>
            <span class="mx-1">healers</span>
            <v-icon :class="expandedClass(expanded['healers'])">$vuetify.icons.expand</v-icon>
          </v-layout>
        </v-chip>

        <v-list v-if="expanded['healers']">
          <v-list-item v-for="(healer, i) in healers" :key="i" class="player-select">
            <v-chip outlined small label class="mr-1">{{ specCount(healer.className, healer.specName) }}</v-chip>
            <v-icon @click="addPlayer(healer.className, healer.specName)" class="mx-1" small>mdi-account-plus</v-icon>
            <WowIcon :className="healer.className" />
            <WowIcon :className="healer.className" :specName="healer.specName" />
            <span :style="{color: classColour(healer.className)}" class="mx-1">{{healer.className}} - {{ healer.specName }}</span>
          </v-list-item>
        </v-list>
      </v-card>
    </v-flex>
  </v-layout>

  <v-layout wrap>
    <v-flex v-for="(classInfo, className) in classes" :key="className" class="pa-1">
      <v-card outlined tile>
        <v-chip
          label
          :text-color="classColour(className)"
          color="transparent"
          @click="expanded = {...expanded, [className]: !expanded[className]}"
          style="width: 100%"
        >
          <v-layout align-center fill-height>
            <WowIcon :className="className" />
            <span class="mx-1">{{ className }}</span>
            <v-icon :class="expandedClass(expanded[className])">$vuetify.icons.expand</v-icon>
          </v-layout>
        </v-chip>
        <v-list v-if="expanded[className]">
          <v-list-item class="player-select">
            <v-chip outlined small label class="mr-1">{{ classCount(className) }}</v-chip>
            <v-icon @click="addPlayer(className)" class="mx-1" small>mdi-account-plus</v-icon>
            <span :style="{color: classColour(className)}" class="mx-1">{{ className }}</span>
          </v-list-item>
          <v-divider class="mx-1" />
          <v-list-item v-for="(specInfo, specName) in classInfo.specs" :key="specName" class="player-select">
            <v-chip outlined small label class="mr-1">{{ specCount(className, specName) }}</v-chip>
            <v-icon @click="addPlayer(className, specName)" class="mx-1" small>mdi-account-plus</v-icon>
            <WowIcon :className="className" :specName="specName" />
            <span :style="{color: classColour(className)}" class="mx-1">{{ specName }}</span>
          </v-list-item>
        </v-list>
      </v-card>
    </v-flex>
  </v-layout>

</v-container>
</template>
<script>
import Assignee from './Assignee'
import WowIcon from './WowIcon'

import { mapGetters } from 'vuex'
import {classes, classColour, classIcon, specIcon, spec, healers} from './wow_info'

export default {
  data: () => ({
    classes,
    healers,
    expanded: {
      healers: true,
    },
  }),

  props: {
  },

  mounted() {
    this.$store.commit('assigns/set', {name: "Yorman", className: 'druid', specName: 'restoration', id: 1})
    this.$store.commit('assigns/set', {name: "Yellowy", className: 'priest', specName: 'holy', id: 2})
    this.$store.commit('assigns/set', {name: "Toshpal", className: 'paladin', specName: 'holy', id: 3})
  },

  computed: {
    ...mapGetters('assigns', ['players', 'abilities']),

    indexedPlayers() {
      return this.players.reduce((ms, m) => {
        let c = ms[m.className]
        if (!c) {
          c = {}
          ms[m.className] = c
        }
        let s = c[m.specName]
        if (!s) {
          s = []
          c[m.specName] = s
        }
        s.push(m)
        return ms
      }, {})
    },

    sortedPlayers() {
      return [...this.players].sort((a, b) => {
        let c = a.className.localeCompare(b.className)
        if (c != 0) return c
        c = a.specName.localeCompare(b.specName)
        return c
      })
    },
  },

  methods: {
    classIcon,
    specIcon,
    spec,
    classColour,

    classCount(className) {
      return Object.values(this.indexedPlayers[className] || {}).reduce((sum, spec) => sum + spec.length, 0)
    },

    specCount(className, specName) {
      return ((this.indexedPlayers[className] || {})[specName] || []).length
    },

    expandedClass(expanded) {
      const c = ['v-data-table__expand-icon']
      if (expanded) {
        c.push('v-data-table__expand-icon--active')
      }
      return c
    },

    addPlayer(className, specName) {
      this.$store.commit('assigns/set', {className, specName})
    },
  },

  components: {
    Assignee,
    WowIcon,
  },
};
</script>
<style>
.player-select {
  min-height: 20px !important;
  padding-left: 4px !important;
  padding-right: 4px !important;
}
</style>
