import Color from 'color'

export function toColor(c) {
  return c ? Color.rgb({r: c.r, g: c.g, b: c.b, alpha: c.a}) : undefined
}

export function toRGBA(c) {
  const o = c.rgb().object()
  return {r: o.r, g: o.g, b: o.b, a: c.alpha()}
}
