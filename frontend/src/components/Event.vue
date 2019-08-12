<template>
<tr @mouseenter="hover = true" @mouseleave="hover = false" :id="rowId">
  <td class="event-time">
    <EventTextField v-model="timeStr" placeholder="0:0" />
  </td>

  <td class="event-label">
    <EventTextField v-model="label" :placeholder="eventId ? 'unnamed' : 'new'" />
  </td>

  <td class="event-assignments">
    <v-layout fluid wrap fill-height align-center justify-start
      class="assignments"
      @drop="handleDrop"
      @dragover.prevent="handleDragOver"
      @dragleave.prevent="handleDragLeave"
    >
      <Assignment :eventId="eventId" :index="index" v-for="(assign, index) in assignments" :key="index" />
      <InsertAssignment :eventId="eventId" :draggedAssign="draggedAssign" class="ml-1" />
    </v-layout>
  </td>

  <td class="event-actions">
    <div v-if="showActions">
    <v-tooltip top>
      <template #activator="{ on }">
        <v-btn tile x-small icon tabindex="-1" @click="clone" v-on="on">
          <v-icon small>mdi-content-copy</v-icon>
        </v-btn>
      </template>
      <span>Copy</span>
    </v-tooltip>
    <v-tooltip top>
      <template #activator="{ on }">
        <v-btn v-if="showActions" tile x-small icon tabindex="-1" @click="clear" v-on="on">
          <v-icon small>mdi-arrow-collapse-left</v-icon>
        </v-btn>
      </template>
      <span>Clear assignments</span>
    </v-tooltip>
    <v-tooltip top>
      <template #activator="{ on }">
        <v-btn v-if="showActions" tile x-small icon tabindex="-1" @click="remove" v-on="on">
          <v-icon small>mdi-delete</v-icon>
        </v-btn>
      </template>
      <span>Delete row</span>
    </v-tooltip>
    </div>
  </td>
</tr>
</template>
<script>
import Assignment from './Assignment'
import EventTextField from './EventTextField'
import InsertAssignment from './InsertAssignment'

import moment from 'moment'
import { eventProps } from '../store/utils';

function padTime(t) {
  return `${t < 10 ? '0' : ''}${t}`
}

const endingNum = /(\d+)$/

export default {
  data: () => ({
    hover: false,
    draggedAssign: null,
  }),

  props: {
    eventId: {
      required: true,
    },
  },

  mounted() {
  },

  computed: {
    ...eventProps(['time', 'label', 'assignments']),

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

    showActions() {
      return this.hover && !this.draggedOver && this.eventId !== undefined
    }
  },

  methods: {
    handleDragOver(event) {
      const assignId = event.dataTransfer.getData("assignId")
      if (assignId) {
        const sourceId = event.dataTransfer.getData("eventId")
        const sourceIdx = event.dataTransfer.getData("assignIndex")
        if (sourceId) {
          event.dataTransfer.dropEffect = "move"
        } else {
          event.dataTransfer.dropEffect = "link"
        }
        if (sourceId != this.eventId || sourceIdx == "" || +sourceIdx < this.assignments.length-1) {
        this.draggedAssign = {assignId, sourceId, assignIndex: sourceIdx}
        }
      }
    },

    handleDragLeave() {
      this.draggedAssign = null
    },

    handleDrop(event) {
      const assignId = event.dataTransfer.getData("assignId")
      if (assignId) {
        event.preventDefault()
        event.stopPropagation();
        const sourceId = event.dataTransfer.getData("eventId")
        const sourceIdx = event.dataTransfer.getData("assignIndex")
        if (sourceId && sourceIdx) {
          this.$store.commit('events/removeAssignment', {id: sourceId, index: sourceIdx})
        }
        this.$store.commit('events/addAssignment', {id: this.eventId, assignId: assignId})
        this.draggedAssign = null
      }
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
      this.$store.commit('events/set', {id: this.eventId, time: this.time, label: this.label})
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
  white-space: nowrap;
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

.assignments .assignment {
  padding-left: 2px;
  padding-right: 2px;
}
</style>
