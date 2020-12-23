import {Module} from 'vuex'
import { State } from '..'

export interface DragAssignState {
  draggedAssign?: number
}

export default {
  namespaced: true,

  state: {
  },
  
  mutations: {
    start(state, assignId) {
      state.draggedAssign = assignId
    },
  
    stop(state) {
      state.draggedAssign = undefined
    },
  },
} as Module<DragAssignState, State>
