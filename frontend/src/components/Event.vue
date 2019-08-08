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
      @drop="handleDrop"
      @dragover.prevent="handleDragOver"
      @dragleave.prevent="handleDragLeave"
    >
      <Assignment :eventId="eventId" :index="index" v-for="(assign, index) in assignments" v-bind:key="index" />
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
  },

  methods: {
    classIcon,
    specIcon,
    spec,

    handleDragOver(event) {
      const assignId = event.dataTransfer.getData("assignId")
      if (assignId) {
        const sourceId = event.dataTransfer.getData("eventId")
        if (sourceId) {
          event.dataTransfer.dropEffect = "move"
        } else {
          event.dataTransfer.dropEffect = "link"
        }
      }
    },

    handleDragLeave() {
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
      }
    },
  },

  components: {
    Assignment,
  },
};
</script>
<style scoped>
td {
  white-space: nowrap;
}
</style>
