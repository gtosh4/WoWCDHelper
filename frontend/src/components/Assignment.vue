<template>
<v-chip
  label
  :close="showHover"
  @mouseenter="hover = true"
  @mouseleave="hover = false"
  :class="chipClass"
  :color="classColour"
  @click:close="remove"
>
  <v-icon v-if="moveable" class="handle" :style="{display: showHover ? '' : 'none'}">mdi-drag</v-icon>
  <span v-if="spell">{{ player.name }} <Spell :spell="spell" :showname="false" /></span>
  <span v-else>{{ name }}</span>
</v-chip>
</template>
<script>
import Spell from './Spell'

import Color from 'color'
import {classes, classIcon, specIcon, spec} from './wow_info'
import {assignProps, dragAssignProps, spell} from '../store/utils'

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
    moveable: {
      default: true,
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
      e.dataTransfer.setData("assignId", this.assignId)
      this.draggedAssign = {
        assignId: this.assignId,
        sourceId: this.eventId,
        sourceIndex: this.index,
      }
    }
    chip.ondragend = () => {
      chip.setAttribute('draggable', 'false')
      this.draggedAssign = null
    }
  },

  computed: {
    ...dragAssignProps(),

    assignId() {
      return this.$store.state.events.events[this.eventId].assignments[this.index]
    },
    ...assignProps(['name', 'className', 'specName', 'playerId']),
    ...spell(),

    player() {
      return this.$store.state.assigns.assigns[this.playerId]
    },

    classColour() {
      const c = Color(classes[this.player.className].colour)
      return c.mix(Color('rgb(66, 66, 66)'), 0.4).string()
    },

    showHover() {
      return this.moveable && this.hover && !this.draggedOver
    },

    chipClass() {
      return this.spell != null ? "spell" : "player"
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
      if (!this.draggedAssign) return
      event.stopPropagation()
      

      if (this.draggedAssign.sourceId !== undefined) {
        const isSameEvent = this.draggedAssign.sourceId == this.eventId
        const isReorder = this.draggedAssign.sourceIndex == this.index || +this.draggedAssign.sourceIndex+1 == this.index
        if (isSameEvent && isReorder) return

        event.dataTransfer.dropEffect = "move"
      } else {
        event.dataTransfer.dropEffect = "link"
      }
      this.draggedOver = true
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
          to: {id: this.eventId, index: this.index}
        })
      } else {
        this.$store.commit('events/addAssignment', {id: this.eventId, assignId: this.draggedAssign.assignId, index: this.index})
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
  padding: 0 0 0 0;
  margin: 0 0 0 0;
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
