package copier

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	// ObjectIDToStringConverter is converter for copier,
	// which realize mongo ObjectID to string convertation.
	ObjectIDToStringConverter = &Converter{
		Src: primitive.ObjectID{},
		Dst: string(""),
		Convert: func(src interface{}) (interface{}, error) {
			value, ok := src.(primitive.ObjectID)
			if !ok {
				return nil, fmt.Errorf("value is not primitive.ObjectID: %T", src)
			}
			return value.Hex(), nil
		},
	}

	// StringToObjectIDConverter is converter for copier,
	// which realize string to mongo ObjectID convertation.
	StringToObjectIDConverter = &Converter{
		Src: string(""),
		Dst: primitive.ObjectID{},
		Convert: func(src interface{}) (interface{}, error) {
			value, ok := src.(string)
			if !ok {
				return nil, fmt.Errorf("value is not string: %T", src)
			}
			if value == "" {
				return primitive.NilObjectID, nil
			}
			id, err := primitive.ObjectIDFromHex(value)
			if err != nil {
				return nil, fmt.Errorf("failed get ObjectID from string: %v", err)
			}
			return id, nil
		},
	}
)
