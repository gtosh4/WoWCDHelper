<template>
<v-col cols="auto" id="cd-palette">
  <PlayerAssigneeGroup />
  <QuickAddHealer />
  <AddPlayers />
</v-col>
</template>
<script>
import PlayerAssigneeGroup from './PlayerAssigneeGroup'
import QuickAddHealer from './QuickAddHealer'
import AddPlayers from './AddPlayers'

import { spec } from './wow_info'

export default {
  data: () => ({
  }),

  props: {
  },

  mounted() {
  },

  computed: {
  },

  methods: {
    addPlayer(className, specName, name, id) {
      const player = {className, specName, name, id}
      this.$store.commit('assigns/set', player)
      if (specName !== undefined) {
        spec(className, specName).spells.forEach(spell => {
          this.$store.commit('assigns/set', {spell, playerId: player.id, id: `${player.id}.${spell.id}`})
        })
      }
    },
  },

  components: {
    PlayerAssigneeGroup,
    QuickAddHealer,
    AddPlayers,
  },
};
</script>
<style>
#cd-palette {
  min-width: 20vw;
  max-width: 40vw;
}
</style>
