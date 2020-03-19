<template>
  <v-chip-group
    column
    class="assignmentGroup"
    @drop="handleDrop"
    @dragover.prevent="handleDragOver"
    @dragleave.prevent="handleDragLeave"
  >
    <Assignment :eventId="eventId" :index="index" v-for="(assign, index) in assignments" :key="index" class="mr-1" />
    <InsertAssignment v-if="draggedOver" />
  </v-chip-group>
</template>
<script>
import Assignment from './Assignment'
import InsertAssignment from './InsertAssignment'
import { eventProps, dragAssignProps } from '../store/utils'

export default {
  data: () => ({
    draggedOver: false,
  }),

  props: {
    eventId: {
    },
  },

  computed: {
    ...eventProps(['assignments']),
    ...dragAssignProps(),
  },

  methods: {
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
  },

  components: {
    Assignment,
    InsertAssignment,
  },
}
</script>
