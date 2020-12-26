<template>
  <v-card>
    <v-tabs v-model="tab">
      <v-tab v-for="key in tabOrder" :key="key" :href="`#${key}`">
        <span>{{ formats[key].name }}</span>
      </v-tab>
    </v-tabs>

    <v-row class="mx-0">
      <v-checkbox v-model="ignoreEmpty" label="Ignore Empty" />
    </v-row>

    <v-select
      v-if="tab == 'ertp'"
      v-model="selectedPlayer"
      :items="players"
      clearable
      dense
      label="Select a Player"
      class="mx-2"
    />

    <v-card-text class="export-content">
      <v-textarea
        auto-grow
        readonly
        full-width
        outlined
        :value="text"
      />
    </v-card-text>

    <v-card-actions>
      <v-spacer />
      <v-btn text @click="$emit('close')">
        Close
      </v-btn>
    </v-card-actions>
  </v-card>
</template>
<script>
import { sortEvents } from '../store/modules/events'

import {formatDuration} from './duration_utils'
import {classes} from './wow_info'
import Color from 'color'

function formatRows(formatLabel, formatPlayer, formatSpell) {
  return (events, assigns, name, config) => {
    var fmtEvts = sortEvents(events)
    if (config.ignoreEmpty) {
      fmtEvts = fmtEvts.filter(e => e.assignments && e.assignments.length > 0)
    }
    const rows = fmtEvts.map(event => {
      const assignments = [...event.assignments.map(a => assigns[a])]
      return `${formatLabel(event)}\t` + assignments.reduce((evtTxts, a, idx) => {
        const player = assigns[a.playerId]
        if (a.spell) {
          if (player == null) {
            evtTxts.push(formatSpell(a))
          } else {
            if (idx > 0 && evtTxts.length > 0 && assignments[idx-1].playerId == a.playerId) {
              evtTxts[evtTxts.length-1] += `+${formatSpell(a)}`
            } else {
              evtTxts.push(`${formatPlayer(player)}${formatSpell(a)}`)
            }
          }
        } else {
          evtTxts.push(formatPlayer(a))
        }
        return evtTxts
      }, []).join("\t")
    })
    return rows.join("\n")
  }
}

function colouredLabel(event) {
  const color = event.colour ? Color.rgb(event.colour.r, event.colour.g, event.colour.b) : null
  return color ? `${color.hex().replace(/#/, '|cFF')}${event.label}|r` : event.label
}

export default {
  data: () => ({
    formats: {
      aa: {
        name: 'Angry Assignments',
        run: formatRows(
          colouredLabel,
          player => (player && player.name) ? `|c${player.className}${player.name}|r` : '',
          assign => `{icon ${assign.spell.id}}`
        ),
      },
      ert: {
        name: 'Exorsus Raid Tools',
        run: formatRows(
          colouredLabel,
          player => (player && player.name) ? `${Color.rgb(classes[player.className].colour).hex().replace(/#/, '|cFF')}${player.name}|r` : '',
          assign => `{spell:${assign.spell.id}}`,
        )
      },
      ertp: {
        name: 'ERT Personal',
        run: (events, assigns, name, config) => {
          const playerAssigns = new Set(Object.values(assigns).filter(a => a.playerId == config.selectedPlayer).map(a => a.id))
          const playerEvents = Object.fromEntries(Object.entries(events)
            .map(([k, v]) => [k, {...v, assignments: v.assignments.filter(a => playerAssigns.has(a))}])
            .filter(([, v]) => v.assignments.length > 0))
            
          const fmt = formatRows(
            event => `{time:${formatDuration(event.time)}}\t${colouredLabel(event)}\t`,
            () => '',
            assign => `{spell:${assign.spell.id}}`,
          )
          return fmt(playerEvents, assigns, name, config)
        }
      },
      json: {
        name: 'JSON',
        run: (events, assigns, name) => JSON.stringify({name, events, assigns}, null, 2),
      },
    },
    tabOrder: ['aa', 'ert', 'ertp', 'json'],
    tab: "aa",

    selectedPlayer: null,
    ignoreEmpty: true,
  }),

  computed: {
    text() {
      const format = this.formats[this.tab]
      if (format === undefined) return ""
      if (format.run === undefined) {
        return `Not implemented: ${format.name}`
      }
      return format.run(
        this.$store.state.events.events,
        this.$store.state.assigns.assigns,
        this.$store.state.name,
        this.config)
    },

    players() {
      return [...this.$store.getters['assigns/players'].map(a => ({text: a.name, value: a.id}))]
    },

    config() {
      return {
        ignoreEmpty: this.ignoreEmpty,
        selectedPlayer: this.selectedPlayer,
      }
    },
  },

  mounted() {
  },

  methods: {
  },
}
</script>
<style>
.export-content {
  overflow-y: auto;
  max-height: 70vh;
}
</style>
