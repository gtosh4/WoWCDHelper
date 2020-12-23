import { Module } from 'vuex'
import { State } from '..'
import Vue from 'vue'

export interface BaseAssign {
  id: string
  playerId: string
  type: string
}

export interface PlayerAssign extends BaseAssign {
  type: "player"
  name: string
  className: string
  specName?: string
}

export interface SpellAssign extends BaseAssign {
  type: "spell"
  spell: {
    id: number
    cfg: any
  }
}

export type Assign = PlayerAssign | SpellAssign

export interface AssignState {
  assigns: {[id: string]: Assign}
}

function nextAssignId(state: AssignState): string {
  let next = 1
  while (state.assigns[`${next}`] !== undefined) {
    next++
  }
  return `${next}`
}

export default {
  namespaced: true,

  state: {
    assigns: {},
  },

  getters: {
    players(state) {
      return Object.values(state.assigns).filter(a => a.type == "player")
    },
  
    spells(state) {
      return Object.values(state.assigns).filter(a => a.type == "spell")
    },
  
    nextAssignId,
  },

  mutations: {
    import(state, assigns) {
      state.assigns = {...assigns}
    },
  
    add(state, assigns: Assign[]) {
      var lastPlayerId: string
      assigns.forEach(assign => {
        if (assign.type == "player") {
          if (assign.id === undefined) {
            do {
              assign.id = nextAssignId(state)
            } while (state.assigns[assign.id] !== undefined)
          }
          
          assign.playerId = assign.id
          lastPlayerId = assign.id
        } else if (assign.type == "spell") {
          if (assign.playerId === undefined) {
            assign.playerId = lastPlayerId
          }
          if (assign.id === undefined) {
            assign.id = `${assign.playerId}.${assign.spell.id}`
          }
        }
        Vue.set(state.assigns, assign.id, assign)
      })
    },

    set(state, assign: Assign) {
      if (assign.id === undefined) {
        do {
          assign.id = nextAssignId(state)
        } while (state.assigns[assign.id] !== undefined)
      }
      Vue.set(state.assigns, assign.id, assign)
    },

    remove(state, id) {
      Vue.delete(state.assigns, id)
    },
  },

  actions: {
    remove({commit}, id) {
      commit('events/clearAssignment', id, {root: true})
      commit('remove', id)
    },
  },
} as Module<AssignState, State>
