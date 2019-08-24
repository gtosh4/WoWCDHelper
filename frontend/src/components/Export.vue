<template>
<v-card>
  <v-tabs v-model="tab">
    <v-tab v-for="key in tabOrder" :key="key" :href="`#${key}`">
      <span>{{ formats[key].name }}</span>
    </v-tab>
  </v-tabs>

  <v-select v-if="tab == 'ertp'" v-model="selectedPlayer" :items="players" dense />

  <v-card-text class="export-content">
    <v-textarea auto-grow readonly full-width outlined :value="text" />
  </v-card-text>

  <v-card-actions>
    <v-spacer />
    <v-btn text @click="$emit('close')">Close</v-btn>
  </v-card-actions>
</v-card>
</template>
<script>
import { sortEvents } from '../store/modules/events'
import {classes} from './wow_info'
import Color from 'color'

function formatRows(formatLabel, formatPlayer, formatSpell) {
  return (events, assigns) => {
    const rows = sortEvents(events).map(event => {
      return `${formatLabel(event)}\t` + event.assignments.map(assignId => {
        const assign = assigns[assignId]
        const player = assigns[assign.playerId]
        if (player !== undefined && assign.id != player.id) {
          return `${formatPlayer(player)}${formatSpell(assign)}`
        } else {
          return formatPlayer(assign)
        }
      }).join("\t")
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
          player => (player && player.name) ? `|c${Color.rgb(classes[player.className].colour).hex().replace(/#/, 'ff')}${player.name}|r` : '',
          assign => `{spell:${assign.spell.id}}`,
        )
      },
      ertp: {
        name: 'ERT Personal',
        run: (events, assigns, name, selectedPlayer) => {
          const fmt = formatRows(
            event => `{time:${event.time.minutes()}:${event.time.seconds()}}\t${colouredLabel(event)}\t`,
            () => '',
            assign => `{spell:${assign.spell.id}}`,
          )
          const playerAssigns = new Set(Object.values(assigns).filter(a => a.playerId == selectedPlayer).map(a => a.id))
          const playerEvents = Object.fromEntries(Object.entries(events)
            .map(([k, v]) => [k, {...v, assignments: v.assignments.filter(a => playerAssigns.has(a))}])
            .filter(([, v]) => v.assignments.length > 0))
          return fmt(playerEvents, assigns)
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
  }),

  props: {
  },

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
        this.selectedPlayer)
    },

    players() {
      return [...this.$store.getters['assigns/players'].map(a => ({text: a.name, value: a.id}))]
    },
  },

  methods: {
  },

  mounted() {
  },

  components: {
  },
};
</script>
<style>
.export-content {
  overflow-y: auto;
  max-height: 70vh;
}
</style>
