<template>
<tr>
  <td class="text-xs-right">
    <v-edit-dialog
      @update:return-value="updateTime"
    > {{ timeStr }}
      <template v-slot:input>
        <v-text-field
          :value="timeStr"
          single-line
        ></v-text-field>
      </template>
    </v-edit-dialog>
  </td>
  <td>{{ label }}</td>
  <td></td>
  <td></td>
</tr>
</template>
<script>
import moment from 'moment'

export default {
  data: () => ({
  }),

  props: {
    label: {
      type: String,
      required: true,
    },
    time: {
      type: Object,
      required: true,
    },
    assigns: {
      type: Array,
      default() { return [] },
    },
  },

  computed: {
    timeStr() {
      return `${this.time.minutes()}:${this.time.seconds()}`
    },

  },
  methods: {
    updateTime(v) {
      const [mins, secs] = v.split(":")
      this.$emit("update:time", moment.duration(+mins, 'minutes').add(+secs, 'seconds'))
    },
  },
};
</script>
