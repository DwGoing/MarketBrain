package enum

import "errors"

type RechargeStatus uint8

const (
	RechargeStatus_UNPAID    RechargeStatus = 1
	RechargeStatus_PAID      RechargeStatus = 2
	RechargeStatus_CANCELLED RechargeStatus = 3
)

func (e RechargeStatus) String() string {
	switch e {
	case RechargeStatus_UNPAID:
		return "UNPAID"
	case RechargeStatus_PAID:
		return "PAID"
	case RechargeStatus_CANCELLED:
		return "CANCELLED"
	default:
		return "UNKNOWN"
	}
}

func (e RechargeStatus) Parse(str string) (RechargeStatus, error) {
	switch str {
	case "UNPAID":
		return RechargeStatus_UNPAID, nil
	case "PAID":
		return RechargeStatus_PAID, nil
	case "CANCELLED":
		return RechargeStatus_CANCELLED, nil
	default:
		return 0, errors.New("unknown status")
	}
}
