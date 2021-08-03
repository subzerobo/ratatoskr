package handlers

type HealthCheckResponse struct {
	Status       string  `json:"status"`
	Container    string  `json:"container"`
	GitCommit    string  `json:"git_commit"`
	Version      string  `json:"go_version"`
	Uptime       string  `json:"kernel_uptime"`
	BinaryUptime string  `json:"binary_uptime"`
	BuildTime    string  `json:"build_time"`
	LAOne        float64 `json:"load_average_one"`
	LAFive       float64 `json:"load_average_five"`
	LAFifteen    float64 `json:"load_average_fifteen"`
}
