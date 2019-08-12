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
  data: () => ({
    internalValue: ""
  }),

  props: {
    value: String,
    placeholder: String,
  },

  watch: {
    value: {
      handler() { this.internalValue = this.value },
      immediate: true,
    },
  },

  mounted() {
  },

  computed: {
  },

  methods: {
    cancel() {
      this.internalValue = this.value
    },

    save() {
      if (this.value == "" && this.internalValue == "") return
      
      this.$emit('input', this.internalValue)
    },
  },

  components: {
  },
};
</script>
<style scoped>
.v-input {
  margin-top: 0px;
  padding-top: 0px;
  min-width: 50px;
}
input {
  text-align: end;
}
</style>
