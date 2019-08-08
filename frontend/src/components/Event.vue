<template>
<tr>
  <td class="text-right">
    <v-edit-dialog :return-value.sync="timeStr">{{ timeStr }}
      <template v-slot:input>
        <v-text-field
          v-model="timeStr"
          single-line
        />
      </template>
    </v-edit-dialog>
  </td>
  <td style="border-right: solid 1px rgba(255, 255, 255, 0.12); white-space: nowrap;">
    <v-edit-dialog :return-value.sync="label">{{ label || 'unnamed' }}
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
      <Assignment :eventId="eventId" :index="index" v-for="(assign, index) in assignments" :key="index" />
      <v-chip v-if="draggedOver" label disabled class="grey lighten-4" />
    </v-layout>
  </td>
  <td style="border-left: solid 1px rgba(255, 255, 255, 0.12); white-space: nowrap;">
    <v-icon small @click="clear">mdi-arrow-collapse-left</v-icon>
    <v-icon small class="pl-1" @click="remove">mdi-delete</v-icon>
  </td>
</tr>
</template>
<script>
import Assignment from './Assignment'

import moment from 'moment'
import {classIcon, specIcon, spec} from './wow_info'
import { eventProps } from '../store/utils';

export default {
  data: () => ({
    draggedOver: false,
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
        const sourceIdx = event.dataTransfer.getData("assignIndex")
        if (sourceId) {
          event.dataTransfer.dropEffect = "move"
        } else {
          event.dataTransfer.dropEffect = "link"
        }
        if (sourceIdx == "" || +sourceIdx < this.assignments.length-1) {
          this.draggedOver = true
        }
      }
    },

    handleDragLeave() {
      this.draggedOver = false
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
        this.draggedOver = false
      }
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
  },
};
</script>
<style>
</style>
