package delivery

type UpdateStateRequest struct {
	NewState string `json:"new_state" valid:"required"`
}
