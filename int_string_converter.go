package copier

import (
	"fmt"
	"strconv"
)

var (
	// IntToStringConverter is converter for copier,
	// which realize int to string convertation.
	IntToStringConverter = Converter{
		Src: int(0),
		Dst: string(""),
		Convert: func(src interface{}) (interface{}, error) {
			value, ok := src.(int)
			if !ok {
				return nil, fmt.Errorf("value is not int: %T", src)
			}
			return strconv.Itoa(value), nil
		},
	}

	// StringToIntConverter is converter for copier,
	// which realize string to int convertation.
	StringToIntConverter = Converter{
		Src: string(""),
		Dst: int(0),
		Convert: func(src interface{}) (interface{}, error) {
			value, ok := src.(string)
			if !ok {
				return nil, fmt.Errorf("value is not string: %T", src)
			}
			return strconv.Atoi(value)
		},
	}
)
