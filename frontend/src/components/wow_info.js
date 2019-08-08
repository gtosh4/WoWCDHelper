export const classes = {
  deathknight: {
    colour: {r: 196, g: 31, b: 59},
    specs: {
      blood: {
        icon: "spell_deathknight_bloodpresence",
      },
      frost: {
        icon: "spell_deathknight_frostpresence.gif",
      },
      unholy: {
        icon: "spell_deathknight_unholypresence.gif",
      },
    },
  },
  demonhunter: {
    colour: {r: 163, g: 48, b: 201},
    specs: {
      havoc: {
        icon: "ability_demonhunter_specdps",
      },
      vengeance: {
        icon: "ability_demonhunter_spectank",
      },
    },
  },
  druid: {
    colour: {r: 255, g: 125, b: 10},
    specs: {
      balance: {
        icon: "spell_nature_starfall",
      },
      feral: {
        icon: "ability_druid_catform",
      },
      guardian: {
        icon: "ability_racial_bearform",
      },
      restoration: {
        icon: "spell_nature_healingtouch",
      },
    },
  },
  hunter: {
    colour: {r: 171, g: 212, b: 115},
    specs: {
      beastmastery: {
        icon: "ability_hunter_bestialdiscipline",
      },
      marksmanship: {
        icon: "ability_hunter_focusedaim",
      },
      survival: {
        icon: "ability_hunter_camouflage",
      },
    },
  },
  mage: {
    colour: {r: 64, g: 199, b: 235},
    specs: {
      arcane: {
        icon: "spell_holy_magicalsentry",
      },
      fire: {
        icon: "spell_fire_firebolt02",
      },
      frost: {
        icon: "spell_frost_frostbolt02",
      },
    },
  },
  monk: {
    colour: {r: 0, g: 255, b: 150},
    specs: {
      brewmaster: {
        icon: "monk_stance_drunkenox",
      },
      mistweaver: {
        icon: "monk_stance_wiseserpent",
      },
      windwalker: {
        icon: "",
      },
    },
  },
  paladin: {
    colour: {r: 245, g: 140, b: 186},
    specs: {
      holy: {
        icon: "spell_holy_holybolt",
      },
      protection: {
        icon: "ability_paladin_shieldofthetemplar",
      },
      retribution: {
        icon: "spell_holy_auraoflight",
      },
    },
  },
  priest: {
    colour: {r: 255, g: 255, b: 255},
    specs: {
      discipline: {
        icon: "spell_holy_powerwordshield",
      },
      holy: {
        icon: "spell_holy_guardianspirit",
      },
      shadow: {
        icon: "spell_shadow_shadowwordpain",
      },
    },
  },
  rogue: {
    colour: {r: 255, g: 245, b: 105},
    specs: {
      assassination: {
        icon: "ability_rogue_eviscerate",
      },
      outlaw: {
        icon: "ability_backstab",
      },
      subtlety: {
        icon: "ability_stealth",
      },
    },
  },
  shaman: {
    colour: {r: 0, g: 112, b: 222},
    specs: {
      elemental: {
        icon: "spell_nature_lightning",
      },
      enhancement: {
        icon: "spell_nature_lightningshield",
      },
      restoration: {
        icon: "spell_nature_magicimmunity",
      },
    },
  },
  warlock: {
    colour: {r: 135, g: 135, b: 237},
    specs: {
      affliction: {
        icon: "spell_shadow_deathcoil",
      },
      demonology: {
        icon: "spell_shadow_metamorphosis",
      },
      destruction: {
        icon: "spell_shadow_rainoffire",
      },
    },
  },
  warrior: {
    colour: {r: 199, g: 156, b: 110},
    specs: {
      arms: {
        icon: "ability_warrior_savageblow",
      },
      fury: {
        icon: "ability_warrior_innerrage",
      },
      protection: {
        icon: "ability_warrior_defensivestance",
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
