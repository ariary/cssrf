package cssrf

type Config struct {
	Length             int
	Port               string
	ExternalUrl        string
	Charset            string
	BackgroundTemplate string
	ImportTemplate     string
	Channels           []chan string
}
