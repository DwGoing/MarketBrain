package enum

import "errors"

type ChainType uint8

const (
	ChainType_TRON ChainType = 1
)

func (e ChainType) String() string {
	switch e {
	case ChainType_TRON:
		return "TRON"
	default:
		return "UNKNOWN"
	}
}

func (e ChainType) Parse(str string) (ChainType, error) {
	switch str {
	case "TRON":
		return ChainType_TRON, nil
	default:
		return 0, errors.New("unknown type")
	}
}
