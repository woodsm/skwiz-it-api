package config

type S3 struct {
	AccessKey    string `json:"key"`
	AccessSecret string `json:"secret"`
	Bucket       string `json:"bucket"`
	Region       string `json:"region"`
}
