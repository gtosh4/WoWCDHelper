<template>
<v-chip-group column class="assignment-group"
  @dragenter.native.prevent="event => handleDrag(event, true)"
  @dragleave.native.prevent="event => handleDrag(event, false)"
  @dragover.native.prevent
  @drop.native.prevent="draggedOver = false; handleDrop()"
>
  <template v-for="(assign, index) in assignments">
    <InsertAssignment :active="assignmentDraggedOver[index]" :key="`insert-${index}`"
      @dragenter.native.prevent.stop="assignmentDraggedOver[index] = true"
      @dragleave.native.prevent.stop="assignmentDraggedOver[index] = false"
      @dragover.native.prevent.stop
      @drop.native.prevent.stop="assignmentDraggedOver[index] = false; handleDrop(index)"
    />

    <Assignment :eventId="eventId" :index="index" :key="`assign-${index}`"
      @dragenter.native.prevent.stop="assignmentDraggedOver[index] = true"
      @dragleave.native.prevent.stop="assignmentDraggedOver[index] = false"
      @dragover.native.prevent.stop
      @drop.native.prevent.stop="assignmentDraggedOver[index] = false; handleDrop(index)"
    />
  </template>
  <InsertAssignment :active="draggedOver" />
</v-chip-group>
</template>
<script>
import Assignment from './Assignment'
import InsertAssignment from './InsertAssignment'
import { eventProps, dragAssignProps } from '../store/utils'

export default {
  data: () => ({
    draggedOver: false,
    assignmentDraggedOver: [],
  }),

  props: {
    eventId: {},
  },

  computed: {
    ...eventProps(['assignments']),
    ...dragAssignProps(),

    insertClass() {
      return this.draggedOver ? "insert-assign-active" : "insert-assign-inactive"
    },
  },

  watch: {
    assignments: {
      handler(v) {
        this.assignmentDraggedOver = [...v.map(() => false)]
      },
      immediate: true,
    },
  },

  mounted() {
  },

  methods: {
    handleDrag(event, enter) {
      let target = event.target
      if (typeof(target.matches) != "function") {
        target = target.parentNode
      }

      if (!target.matches(".assignment,.assignment *,.insert-assign,.insert-assign *")) {
        this.draggedOver = enter
      }
    },

    handleDrop(index) {
      if (index == null) {
        index = this.assignments.length
      }
      if (this.draggedAssign.sourceId !== undefined && this.draggedAssign.sourceIndex !== undefined) {
        this.$store.commit('events/moveAssignment', {
          from: {id: this.draggedAssign.sourceId, index: this.draggedAssign.sourceIndex},
          to: {id: this.eventId, index}
        })
      } else {
        this.$store.commit('events/addAssignment', {id: this.eventId, assignId: this.draggedAssign.assignId, index})
      }
    },
  },

  components: {
    Assignment,
    InsertAssignment,
  },
}
</script>
<style>
.v-chip-group.assignment-group .v-chip {
  margin-top: 0;
  margin-bottom: 0;
  margin-right: 4px;
}
</style>
