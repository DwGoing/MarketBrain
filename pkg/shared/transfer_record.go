package shared

// @title 充值状态
type TransferStatus int8

const (
	TransferStatus_SUCCESS TransferStatus = 1
	TransferStatus_FAILED  TransferStatus = 2
)

func (e TransferStatus) ToString() string {
	switch e {
	case TransferStatus_SUCCESS:
		return "SUCCESS"
	case TransferStatus_FAILED:
		return "FAILED"
	default:
		return "UNKNOWN"
	}
}

type TransferRecord struct {
	Record
	FromIndex   uint32         `gorm:"column:FROM_INDEX"`
	FromAddress string         `gorm:"column:FROM_ADDRESS"`
	To          string         `gorm:"column:TO"`
	Token       string         `gorm:"column:TOKEN"`
	Amount      string         `gorm:"column:AMOUNT"`
	Status      TransferStatus `gorm:"column:STATUS"`
	Error       string         `gorm:"column:ERROR"`
}
