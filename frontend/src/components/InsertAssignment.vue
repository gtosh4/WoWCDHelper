<template>
  <v-chip label color="grey lighten-4" :class="chipClass" />
</template>
<script>
export default {
  props: {
    active: {
      type: Boolean,
    },
  },

  computed: {
    chipClass() {
      return ["insert-assign", `insert-assign-${this.active ? 'active': 'inactive'}`]
    },
  },
}
</script>
<style>
.v-chip.insert-assign {
  width: 25px;
  opacity: 0.8;
  padding: 0 0 0 0;
  margin: 0 0 0 0;
}
.v-chip.insert-assign.insert-assign-inactive {
  display: none;
}
.v-chip.insert-assign.insert-assign-active {
  display: inherit;
}
</style>
