package enum

type RechargeStatus uint8

const (
	RechargeStatus_UNPAID RechargeStatus = 1
	RechargeStatus_PAID   RechargeStatus = 2
)

func (e RechargeStatus) String() string {
	switch e {
	case RechargeStatus_UNPAID:
		return "UNPAID"
	case RechargeStatus_PAID:
		return "PAID"
	default:
		return "UNKNOWN"
	}
}
