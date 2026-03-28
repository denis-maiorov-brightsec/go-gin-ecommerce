package api

type MessageResponse struct {
	Message string `json:"message" example:"The unversioned root route is deprecated. Migrate to /v1/health."`
}

type StatusResponse struct {
	Status string `json:"status" example:"ok"`
}
