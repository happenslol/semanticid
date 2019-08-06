package semanticid

import (
	"encoding/json"
)

var _ json.Marshaler = &SemanticID{}
var _ json.Unmarshaler = &SemanticID{}

// MarshalJSON implements the json.Marshaler interface for SemanticID
func (sid SemanticID) MarshalJSON() ([]byte, error) {
	if sid.IsNil() {
		return json.Marshal(nil)
	}

	str := sid.String()
	return json.Marshal(str)
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
