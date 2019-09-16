import router from '../../plugins/router'
import jsonurl from 'json-url'
import moment from 'moment'

const codec = jsonurl('lzma')

function shortenEvent(event) {
  return {
    id: event.id,
    t: event.time.asSeconds(),
    a: event.assignments,
    l: event.label,
    c: event.colour,
  }
}

function expandEvent(e) {
  return {
    id: e.id,
    time: moment.duration(e.t, 'seconds'),
    assignments: e.a,
    label: e.l,
    colour: e.c,
  }
}

function shortenAssign(assign) {
  const a = {id: assign.id}
  if (assign.id == assign.playerId) {
    a.n = assign.name
    a.c = assign.className
    if (assign.specName) {
      a.s = assign.specName
    }
  } else {
    // ONLY set `p` for non-player assigns
    a.p = assign.playerId
    if (assign.spell) {
      a.s = {id: assign.spell.id, c: assign.spell.cfg}
    }
  }
  return a
}

function expandAssign(a) {
  const assign = {id: a.id}
  if (a.p) {
    assign.playerId = a.p
    if (a.s) {
      assign.spell = {
        id: a.s.id,
        cfg: a.s.c,
      }
    }
  } else {
    assign.name = a.n
    assign.playerId = a.id
    assign.className = a.c
    assign.specName = a.s
  }
  return assign
}

function shortenState(state) {
  const s = {
    n: state.name,
    u: state.logURL,
  }
  if (state.assigns) {
    s.a = Object.values(state.assigns.assigns).map(shortenAssign)
  }
  if (state.events) {
    s.e = Object.values(state.events.events).map(shortenEvent)
  }
  return s
}

function expandState(s) {
  const state = {
    name: s.n,
    logURL: s.u,
    events: {events: {}},
    assigns: {assigns: {}},
    dragassign: {info: null},
  }
  if (s.a) {
    s.a.forEach(a => {
      const assign = expandAssign(a)
      state.assigns.assigns[assign.id] = assign
    })
  }
  if (s.e) {
    s.e.forEach(e => {
      const event = expandEvent(e)
      state.events.events[event.id] = event
    })
  }
  return state
}

export function loadFromURL(store, route) {
  // use replaceState instead of mutations so that we don't trigger updateURL and stomp our changes

  const data = route.params.data
  if (!data) {
    store.replaceState({
      name: "",
      logURL: null,
      events: {events: {}},
      assigns: {assigns: {}},
      dragassign: {info: null},
    })
    return
  }

  codec.decompress(data).then(short => {
    const json = JSON.parse(short)
    const state = expandState(json)
    store.replaceState(state)
  }).catch(err => console.log("Error decompressing", {data, err}))
}

export function updateURLPlugin(store) {
  const ignorePatterns = [
    /^dragassign/,
  ]
  store.subscribe((mutation, state) => {
    const ignored = ignorePatterns.some(p => p.test(mutation.type))
    if (ignored) {
      return
    }
    const short = shortenState(state)
    const str = JSON.stringify(short)
    codec.compress(str).then(data => {
      if (data == router.currentRoute.params.data) {
        return
      }
      const r = {
        path: `/${data}`,
        query: router.currentRoute.query,
      }
      router.push(r)
    }).catch(err => console.log("Error compressing", {state, short, str, err}))
  })
}
