package components

type Config struct {
	DryRun         bool
	Excludes       map[string]bool
	EnableFeatures map[string]bool
	Demo           bool
}

func NewConfig(dryRun bool, excludes []string, enables []string, demo bool) Config {
	return Config{
		DryRun:         dryRun,
		Excludes:       getBoolMapFromArray(excludes),
		EnableFeatures: getBoolMapFromArray(enables),
		Demo:           demo,
	}
}

func getBoolMapFromArray(input []string) map[string]bool {
	result := map[string]bool{}

	for _, key := range input {
		result[key] = true
	}

	return result
}
