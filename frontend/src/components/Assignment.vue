<template>
  <v-chip
    label
    :close="showHover"
    :class="chipClass"
    :color="classColour"
    @mouseenter="hover = true"
    @mouseleave="hover = false"
    @click:close="remove"
  >
    <v-icon
      v-if="moveable"
      class="handle"
      :style="{display: showHover ? '' : 'none'}"
    >
      drag_indicator
    </v-icon>
    <span class="assignment-content">{{ spell ? name : player.name }}<Spell v-if="spell" :spell="spell" :showname="false" /></span>
  </v-chip>
</template>
<script>
import Spell from './Spell'

import Color from 'color'
import {classes, classIcon, specIcon, spec} from './wow_info'
import {assignProps, dragAssignProps, player, spell} from '../store/utils'

export default {
  components: {
    Spell,
  },

  props: {
    eventId: {
      type: Number,
      required: true,
    },
    index: {
      type: Number,
      required: true,
    },
    moveable: {
      type: Boolean,
      default: true,
    },
  },

  data: () => ({
    hover: false,
    draggedOver: false,
  }),

  computed: {
    ...dragAssignProps(),

    assignId() {
      return this.$store.state.events.events[this.eventId].assignments[this.index]
    },
    ...assignProps(['name', 'className', 'specName', 'playerId']),
    ...player(),
    ...spell(),

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

  methods: {
    classIcon,
    specIcon,
    spec,

    remove() {
      this.$store.commit('events/removeAssignment', {id: this.eventId, index: this.index})
    },
  },
}
</script>
<style>
.assignment .handle {
  cursor: grab;
}

/* .assignment .assignment-content {
  pointer-events: none;
} */
</style>
