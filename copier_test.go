package copier

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Copy_InvalidDestination(t *testing.T) {
	err := Copy(nil, nil)
	assert.EqualError(t, err, ErrInvalidDestination.Error())
}

func Test_Copy_ErrInvalidSource(t *testing.T) {
	var dst string
	err := Copy(&dst, nil)
	assert.EqualError(t, err, ErrInvalidSource.Error())
}

func Test_Copy_StringToString(t *testing.T) {
	var dst string
	var src = "Lorem ipsum dolor sit amet"
	err := Copy(&dst, &src)
	assert.NoError(t, err)
	assert.Equal(t, "Lorem ipsum dolor sit amet", dst)
}

func Test_Copy_IntToString(t *testing.T) {
	var dst string
	var src = 100
	err := Copy(&dst, &src)
	assert.EqualError(t, err, "src and dst fields has different types: expected int, actual string")
	assert.Equal(t, "", dst)
}

func Test_Copy_IntToStringWithConverter(t *testing.T) {
	var dst string
	var src = 100
	err := Copy(&dst, &src, IntToStringConverter)
	assert.NoError(t, err)
	assert.Equal(t, "100", dst)
}

func Test_Copy_StringToIntWithConverter(t *testing.T) {
	var dst int
	var src = "100"
	err := Copy(&dst, &src, StringToIntConverter)
	assert.NoError(t, err)
	assert.Equal(t, 100, dst)
}

func Test_Copy_StringToIntWithConverterInvalidNumber(t *testing.T) {
	var dst int
	var src = "Lorem"
	err := Copy(&dst, &src, StringToIntConverter)
	assert.EqualError(t, err, `strconv.Atoi: parsing "Lorem": invalid syntax`)
	assert.Equal(t, 0, dst)
}

func Test_Copy_SliceStringToSliceString(t *testing.T) {
	var dst []string
	var src = []string{"Lorem", "ipsum", "dolor", "sit", "amet"}
	err := Copy(&dst, &src, StringToIntConverter)
	assert.NoError(t, err)
	assert.Equal(t, []string{"Lorem", "ipsum", "dolor", "sit", "amet"}, dst)
}
func Test_Copy_EmptySliceToSliceString(t *testing.T) {
	var dst []string
	var src = []string{}
	err := Copy(&dst, &src, StringToIntConverter)
	assert.NoError(t, err)
	assert.Equal(t, []string(nil), dst)
}

func Test_Copy_SliceStringToSliceInt(t *testing.T) {
	var dst []int
	var src = []string{"Lorem", "ipsum", "dolor", "sit", "amet"}
	err := Copy(&dst, &src)
	assert.EqualError(t, err, "src and dst fields has different types: expected string, actual int")
	assert.Equal(t, []int(nil), dst)
}

func Test_Copy_SliceStringToStringError(t *testing.T) {
	var dst string
	var src = []string{"Lorem", "ipsum", "dolor", "sit", "amet"}
	err := Copy(&dst, &src)
	assert.EqualError(t, err, "src and dst fields has different types: expected slice, actual string")
	assert.Equal(t, string(""), dst)
}

func Test_Copy_SliceStringToSliceIntWithConverter(t *testing.T) {
	var dst []int
	var src = []string{"100", "200", "300", "400", "500"}
	err := Copy(&dst, &src, StringToIntConverter)
	assert.NoError(t, err)
	assert.Equal(t, []int{100, 200, 300, 400, 500}, dst)
}

func Test_Copy_SliceStringToSliceIntWithConverterInvalidNumber(t *testing.T) {
	var dst []int
	var src = []string{"Lorem", "ipsum", "dolor", "sit", "amet"}
	err := Copy(&dst, &src, StringToIntConverter)
	assert.EqualError(t, err, `strconv.Atoi: parsing "Lorem": invalid syntax`)
	assert.Equal(t, []int(nil), dst)
}

func Test_Copy_SliceStringToArrayString(t *testing.T) {
	var dst [5]int
	var src = []string{"100", "200", "300", "400", "500"}
	err := Copy(&dst, &src, StringToIntConverter)
	assert.NoError(t, err)
	assert.Equal(t, [5]int{100, 200, 300, 400, 500}, dst)
}

func Test_Copy_MapStringStringToMapStringInt(t *testing.T) {
	var dst map[string]int
	var src = map[string]string{"Lorem": "100", "ipsum": "200", "dolor": "300", "sit": "400", "amet": "500"}
	err := Copy(&dst, &src)
	assert.EqualError(t, err, "src and dst fields has different types: expected string, actual int")
	expRes := map[string]int{}
	assert.Equal(t, expRes, dst)
}

func Test_Copy_MapStringStringToMapStringIntWithConverter(t *testing.T) {
	var dst map[string]int
	var src = map[string]string{"Lorem": "100", "ipsum": "200", "dolor": "300", "sit": "400", "amet": "500"}
	err := Copy(&dst, &src, StringToIntConverter)
	assert.NoError(t, err)
	expRes := map[string]int{"Lorem": 100, "ipsum": 200, "dolor": 300, "sit": 400, "amet": 500}
	assert.Equal(t, expRes, dst)
}

