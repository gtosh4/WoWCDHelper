<template>
<v-card>
  <v-card-text class="import-content">
    <v-textarea auto-grow full-width outlined @change="v => content = v" />
  </v-card-text>

  <v-card-actions>
    <v-btn text @click="doImport">Import</v-btn>
    <v-spacer />
    <v-btn text @click="$emit('close')">Close</v-btn>
  </v-card-actions>
</v-card>
</template>
<script>
import moment from 'moment'

export default {
  data: () => ({
    content: "",
  }),

  props: {
  },

  computed: {
  },

  methods: {
    doImport() {
      const t = JSON.parse(this.content) || {}
      const assigns = t.assigns || {}
      const events = t.events || {}
      Object.values(events).forEach((v) => {
        v.time = moment.duration(v.time)
      })

      this.$store.commit("assigns/import", assigns)
      this.$store.commit("events/import", events)
      this.$emit("close")
    },
  },

  mounted() {
  },

  components: {
  },
};
</script>
<style>
.import-content {
  overflow-y: auto;
  max-height: 70vh;
}
</style>
