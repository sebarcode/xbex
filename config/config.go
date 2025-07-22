package config

type ModConfig struct {
	ShaKey  string `yaml:"sha_key" json:"sha_key"`
	JwtSalt string `yaml:"jwt_salt" json:"jwt_salt"`
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

func SetConfig(c *ModConfig) {
	if c == nil {
		return
	}
	cfg = c
}
