<template>
  <v-col id="cd-palette" cols="auto">
    <PlayerAssigneeGroup />
    <QuickAddHealer @addPlayer="addPlayer" />
    <AddPlayers @addPlayer="addPlayer" />
  </v-col>
</template>
<script lang="ts">
import Vue from 'vue'
import PlayerAssigneeGroup from './PlayerAssigneeGroup.vue'
import QuickAddHealer from './QuickAddHealer.vue'
import AddPlayers from './AddPlayers.vue'
import { Assign, PlayerAssign, SpellAssign } from '../store/modules/assigns'

import { spec } from './wow_info'

export default Vue.extend({
  components: {
    PlayerAssigneeGroup,
    QuickAddHealer,
    AddPlayers,
  },

  methods: {
    addPlayer(className: string, specName?: string, name?: string, id?: string) {
      const player = {className, specName, name, id, type: "player"} as PlayerAssign
      
      const assigns = [player] as Assign[]
      if (specName !== undefined) {
        spec(className, specName).spells.forEach(spell => {
          assigns.push({spell, type: "spell"} as SpellAssign)
        })
      }
      this.$store.commit('assigns/add', assigns)
    },
  },
})
</script>

<style>
#cd-palette {
  min-width: 120px;
  max-width: 40vw;
}
</style>
