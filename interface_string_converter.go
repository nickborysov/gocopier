package copier

import "fmt"

var (
	// InterfaceToStringConverter is converter for copier,
	// which realize interface to string convertation.
	InterfaceToStringConverter = Converter{
		Src: nil,
		Dst: string(""),
		Convert: func(src interface{}) (interface{}, error) {
			value, ok := src.(string)
			if !ok {
				return nil, fmt.Errorf("value is not string: %T", src)
			}
			return value, nil
		},
	}

	// StringToInterfaceConverter is converter for copier,
	// which realize string to interface convertation.
	StringToInterfaceConverter = Converter{
		Src: string(""),
		Dst: nil,
		Convert: func(src interface{}) (interface{}, error) {
			value, ok := src.(string)
			if !ok {
				return nil, fmt.Errorf("value is not string: %T", src)
			}
			return value, nil
		},
	}
)
