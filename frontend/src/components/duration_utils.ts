require("moment-duration-format")

export function formatDuration(t) {
  return t ? t.format('mm:ss', { trim: false }) : ''
}
