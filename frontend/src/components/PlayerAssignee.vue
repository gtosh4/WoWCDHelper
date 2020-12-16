<template>
<v-card tile :color="classColour" class="assignee">
  <v-card-title ref="player" class="player">
    <v-icon class="handle">drag_indicator</v-icon>
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
      @keydown.esc.stop="nameTmp = name"
      @keydown.enter="name = nameTmp"
      @blur="name = nameTmp"
      class="name-field"
    />
    <v-menu>
      <template #activator="{ on, attrs }">
        <v-btn
          v-bind="attrs"
          v-on="on"
          icon
          tile
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
  </v-card-title>

  <v-card-text class="playerspells">
    <v-list>
      <SpellAssignee v-for="(spell, i) in playerSpells" :key="i" :assignId="spell.id" />
    </v-list>
  </v-card-text>
</v-card>
</template>
<script>
import SpellAssignee from './SpellAssignee'
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
    const chip = this.$refs.player
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

    playerSpells() {
      return Object.values(this.$store.state.assigns.assigns).filter(a => a.playerId == this.player.id && a.spell != null).sort((a, b) => {
        return a.spell.id < b.spell.id
      })
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
    SpellAssignee,
  }
}
</script>

<style lang="scss">
.assignee {
  width: 100%;

  .player {
    padding: 0;

    .name-field .v-input__control {
      min-height: 18px;

      input {
        padding: 0;
      }
    }
  }

  .playerspells {
    padding-bottom: 0;

    &>.v-list {
      margin: 0;
      background-color: transparent;
      padding: 0;
    }
  }

  .handle {
    cursor: grab;
  }
}
</style>
