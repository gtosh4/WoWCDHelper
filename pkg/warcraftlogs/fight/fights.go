package fight

var FightNotFound = Fight{}

func (f *Fights) Fight(num int64) (fight *Fight, friendlies []Friendly, enemies []Enemy, phases *Phases) {
	fight = &f.Fights[num-1]
	if fight.Id != num {
		for _, fi := range f.Fights {
			if fi.Id == num {
				fightSnap := fi
				fight = &fightSnap
				break
			}
		}
	}
	if fight.Id != num {
		fight = &FightNotFound
		return
	}

	for _, friend := range f.Friendlies {
		if friend.Type != "NPC" {
			for _, enc := range friend.Fights {
				if enc.Id == fight.Id {
					friendlies = append(friendlies, friend)
					break
				}
			}
		}
	}

	for _, enemy := range f.Enemies {
		for _, enc := range enemy.Fights {
			if enc.Id == fight.Id {
				enemies = append(enemies, enemy)
				break
			}
		}
	}

	for _, phase := range f.Phases {
		if phase.Boss == fight.Boss {
			phaseSnap := phase
			phases = &phaseSnap
			break
		}
	}

	return
}