func Test_Copy_MapStringStingToMapIntInt(t *testing.T) {
	var dst map[int]int
	var src = map[string]string{"1": "100", "4": "200", "20": "300", "30": "400", "100": "500"}
	err := Copy(&dst, &src)
	assert.EqualError(t, err, "src and dst fields has different types: expected string, actual int")
	expRes := map[int]int{}
	assert.Equal(t, expRes, dst)
}

func Test_Copy_MapStringStingToMapIntIntWithConverter(t *testing.T) {
	var dst map[int]int
	var src = map[string]string{"1": "100", "4": "200", "20": "300", "30": "400", "100": "500"}
	err := Copy(&dst, &src, StringToIntConverter)
	assert.NoError(t, err)
	expRes := map[int]int{1: 100, 4: 200, 20: 300, 30: 400, 100: 500}
	assert.Equal(t, expRes, dst)
}

func Test_Copy_StructAToStructBWithConverter(t *testing.T) {
	type A struct {
		ID   string
		Name string
		Type string
	}
	type B struct {
		ID   int
		Name string
	}
	var dst B
	var src = A{
		ID:   "100",
		Name: "Jonh",
		Type: "skipped",
	}
	expResult := B{
		ID:   100,
		Name: "Jonh",
	}
	err := Copy(&dst, &src, StringToIntConverter)
	assert.NoError(t, err)
	assert.Equal(t, expResult, dst)
}

func Test_Copy_StructWithSliceToStructWithSliceWithConverter(t *testing.T) {
	type A struct {
		Name   string
		Values []string
	}
	type B struct {
		Name   string
		Values []int
	}
	var dst B
	var src = A{
		Name:   "Jonh",
		Values: []string{"100", "200", "300"},
	}
	expResult := B{
		Name:   "Jonh",
		Values: []int{100, 200, 300},
	}
	err := Copy(&dst, &src, StringToIntConverter)
	assert.NoError(t, err)
	assert.Equal(t, expResult, dst)
}

func Test_Copy_StructWithMapToStructWithMapWithConverter(t *testing.T) {
	type A struct {
		Name   string
		Values map[string]string
	}
	type B struct {
		Name   string
		Values map[string]int
	}
	var dst B
	var src = A{
		Name:   "Jonh",
		Values: map[string]string{"Lorem": "100", "ipsum": "200", "dolor": "300", "sit": "400", "amet": "500"},
	}
	expResult := B{
		Name:   "Jonh",
		Values: map[string]int{"Lorem": 100, "ipsum": 200, "dolor": 300, "sit": 400, "amet": 500},
	}
	err := Copy(&dst, &src, StringToIntConverter)
	assert.NoError(t, err)
	assert.Equal(t, expResult, dst)
}

func Test_Copy_SliceOfStructWithMapToSliceOfStructWithMapWithConverter(t *testing.T) {
	type A struct {
		Name   string
		Values map[string]string
	}
	type B struct {
		Name   string
		Values map[string]int
	}
	var dst []B
	var src = []A{
		{
			Name:   "Jonh",
			Values: map[string]string{"Lorem": "100", "ipsum": "200", "dolor": "300", "sit": "400", "amet": "500"},
		},
	}
	expResult := []B{
		{
			Name:   "Jonh",
			Values: map[string]int{"Lorem": 100, "ipsum": 200, "dolor": 300, "sit": 400, "amet": 500},
		},
	}
	err := Copy(&dst, &src, StringToIntConverter)
	assert.NoError(t, err)
	assert.Equal(t, expResult, dst)
}

func Test_Copy_StructWithPointersToStructWithConverter(t *testing.T) {
	type A struct {
		Name  *string
		Value *string
	}
	type B struct {
		Name  *string
		Value *int
	}
	var dst B
	name := "Jonh"
	valueString := "100"
	var src = A{
		Name:  &name,
		Value: &valueString,
	}
	valueInt := 100
	var expResult = B{
		Name:  &name,
		Value: &valueInt,
	}
	err := Copy(&dst, &src, StringToIntConverter)
	assert.NoError(t, err)
	assert.Equal(t, expResult, dst)
}

func Test_Copy_StructWithPointersToStructWithoutPointerWithConverter(t *testing.T) {
	type A struct {
		Name  *string
		Value *string
	}
	type B struct {
		Name  string
		Value int
	}
	var dst B
	name := "Jonh"
	valueString := "100"
	var src = A{
		Name:  &name,
		Value: &valueString,
	}
	var expResult = B{
		Name:  "Jonh",
		Value: 100,
	}
	err := Copy(&dst, &src, StringToIntConverter)
	assert.NoError(t, err)
	assert.Equal(t, expResult, dst)
}
