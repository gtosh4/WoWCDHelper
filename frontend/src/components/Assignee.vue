<template>
<v-layout>
  <v-chip v-if="!ability" label :color="classColour" class="assignee" >
    <v-icon class="handle">mdi-drag</v-icon>
    <WowIcon :className="className" />
    <WowIcon v-if="specName" :className="className" :specName="specName" />
    <v-text-field v-model="name" solo flat placeholder="Name" hide-details background-color="transparent" width="100%" />
  </v-chip>
  <v-chip v-else label :color="classColour" draggable class="assignee">
    <v-icon class="handle">mdi-drag</v-icon>
    Empty
  </v-chip>
  <v-icon @click="deleteAssign">mdi-delete</v-icon>
</v-layout>
</template>
<script>
import WowIcon from './WowIcon'

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

  components: {
    WowIcon,
  },
}
</script>
<style>
</style>

