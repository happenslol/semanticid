package semanticid

import (
	"encoding/json"
)

var _ json.Marshaler = &SemanticID{}
var _ json.Unmarshaler = &SemanticID{}

// MarshalJSON implements the json.Marshaler interface for SemanticID
func (sid SemanticID) MarshalJSON() ([]byte, error) {
	str := sid.String()
	return []byte(str), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for SemanticID
func (sid *SemanticID) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}

	parsed, err := FromString(str)
	if err != nil {
		return err
	}

	*sid = parsed
	return nil
}
