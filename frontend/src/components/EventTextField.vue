<template>
  <v-text-field
    v-model="internalValue"
    single-line
    hide-details
    :placeholder="placeholder"
    @keydown.esc.stop="cancel"
    @keydown.enter="save"
    @blur="save"
  />
</template>
<script>
export default {

  components: {
  },

  props: {
    value: {
      type: String,
      default: undefined,
    },

    placeholder: {
      type: String,
      default: "",
    },
  },

  data: () => ({
    internalValue: ""
  }),

  watch: {
    value: {
      handler() { this.internalValue = this.value },
      immediate: true,
    },
  },

  methods: {
    cancel() {
      this.internalValue = this.value
    },

    save() {
      if (this.value == "" && this.internalValue == "") return
      if (this.value == this.internalValue) return
      
      this.$emit('input', this.internalValue)
    },
  },
}
</script>
<style scoped>
.v-input {
  margin-top: 0px;
  padding-top: 0px;
  width: 100%;
}
</style>
