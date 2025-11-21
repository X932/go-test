package response

type Body struct {
	Message string `json:"message"`
	Payload any    `json:"payload,omitempty"`
}
