package enrichers

import "sqser/models"

type EnricherCreator func(enricherConfig *models.EnricherConfig) models.Enricher

var Enrichers = map[string]EnricherCreator{}

func Add(name string, enricherCreator EnricherCreator) {
	Enrichers[name] = enricherCreator
}
