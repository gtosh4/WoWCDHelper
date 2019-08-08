export const classes = {
  deathknight: {
    colour: {r: 196, g: 31, b: 59},
    specs: {
      blood: {
        icon: "spell_deathknight_bloodpresence",
        spells: [],
      },
      frost: {
        icon: "spell_deathknight_frostpresence.gif",
        spells: [],
      },
      unholy: {
        icon: "spell_deathknight_unholypresence.gif",
        spells: [],
      },
    },
  },
  demonhunter: {
    colour: {r: 163, g: 48, b: 201},
    specs: {
      havoc: {
        icon: "ability_demonhunter_specdps",
        spells: [],
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
        spells: [],
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
        spells: [],
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
        spells: [],
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
        spells: [],
      },
      holy: {
        icon: "spell_holy_guardianspirit",
        spells: [],
      },
      shadow: {
        icon: "spell_shadow_shadowwordpain",
        spells: [],
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
        spells: [],
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
        spells: [],
      },
      fury: {
        icon: "ability_warrior_innerrage",
        spells: [],
      },
      protection: {
        icon: "ability_warrior_defensivestance",
        spells: [],
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
