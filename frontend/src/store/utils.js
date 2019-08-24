export function assignProps(props) {
  const p = {}
  props.forEach(prop => {
    p[prop] = {
      get() {
        if (this.assignId === undefined) return undefined

        const assign = this.$store.state.assigns.assigns[this.assignId]
        return assign !== undefined ? assign[prop] : undefined
      },
      set(v) {
        const assign = this.assignId === undefined ? {} : this.$store.state.assigns.assigns[this.assignId]
        this.$store.commit('assigns/set', {...assign, [prop]: v})
      },
    }
  });
  return p
}

export function player() {
  return {
    player() {
      const assign = this.$store.state.assigns.assigns[this.assignId]
      return this.$store.state.assigns.assigns[assign.playerId]
    }
  }
}

export function eventProps(props) {
  const p = {}
  props.forEach(prop => {
    p[prop] = {
      get() {
        if (this.eventId === undefined) return undefined

        const event = this.$store.state.events.events[this.eventId]
        return event !== undefined ? event[prop] : undefined
      },
      set(v) {
        const event = this.eventId === undefined ? {} : this.$store.state.events.events[this.eventId]
        this.$store.commit('events/set', {...event, [prop]: v})
      },
    }
  });
  return p
}

export function dragAssignProps() {
  return {
    draggedAssign: {
      get() {
        return this.$store.state.dragassign.info
      },

      set(info) {
        if (!info || info == {}) {
          this.$store.commit('dragassign/stop')
        } else {
          this.$store.commit('dragassign/start', info)
        }
      },
    },
  }
}
