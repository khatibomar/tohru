package angoslayer

type Config struct {
	clientID     string
	clientSecret string
}

func NewConfig(clientID, clientSecret string) *Config {
	return &Config{clientID: clientID, clientSecret: clientSecret}
}
