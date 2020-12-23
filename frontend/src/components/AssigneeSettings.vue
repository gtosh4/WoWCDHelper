<template>
  <v-card>
    <v-card-title>
      <WowIcon :class-name="className" />
      <WowIcon v-if="specName" :class-name="className" :spec-name="specName" />
      {{ player.name }}{{ spell ? `: ${spell.name}` : '' }}
    </v-card-title>
    <v-list v-if="spell">
      <v-list-item v-for="(opt, i) in spell.options" :key="i">
        <v-checkbox
          v-if="opt.type == 'bool'"
          :input-value="getValue(opt)"
          @change="v => setOpt(opt, v)"
        />
        <span>{{ opt.text }}</span>
      </v-list-item>
    </v-list>

    <v-card-actions>
      <v-btn @click="$emit('close')">
        Close
      </v-btn>
    </v-card-actions>
  </v-card>
</template>
<script lang="ts">
import { assignProps, player, spell } from '../store/utils'
import WowIcon from './WowIcon.vue'
import Vue from 'vue'
import { SpellOptions } from './wow_info'

export default Vue.extend({
  components: {
    WowIcon,
  },
  
  props: {
    assignId: {
      type: String,
      required: true,
    },
  },

  computed: {
    ...assignProps(['name', 'className', 'specName', 'playerId']),
    ...player(),
    ...spell(),
  },

  methods: {
    getValue(opt: SpellOptions<any>): any {
      const cfg = this.spell.cfg
      return cfg?.[opt.prop] || opt.default
    },

    setOpt(opt: SpellOptions<any>, v: any) {
      Vue.set(this.spell.cfg, opt.prop, v)
    },
  },
})
</script>
<style>
</style>
