<template>
<v-avatar tile class="mx-1" size="18">
  <img v-if="src" :src="src" class="wowicon" draggable="false" />
  <v-icon v-else>mdi-alert-circle</v-icon>
</v-avatar>
</template>
<script>
import {classIcon, specIcon, abilityIcon, spec} from './wow_info'

export default {
  data: () => ({
  }),

  props: {
    className: {
      type: String,
    },
    specName: {
      type: String,
    },
    abilityName: {
      type: String,
    },
    size: {
      type: Number,
      default: 18,
    },
  },

  computed: {
    src() {
      if (this.abilityName) {
        return abilityIcon(this.abilityName)
      } else if (this.specName && this.className) {
        return specIcon(spec(this.className, this.specName))
      } else if (this.className) {
        return classIcon(this.className)
      } else {
        return undefined
      }
    }
  },

  methods: {
  },

  mounted() {
  },

  components: {
  },
};
</script>
<style>
</style>
