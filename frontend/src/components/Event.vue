<template>
<tr @mouseenter="hover = true" @mouseleave="hover = false" :id="rowId" :style="style">
  <td class="event-time">
    <EventTextField v-model="timeStr" placeholder="0:0" />
  </td>

  <td class="event-label">
    <EventTextField v-model="label" :placeholder="eventId ? 'unnamed' : 'new'" />
  </td>

  <td class="event-assignments">
    <v-row no-gutters>
      <v-col>
        <AssignmentGroup v-if="eventId !== undefined"
          :eventId="eventId"
        />
      </v-col>
      <v-col align-right cols="auto">
        <v-container pa-0 ma-0 fill-height>
          <EventActions v-if="showActions"
            class="event-actions"
            @clone="clone"
            @clear="clear"
            @remove="remove"
            @colour="colourPicker = colour ? colour : toRGBA(defaultColour); showColour = true"
          />
        </v-container>
      </v-col>
    </v-row>
  </td>
</tr>
</template>
<script>
import EventTextField from './EventTextField'
import EventActions from './EventActions'
import AssignmentGroup from './AssignmentGroup'

import moment from 'moment'
import Color from 'color'
import { eventProps } from '../store/utils'
import {toColor, toRGBA} from './colour_utils'
import {formatDuration} from './duration_utils'

const endingNum = /(\d+)$/
const defaultColour = Color('rgb(66, 66, 66)')

export default {
  data: () => ({
    defaultColour,

    hover: false,
    draggedOver: false,

    showColour: false,
    colourPicker: null,
  }),

  props: {
    eventId: {
      required: true,
    },
  },

  mounted() {
  },

  computed: {
    ...eventProps(['time', 'label', 'assignments', 'colour']),

    timeStr: {
      get() {
        return formatDuration(this.time)
      },
      set(v) {
        if (v === undefined || v == "") {
          this.time = undefined
        } else {
          const [mins, secs] = v.split(":")
          this.time = moment.duration(+mins, 'minutes').add(+secs, 'seconds')
        }
      },
    },

    rowId() {
      return `event-${this.eventId === undefined ? 'new': this.eventId}`
    },

    style() {
      var colour = toColor(this.colour)
      colour = colour ? colour.darken(0.5) : defaultColour
      if (this.hover) colour = colour.lighten(0.5)
      
      return {
        "background-color": colour.string(),
      }
    },

    prevDragSpellTime() {
      if (!this.draggedAssign || !this.time) return
      const assign = this.$store.state.assigns.assigns[this.draggedAssign.assignId]
      if (!assign || !assign.spell) return

      const events = this.$store.getters["events/orderedEvents"].filter(e => 
        (assign.eventId === undefined || e.eventId != assign.eventId) &&
        e.assignments.indexOf(assign.id) >= 0 &&
        (e.time && e.time.asSeconds() < this.time.asSeconds())
      )
      if (events.length == 0) return
      return events[events.length-1].time
    },

    showActions() {
      return this.hover && this.eventId !== undefined && !this.draggedAssign
    },

    swatches() {
      return this.$store.getters["events/allEventColours"].map(c => [c])
    },
  },

  methods: {
    toRGBA,

    clone() {
      let label = this.label
      const num = label.match(endingNum)
      if (num != null) {
        label = label.replace(endingNum,  `${(+num[0])+1}`)
      }
      this.$store.commit('events/set', {time: this.time, label, colour: this.colour, assignments: [...this.assignments]})
    },

    clear() {
      const event = this.$store.state.events.events[this.eventId]
      this.$store.commit('events/set', {...event, assignments: []})
    },

    remove() {
      this.$store.commit('events/delete', this.eventId)
    },
  },

  components: {
    AssignmentGroup,
    EventTextField,
    EventActions,
  },
};
</script>
<style>

.v-data-table td.event-time {
  padding-left: 8px;
  padding-right: 8px;
  width: 8ch;
}
.v-data-table td.event-time input {
  text-align: end;
}

.v-data-table td.event-label {
  border-right: solid 1px rgba(255, 255, 255, 0.3);
  white-space: nowrap;
  padding-left: 8px;
  padding-right: 8px;
  width: 20ch;
}

.v-data-table td.event-assignments {
  padding-left: 8px;
  padding-right: 8px;
}

.v-data-table td.event-assignments .event-actions {
  padding-left: 4px;
  padding-right: 4px;
}

</style>
