<template>
<tr @mouseenter="hover = true" @mouseleave="hover = false" :id="rowId" :style="style">
  <td class="event-time">
    <EventTextField v-model="timeStr" placeholder="0:0" />
  </td>

  <td class="event-label">
    <EventTextField v-model="label" :placeholder="eventId ? 'unnamed' : 'new'" />
  </td>

  <td class="event-assignments" :class="dragDestClass">
    <v-dialog persistent v-model="showColour" @keydown.esc.stop="showColour = false" max-width="300px">
      <v-card tile>
        <v-color-picker v-model="colourPicker" :swatches="swatches" show-swatches flat hide-mode-switch />

        <v-card-actions>
          <v-btn @click="showColour = false; colour = null">Clear</v-btn>
          <v-spacer />
          <v-btn @click="showColour = false">Discard</v-btn>
          <v-btn @click="colour = colourPicker; showColour = false">Save</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-container pa-0 ma-0 justify-start>
      <v-layout>
        <v-flex v-if="eventId !== undefined"
          grow
          class="assignments"
          @drop="handleDrop"
          @dragover.prevent="handleDragOver"
          @dragleave.prevent="handleDragLeave"
        >
          <v-layout>
            <Assignment :eventId="eventId" :index="index" v-for="(assign, index) in assignments" :key="index" />
            <InsertAssignment v-if="draggedOver" class="ml-1" />
          </v-layout>
        </v-flex>

        <v-flex v-if="showActions" shrink>
          <v-tooltip top>
            <template #activator="{ on }">
              <v-btn v-on="on" @click.stop="colourPicker = colour ? colour : toRGBA(defaultColour); showColour = true" tile x-small icon>
                <v-icon>mdi-palette</v-icon>
              </v-btn>
            </template>
            <span>Set Colour</span>
          </v-tooltip>

          <v-tooltip top>
            <template #activator="{ on }">
              <v-btn tile x-small icon tabindex="-1" @click="clone" v-on="on">
                <v-icon>mdi-content-copy</v-icon>
              </v-btn>
            </template>
            <span>Copy</span>
          </v-tooltip>

          <v-tooltip top>
            <template #activator="{ on }">
              <v-btn tile x-small icon tabindex="-1" @click="clear" v-on="on">
                <v-icon>mdi-backspace</v-icon>
              </v-btn>
            </template>
            <span>Clear assignments</span>
          </v-tooltip>

          <v-tooltip top>
            <template #activator="{ on }">
              <v-btn tile x-small icon tabindex="-1" @click="remove" v-on="on">
                <v-icon>mdi-delete</v-icon>
              </v-btn>
            </template>
            <span>Delete row</span>
          </v-tooltip>
        </v-flex>

      </v-layout>
      <v-layout style="height: 4px" />
    </v-container>
  </td>
</tr>
</template>
<script>
import Assignment from './Assignment'
import EventTextField from './EventTextField'
import InsertAssignment from './InsertAssignment'

import moment from 'moment'
import Color from 'color'
import { eventProps, dragAssignProps } from '../store/utils'
import {toColor, toRGBA} from './colour_utils'

function padTime(t) {
  return `${t < 10 ? '0' : ''}${t}`
}

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
    ...dragAssignProps(),

    timeStr: {
      get() {
        const t = this.time
        return t !== undefined ? `${t.minutes()}:${padTime(t.seconds())}` : ""
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
      var colour = toColor(this.colour) || defaultColour
      if (this.hover) {
        colour = colour.isLight() ? colour.darken(0.5) : colour.lighten(0.5)
      }
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

    dragDestClass() {
      if (!this.time || !this.draggedAssign) return

      const spell = this.$store.state.assigns.assigns[this.draggedAssign.assignId].spell

      if (!spell || !spell.cd) return

      const onCD = this.prevDragSpellTime
        ? (this.time.asSeconds() - this.prevDragSpellTime.asSeconds()) < spell.cd
        : false
      
      return onCD ? 'drag-on-cd' : 'drag-off-cd'
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

    handleDragOver(event) {
      if (!this.draggedAssign) return

      if (this.draggedAssign.sourceId) {
        event.dataTransfer.dropEffect = "move"
      } else {
        event.dataTransfer.dropEffect = "link"
      }
      const isNewAssignment = this.draggedAssign.sourceId === undefined
      const isSameEvent = !isNewAssignment && this.draggedAssign.sourceId == this.eventId
      const isReorder = +this.draggedAssign.sourceIndex < this.assignments.length-1

      if (isNewAssignment || !isSameEvent || !isReorder) {
        this.draggedOver = true
      }
    },

    handleDragLeave() {
      this.draggedOver = false
    },

    handleDrop(event) {
      if (!this.draggedAssign) return

      event.preventDefault()
      event.stopPropagation();

      if (this.draggedAssign.sourceId !== undefined && this.draggedAssign.sourceIndex !== undefined) {
        this.$store.commit('events/moveAssignment', {
          from: {id: this.draggedAssign.sourceId, index: this.draggedAssign.sourceIndex},
          to: {id: this.eventId, index: this.assignments.length}
        })
      } else {
        this.$store.commit('events/addAssignment', {id: this.eventId, assignId: this.draggedAssign.assignId})
      }
      this.draggedOver = false
    },

    clone() {
      let label = this.label
      const num = label.match(endingNum)
      if (num != null) {
        label = label.replace(endingNum,  `${(+num[0])+1}`)
      }
      this.$store.commit('events/set', {time: this.time, label, assignments: [...this.assignments]})
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
    Assignment,
    EventTextField,
    InsertAssignment,
  },
};
</script>
<style>
.event-label {
  border-right: solid 1px rgba(255, 255, 255, 0.12);
  white-space: nowrap;
}
.v-data-table td.event-actions {
  padding-left: 4px;
  padding-right: 4px;
}

.v-data-table td.event-time {
  padding-left: 8px;
  padding-right: 8px;
  min-width: 50px;
}
.v-data-table td.event-time input {
  text-align: end;
}

.v-data-table td.event-label {
  padding-left: 8px;
  padding-right: 8px;
  min-width: 100px;
}

.v-data-table td.event-assignments {
  padding-left: 8px;
  padding-right: 8px;
}

.assignments {
  flex-grow: 1;
}

.assignments .assignment {
  padding-left: 2px;
  padding-right: 2px;
}

.drag-on-cd {
  background-color: rgba(128, 25, 35, 0.8);
}

.drag-off-cd {
  background-color: rgba(102, 128, 103, 0.8);
}
</style>
