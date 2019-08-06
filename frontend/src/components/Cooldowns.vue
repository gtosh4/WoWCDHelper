<template>
<v-container fluid>
  <v-layout wrap>
    <v-flex v-for="(classInfo, className) in classes" v-bind:key="className" xs3 class="pa-1">
      <v-card outlined>
        <v-chip
          label
          :text-color="classInfo.colour"
          color="transparent"
          @click="expanded = {...expanded, [className]: !expanded[className]}"
          style="width: 100%"
        >
          <img :src="classIcon(className)" class="mr-2" />
          {{ className }}
          <v-spacer />
          <v-icon :class="expandedClass(expanded[className])">$vuetify.icons.expand</v-icon>
        </v-chip>
        <v-list v-if="expanded[className]">
          <v-list-item v-for="(specInfo, specName) in classInfo.specs" v-bind:key="specName">
            <v-chip outlined>{{ specCount(className, specName) }}</v-chip>
            <v-icon @click="addSpec(className, specName)" class="ml-2">mdi-account-plus</v-icon>
            <img :src="specIcon(specInfo)" class="mx-2" v-if="specInfo.icon" />
            <v-icon v-else>mdi-account-alert-outline</v-icon>
            <span :style="{color: classInfo.colour}">{{ specName }}</span>
          </v-list-item>
        </v-list>
      </v-card>
    </v-flex>
  </v-layout>
</v-container>
</template>
<script>

const classes = {
  deathknight: {
    colour: "#C41F3B",
    specs: {
    },
  },
  demonhunter: {
    colour: "#A330C9",
    specs: {
    },
  },
  druid: {
    colour: "#FF7D0A",
    specs: {
      resto: {
        icon: "talentspec_druid_restoration",
      },
    },
  },
  hunter: {
    colour: "#ABD473",
    specs: {
    },
  },
  mage: {
    colour: "#40C7EB",
    specs: {
    },
  },
  monk: {
    colour: "#00FF96",
    specs: {
    },
  },
  paladin: {
    colour: "#F58CBA",
    specs: {
      holy: {
        icon: "spell_holy_holybolt",
      },
    },
  },
  priest: {
    colour: "#FFFFFF",
    specs: {
    },
  },
  rogue: {
    colour: "#FFF569",
    specs: {
    },
  },
  shaman: {
    colour: "#0070DE",
    specs: {
    },
  },
  warlock: {
    colour: "#8787ED",
    specs: {
    },
  },
  warrior: {
    colour: "#C79C6E",
    specs: {
    },
  },
}

export default {
  data: () => ({
    classes,
    members: [],
    expanded: {},
  }),

  props: {
    assigns: {
      type: Array,
      required: true,
    },
  },

  computed: {
    indexedMembers() {
      return this.members.reduce((ms, m) => {
        let c = ms[m.className]
        if (!c) {
          c = {}
          ms[m.className] = c
        }
        let s = c[m.specName]
        if (!s) {
          s = []
          c[m.specName] = s
        }
        s.push(m)
        return ms
      }, [])
    },
  },

  methods: {
    classIcon(className, size="small") {
      return `https://wow.zamimg.com/images/wow/icons/${size}/class_${className}.jpg`
    },

    specIcon(specInfo, size="tiny") {
      return `https://wow.zamimg.com/images/wow/icons/${size}/${specInfo.icon}.gif`
    },

    specCount(className, specName) {
      return ((this.indexedMembers[className] || {})[specName] || []).length
    },

    addSpec(className, specName) {
      this.members = [...this.members, {className, specName}]
    },

    expandedClass(expanded) {
      const c = ['v-data-table__expand-icon']
      if (expanded) {
        c.push('v-data-table__expand-icon--active')
      }
      return c
    }
  },
};
</script>
