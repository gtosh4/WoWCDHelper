<template>
<div
  @drop="handleDrop"
  @dragover.prevent="handleDragOver"
  @dragleave.prevent="handleDragLeave"
>
<v-chip v-if="draggedOver" label disabled class="grey lighten-4" />
<v-chip
  v-if="!spell"
  label
  :close="hover && !draggedOver"
  @mouseenter="hover = true"
  @mouseleave="hover = false"
  class="assignment player"
  :color="classColour"
  @click:close="remove"
>
  <v-icon class="handle">mdi-drag</v-icon>
  <span>{{ name }}</span>
</v-chip>
<v-chip
  v-else
  label
  :close="hover && !draggedOver"
  @mouseenter="hover = true"
  @mouseleave="hover = false"
  class="assignment spell"
  :color="classColour"
  @click:close="remove"
>
  <v-icon class="handle" :style="{display: hover ? 'flex' : 'none'}">mdi-drag</v-icon>
  <span>{{player.name}} <Spell :spell="spell" :showname="false" class="ml-1" /></span>
</v-chip>
</div>
</template>
<script>
import Spell from './Spell'

import {classes, classIcon, specIcon, spec} from './wow_info'
import {assignProps} from '../store/utils'

export default {
  data: () => ({
    hover: false,
    draggedOver: false,
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
    const chip = this.$el
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
    }
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
        this.draggedOver = true
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
          this.$store.commit('events/moveAssignment', {from: {id: sourceId, index: sourceIdx}, to: {id: this.eventId, index: this.index}})
        } else {
          this.$store.commit('events/addAssignment', {id: this.eventId, assignId: assignId, index: this.index})
        }
      }
      this.draggedOver = false
    },
  },

  components: {
    Spell,
  },
};
</script>
<style>
.assignment {
  margin-left: 4px;
  margin-right: 4px;
}
</style>
