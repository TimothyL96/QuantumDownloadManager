package manager

// FlagState indicates some download states:
// 0 = unknown
// 1 = allowed
// 2 = not allowed
type FlagState int

const (
	unknown FlagState = iota
	allowed
	notAllowed
)

func (s FlagState) String() string {
	stateStr := ""

	switch s {
	case 0:
		stateStr = "unknown"
	case 1:
		stateStr = "allowed"
	case 2:
		stateStr = "not allowed"
	}

	return stateStr
}
