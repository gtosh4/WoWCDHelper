import moment from "moment"
import 'moment-duration-format'

export function formatDuration(t?: moment.Duration) {
  return t?.format('mm:ss', { trim: false }) || ''
}
