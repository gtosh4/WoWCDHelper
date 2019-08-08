<template>
<v-layout>
  <v-chip v-if="!ability" label :color="classColour" class="assign-chip" >
    <v-icon class="handle">mdi-drag</v-icon>
    <v-avatar tile class="mx-1" size="18"><img :src="classIcon(className)" /></v-avatar>
    <v-avatar tile class="mx-1" size="18"><img v-if="specName" :src="specIcon(spec(className, specName))" /></v-avatar>
    <v-text-field v-model="name" solo flat placeholder="Name" hide-details background-color="transparent" />
  </v-chip>
  <v-chip v-else label :color="classColour" draggable class="assign-chip">
    <v-icon class="handle">mdi-drag</v-icon>
    Empty
  </v-chip>
  <v-icon @click="deleteAssign">mdi-delete</v-icon>
</v-layout>
</template>
<script>
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
    const chip = this.$el.querySelector(".assign-chip")
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
    ...assignProps(['name', 'className', 'specName', 'ability']),

    classColour() {
      const c = classes[this.className].colour
      return `rgba(${c.r}, ${c.g}, ${c.b}, 0.5)`
    }
  },

  methods: {
    classIcon,
    specIcon,
    spec,

    deleteAssign() {
      this.$store.commit('assigns/delete', this.assignId)
    },
  },
}
</script>
<style>
img {
  pointer-events: none;
}
</style>

