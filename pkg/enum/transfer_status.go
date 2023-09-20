package enum

import "errors"

type TransferStatus int8

const (
	TransferStatus_SUCCESS TransferStatus = 1
	TransferStatus_FAILED  TransferStatus = 2
)

func (e TransferStatus) String() string {
	switch e {
	case TransferStatus_SUCCESS:
		return "SUCCESS"
	case TransferStatus_FAILED:
		return "FAILED"
	default:
		return "UNKNOWN"
	}
}

func (e TransferStatus) Parse(str string) (TransferStatus, error) {
	switch str {
	case "SUCCESS":
		return TransferStatus_SUCCESS, nil
	case "FAILED":
		return TransferStatus_FAILED, nil
	default:
		return 0, errors.New("unknown status")
	}
}
