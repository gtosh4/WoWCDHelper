import {spells, spec} from '../components/wow_info'
import { DragAssignState } from './modules/dragassign'

interface HasAssign {
  assignId?: number
  $store: any
}

interface Prop<T> {
  get(): T
  set(v: T): void // eslint-disable-line
}

export function assignProps(props: string[]) {
  const p = {} as {[prop: string]: Prop<any>}

  props.forEach(prop => {
    p[prop] = {
      get(this: HasAssign) {
        if (this.assignId === undefined) return undefined

        const assign = this.$store.state.assigns.assigns[this.assignId]
        if (assign === undefined) return undefined

        if (assign[prop] !== undefined) {
          return assign[prop]
        }
        
        if (assign.id != assign.playerId) {
          // Fallback to player assign (eg for specName, className)
          const pAssign = this.$store.state.assigns.assigns[assign.playerId]
          if (pAssign !== undefined && pAssign[prop] !== undefined) {
            return pAssign[prop]
          }
        }
        return undefined
      },

      set(this: HasAssign, v) {
        const assign = this.assignId === undefined ? {} : this.$store.state.assigns.assigns[this.assignId]
        this.$store.commit('assigns/set', {...assign, [prop]: v})
      },
    }
  })
  return p
}

export function spell() {
  return {
    spell: {
      get(this: HasAssign) {
        if (this.assignId === undefined) return undefined
        const assign = this.$store.state.assigns.assigns[this.assignId]
        if (!assign.spell || !assign.spell.id) return undefined

        const s = {
          ...spells[assign.spell.id],
          ...assign.spell,
        }
        if (!s.cfg) {
          s.cfg = {}
        }
        
        const pAssign = this.$store.state.assigns.assigns[assign.playerId]
        const specSpell = spec(pAssign.className, pAssign.specName).spells.find(ss => ss.id == s.id)
        if (specSpell) {
          s.options = specSpell.options
          if (typeof(specSpell.configure) == 'function') {
            specSpell.configure(s)
          }
        }

        return s
      },
      set(this: HasAssign, v: {id: number, cfg: any}) {
        if (this.assignId === undefined) return
        const s = {
          id: v.id,
          cfg: v.cfg,
        }
        const assign = this.$store.state.assigns.assigns[this.assignId]
        this.$store.commit('assigns/set', {...assign, spell: s})
      },
    }
  }
}

export function player() {
  return {
    player(this: HasAssign) {
      if (this.assignId === undefined) return undefined
      const assign = this.$store.state.assigns.assigns[this.assignId]
      return this.$store.state.assigns.assigns[assign.playerId]
    }
  }
}

interface HasEvent {
  eventId?: number
  $store: any
}

export function eventProps(props: string[]) {
  const p = {} as {[prop: string]: Prop<any>}

  props.forEach(prop => {
    p[prop] = {
      get(this: HasEvent) {
        if (this.eventId === undefined) return undefined

        const event = this.$store.state.events.events[this.eventId]
        return event !== undefined ? event[prop] : undefined
      },
      set(this: HasEvent, v) {
        const event = this.eventId === undefined ? {} : this.$store.state.events.events[this.eventId]
        this.$store.commit('events/set', {...event, [prop]: v})
      },
    }
  })
  return p
}

export function dragAssignProps() {
  return {
    draggedAssign: {
      get(this: {$store: any}) {
        const state = this.$store.state.dragassign as DragAssignState
        return state.draggedAssign
      },

      set(this: {$store: any}, assignId?: number) {
        if (assignId) {
          this.$store.commit('dragassign/start', assignId)
        } else {
          this.$store.commit('dragassign/stop')
        }
      },
    },
  }
}
