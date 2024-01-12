package handler

type RackRequest struct {
	Name string `json:"name" form:"name"`
}

type RackResponse struct {
	ID   uint   `json:"rack_id"`
	Name string `json:"name"`
}
