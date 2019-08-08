<template>
<tr>
  <td class="text-right">
    <v-edit-dialog>{{ timeStr }}
      <template v-slot:input>
        <v-text-field
          v-model="timeStr"
          single-line
        />
      </template>
    </v-edit-dialog>
  </td>
  <td style="border-right: solid 1px rgba(255, 255, 255, 0.12)">
    <v-edit-dialog>{{ label }}
      <template v-slot:input>
        <v-text-field
          v-model="label"
          single-line
        />
      </template>
    </v-edit-dialog>
  </td>
  <td>
    <v-layout fluid wrap fill-height align-center justify-start
      class="assignments"
      v-sortable="sortable"
      @drop="handleDrop"
      @dragover.prevent
    >
      <Assignment :eventId="eventId" :assignId="assign" v-for="(assign, index) in assignments" v-bind:key="index" />
    </v-layout>
  </td>
  <td></td>
</tr>
</template>
<script>
import Assignment from './Assignment'

import moment from 'moment'
import {classIcon, specIcon, spec} from './wow_info'
import { eventProps } from '../store/utils';

export default {
  data: () => ({
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
        return `${this.time.minutes()}:${this.time.seconds()}`
      },
      set(v) {
        const [mins, secs] = v.split(":")
        this.time = moment.duration(+mins, 'minutes').add(+secs, 'seconds')
      },
    },

    sortable() {
      const event = this
      return {
        group: {
          name: `${this.label}-assignment`,
          pull: to => {
            return to.el.classList.contains("assignments")
          },
          put: to => {
            return to.el.classList.contains("assignments")
          },
        },
        handle: ".handle",
        setData(dt, el) {
          dt.setData("eventId", event.eventId)
          dt.setData("assignId", el.__vue__.assignId)
        },
        onAdd(evt) {
          const sourceId = evt.originalEvent.dataTransfer.getData("eventId")
          const assignId = evt.originalEvent.dataTransfer.getData("assignId")
          event.$store.commit('events/moveAssignment', {fromId: sourceId, toId: event.eventId, assignId: assignId})
        },
        onUpdate(evt) {
          const oldI = evt.oldIndex
          const newI = evt.newIndex
          const assignments = [...event.assignments]
          const old = assignments.splice(oldI, 1)[0]
          assignments.splice(newI, 0, old)
          event.assignments = assignments
        },
      }
    },
  },

  methods: {
    classIcon,
    specIcon,
    spec,

    handleDrop(event) {
      const assignId = event.dataTransfer.getData("assignId")
      const sourceId = event.dataTransfer.getData("eventId")
      if (assignId && !sourceId) {
        event.preventDefault()
        this.$store.commit('events/addAssignment', {id: this.eventId, assignId: assignId})
      }
    },
  },

  components: {
    Assignment,
  },
};
</script>
