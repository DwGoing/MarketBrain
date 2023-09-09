package enum

import "errors"

type ServiceType uint8

const (
	ServiceType_FUNDS ServiceType = 1
	ServiceType_DATA  ServiceType = 2
)

func (e ServiceType) String() string {
	switch e {
	case ServiceType_FUNDS:
		return "FUNDS"
	case ServiceType_DATA:
		return "DATA"
	default:
		return "UNKNOWN"
	}
}

func (e ServiceType) Parse(str string) (ServiceType, error) {
	switch str {
	case "FUNDS":
		return ServiceType_FUNDS, nil
	case "DATA":
		return ServiceType_DATA, nil
	default:
		return 0, errors.New("unknown type")
	}
}
