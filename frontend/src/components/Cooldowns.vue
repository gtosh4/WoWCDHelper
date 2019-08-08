<template>
<v-container fluid>

  <v-layout wrap>
    <v-flex v-for="(classInfo, className) in classes" v-bind:key="className" class="pa-1">
      <v-card outlined>
        <v-chip
          label
          :text-color="classColour(className)"
          color="transparent"
          @click="expanded = {...expanded, [className]: !expanded[className]}"
          style="width: 100%"
        >
          <v-layout align-center fill-height>
            <v-avatar tile class="mr-1" size="24"><img :src="classIcon(className)" draggable="false" ondragstart="return false;" /></v-avatar>
            <span class="mx-1">{{ className }}</span>
            <v-icon :class="expandedClass(expanded[className])">$vuetify.icons.expand</v-icon>
          </v-layout>
        </v-chip>
        <v-list v-if="expanded[className]" dense>
          <v-list-item class="player-select">
            <v-chip outlined small label class="mr-1">{{ classCount(className) }}</v-chip>
            <v-icon @click="addPlayer(className)" class="mx-1" small>mdi-account-plus</v-icon>
            <span :style="{color: classColour(className)}" class="mx-1">{{ className }}</span>
          </v-list-item>
          <v-divider class="mx-1" />
          <v-list-item v-for="(specInfo, specName) in classInfo.specs" v-bind:key="specName" class="player-select">
            <v-chip outlined small label class="mr-1">{{ specCount(className, specName) }}</v-chip>
            <v-icon @click="addPlayer(className, specName)" class="mx-1" small>mdi-account-plus</v-icon>
            <v-avatar tile class="mx-1" size="18"><img :src="specIcon(specInfo)" draggable="false" ondragstart="return false;" /></v-avatar>
            <span :style="{color: classColour(className)}" class="mx-1">{{ specName }}</span>
          </v-list-item>
        </v-list>
      </v-card>
    </v-flex>
  </v-layout>

  <v-layout>
    <v-flex grow class="pa-1">
      <v-card outlined>

        <v-list dense>
          <v-list-item v-for="(player, index) in sortedPlayers" v-bind:key="index" min-height="20px">
            <Assignee :assignId="player.id" />
          </v-list-item>
        </v-list>

      </v-card>
    </v-flex>
  </v-layout>

</v-container>
</template>
<script>
import Assignee from './Assignee'

import { mapGetters } from 'vuex'
import {classes, classColour, classIcon, specIcon, spec} from './wow_info'

export default {
  data: () => ({
    classes,
    expanded: {},
  }),

  props: {
  },

  mounted() {
    this.$store.commit('assigns/set', {name: "Yorman", className: 'druid', specName: 'resto', id: 1})
    this.$store.commit('assigns/set', {name: "Test", className: 'druid', specName: 'resto', id: 2})
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
  },
};
</script>
<style>
.player-select {
  min-height: 20px;
  padding-left: 4px;
  padding-right: 4px;
}
</style>
