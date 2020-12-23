import Color from 'color'

export function toColor(c?: {r: number, g: number, b: number, a?: number}) {
  return c ? Color.rgb({r: c.r, g: c.g, b: c.b, alpha: c.a}) : undefined
}

export function toRGBA(c: Color): {r: number, g: number, b: number, a?: number} {
  const o = c.rgb().object()
  return {r: o.r, g: o.g, b: o.b, a: c.alpha()}
}
