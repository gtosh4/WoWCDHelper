<template>
<span :data-tooltip-href="`//www.wowdb.com/spells/${spell.id}`">
  <v-avatar tile class="mr-1" size="18">
    <img v-if="spell.icon" :src="spell.icon" draggable="false" ondragstart="return false" />
  </v-avatar>
  <span v-if="showname">{{spell.name}}</span>
</span>
</template>
<script>
export default {
  data: () => ({
  }),

  props: {
    spell: {
      type: Object,
      required: true,
    },
    showname: {
      type: Boolean,
      default: true,
    }
  },

  computed: {
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
