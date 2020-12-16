<template>
<v-list-item :color="classColour" class="spell">
  <v-icon class="handle">drag_indicator</v-icon>
  <span>{{ assignmentCount }}</span>
  <Spell :spell="spell" />

  <v-menu>
    <template #activator="{ on, attrs }">
      <v-btn
        v-bind="attrs"
        v-on="on"
        icon
        tile
        small
      >
        <v-icon>more_vert</v-icon>
      </v-btn>
    </template>

    <v-list dense>
      <v-list-item @click="$emit('config')">
        <v-list-item-icon><v-icon>settings</v-icon></v-list-item-icon>
        Settings
      </v-list-item>

      <v-list-item @click="clearAssign">
        <v-list-item-icon><v-icon>backspace</v-icon></v-list-item-icon>
        Clear assignments
      </v-list-item>

      <v-list-item @click="deleteAssign">
        <v-list-item-icon><v-icon>delete</v-icon></v-list-item-icon>
        Delete
      </v-list-item>
    </v-list>
  </v-menu>
</v-list-item>

</template>
<script>
import Spell from './Spell'

import Color from 'color'
import {classes, classIcon, specIcon, spec} from './wow_info'
import {assignProps, dragAssignProps, player, spell} from '../store/utils'

export default {
  data: () => ({
    nameTmp: "",
  }),

  props: {
    assignId: {
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
      e.dataTransfer.setData("assignId", this.assignId)
      this.draggedAssign = {
        assignId: this.assignId,
      }
    }
    chip.ondragend = () => {
      chip.setAttribute('draggable', 'false')
      this.draggedAssign = null
    }
  },

  computed: {
    ...dragAssignProps(),
    ...assignProps(['name', 'className', 'specName', 'playerId']),
    ...spell(),
    ...player(),

    classColour() {
      const c = Color(classes[this.player.className].colour)
      return c.mix(Color('rgb(66, 66, 66)'), 0.4).string()
    },

    assignmentCount() {
      return Object.values(this.$store.state.events.events).map(e => e.assignments.filter(a => a == this.assignId)).flat().length
    },
  },

  watch: {
    name: {
      handler() { this.nameTmp = this.name},
      immediate: true,
    },
  },

  methods: {
    classIcon,
    specIcon,
    spec,

    deleteAssign() {
      this.$store.commit('deleteAssign', this.assignId)
    },

    clearAssign() {
      this.$store.commit('clearAssign', this.assignId)
    },
  },

  components: {
    Spell,
  },
}
</script>

<style lang="scss">
.spell.v-list-item {
  border-top: 1px solid black;
  min-height: 24px;

  span {
    margin-right: 4px;
  }

  .spelldetails {
    width: 100%;
  }
}
</style>

