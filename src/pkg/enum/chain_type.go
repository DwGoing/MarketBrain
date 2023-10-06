package enum

import "errors"

type ChainType uint8

const (
	ChainType_Tron ChainType = 1
)

func (e ChainType) String() string {
	switch e {
	case ChainType_Tron:
		return "TRON"
	default:
		return "UNKNOWN"
	}
}

func (e ChainType) Parse(str string) (ChainType, error) {
	switch str {
	case "TRON":
		return ChainType_Tron, nil
	default:
		return 0, errors.New("unknown type")
	}
}
