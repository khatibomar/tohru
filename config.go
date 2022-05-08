package tohru

type Config struct {
	clientID          string
	clientSecret      string
	backupLinksSecret string
}

func NewConfig(clientID, clientSecret, backupLinksSecret string) *Config {
	return &Config{clientID: clientID, clientSecret: clientSecret, backupLinksSecret: backupLinksSecret}
}
