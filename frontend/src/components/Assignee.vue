<template>
<v-layout class="assignee">
  <v-chip v-if="!spell" label :color="classColour" class="player" >
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
  </v-chip>
  <v-chip v-else label :color="classColour" draggable class="spell">
    <v-icon class="handle">mdi-drag</v-icon>
    <span class="mr-1">{{ assignmentCount }}</span>
    <Spell :spell="spell" />
  </v-chip>
  <v-dialog v-model="showSettings" width="650px">
    <template #activator="{ on: dialog }">
      <v-tooltip top>
        <template #activator="{ on: tooltip }">
          <v-btn tile x-small icon tabindex="-1" v-on="{...tooltip, ...dialog}">
            <v-icon>mdi-settings</v-icon>
          </v-btn>
        </template>
        <span>Settings</span>
      </v-tooltip>
    </template>

    <SpellSettings v-if="spell" :assignId="assignId" @close="showSettings = false" />
    <PlayerSettings v-else :assignId="assignId" @close="showSettings = false" />
  </v-dialog>
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
    WowIcon,
    Spell,
    SpellSettings,
    PlayerSettings,
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

