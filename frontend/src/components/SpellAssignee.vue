<template>
<v-chip label :color="classColour" class="spell assignee">
  <v-icon class="handle">mdi-drag</v-icon>
  <span class="mr-1">{{ assignmentCount }}</span>
  <Spell :spell="spell" />
  <v-btn-toggle group>
    <v-tooltip top>
      <template #activator="{ on }">
        <v-btn tile small icon tabindex="-1" @click="$emit('config')" v-on="on">
          <v-icon>mdi-settings</v-icon>
        </v-btn>
      </template>
      <span>Settings</span>
    </v-tooltip>
    <v-tooltip top>
      <template #activator="{ on }">
        <v-btn tile small icon tabindex="-1" @click="clearAssign" v-on="on">
          <v-icon>mdi-backspace</v-icon>
        </v-btn>
      </template>
      <span>Clear assignments</span>
    </v-tooltip>
    <v-tooltip top>
      <template #activator="{ on }">
        <v-btn tile small icon tabindex="-1" @click="deleteAssign" v-on="on">
          <v-icon>mdi-delete</v-icon>
        </v-btn>
      </template>
      <span>Delete</span>
    </v-tooltip>

  </v-btn-toggle>
</v-chip>

</template>
<script>
import WowIcon from './WowIcon'
import Spell from './Spell'
import SpellSettings from './SpellSettings'
import PlayerSettings from './PlayerSettings'

import Color from 'color'
import {classes, classIcon, specIcon, spec} from './wow_info'
import {assignProps, dragAssignProps, player, spell} from '../store/utils'

export default {
  data: () => ({
    showSettings: false,
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
<style>
.spell.assignee {
  width: 100%;
}

.spell.assignee .v-chip__content {
  width: 100%;
}

.spell.assignee .spelldetails {
  width: 100%;
}

.v-chip .handle {
  cursor: grab;
}
</style>

