<template>
<v-btn-toggle group>
  <v-tooltip top>
    <template #activator="{ on }">
      <v-btn tile small icon tabindex="-1" @click="$emit('config')" v-on="on">
        <v-icon>mdi-settings</v-icon>
      </v-btn>
    </template>
    <span>Configure</span>
  </v-tooltip>

  <v-tooltip top>
    <template #activator="{ on }">
      <v-btn tile small icon tabindex="-1" @click="$emit('clone')" v-on="on">
        <v-icon>mdi-content-copy</v-icon>
      </v-btn>
    </template>
    <span>Copy</span>
  </v-tooltip>

  <v-tooltip top>
    <template #activator="{ on }">
      <v-btn tile small icon tabindex="-1" @click="$emit('clear')" v-on="on">
        <v-icon>mdi-backspace</v-icon>
      </v-btn>
    </template>
    <span>Clear assignments</span>
  </v-tooltip>

  <v-tooltip top>
    <template #activator="{ on }">
      <v-btn tile small icon tabindex="-1" @click="$emit('remove')" v-on="on">
        <v-icon>mdi-delete</v-icon>
      </v-btn>
    </template>
    <span>Delete row</span>
  </v-tooltip>
</v-btn-toggle>
</template>
<script>
export default {
  data: () => ({}),

  props: {
    
  }
}
</script>
