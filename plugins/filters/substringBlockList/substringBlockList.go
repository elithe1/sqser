package enrichers

import (
	"go.uber.org/zap"
	"sqser/models"
	"sqser/plugins/filters"
	"sqser/sqsercore"
	"strings"
)

type SubstringBlockList struct {
	Config *Config
	Log    *zap.SugaredLogger
}

type Config []string

func newSubstringBlockList(config *Config, logger *zap.SugaredLogger) *SubstringBlockList {
	return &SubstringBlockList{Config: config, Log: logger}
}

func (sbl *SubstringBlockList) ApplyFilter(liad *models.ListItemsActionData) []*sqsercore.QueueCount {
	var ret []*sqsercore.QueueCount
	for _, item := range liad.FilteredQueueItems {
		shouldFilterOut := false
		for _, allowListStr := range *sbl.Config {
			if strings.Contains(item.QueueName, allowListStr) {
				shouldFilterOut = true
				break
			}
		}
		if !shouldFilterOut {
			ret = append(ret, item)

		}
	}
	return ret
}

func init() {
	pluginName := "substringBlockList"

	filters.Add(pluginName, func(filterConfig *models.FilterConfig) models.Filter {
		logger := zap.S().With(zap.String("pluginType", "filter"), zap.String("pluginName", pluginName))

		c := Config{}
		err := filterConfig.Values.Unmarshal(&c)
		if err != nil {
			zap.S().Error("couldn't load filter config")
		}
		s := newSubstringBlockList(&c, logger)
		logger.Infof("loaded successfully")
		return s
	})
}
