package deconzauth

type AuthRequest struct {
	Username   string `json:"username,omitempty"`
	Devicetype string `json:"devicetype"`
}

type AuthSuccessResponse struct {
	Success struct {
		Username string `json:"username"`
	} `json:"success"`
}

type AuthErrorResponse struct {
	Error struct {
		Address     string `json:"address"`
		Description string `json:"description"`
		Type        int    `json:"type"`
	} `json:"error"`
}
