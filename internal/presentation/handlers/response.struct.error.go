package handlers

type ResponseError struct {
	StatusCode  int         `json:"-"`
	Code        string      `json:"code" example:"ResourceNotFound"`
	Message     string      `json:"message" example:"license is not found"`
	Target      string      `json:"target,omitempty" example:"subscription"`
	Details     interface{} `json:"details,omitempty" swaggerignore:"true"`
	BatchResult interface{} `json:"batch_result,omitempty" swaggerignore:"true"`
}
