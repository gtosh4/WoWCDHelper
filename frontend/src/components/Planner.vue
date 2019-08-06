ret<template>
<v-container fluid>
  <v-layout>
    <v-flex>
    <v-data-table
      :headers="headers"
      :items="items"
      :items-per-page="-1"
      hide-default-header
      hide-default-footer
      :custom-sort="sort"
    >

      <template v-slot:item="{ item }">
        <Event
          :label="item.label" :time="item.time" :assigns="item.assigns"
          @update:label="(l) => item.label = l"
          @update:time="(t) => item.time = t"
          @update:assigns="(a) => item.assigns = a"
        />
      </template>

    </v-data-table>
    </v-flex>
    <v-flex xs4>
      <Cooldowns :assigns="allAssigns" />
    </v-flex>
  </v-layout>
</v-container>
</template>
<script>
import Cooldowns from './Cooldowns'
import Event from './Event'

import moment from 'moment'

export default {
  data: () => ({
    headers: [
        { text: 'Time',        value: 'time',    align: 'right', width: '100px'  },
        { text: 'Label',       value: 'label',   align: 'left'  },
        { text: 'Assignments', value: 'assigns', align: 'left'  },
        { text: '',            value: 'clear',   align: 'right',  sortable: false, width: '1px' },
    ],
    items: [
      {time: moment.duration(0, 'seconds'), label: "Event 1", assigns: []},
      {time: moment.duration(15, 'seconds'), label: "Event 2", assigns: []},
      {time: moment.duration(1.5, 'minutes'), label: "Event 2", assigns: []},
    ]
  }),

  computed: {
    allAssigns() {
      return this.items.map(i => i.assigns).flat()
    },
  },

  methods: {
    sort(items) {
      return items.sort((a, b) => {
        const t = a.time.asSeconds() - b.time.asSeconds()
        if (t != 0) return t
        return a.label > b.label
      })
    },
  },

  mounted() {
  },

  components: {
    Event,
    Cooldowns,
  },
};
</script>
