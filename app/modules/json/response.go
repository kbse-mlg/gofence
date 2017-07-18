package json

// JSONResponse model
type JSONResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
}

// MakeJSON create json data
func MakeJSON(success bool, data interface{}, err string) JSONResponse {
	return JSONResponse{
		Success: success,
		Data:    data,
		Error:   err,
	}
}

// OK json response
func OK() JSONResponse {
	return JSONResponse{true, nil, ""}
}

// ERROR json response
func ERROR(err string) JSONResponse {
	return JSONResponse{false, nil, err}
}

// DATA OK with response data
func DATA(data interface{}) JSONResponse {
	return JSONResponse{true, data, ""}
}
