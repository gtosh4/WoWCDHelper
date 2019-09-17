<template>
<v-card>
  <v-card-title>
    {{ player.name }}: {{ spell.name }}
  </v-card-title>
  <v-list>
    <v-list-item v-for="(opt, i) in spell.options" :key="i">
      <v-checkbox v-if="opt.type == Boolean"
        :input-value="getValue(opt)"
        @change="v => setOpt(opt, v)"
      />
      <span>{{opt.text}}</span>
    </v-list-item>
  </v-list>
</v-card>
</template>
<script>
import { assignProps, player, spell } from '../store/utils'

export default {
  data: () => ({
  }),

  props: {
    assignId: {
      required: true,
    },
  },

  computed: {
    ...assignProps(['name', 'className', 'specName', 'playerId']),
    ...spell(),
    ...player(),
  },

  methods: {
    getValue(opt) {
      const cfg = this.spell.cfg
      return cfg && cfg[opt.prop] ? cfg[opt.prop] : opt.default
    },

    setOpt(opt, v) {
      const s = {
        ...this.spell,
      }
      if (!s.cfg) {
        s.cfg = {}
      }
      s.cfg = {...s.cfg, [opt.prop]: v}
      
      this.spell = s
    },
  },

  mounted() {
  },

  components: {
  },
};
</script>
<style>
</style>
