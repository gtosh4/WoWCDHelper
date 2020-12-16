export const classes = {
  deathknight: {
    colour: {r: 196, g: 31, b: 59},
    specs: {
      blood: {
        icon: "spell_deathknight_bloodpresence",
        spells: [
          {id: 51052}, // Anti-magic Zone
        ],
      },
      frost: {
        icon: "spell_deathknight_frostpresence",
        spells: [
          {id: 51052}, // Anti-magic Zone
        ],
      },
      unholy: {
        icon: "spell_deathknight_unholypresence",
        spells: [
          {id: 51052}, // Anti-magic Zone
        ],
      },
    },
  },
  demonhunter: {
    colour: {r: 163, g: 48, b: 201},
    specs: {
      havoc: {
        icon: "ability_demonhunter_specdps",
        spells: [
          {id: 196718}, // Darkness
        ],
      },
      vengeance: {
        icon: "ability_demonhunter_spectank",
        spells: [],
      },
    },
  },
  druid: {
    colour: {r: 255, g: 125, b: 10},
    specs: {
      balance: {
        icon: "spell_nature_starfall",
        spells: [],
      },
      feral: {
        icon: "ability_druid_catform",
        spells: [],
      },
      guardian: {
        icon: "ability_racial_bearform",
        spells: [],
      },
      restoration: {
        icon: "spell_nature_healingtouch",
        spells: [
          { // Tranquility
            id: 740,
            options: [
              {text: "Inner Peace", type: Boolean, default: false, prop: 'inner_peace'},
            ],
            configure(spell) {
              if (spell.cfg.inner_peace) {
                spell.cd -= 60
              }
            },
          },
          {id: 33891}, // Incarnation: Tree of Life
        ],
      },
    },
  },
  hunter: {
    colour: {r: 171, g: 212, b: 115},
    specs: {
      beastmastery: {
        icon: "ability_hunter_bestialdiscipline",
        spells: [],
      },
      marksmanship: {
        icon: "ability_hunter_focusedaim",
        spells: [],
      },
      survival: {
        icon: "ability_hunter_camouflage",
        spells: [],
      },
    },
  },
  mage: {
    colour: {r: 64, g: 199, b: 235},
    specs: {
      arcane: {
        icon: "spell_holy_magicalsentry",
        spells: [],
      },
      fire: {
        icon: "spell_fire_firebolt02",
        spells: [],
      },
      frost: {
        icon: "spell_frost_frostbolt02",
        spells: [],
      },
    },
  },
  monk: {
    colour: {r: 0, g: 255, b: 150},
    specs: {
      brewmaster: {
        icon: "monk_stance_drunkenox",
        spells: [],
      },
      mistweaver: {
        icon: "monk_stance_wiseserpent",
        spells: [
          {id: 115310}, // Revival
        ],
      },
      windwalker: {
        icon: "monk_stance_whitetiger",
        spells: [],
      },
    },
  },
  paladin: {
    colour: {r: 245, g: 140, b: 186},
    specs: {
      holy: {
        icon: "spell_holy_holybolt",
        spells: [
          {id: 31821}, // Aura Mastery
          { // Avenging Wrath
            id: 31884,
            options: [
            ],
          },
        ],
      },
      protection: {
        icon: "ability_paladin_shieldofthetemplar",
        spells: [],
      },
      retribution: {
        icon: "spell_holy_auraoflight",
        spells: [],
      },
    },
  },
  priest: {
    colour: {r: 255, g: 255, b: 255},
    specs: {
      discipline: {
        icon: "spell_holy_powerwordshield",
        spells: [
          {id: 62618}, // Power Word: Barrier
          {id: 246287}, // Evangelism
          {id: 47536}, // Rapture
        ],
      },
      holy: {
        icon: "spell_holy_guardianspirit",
        spells: [
          { // Divine Hymn
            id: 64843,
            options: [
            ],
          }, 
          {id: 265202}, // Holy Word: Salvation
        ],
      },
      shadow: {
        icon: "spell_shadow_shadowwordpain",
        spells: [
          {id: 15286}, // Vampiric Embrace
        ],
      },
    },
  },
  rogue: {
    colour: {r: 255, g: 245, b: 105},
    specs: {
      assassination: {
        icon: "ability_rogue_eviscerate",
        spells: [],
      },
      outlaw: {
        icon: "ability_backstab",
        spells: [],
      },
      subtlety: {
        icon: "ability_stealth",
        spells: [],
      },
    },
  },
  shaman: {
    colour: {r: 0, g: 112, b: 222},
    specs: {
      elemental: {
        icon: "spell_nature_lightning",
        spells: [],
      },
      enhancement: {
        icon: "spell_nature_lightningshield",
        spells: [],
      },
      restoration: {
        icon: "spell_nature_magicimmunity",
        spells: [
          { // Healing Tide Totem
            id: 108280,
          },
          {id: 98008}, // Spirit Link Totem
        ],
      },
    },
  },
  warlock: {
    colour: {r: 135, g: 135, b: 237},
    specs: {
      affliction: {
        icon: "spell_shadow_deathcoil",
        spells: [],
      },
      demonology: {
        icon: "spell_shadow_metamorphosis",
        spells: [],
      },
      destruction: {
        icon: "spell_shadow_rainoffire",
        spells: [],
      },
    },
  },
  warrior: {
    colour: {r: 199, g: 156, b: 110},
    specs: {
      arms: {
        icon: "ability_warrior_savageblow",
        spells: [
          {id: 97462}, // Rallying Cry
        ],
      },
      fury: {
        icon: "ability_warrior_innerrage",
        spells: [
          {id: 97462}, // Rallying Cry
        ],
      },
      protection: {
        icon: "ability_warrior_defensivestance",
        spells: [
          {id: 97462}, // Rallying Cry
        ],
      },
    },
  },
}

