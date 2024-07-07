package proofs

// SiblingType encodes, in the context of a Merkle tree, if
// a sibling is on the left or on the right; this information
// is used by the client and the server to properly communicate
// the proofs
type SiblingType int

const (
	NoSibling SiblingType = iota
	LeftSibling
	RightSibling
	Unknown
)

func (s SiblingType) String() string {
	switch s {
	case NoSibling:
		return "none"
	case LeftSibling:
		return "left"
	case RightSibling:
		return "right"
	default:
		return "unknown"
	}
}

func GetSiblingType(s string) SiblingType {
	switch s {
	case "none":
		return NoSibling
	case "left":
		return LeftSibling
	case "right":
		return RightSibling
	default:
		return Unknown
	}
}

type ProofPart struct {
	SiblingType SiblingType
	SiblingHash string
}
