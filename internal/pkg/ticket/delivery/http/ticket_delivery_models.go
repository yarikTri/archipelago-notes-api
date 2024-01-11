package delivery

type UpdateStateRequest struct {
	NewState string `json:"new_state" valid:"required"`
}

type FinalizeWritingRequest struct {
	State string `json:"state" valid:"required"`
}
