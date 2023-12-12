package utils

// Use https://github.com/omniti-labs/jsend for response format

const (
	StatusError   = "error"
	StatusFail    = "fail"
	StatusSuccess = "success"
)

type Body struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Code    int         `json:"code,omitempty"`
}

func FormatSuccess(data interface{}) *Body {
	return &Body{
		Status: StatusSuccess,
		Data:   data,
	}
}

func FormatFail(data interface{}) *Body {
	return &Body{
		Status: StatusFail,
		Data:   data,
	}
}

func FormatError(message string) *Body {
	return &Body{
		Status:  StatusError,
		Message: message,
	}
}
