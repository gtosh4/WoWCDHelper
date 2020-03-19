<template>
<v-card outlined tile>
  <v-row>
    <v-text-field
      v-model="logURL"
      :loading="loading"
      outlined
      single-line
      clearable
      placeholder="log URL"
      hint="https://www.warcraftlogs.com/reports/aaaaaaa#fight=1"
      class="logUrlInput"
      :rules="[validateLogURL]"
    >
    <template #append>
      <v-icon @click="loadLog">mdi-send</v-icon>
    </template>
    </v-text-field>
  </v-row>
  <svg :viewBox="`0 0 ${width} ${height}`">
    <g :transform="`translate(${margin.left}, ${margin.top})`" id="chart">
      <g ref="hpaxis" />
      <g ref="dtpsaxis" />
      <g ref="xaxis" />

      <path :d="hpLine" fill="slategrey" />
      <path :d="dtpsLine" fill="none" stroke="crimson" stroke-width="2" />
      <path v-if="showHealing" :d="hpsLine" fill="none" stroke="springgreen" stroke-width="2" />

      <line v-for="event in events" :key="event.id"
        :x1="xForEvent(event)" :x2="xForEvent(event)"
        y1="0" :y2="height"
        stroke-width="1.5"
        :stroke="colourForEvent(event)"
        :id="`line-event-${event.id}`"
      />
    </g>
  </svg>
  <v-container ma-0 pa-0 justify-start align-start class="timeline-events">
    <v-row row justify-start align-start ma-0>
      <v-row v-for="event in events" :key="event.id" column justify-start align-start>
        <Assignment v-for="(assign, index) in event.assignments" :key="index" :eventId="event.id" :index="index" :moveable="false" />
      </v-row>
    </v-row>
  </v-container>
</v-card>
</template>
<script>
import Assignment from './Assignment'

import moment from 'moment'
import * as d3 from 'd3'
import Color from 'color'

import { BackendAPI } from '../api/backend'
import { toColor } from './colour_utils';
import {formatDuration} from './duration_utils'

const defaultColour = Color('rgb(200, 200, 200)')
const backend = new BackendAPI()

const logURLParse = /reports\/(\w+)#.*(fight=(\d+))/

function clearChildren(node) {
  while (node.firstChild) {
    node.remove(node.firstChild)
  }
}

export default {
  data: () => ({
    margin: {top: 10, right: 50, bottom: 30, left: 50},
    height: 480,
    width: 1280,

    hpAxis: null,
    dtpsAxis: null,
    xAxis: null,
    legend: null,

    loading: null,
    raidHealth: null,
  }),

  props: {
    showHealing: {
      default: false,
    },
  },

  created() {
  },

  mounted() {
    this.hpAxis = d3.select(this.$refs.hpaxis)
    this.dtpsAxis = d3.select(this.$refs.dtpsaxis)
    this.xAxis = d3.select(this.$refs.xaxis)
    this.legend = d3.select(this.$refs.legend)

    if (this.logURL) {
      this.loadLog()
    }
  },

  computed: {
    logURL: {
      get() {
        return this.$store.state.logURL
      },

      set(v) {
        this.$store.commit("setLogURL", v)
      },
    },

    maxTime() {
      const maxEvent =  this.events.length > 0 ? this.events[this.events.length - 1].time : moment.duration(0)
      if (this.raidHealth) {
        const maxHealthEvent = this.raidHealth[this.raidHealth.length-1].timestamp
        return maxEvent.asSeconds() > maxHealthEvent.asSeconds() ? maxEvent : maxHealthEvent
      }
      return maxEvent
    },

    events() {
      return this.$store.getters['events/orderedEvents']
    },

    x() {
      return d3.scaleLinear()
        .domain([0, this.maxTime.asSeconds()])
        .range([0, this.chartSize.width])
    },

    yRaidHP() {
      return d3.scaleLinear()
        .domain([0, 1])
        .range([this.chartSize.height, 0])
    },

    yRaidDTPS() {
      var maxDTPS = 0
      if (this.raidHealth) {
        maxDTPS = Math.max(...this.raidHealth.map(e => e.damage_taken))
      }
      return d3.scaleLinear()
        .domain([0, maxDTPS])
        .range([this.chartSize.height, 0])
    },

    chartSize() {
      return {
        width: this.width - (this.margin.left + this.margin.right),
        height: this.height - (this.margin.top + this.margin.bottom),
      }
    },

    hpLine() {
      if (!this.raidHealth) return undefined

      const area = d3.area()
        .x(d => this.x(d.time_sec))
        .y0(this.yRaidHP(0))
        .y1(d => this.yRaidHP(d.current / d.max))
      return area(this.raidHealth)
    },

    dtpsLine() {
      if (!this.raidHealth) return undefined

      const line = d3.line()
        .x(d => this.x(d.time_sec))
        .y(d => this.yRaidDTPS(d.damage_taken))
      return line(this.raidHealth)
    },

    hpsLine() {
      if (!this.raidHealth) return undefined

      const line = d3.line()
        .x(d => this.x(d.time_sec))
        .y(d => this.yRaidDTPS(d.healing))
      return line(this.raidHealth)
    }
  },

  watch: {
    x() {
      this.updateAxis()
    },
    yRaidHP() {
      this.updateAxis()
    },
    yRaidDTPS() {
      this.updateAxis()
    },
  },

  methods: {
    updateAxis() {
      if (this.xAxis) {
        clearChildren(this.xAxis)
        const ticks = Array.from({length: Math.floor(this.maxTime.asSeconds() / 30)+1}, (_,i) => i * 30)
        this.xAxis
          .attr("transform", `translate(0, ${this.chartSize.height})`)
          .call(
            d3.axisBottom(this.x)
              .tickValues(ticks)
              .tickFormat(s => formatDuration(moment.duration(s, "seconds")))
          )
      }

      if (this.hpAxis) {
        clearChildren(this.hpAxis)
        this.hpAxis
          .call(
            d3.axisLeft(this.yRaidHP)
              .tickFormat(d3.format(".0%"))
          )
      }
      
      if (this.dtpsAxis) {
        clearChildren(this.dtpsAxis)
        this.dtpsAxis
          .attr("transform", `translate(${this.chartSize.width}, 0)`)
          .call(
            d3.axisRight(this.yRaidDTPS)
              .tickFormat(d3.format("~s"))
          )
      }
    },

    xForEvent(event) {
      return this.x(event.time.asSeconds())
    },

    colourForEvent(event) {
      return (toColor(event.colour) || defaultColour).string()
    },

    loadLog() {
      const match = (this.logURL || '').match(logURLParse)
      if (!match) return
      const [, id, , fight] = match
      this.loading = true
      backend.getRaidHealth(id, fight)
        .then(d => {
          this.loading = false
          this.raidHealth = [...d.map(d => {
            return {
              ...d,
              timestamp: moment.duration(d.time_sec, 'second'),
            }
          })]
        })
        .catch(e => console.log("getRaidHealth error", {e, match, id, fight}))
    },

    validateLogURL() {
      return !!this.logURL && (!!this.logURL.match(logURLParse) || 'URL does not contain log id and fight number')
    },
  },

  components: {
    Assignment,
  },
};
</script>
<style>
.timeline-events .assignment {
  margin-bottom: 4px;
}
.logUrlInput {
  min-height: 0px;
}
</style>
