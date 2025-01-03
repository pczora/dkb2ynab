package config

type FromKeyringConfig struct {
	Key string `mapstructure:"key"`
}

type PasswordConfig struct {
	FromKeyring FromKeyringConfig `mapstructure:"fromKeyring"`
}

type CredentialConfig struct {
	Username string         `mapstructure:"username"`
	Password PasswordConfig `mapstructure:"password"`
}

type BankConfig struct {
	Name        string `mapstructure:"name"`
	Bank        string `mapstructure:"bank"`
	Credentials CredentialConfig
}
