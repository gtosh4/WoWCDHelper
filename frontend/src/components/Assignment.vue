<template>
<v-flex shrink pa-0 ma-0
  @drop="handleDrop"
  @dragover.prevent="handleDragOver"
  @dragleave.prevent="handleDragLeave"
  class="assignment"
>
  <InsertAssignment :eventId="eventId" :draggedAssign="draggedAssign" />
  <v-chip
    v-if="!spell"
    label
    :close="showHover"
    @mouseenter="hover = true"
    @mouseleave="hover = false"
    class="player"
    :color="classColour"
    @click:close="remove"
  >
    <v-icon class="handle" :style="{display: showHover ? '' : 'none'}">mdi-drag</v-icon>
    <span>{{ name }}</span>
  </v-chip>
  <v-chip
    v-else
    label
    :close="showHover"
    @mouseenter="hover = true"
    @mouseleave="hover = false"
    class="spell"
    :color="classColour"
    @click:close="remove"
  >
    <v-icon class="handle" :style="{display: showHover ? '' : 'none'}">mdi-drag</v-icon>
    <span>{{player.name}} <Spell :spell="spell" :showname="false" /></span>
  </v-chip>
</v-flex>
</template>
<script>
import Spell from './Spell'
import InsertAssignment from './InsertAssignment'

import {classes, classIcon, specIcon, spec} from './wow_info'
import {assignProps} from '../store/utils'

export default {
  data: () => ({
    hover: false,
    draggedAssign: null,
  }),

  props: {
    eventId: {
      required: true,
    },
    index: {
      required: true,
    },
  },

  mounted() {
    const chip = this.$el.querySelector(".assignment > .v-chip")
    const handle = chip.querySelector(".handle")

    handle.onmousedown = () => {
      chip.setAttribute('draggable', 'true')
    }
    handle.onmouseup = () => {
      chip.setAttribute('draggable', 'false')
    }

    chip.ondragstart = (e) => {
      e.dataTransfer.setData("eventId", this.eventId)
      e.dataTransfer.setData("assignId", this.assignId)
      e.dataTransfer.setData("assignIndex", this.index)
    }
    chip.ondragend = () => {
      chip.setAttribute('draggable', 'false')
    }
  },

  computed: {
    assignId() {
      return this.$store.state.events.events[this.eventId].assignments[this.index]
    },
    ...assignProps(['name', 'className', 'specName', 'spell', 'playerId']),

    player() {
      return this.$store.state.assigns.assigns[this.playerId]
    },

    classColour() {
      const c = classes[this.player.className].colour
      return `rgba(${c.r}, ${c.g}, ${c.b}, 0.5)`
    },

    showHover() {
      return this.hover && !this.draggedAssign
    },
  },

  methods: {
    classIcon,
    specIcon,
    spec,

    remove() {
      this.$store.commit('events/removeAssignment', {id: this.eventId, index: this.index})
    },

    handleDragOver(event) {
      const assignId = event.dataTransfer.getData("assignId"),
            assignIndex = event.dataTransfer.getData("assignIndex"),
            sourceId = event.dataTransfer.getData("eventId")
      event.stopPropagation();
      
      if (sourceId == this.eventId && (assignIndex == this.index || +assignIndex+1 == this.index)) return
      if (assignId) {
        if (sourceId) {
          event.dataTransfer.dropEffect = "move"
        } else {
          event.dataTransfer.dropEffect = "link"
        }
        this.draggedAssign = {assignId, sourceId, assignIndex}
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
          this.$store.commit('events/moveAssignment', {from: {id: sourceId, index: sourceIdx}, to: {id: this.eventId, index: this.index}})
        } else {
          this.$store.commit('events/addAssignment', {id: this.eventId, assignId: assignId, index: this.index})
        }
      }
      this.draggedAssign = null
    },
  },

  components: {
    Spell,
    InsertAssignment,
  },
};
</script>
<style>
.assignment {
  display: inline-flex;
}
.assignment .v-chip {
  margin-left: 4px;
  margin-right: 0px;
}

.assignment .handle {
  cursor: grab;
}
</style>
