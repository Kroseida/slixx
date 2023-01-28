package dto_test

import (
	"github.com/stretchr/testify/assert"
	"kroseida.org/slixx/pkg/dto"
	"testing"
)

type TestObject struct {
	Name        string
	Age         int
	Blub        map[string]string
	AnotherBlub map[string]string
}

type TestObjectDto struct {
	Name string
	Age  int
	Blub map[string]string
}

func Test_Map(t *testing.T) {
	testObject := TestObject{
		Name: "Test",
		Age:  42,
		Blub: map[string]string{
			"blub": "blub",
		},
		AnotherBlub: map[string]string{
			"anotherBlub": "anotherBlub",
		},
	}

	var testObjectDto TestObjectDto
	dto.Map(&testObject, &testObjectDto)

	assert.Equal(t, testObject.Name, testObjectDto.Name)
	assert.Equal(t, testObject.Age, testObjectDto.Age)
	assert.Equal(t, testObject.Blub, testObjectDto.Blub)
	assert.Equal(t, testObject.Blub["blub"], testObjectDto.Blub["blub"])
	assert.Equal(t, "", testObjectDto.Blub["anotherBlub"])
}
