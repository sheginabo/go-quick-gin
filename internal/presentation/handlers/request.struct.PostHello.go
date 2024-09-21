package handlers

type PostHelloRequest struct {
	Message string `json:"message" validate:"required"`
}
