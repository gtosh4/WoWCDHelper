export function assignProps(props) {
  const p = {}
  props.forEach(prop => {
    p[prop] = {
      get() {
        return this.$store.state.assigns.assigns[this.assignId][prop]
      },
      set(v) {
        this.$store.commit('assigns/set', {...this.$store.state.assigns.assigns[this.assignId], [prop]: v})
      },
    }
  });
  return p
}

export function eventProps(props) {
  const p = {}
  props.forEach(prop => {
    p[prop] = {
      get() {
        return this.$store.state.events.events[this.eventId][prop]
      },
      set(v) {
        this.$store.commit('events/set', {...this.$store.state.events.events[this.eventId], [prop]: v})
      },
    }
  });
  return p
}
