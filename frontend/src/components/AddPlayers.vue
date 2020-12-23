<template>
  <v-card
    outlined
    tile
    class="add-players-all"
  >
    <v-btn
      tile
      text
      left
      width="100%"
      @click="expanded = !expanded"
    >
      <v-icon :class="expandedClass">
        $vuetify.icons.expand
      </v-icon>
      <span>all</span>
    </v-btn>
    <v-lazy>
      <v-row no-gutters justify="start" :style="{display: expanded ? '' : 'none'}">
        <v-col v-for="(classInfo, className) in classes" :key="className" cols="4">
          <v-card outlined tile class="add-player-class">

            <v-card-title>
              <WowIcon :class-name="className" />
              <span :style="{color: classColour(className)}">{{ className }}</span>
              <div class="add-player-action">
                <span>{{ counts[className].total }}</span>
                <v-btn
                  tile
                  x-small
                  icon
                  @click="$emit('addPlayer', className)"
                >
                  <v-icon>mdi-account-plus</v-icon>
                </v-btn>
              </div>
            </v-card-title>

            <v-divider />
            <v-list>
              <v-list-item v-for="(specInfo, specName) in classInfo.specs" :key="specName" class="player-select">
                <WowIcon :class-name="className" :spec-name="specName" />
                <span :style="{color: classColour(className)}">{{ specName }}</span>
                <div class="add-player-action">
                  <span>{{ counts[className][specName] }}</span>
                  <v-btn
                    tile
                    x-small
                    icon
                    @click="$emit('addPlayer', className, specName)"
                  >
                    <v-icon>mdi-account-plus</v-icon>
                  </v-btn>
                </div>
              </v-list-item>
            </v-list>
          </v-card>
        </v-col>
      </v-row>
    </v-lazy>
  </v-card>
</template>
<script>
import WowIcon from './WowIcon'

import { mapGetters } from 'vuex'
import { classes, classColour } from './wow_info'

export default {
  components: {
    WowIcon,
  },

  data: () => ({
    classes,
    expanded: false,
  }),

  computed: {
    ...mapGetters('assigns', ['players']),

    expandedClass() {
      const c = ['v-data-table__expand-icon']
      if (this.expanded) {
        c.push('v-data-table__expand-icon--active')
      }
      return c
    },

    counts() {
      const cnt = Object.entries(classes).reduce((c, [className, classObj]) => {
        if (!c[className]) c[className] = {total: 0}
        Object.keys(classObj.specs).forEach(specName => c[className][specName] = 0)
        return c
      }, {total: 0})
      return this.players.reduce((cnt, {className, specName}) => {
        cnt.total++

        const classCnt = cnt[className]
        if (classCnt != null) {
          classCnt.total++

          if (classCnt[specName] != null) {
            classCnt[specName]++
          }
        }
        return cnt
      }, cnt)
    },
  },

  methods: {
    classColour,

    classCount(className) {
      return Object.values(this.players[className] || {}).reduce((sum, spec) => sum + spec.length, 0)
    },

    specCount(className, specName) {
      return ((this.players[className] || {})[specName] || []).length
    },
  },
}
</script>
<style>
.add-players-all {
  max-width: 20vw;
}

.add-players-all > .v-chip {
  width: 100%;
}

.add-players-all .v-list {
  padding-top: 0;
}

.add-players-all .v-list-item {
  min-height: 20px;
  padding-left: 4px;
  padding-right: 4px;
}

.add-player-class .v-card__title {
  padding: 4px;
}

.add-player-class .v-divider {
  margin-left: 4px;
  margin-right: 4px;
}

.add-players-all .add-player-action {
  margin-left: auto;
}

.add-players-all .add-player-action span {
  margin-right: 8px;
  white-space: nowrap;
}
</style>
