<template>
<v-chip label disabled :color="color" />
</template>
<script>
export default {
  data: () => ({
  }),

  props: {
    value: {
    },
  },

  computed: {
    color() {
      return "grey lighten-4"
    },
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
