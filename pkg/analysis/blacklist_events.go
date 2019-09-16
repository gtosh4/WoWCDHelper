package analysis

/*
The purpose of these blacklists are to filter out sources of avoidable damage that is always a mistake to take.
*/

var ignoredAbilitiesByBoss = map[int64]map[int64]struct{}{
	2293: toSet(292565),
}

func toSet(ids ...int64) map[int64]struct{} {
	m := make(map[int64]struct{})
	for _, id := range ids {
		m[id] = struct{}{}
	}
	return m
}
