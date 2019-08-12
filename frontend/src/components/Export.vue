<template>
<v-card>
  <v-tabs v-model="tab">
    <v-tab v-for="key in tabOrder" :key="key" :href="`#${key}`">
      <span>{{ formats[key].name }}</span>
    </v-tab>
  </v-tabs>

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
function formatRows(formatPlayer, formatSpell) {
  return (events, assigns) => {
    const rows = Object.values(events).map(event => {
      return `${event.label}\t` + event.assignments.map(assignId => {
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

export default {
  data: () => ({
    formats: {
      aa: {
        name: 'Angry Assignments',
        run: formatRows(
          player => `|c${player.className}${player.name}|r`,
          assign => `{icon ${assign.spell.id}}`
        ),
      },
      ert: {
        name: 'Exorsus Raid Tools',
      },
      ertp: {
        name: 'ERT Personal',
      },
      json: {
        name: 'JSON',
        run: (events, assigns) => JSON.stringify({events, assigns}, null, 2),
      },
    },
    tabOrder: ['aa', 'ert', 'ertp', 'json'],
    tab: "aa",
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
      return format.run(this.$store.state.events.events, this.$store.state.assigns.assigns)
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
