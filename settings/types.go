package settings

// Settings contain all the possible settings for the application,
type Settings struct {
	Plugin  Plugin  `mapstructure:"plugin"`
	Webhook Webhook `mapstructure:"webhook"`
}

// Plugin settings
type Plugin struct {
	Enable   bool     `mapstructure:"enable"`
	Disabled []string `mapstructure:"disabled"`
}

// Webhook settings
type Webhook struct {
	Port int `mapstructure:"port"`
}
