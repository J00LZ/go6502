// Code generated by "enumer -type=DeviceType -json -yaml"; DO NOT EDIT.

//
package deviceinfo

import (
	"encoding/json"
	"fmt"
)

const _DeviceTypeName = "ROMRAMPPU"

var _DeviceTypeIndex = [...]uint8{0, 3, 6, 9}

func (i DeviceType) String() string {
	if i < 0 || i >= DeviceType(len(_DeviceTypeIndex)-1) {
		return fmt.Sprintf("DeviceType(%d)", i)
	}
	return _DeviceTypeName[_DeviceTypeIndex[i]:_DeviceTypeIndex[i+1]]
}

var _DeviceTypeValues = []DeviceType{0, 1, 2}

var _DeviceTypeNameToValueMap = map[string]DeviceType{
	_DeviceTypeName[0:3]: 0,
	_DeviceTypeName[3:6]: 1,
	_DeviceTypeName[6:9]: 2,
}

// DeviceTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func DeviceTypeString(s string) (DeviceType, error) {
	if val, ok := _DeviceTypeNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to DeviceType values", s)
}

// DeviceTypeValues returns all values of the enum
func DeviceTypeValues() []DeviceType {
	return _DeviceTypeValues
}

// IsADeviceType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i DeviceType) IsADeviceType() bool {
	for _, v := range _DeviceTypeValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for DeviceType
func (i DeviceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for DeviceType
func (i *DeviceType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("DeviceType should be a string, got %s", data)
	}

	var err error
	*i, err = DeviceTypeString(s)
	return err
}

// MarshalYAML implements a YAML Marshaler for DeviceType
func (i DeviceType) MarshalYAML() (interface{}, error) {
	return i.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for DeviceType
func (i *DeviceType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	var err error
	*i, err = DeviceTypeString(s)
	return err
}
