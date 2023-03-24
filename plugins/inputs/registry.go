package inputs

import "sqser/models"

type InputCreator func() models.Input

var Inputs = map[string]InputCreator{}

func Add(name string, inputCreator InputCreator) {
	Inputs[name] = inputCreator
}
