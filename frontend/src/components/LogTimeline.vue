<template>
  <v-card outlined tile>
    <v-card-title>
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
          <v-icon @click="loadLog">
            mdi-send
          </v-icon>
        </template>
      </v-text-field>
    </v-card-title>
    <svg :viewBox="`0 0 ${width} ${height+assigmentsMaxHeight}`" class="log-timeline">
      <g id="chart" :transform="`translate(${margin.left}, ${margin.top})`">
        <path :d="hpLine" fill="slategrey" />
        <path
          :d="dtpsLine"
          fill="none"
          stroke="crimson"
          stroke-width="2"
        />
        <path
          v-if="showHealing"
          :d="hpsLine"
          fill="none"
          stroke="springgreen"
          stroke-width="2"
        />

        <g v-for="event in events" :key="event.id" :transform="`translate(${xForEvent(event)}, 0)`">
          <line 
            :id="`line-event-${event.id}`"
            :x1="0"
            :x2="0"
            y1="0"
            :y2="chartSize.height + xAxisHeight + assignMargin"
            stroke-width="1.5"
            :stroke="colourForEvent(event)"
          />
          <text
            :id="`text-event-${event.id}`"
            :x="-chartSize.height/2"
            :y="-4"
            :fill="colourForEvent(event)"
            transform="rotate(-90)"
            class="label-text"
            text-anchor="middle"
          >
            {{ event.label || "" }}
          </text>
          <g :transform="`translate(${-(assignImageSize/2)}, ${xAxisHeight + chartSize.height + assignMargin + 2})`">
            <image
              v-for="(assign, index) in event.assignments"
              :key="index"
              x="0"
              :y="index * (assignImageSize+1)"
              :height="`${assignImageSize}px`"
              :width="`${assignImageSize}px`"
              :href="imageForAssign(assign)"
            />
          </g>
        </g>
        <g ref="hpaxis" />
        <g ref="dtpsaxis" />
        <g ref="xaxis" />
      </g>
    </svg>
  </v-card>
</template>
<script>
import moment from 'moment'
import * as d3 from 'd3'
import Color from 'color'

import { BackendAPI } from '../api/backend'
import { toColor } from './colour_utils'
import {formatDuration} from './duration_utils'
import {spells, spec, specIcon, classIcon} from './wow_info'

const defaultColour = Color('rgb(200, 200, 200)')
const backend = new BackendAPI()

const reReportPath = /reports\/(\w+)/
const reFightHash = /fight=(\d+)/

function clearChildren(node) {
  while (node.firstChild) {
    node.remove(node.firstChild)
  }
}

export default {
  props: {
    showHealing: {
      type: Boolean,
      default: false,
    },
  },

  data: () => ({
    margin: {top: 10, right: 50, bottom: 30, left: 50},
    height: 480,
    width: 1280,
    assignImageSize: 18,
    assignMargin: 4,

    hpAxis: null,
    dtpsAxis: null,
    xAxis: null,
    legend: null,

    loading: null,
    raidHealth: null,

    // "hack" to add reactivity for when the axis are updated
    axisComputation: 0,
  }),

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
      return this.$store.getters['events/ordered']
    },

    assigmentsMaxHeight() {
      return this.events.reduce((max, event) => max > event.assignments.length ? max : event.assignments.length, 0) * this.assignImageSize
    },

    xAxisHeight() {
      this.axisComputation
      const node = this.xAxis ? this.xAxis.node() : null
      return node ? node.getBBox().height : 0
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
        .y1(d => this.yRaidHP(d.max > 0 ? d.current / d.max : 1))
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

  mounted() {
    this.hpAxis = d3.select(this.$refs.hpaxis)
    this.dtpsAxis = d3.select(this.$refs.dtpsaxis)
    this.xAxis = d3.select(this.$refs.xaxis)
    this.legend = d3.select(this.$refs.legend)

    if (this.logURL) {
      this.loadLog()
    }
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
      this.axisComputation++
    },

    xForEvent(event) {
      return this.x(event.time.asSeconds())
    },

    colourForEvent(event) {
      return (toColor(event.colour) || defaultColour).string()
    },

    loadLog() {
      const url = new URL(this.logURL)
      if (url === undefined || url.hostname != "www.warcraftlogs.com") return

      const [, report] = url.pathname.match(reReportPath)
      const [, fight] = url.hash.match(reFightHash)
      if (report == null || fight == null) return

      this.loading = true
      backend.getRaidHealth(report, fight)
        .then(d => {
          this.loading = false
          this.raidHealth = [...d.map(d => {
            return {
              ...d,
              timestamp: moment.duration(d.time_sec, 'second'),
            }
          })]
        })
        .catch(e => {
          this.loading = false
          console.error("getRaidHealth error", {e, report, fight})
        })
    },

    validateLogURL() {
      if (this.logURL == null || this.logURL == "") return true

      const url = new URL(this.logURL)
      if (url === undefined || url.hostname != "www.warcraftlogs.com") return `${url ? url.hostname : ''} not a warcraftlogs URL`

      const [report] = url.pathname.match(reReportPath)
      const [fight] = url.hash.match(reFightHash)

      if (report == null || fight == null) return 'URL does not contain log id and fight number'

      return true
    },

    imageForAssign(assignId) {
      const assign = this.$store.state.assigns.assigns[assignId]
      if (assign.spell != undefined) {
        const spell = spells[assign.spell.id]
        return spell != undefined ? spell.icon : ''
      }
      if (assign.className && assign.specName) {
        return specIcon(spec(assign.className, assign.specName))
      }
      if (assign.className) {
        return classIcon(assign.className)
      }
    },
  },
}
</script>
<style>
.logUrlInput {
  min-height: 0px;
}
.log-timeline .label-text {
    paint-order: stroke;
    stroke: black;
    stroke-width: 1px;
    stroke-linecap: butt;
    stroke-linejoin: miter;
}
</style>
