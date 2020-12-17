import axios from 'axios'

const client = axios.create({
  // baseURL: 'https://www.wowdb.com',
  baseURL: '//www.wowdb.com',
})

function parseJson(body) {
  // For whatever reason, wowdb returns JSON wrapped in ()
  return JSON.parse(body.substring(1, body.length-1))
}

export default {
  getSpellTooltip(spell) {
    return client.get(`/spells/${spell}/tooltip`).then(response => parseJson(response.data).Tooltip)
  },

  getSpellInfo(spell) {
    return client.get(`/api/spell/${spell}`).then(response => parseJson(response.data))
  },
}
