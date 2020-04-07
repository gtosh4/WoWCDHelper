<template>
<tr :id="rowId" class="event" :style="rowStyle"
  @mouseenter="hovered=true"
  @mouseleave="hovered=false"
>
  <td class="event-time" :style="timeSyle">
    <EventTextField v-model="timeStr" placeholder="0:0" />
  </td>

  <td class="event-label" :style="labelStyle">
    <EventTextField v-model="label" :placeholder="eventId ? 'unnamed' : 'new'" />
  </td>

  <td class="event-assignments" :style="assignmentsStyle">
    <v-row no-gutters>
      <v-col>
        <AssignmentGroup v-if="eventId !== undefined"
          :eventId="eventId"
        />
      </v-col>
      <v-col align-right cols="auto">
        <v-container pa-0 ma-0 fill-height>
          <EventActions v-if="eventId !== undefined"
            class="event-actions"
            @clone="clone"
            @clear="clear"
            @remove="remove"
            @config="$emit('config')"
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
import { eventProps } from '../store/utils'
import {toColor} from './colour_utils'
import {formatDuration} from './duration_utils'

const endingNum = /(\d+)$/

export default {
  data: () => ({
    draggedOver: false,
    hovered: false,
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

    rowStyle() {
      const c = this.colour ? toColor(this.colour) : null

      const s = {}
      if (c != null) {
        s["background-color"] = c.alpha(this.hovered ? 0.6 : 0.1).string()
      }
      return s
    },

    timeSyle() {
      const c = this.colour ? toColor(this.colour) : null

      const s = {}
      if (c != null) {
        s["border-left"] = `4px solid ${toColor(this.colour).string()}`
      }
      return s
    },

    labelStyle() {
      const c = this.colour ? toColor(this.colour) : null

      const s = {}
      s["border-right"] = `2px solid ${c != null ? toColor(this.colour).string() : 'rgba(255, 255, 255, 0.3)'}`
      return s
    },

    assignmentsStyle() {
      const c = this.colour ? toColor(this.colour) : null

      const s = {}
      if (c != null) {
        s["border-right"] = `4px solid ${toColor(this.colour).string()}`
      }
      return s
    }
  },

  methods: {
    clone() {
      let label = this.label
      const num = label.match(endingNum)
      if (num != null) {
        label = label.replace(endingNum,  `${(+num[0])+1}`)
      }
      this.$store.commit('events/set', {time: this.time, label, colour: this.colour, assignments: [...this.assignments]})
    },

    clear() {
      this.assignments = []
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
.v-data-table .event td {
  border-top: 2px solid rgba(255,255,255,0);
}

.v-data-table .event td.event-time {
  padding-left: 8px;
  padding-right: 8px;
  width: 8ch;
}
.v-data-table .event td.event-time input {
  text-align: end;
}

.v-data-table .event td.event-label {
  white-space: nowrap;
  padding-left: 8px;
  padding-right: 8px;
  width: 20ch;
}

.v-data-table .event td.event-assignments {
  padding-left: 8px;
  padding-right: 8px;
}

.v-data-table .event td.event-assignments .event-actions {
  padding-left: 4px;
  padding-right: 4px;
  display: none;
}

.v-data-table tr:hover.event td.event-assignments .event-actions {
  display: inherit;
}

</style>
