<template>
<v-layout class="assignee">
  <v-chip v-if="!spell" label :color="classColour" class="player" >
    <v-icon class="handle">mdi-drag</v-icon>
    <WowIcon :className="className" />
    <WowIcon v-if="specName" :className="className" :specName="specName" />
    <WowIcon v-else />
    <v-text-field v-model="name" solo flat placeholder="Name" hide-details background-color="transparent" width="100%" />
  </v-chip>
  <v-chip v-else label :color="classColour" draggable class="spell">
    <v-icon class="handle">mdi-drag</v-icon>
    <span class="mr-1">{{ assignmentCount }}</span>
    <Spell :spell="spell" />
  </v-chip>
    <v-tooltip top>
      <template #activator="{ on }">
        <v-btn tile x-small icon tabindex="-1" @click="clearAssign" v-on="on">
          <v-icon>mdi-backspace</v-icon>
        </v-btn>
      </template>
      <span>Clear assignments</span>
    </v-tooltip>
    <v-tooltip top>
      <template #activator="{ on }">
        <v-btn tile x-small icon tabindex="-1" @click="deleteAssign" v-on="on">
          <v-icon>mdi-delete</v-icon>
        </v-btn>
      </template>
      <span>Delete</span>
    </v-tooltip>
</v-layout>
</template>
<script>
import WowIcon from './WowIcon'
import Spell from './Spell'

import {classes, classIcon, specIcon, spec} from './wow_info'
import {assignProps} from '../store/utils'

export default {
  data: () => ({
  }),

  props: {
    assignId: {
      required: true,
    },
  },

  mounted() {
    const chip = this.$el.querySelector(".assignee > .v-chip")
    const handle = chip.querySelector(".handle")

    handle.onmousedown = () => {
      chip.setAttribute('draggable', 'true')
    }
    handle.onmouseup = () => {
      chip.setAttribute('draggable', 'false')
    }

    chip.ondragstart = (e) => {
      e.dataTransfer.setData("assignId", this.assignId)
    }
    chip.ondragend = () => {
      chip.setAttribute('draggable', 'false')
    }
  },

  computed: {
    ...assignProps(['name', 'className', 'specName', 'spell', 'playerId']),

    player() {
      return this.$store.state.assigns.assigns[this.playerId]
    },

    classColour() {
      const c = classes[this.player.className].colour
      return `rgba(${c.r}, ${c.g}, ${c.b}, 0.5)`
    },

    assignmentCount() {
      return Object.values(this.$store.state.events.events).map(e => e.assignments.filter(a => a == this.assignId)).flat().length
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
      this.$store.commit('events/clearAssignee', this.assignId)
    },
  },

  components: {
    WowIcon,
    Spell,
  },
}
</script>
<style>
.assignee .v-chip {
  width: 100%;
}
.assignee .handle {
  cursor: grab;
}
</style>

