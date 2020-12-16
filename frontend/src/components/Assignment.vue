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
  <v-icon v-if="moveable" class="handle" :style="{display: showHover ? '' : 'none'}">drag_indicator</v-icon>
  <span class="assignment-content">{{ spell ? name : player.name }}<Spell v-if="spell" :spell="spell" :showname="false" /></span>
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
    const chip = this.$el
    const handle = chip.querySelector(".handle")

    if (handle != null) {
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
      return [this.spell != null ? "spell" : "player", "assignment"]
    }
  },

  methods: {
    classIcon,
    specIcon,
    spec,

    remove() {
      this.$store.commit('events/removeAssignment', {id: this.eventId, index: this.index})
    },
  },

  components: {
    Spell,
  },
};
</script>
<style>
.assignment .handle {
  cursor: grab;
}

/* .assignment .assignment-content {
  pointer-events: none;
} */
</style>
