package outputs

import "sqser/models"

type OutputCreator func() models.Output

var Outputs = map[string]OutputCreator{}

func Add(name string, outputCreator OutputCreator) {
	Outputs[name] = outputCreator
}
