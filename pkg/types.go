package pkg

type Request struct {
	InstanceID string `json:"instance_id"`
	ApiKey     string `json:"api_key,omitempty"`
	ApiSecret  string `json:"api_secret,omitempty"`
}
