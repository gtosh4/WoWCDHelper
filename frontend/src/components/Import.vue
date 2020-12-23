<template>
  <v-card>
    <v-card-text class="import-content">
      <v-textarea
        v-model="content"
        auto-grow
        full-width
        outlined
      />
    </v-card-text>

    <v-card-actions>
      <v-btn text @click="doImport">
        Import
      </v-btn>
      <v-spacer />
      <v-btn text @click="$emit('close')">
        Close
      </v-btn>
    </v-card-actions>
  </v-card>
</template>
<script>
import moment from 'moment'

export default {
  data: () => ({
    content: "",
  }),

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
      if (t.name) {
        this.$store.commit("setName", t.name)
      }
      this.content = ""
      this.$emit("close")
    },
  },
}
</script>
<style>
.import-content {
  overflow-y: auto;
  max-height: 70vh;
}
</style>
