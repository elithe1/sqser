package enrichers

import (
	"go.uber.org/zap"
	"sqser/models"
	"sqser/plugins/filters"
	"sqser/sqsercore"
	"strings"
)

type SubstringAllowList struct {
	Config *Config
	Log    *zap.SugaredLogger
}

type Config []string

func newSubstringAllowList(config *Config, logger *zap.SugaredLogger) *SubstringAllowList {
	return &SubstringAllowList{Config: config, Log: logger}
}

func (sal *SubstringAllowList) ApplyFilter(liad *models.ListItemsActionData) []*sqsercore.QueueCount {
	var ret []*sqsercore.QueueCount
	for _, item := range liad.FilteredQueueItems {
		for _, allowListStr := range *sal.Config {
			if strings.Contains(item.QueueName, allowListStr) {
				ret = append(ret, item)
				break
			}
		}
	}
	return ret
}

func init() {
	pluginName := "substringAllowList"
	filters.Add(pluginName, func(filterConfig *models.FilterConfig) models.Filter {
		logger := zap.S().With(zap.String("pluginType", "filter"), zap.String("pluginName", pluginName))
		c := Config{}
		err := filterConfig.Values.Unmarshal(&c)
		if err != nil {
			logger.Error("couldn't load filter config")
		}
		s := newSubstringAllowList(&c, logger)
		logger.Infof("loaded successfully")
		return s
	})
}
