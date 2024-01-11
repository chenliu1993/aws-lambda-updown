package pkg

type Request struct {
	InstanceID string `json:"instance_id"`
	Timezone   string `json:"timezone,omitempty"`
	StartHour  string `json:"start_hour,omitempty"`
	StopHour   string `json:"stop_hour,omitempty"`
	ApiKey     string `json:"api_key,omitempty"`
	ApiSecret  string `json:"api_secret,omitempty"`
}
