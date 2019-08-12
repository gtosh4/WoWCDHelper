<template>
<v-layout>
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
  <v-icon @click="deleteAssign">mdi-delete</v-icon>
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
    const chip = this.$el.querySelector(".assignee")
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
      this.$store.commit('assigns/delete', this.assignId)
    },
  },

  components: {
    WowIcon,
    Spell,
  },
}
</script>
<style>
.v-chip.assignee {
  width: 100%;
}
</style>

