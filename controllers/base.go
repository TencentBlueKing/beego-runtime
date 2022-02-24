package controllers

type BaseResponse struct {
	Result  bool   `json:"result"`
	Message string `json:"message"`
}
