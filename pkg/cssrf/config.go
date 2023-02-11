package cssrf

type Config struct {
	Length             int
	Port               string
	ExternalUrl        string
	Charset            string
	Prefix             string
	BackgroundTemplate string
	ImportTemplate     string
	Elt                string
	Channels           []chan string
}
