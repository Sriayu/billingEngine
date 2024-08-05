package response

import (
	"encoding/json"
	"log"
	"net/http"
)

type Meta struct {
	Page      int     `json:"page,omitempty"`
	TotalPage float64 `json:"totalPage"`
	Total     int     `json:"total"`
}

type ResponseMessage struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type IResponseClient interface {
	ResponseJSON(w http.ResponseWriter, message string, data interface{}, meta *Meta) error
	HttpError(w http.ResponseWriter, err error, errorCode int) error
}

type responseClient struct{}

func NewResponseClient() IResponseClient {
	return &responseClient{}
}

func (r *responseClient) ResponseJSON(
	w http.ResponseWriter,
	message string,
	data interface{},
	meta *Meta,
) error {
	var sendError error

	response := ResponseMessage{
		Message: message,
		Data:    data,
		Meta:    meta,
	}
	resp, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	_, sendError = w.Write(resp)
	return sendError
}

func (r *responseClient) HttpError(
	w http.ResponseWriter,
	err error,
	errorCode int,
) error {
	resp := make(map[string]string)
	resp["message"] = err.Error()
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorCode)

	_, sendError := w.Write(jsonResp)
	return sendError
}
