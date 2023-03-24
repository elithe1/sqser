package filters

import "sqser/models"

type FilterCreator func(filterConfig *models.FilterConfig) models.Filter

var Filters = map[string]FilterCreator{}

func Add(name string, filterCreator FilterCreator) {
	Filters[name] = filterCreator
}
