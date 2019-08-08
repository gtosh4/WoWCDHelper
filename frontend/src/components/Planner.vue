<template>
<v-container fluid>
  <v-layout>
    <v-flex>
      <v-data-table
        :headers="headers"
        :items="items"
        :items-per-page="-1"
        hide-default-header
        hide-default-footer
      >

        <template v-slot:item="{ item }">
          <Event :eventId="item.id" />
        </template>

        <template v-slot:footer>
          <v-footer>
            <v-edit-dialog :return-value.sync="addItem"><v-icon>mdi-plus-circle</v-icon>
              <template v-slot:input>
                <v-text-field
                  placeholder="time"
                  single-line
                />
              </template>
            </v-edit-dialog>
            <v-spacer />
            <v-dialog persistent v-model="showExport">
              <template v-slot:activator="{ on }">
                <v-icon @click="on.click">mdi-export</v-icon>
              </template>

              <Export @close="showExport = false" />
            </v-dialog>
          </v-footer>
        </template>

      </v-data-table>
    </v-flex>
    <v-flex xs4>
      <Cooldowns />
    </v-flex>
  </v-layout>
</v-container>
</template>
<script>
import Cooldowns from './Cooldowns'
import Event from './Event'
import Export from './Export'

import moment from 'moment'

export default {
  data: () => ({
    headers: [
        { text: 'Time',        value: 'time',    align: 'right' },
        { text: 'Label',       value: 'label',   align: 'left'  },
        { text: 'Assignments', value: 'assigns', align: 'left',  width: '100%'  },
        { text: '',            value: 'clear',   align: 'right', width: '1px'   },
    ],
    assignees: {},
    showExport: false,
  }),

  computed: {
    items() {
      return [...Object.values(this.$store.state.events.events)].sort((a, b) => {
        const t = a.time.asSeconds() - b.time.asSeconds()
        if (t != 0) return t
        return a.label > b.label
      })
    },

    addItem: {
      get() {
        return undefined
      },
      set(v) {
        const [mins, secs] = v.split(":")
        this.$store.commit('events/set', {time: moment.duration(+mins, 'minutes').add(+secs, 'seconds')})
      },
    }
  },

  methods: {
  },

  mounted() {
    this.$store.commit('events/set', {time: moment.duration(0, 'seconds'), label: "Event 1", assignments: []})
    this.$store.commit('events/set', {time: moment.duration(15, 'seconds'), label: "Event 2", assignments: []})
    this.$store.commit('events/set', {time: moment.duration(1.5, 'minutes'), label: "Event 2", assignments: []})
  },

  components: {
    Event,
    Cooldowns,
    Export,
  },
};
</script>
