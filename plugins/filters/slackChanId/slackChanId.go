package enrichers

import (
	"go.uber.org/zap"
	"sqser/models"
	"sqser/plugins/filters"
	"sqser/sqsercore"
	"strings"
)

type SlackChanId struct {
	Config *Config
	Log    *zap.SugaredLogger
}

type Config []struct {
	Name string `yaml:"name"`
	Id   string `yaml:"id"`
}

func newSlackChanId(config *Config, logger *zap.SugaredLogger) *SlackChanId {
	return &SlackChanId{Config: config, Log: logger}
}

func (sal *SlackChanId) ApplyFilter(liad *models.ListItemsActionData) []*sqsercore.QueueCount {
	var envSubstring string
	for _, item := range *sal.Config {
		if item.Id == liad.SlackData.ChanId {
			envSubstring = item.Name
			break
		}
	}
	if envSubstring == "" {
		sal.Log.Error("couldn't find slack chan in provided config")
		return liad.FilteredQueueItems
	}
	var ret []*sqsercore.QueueCount
	for _, item := range liad.FilteredQueueItems {
		if strings.Contains(item.QueueName, envSubstring) {
			ret = append(ret, item)
		}

	}
	return ret
}

func init() {
	pluginName := "slackChannelId"
	filters.Add(pluginName, func(filterConfig *models.FilterConfig) models.Filter {
		logger := zap.S().With(zap.String("pluginType", "filter"), zap.String("pluginName", pluginName))
		c := Config{}
		err := filterConfig.Values.Unmarshal(&c)
		if err != nil {
			logger.Error("couldn't load filter config")
		}
		s := newSlackChanId(&c, logger)
		logger.Infof("loaded successfully")
		return s
	})
}
