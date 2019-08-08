export const classes = {
  deathknight: {
    colour: {r: 196, g: 31, b: 59},
    specs: {
    },
  },
  demonhunter: {
    colour: {r: 163, g: 48, b: 201},
    specs: {
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
      resto: {
        icon: "spell_nature_healingtouch",
      },
    },
  },
  hunter: {
    colour: {r: 171, g: 212, b: 115},
    specs: {
    },
  },
  mage: {
    colour: {r: 64, g: 199, b: 235},
    specs: {
    },
  },
  monk: {
    colour: {r: 0, g: 255, b: 150},
    specs: {
    },
  },
  paladin: {
    colour: {r: 245, g: 140, b: 186},
    specs: {
      holy: {
        icon: "spell_holy_holybolt",
      },
    },
  },
  priest: {
    colour: {r: 255, g: 255, b: 255},
    specs: {
    },
  },
  rogue: {
    colour: {r: 255, g: 245, b: 105},
    specs: {
    },
  },
  shaman: {
    colour: {r: 0, g: 112, b: 222},
    specs: {
    },
  },
  warlock: {
    colour: {r: 135, g: 135, b: 237},
    specs: {
    },
  },
  warrior: {
    colour: {r: 199, g: 156, b: 110},
    specs: {
    },
  },
}

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

export function spec(className, specName) {
  return (classes[className] || {}).specs[specName] || {}
}
