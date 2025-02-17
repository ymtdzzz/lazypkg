package components

type Config struct {
	DryRun   bool
	Excludes map[string]bool
}

func NewConfig(dryRun bool, excludes []string) Config {
	return Config{
		DryRun:   dryRun,
		Excludes: getExcludeMap(excludes),
	}
}

func getExcludeMap(excludes []string) map[string]bool {
	result := map[string]bool{}

	for _, exclude := range excludes {
		result[exclude] = true
	}

	return result
}
