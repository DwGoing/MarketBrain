package enum

import "errors"

type RechargeStatus uint8

const (
	RechargeStatus_UNPAID        RechargeStatus = 1
	RechargeStatus_PAID          RechargeStatus = 2
	RechargeStatus_CANCELLED     RechargeStatus = 3
	RechargeStatus_NOTIFY_FAILED RechargeStatus = 4
	RechargeStatus_NOTIFY_OK     RechargeStatus = 5
)

func (e RechargeStatus) String() string {
	switch e {
	case RechargeStatus_UNPAID:
		return "UNPAID"
	case RechargeStatus_PAID:
		return "PAID"
	case RechargeStatus_CANCELLED:
		return "CANCELLED"
	case RechargeStatus_NOTIFY_FAILED:
		return "NOTIFY_FAILED"
	case RechargeStatus_NOTIFY_OK:
		return "NOTIFY_OK"
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
	case "NOTIFY_FAILED":
		return RechargeStatus_NOTIFY_FAILED, nil
	case "NOTIFY_OK":
		return RechargeStatus_NOTIFY_OK, nil
	default:
		return 0, errors.New("unknown status")
	}
}
