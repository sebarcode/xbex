package config

type ModConfig struct {
}

var (
	cfg *ModConfig
)

func Config() *ModConfig {
	if cfg == nil {
		cfg = &ModConfig{}
	}
	return cfg
}