export const spells = {
  196718: {id: 196718, name: "Darkness", icon: "https://media.wowdb.com/wow/icons/large/ability_demonhunter_darkness.jpg", cd: 180},
  740:    {id: 740,    name: "Tranquility", icon: "https://media.wowdb.com/wow/icons/large/spell_nature_tranquility.jpg", cd: 180},
  33891:  {id: 33891,  name: "Incarnation: Tree of Life", icon: "https://media.wowdb.com/wow/icons/large/ability_druid_treeoflife.jpg", cd: 180},
  115310: {id: 115310, name: "Revival", icon: "https://media.wowdb.com/wow/icons/large/spell_monk_revival.jpg", cd: 180},
  31821:  {id: 31821,  name: "Aura Mastery", icon: "https://media.wowdb.com/wow/icons/large/spell_holy_auramastery.jpg", cd: 180},
  31884:  {id: 31884,  name: "Avenging Wrath", icon: "https://media.wowdb.com/wow/icons/large/spell_holy_avenginewrath.jpg", cd: 120},
  62618:  {id: 62618,  name: "Power Word: Barrier", icon: "https://media.wowdb.com/wow/icons/large/spell_holy_powerwordbarrier.jpg", cd: 180},
  246287: {id: 246287, name: "Evangelism", icon: "https://media.wowdb.com/wow/icons/large/spell_holy_divineillumination.jpg", cd: 90},
  47536:  {id: 47536,  name: "Rapture", icon: "https://media.wowdb.com/wow/icons/large/spell_holy_rapture.jpg", cd: 90},
  64843:  {id: 64843,  name: "Divine Hymn", icon: "https://media.wowdb.com/wow/icons/large/spell_holy_divinehymn.jpg", cd: 180},
  265202: {id: 265202, name: "Holy Word: Salvation", icon: "https://media.wowdb.com/wow/icons/large/ability_priest_archangel.jpg", cd: 300},
  15286:  {id: 15286,  name: "Vampiric Embrace", icon: "https://media.wowdb.com/wow/icons/large/spell_shadow_unsummonbuilding.jpg", cd: 120},
  108280: {id: 108280, name: "Healing Tide Totem", icon: "https://media.wowdb.com/wow/icons/large/ability_shaman_healingtide.jpg", cd: 180},
  98008:  {id: 98008,  name: "Spirit Link Totem", icon: "https://media.wowdb.com/wow/icons/large/spell_shaman_spiritlink.jpg", cd: 180},
  97462:  {id: 97462,  name: "Rallying Cry", icon: "https://media.wowdb.com/wow/icons/large/ability_warrior_rallyingcry.jpg", cd: 180},
  51052:  {id: 51052,  name: "Anti-Magic Zone", icon: "https://icons.wowdb.com/ptr/large/spell_deathknight_antimagiczone.jpg?36734", cd: 120},
}

export const healers = [
  {className: 'druid', specName: 'restoration'},
  {className: 'monk', specName: 'mistweaver'},
  {className: 'paladin', specName: 'holy'},
  {className: 'priest', specName: 'discipline'},
  {className: 'priest', specName: 'holy'},
  {className: 'shaman', specName: 'restoration'},
]

export function classColour(className) {
  const c = classes[className].colour
  return `rgb(${c.r}, ${c.g}, ${c.b})`
}

export function classIcon(className, size="small") {
  return `https://wow.zamimg.com/images/wow/icons/${size}/class_${className}.jpg`
}

export function specIcon(specInfo, size="tiny") {
  return `https://wow.zamimg.com/images/wow/icons/${size}/${specInfo.icon}.gif`
}

export function abilityIcon(abilityName, size="tiny") {
  return `https://wow.zamimg.com/images/wow/icons/${size}/ability_${abilityName}.gif`
}

export function spec(className, specName) {
  return (classes[className] || {}).specs[specName] || {}
}
