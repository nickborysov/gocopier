package copier

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_Converter_SrcType(t *testing.T) {
	converter := Converter{Src: string("")}
	typ := converter.SrcType()
	assert.Equal(t, reflect.String, typ.Kind())
}

func Test_Converter_SrcType_nil(t *testing.T) {
	converter := Converter{}
	typ := converter.SrcType()
	assert.Equal(t, interfaceType, typ)
}

func Test_Converter_DstType(t *testing.T) {
	converter := Converter{Dst: string("")}
	typ := converter.DstType()
	assert.Equal(t, reflect.String, typ.Kind())
}

func Test_Converter_DstType_nil(t *testing.T) {
	converter := Converter{}
	typ := converter.DstType()
	assert.Equal(t, interfaceType, typ)
}

func Test_IntToStringConverter_ConvertInt(t *testing.T) {
	result, err := IntToStringConverter.Convert(100)
	assert.NoError(t, err)
	assert.Equal(t, "100", result)
}

func Test_IntToStringConverter_ConvertInvalid(t *testing.T) {
	result, err := IntToStringConverter.Convert("100")
	assert.Equal(t, errors.New("value is not int: string"), err)
	assert.Equal(t, nil, result)
}

func Test_StringToIntConverter_ConvertString(t *testing.T) {
	result, err := StringToIntConverter.Convert("100")
	assert.NoError(t, err)
	assert.Equal(t, 100, result)
}

func Test_StringToIntConverter_ConvertInvalid(t *testing.T) {
	result, err := StringToIntConverter.Convert(100)
	assert.Equal(t, errors.New("value is not string: int"), err)
	assert.Equal(t, nil, result)
}

func Test_InterfaceToStringConverter_ConvertString(t *testing.T) {
	result, err := InterfaceToStringConverter.Convert(interface{}("Lorem"))
	assert.NoError(t, err)
	assert.Equal(t, "Lorem", result)
}

func Test_InterfaceToStringConverter_ConvertInvalid(t *testing.T) {
	result, err := InterfaceToStringConverter.Convert(nil)
	assert.Equal(t, errors.New("value is not string: <nil>"), err)
	assert.Equal(t, nil, result)
}

func Test_StringToInterfaceConverter_ConvertString(t *testing.T) {
	result, err := StringToInterfaceConverter.Convert("Lorem")
	assert.NoError(t, err)
	assert.Equal(t, "Lorem", result)
}

func Test_StringToInterfaceConverter_ConvertInvalid(t *testing.T) {
	result, err := StringToInterfaceConverter.Convert(100)
	assert.Equal(t, errors.New("value is not string: int"), err)
	assert.Equal(t, nil, result)
}

func Test_ObjectIDToStringConverter_ConvertObjectID(t *testing.T) {
	id := primitive.NewObjectID()
	result, err := ObjectIDToStringConverter.Convert(id)
	assert.NoError(t, err)
	assert.Equal(t, id.Hex(), result)
}

func Test_ObjectIDToStringConverter_ConvertInvalid(t *testing.T) {
	result, err := ObjectIDToStringConverter.Convert(100)
	assert.Equal(t, errors.New("value is not primitive.ObjectID: int"), err)
	assert.Equal(t, nil, result)
}

func Test_StringToObjectIDConverter_ConvertString(t *testing.T) {
	id := primitive.NewObjectID()
	result, err := StringToObjectIDConverter.Convert(id.Hex())
	assert.NoError(t, err)
	assert.Equal(t, id, result)
}

func Test_StringToObjectIDConverter_ConvertEmptyString(t *testing.T) {
	id := primitive.NilObjectID
	result, err := StringToObjectIDConverter.Convert("")
	assert.NoError(t, err)
	assert.Equal(t, id, result)
}

func Test_StringToObjectIDConverter_ConvertInvalidString(t *testing.T) {
	result, err := StringToObjectIDConverter.Convert("invalid")
	assert.Equal(t, errors.New("failed get ObjectID from string: encoding/hex: invalid byte: U+0069 'i'"), err)
	assert.Equal(t, nil, result)
}

func Test_StringToObjectIDConverter_ConvertInvalid(t *testing.T) {
	result, err := StringToObjectIDConverter.Convert(100)
	assert.Equal(t, errors.New("value is not string: int"), err)
	assert.Equal(t, nil, result)
}
