package eventbus

func opt[T any](v []T, def ...T) T {
	if len(v) == 0 {
		var altDef T
		return opt(def, altDef)
	}
	return v[0]
}
