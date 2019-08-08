<template>
<v-card>
  <v-tabs v-model="tab">
    <v-tab v-for="format in formats" :key="format.key" :href="`#export-${format.key}`">
      <span>{{ format.name }}</span>
    </v-tab>

    <v-tab-item v-for="format in formats" :key="format.key" :value="`export-${format.key}`">
      <v-textarea readonly :value="text[format.key]" />
    </v-tab-item>
  </v-tabs>

  <v-card-actions>
    <v-spacer />
    <v-btn @click="$emit('close')">Close</v-btn>
  </v-card-actions>
</v-card>
</template>
<script>
export default {
  data: () => ({
    formats: [
      {name: 'Angry Assignments', key: "aa"},
      {name: 'Exorsus Raid Tools', key: "ert"},
      {name: 'ERT Personal', key: "ertp"},
    ],
    tab: undefined,
  }),

  props: {
  },

  computed: {
    text() {
      return Object.fromEntries(this.formats.map(format => [format.key, `Tab: ${this.tab}`]))
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
