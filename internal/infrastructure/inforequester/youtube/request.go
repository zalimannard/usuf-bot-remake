package youtubeinforequester

type Config interface {
	APIKey() string
}

type Requester struct {
	apiKey string
}

func New(config Config) *Requester {
	return &Requester{
		apiKey: config.APIKey(),
	}
}
