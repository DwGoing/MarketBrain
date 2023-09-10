package enum

type ApiErrorType int64

const (
	ApiErrorType_Ok               ApiErrorType = 200
	ApiErrorType_RequestBindError ApiErrorType = 501
	ApiErrorType_ServiceError     ApiErrorType = 502
)

func (e ApiErrorType) Code() int64 {
	return int64(e)
}
