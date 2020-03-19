<template>
<v-col shrink id="cd-palette"><v-card outlined tile>
  <v-row>
    <v-col>

      <v-list class="pb-0">
        <v-list-item v-for="(player, index) in sortedPlayers" :key="index" class="assignee">
          <v-card outlined tile width="100%" >
            <Assignee :assignId="player.id" />
            <v-list>
              <v-list-item v-for="(spell, i) in playerSpells[player.id]" :key="i" class="px-0">
                <Assignee :assignId="spell.id" />
              </v-list-item>
            </v-list>
          </v-card>
        </v-list-item>
      </v-list>

    </v-col>
  </v-row>

  <v-row wrap>
    <v-col xs12 class="pa-1">
      <v-card outlined tile>
        <v-chip
          label
          color="transparent"
          @click="expanded = {...expanded, healers: !expanded['healers']}"
          style="width: 100%"
        >
          <v-row align-center fill-height>
            <span class="mx-1">healers</span>
            <v-icon :class="expandedClass(expanded['healers'])">$vuetify.icons.expand</v-icon>
          </v-row>
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
    </v-col>
  </v-row>

  <v-row wrap>
    <v-col xs12 class="pa-1">
      <v-card outlined tile>
      <v-chip
        label
        color="transparent"
        @click="expanded = {...expanded, all: !expanded['all']}"
        style="width: 100%"
      >
        <v-row align-center fill-height>
          <span class="mx-1">all</span>
          <v-icon :class="expandedClass(expanded['all'])">$vuetify.icons.expand</v-icon>
        </v-row>
      </v-chip>
      <v-container v-if="expanded.all" grid-list-sm><v-row wrap>
        <v-col v-for="(classInfo, className) in classes" :key="className" class="pa-1">
            <v-chip
              label
              :text-color="classColour(className)"
              color="transparent"
              style="width: 100%"
            >
              <v-row align-center fill-height>
                <WowIcon :className="className" />
                <span class="mx-1">{{ className }}</span>
              </v-row>
            </v-chip>
            <v-list>
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
        </v-col>
      </v-row></v-container>
      </v-card>
    </v-col>
  </v-row>

</v-card></v-col>
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
    draggedAssign: Object,
  },

  mounted() {
  },

  computed: {
    ...mapGetters('assigns', ['players', 'spells']),

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
        if (a.specName && b.specName) {
          c = a.specName.localeCompare(b.specName)
          if (c != 0) return c
        }
        c = a.id - b.id
        return c
      })
    },

    playerSpells() {
      return (this.spells || []).reduce((m, spell) => {
        let pspells = m[spell.playerId]
        if (!pspells) {
          pspells = []
          m[spell.playerId] = pspells
        }
        pspells.push(spell)
        return m
      }, {})
    }
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

    addPlayer(className, specName, name, id) {
      const player = {className, specName, name, id}
      this.$store.commit('assigns/set', player)
      if (specName !== undefined) {
        spec(className, specName).spells.forEach(spell => {
          this.$store.commit('assigns/set', {spell, playerId: player.id, id: `${player.id}.${spell.id}`})
        })
      }
    },
  },

  components: {
    Assignee,
    WowIcon,
  },
};
</script>
<style>
#cd-palette {
  padding: 0 0 0 4px;
  min-width: 400px;
  max-width: 650px;
}
.v-list-item.player-select {
  min-height: 20px;
  padding-left: 4px;
  padding-right: 4px;
}
.v-list-item.assignee {
  min-height: 20px;
  padding-left: 4px;
  padding-right: 4px;
  margin-bottom: 4px;
}
.v-list-item.assignee .v-list {
  padding-left: 28px;
  padding-right: 28px;
}
.v-list-item.assignee .v-list-item {
  min-height: 20px;
  padding-left: 4px;
  margin-bottom: 4px;
}
</style>
