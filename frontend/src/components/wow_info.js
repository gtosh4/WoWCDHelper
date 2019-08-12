export const classes = {
  deathknight: {
    colour: {r: 196, g: 31, b: 59},
    specs: {
      blood: {
        icon: "spell_deathknight_bloodpresence",
        spells: [],
      },
      frost: {
        icon: "spell_deathknight_frostpresence",
        spells: [],
      },
      unholy: {
        icon: "spell_deathknight_unholypresence",
        spells: [],
      },
    },
  },
  demonhunter: {
    colour: {r: 163, g: 48, b: 201},
    specs: {
      havoc: {
        icon: "ability_demonhunter_specdps",
        spells: [
          {id: 196718, name: "Darkness", icon: "https://media.wowdb.com/wow/icons/large/ability_demonhunter_darkness.jpg"},
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
          {id: 740, name: "Tranquility", icon: "https://media.wowdb.com/wow/icons/large/spell_nature_tranquility.jpg"},
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
          {id: 115310, name: "Revival", icon: "https://media.wowdb.com/wow/icons/large/spell_monk_revival.jpg"},
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
          {id: 31821, name: "Aura Mastery", icon: "https://media.wowdb.com/wow/icons/large/spell_holy_auramastery.jpg"},
          {id: 31884, name: "Avenging Wrath", icon: "https://media.wowdb.com/wow/icons/large/spell_holy_avenginewrath.jpg"},
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
          {id: 81782, name: "Power Word: Barrier", icon: "https://media.wowdb.com/wow/icons/large/spell_holy_powerwordbarrier.jpg"},
          {id: 246287, name: "Evangelism", icon: "https://media.wowdb.com/wow/icons/large/spell_holy_divineillumination.jpg"},
          {id: 47536, name: "Rapture", icon: "https://media.wowdb.com/wow/icons/large/spell_holy_rapture.jpg"},
        ],
      },
      holy: {
        icon: "spell_holy_guardianspirit",
        spells: [
          {id: 64843, name: "Divine Hymn", icon: "https://media.wowdb.com/wow/icons/large/spell_holy_divinehymn.jpg"},
          {id: 265202, name: "Holy Word: Salvation", icon: "https://media.wowdb.com/wow/icons/large/ability_priest_archangel.jpg"},
        ],
      },
      shadow: {
        icon: "spell_shadow_shadowwordpain",
        spells: [
          {id: 15286, name: "Vampiric Embrace", icon: "https://media.wowdb.com/wow/icons/large/spell_shadow_unsummonbuilding.jpg"},
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
          {id: 108280, name: "Healing Tide Totem", icon: "https://media.wowdb.com/wow/icons/large/ability_shaman_healingtide.jpg"},
          {id: 98008, name: "Spirit Link Totem", icon: "https://media.wowdb.com/wow/icons/large/spell_shaman_spiritlink.jpg"},
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
          {id: 97462, name: "Rallying Cry", icon: "https://media.wowdb.com/wow/icons/large/ability_warrior_rallyingcry.jpg"},
        ],
      },
      fury: {
        icon: "ability_warrior_innerrage",
        spells: [
          {id: 97462, name: "Rallying Cry", icon: "https://media.wowdb.com/wow/icons/large/ability_warrior_rallyingcry.jpg"},
        ],
      },
      protection: {
        icon: "ability_warrior_defensivestance",
        spells: [
          {id: 97462, name: "Rallying Cry", icon: "https://media.wowdb.com/wow/icons/large/ability_warrior_rallyingcry.jpg"},
        ],
      },
    },
  },
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
