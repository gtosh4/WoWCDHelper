<template>
<v-chip label :color="classColour" class="player assignee">
  <v-icon class="handle">mdi-drag</v-icon>
  <WowIcon :className="className" />
  <WowIcon v-if="specName" :className="className" :specName="specName" />
  <WowIcon v-else />
  <v-text-field
    v-model="nameTmp"
    solo
    flat
    placeholder="Name"
    hide-details
    background-color="transparent"
    width="100%"
    @keydown.esc.stop="nameTmp = name"
    @keydown.enter="name = nameTmp"
    @blur="name = nameTmp"
  />
  <v-btn-toggle group>
    <v-tooltip top>
      <template #activator="{ on }">
        <v-btn tile small icon tabindex="-1" @click="$emit('config')" v-on="on">
          <v-icon>mdi-cog-outline</v-icon>
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

import Color from 'color'

import {classes} from './wow_info'
import {assignProps, dragAssignProps, player} from '../store/utils'

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
    ...assignProps(['name', 'className', 'specName']),
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
    deleteAssign() {
      this.$store.commit('deleteAssign', this.assignId)
    },

    clearAssign() {
      this.$store.commit('clearAssign', this.assignId)
    },
  },

  components: {
    WowIcon,
  }
}
</script>
<style>
.player.assignee {
  width: 100%;
}
.player.assignee .v-chip__content {
  width: 100%;
}

.player.assignee .handle {
  cursor: grab;
}
</style>
