package deconzauth

type AuthRequest struct {
	Username   string `json:"username,omitempty"`
	Devicetype string `json:"devicetype"`
}
type AuthSuccess struct {
	Username string `json:"username"`
}

type AuthSuccessResponse struct {
	Success AuthSuccess `json:"success"`
}

type AuthError struct {
	Address     string `json:"address"`
	Description string `json:"description"`
	Type        int    `json:"type"`
}

type AuthErrorResponse struct {
	Error AuthError `json:"error"`
}
