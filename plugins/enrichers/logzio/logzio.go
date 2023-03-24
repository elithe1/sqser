package enrichers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"sqser/models"
	"sqser/plugins/enrichers"
	"sqser/sqsercore"
	"time"
)

type Logzio struct {
	Config      *Config
	TimeStamp   string
	SearchField string
	Log         *zap.SugaredLogger
}

type Config struct {
	Accounts []struct {
		Name string `yaml:"name"`
		Id   string `yaml:"id"`
	} `yaml:"accounts"`
	EnrichFields struct {
		TimeStampFieldName   string `yaml:"timeStamp"`
		SearchFieldFieldName string `yaml:"searchField"`
	} `yaml:"enrichFields"`
}

func newLogzio(config *Config, logger *zap.SugaredLogger) *Logzio {
	return &Logzio{Config: config, Log: logger}
}

func (l *Logzio) extractEnrichingFields(item *sqsercore.SQSerItem) error {
	messageBody := make(map[string]interface{})
	if err := json.Unmarshal([]byte(*item.MessageBody), &messageBody); err != nil {
		l.Log.Error("couldn't load enricher config")
		return nil
	}

	for item := range messageBody {
		if item == l.Config.EnrichFields.TimeStampFieldName {
			l.TimeStamp = fmt.Sprintf("%v", messageBody[item])
		} else if item == l.Config.EnrichFields.SearchFieldFieldName {
			l.SearchField = fmt.Sprintf("%v", messageBody[item])
		}
	}
	if l.TimeStamp == "" {
		return errors.New("couldn't extract timestamp field")
	}
	if l.TimeStamp == "" || l.SearchField == "" {
		return errors.New("couldn't extract search field field")
	}
	return nil
}

func (l *Logzio) Enrich(data *models.GetItemActionData) *models.EnrichLinkBlock {
	l.Log.Infof("Running Enricher")
	if data == nil {
		return nil
	}
	err := l.extractEnrichingFields(data.SQSerItem)
	if err != nil {
		l.Log.Error("Running Enricher errored %s", err)
		return nil
	}
	logzStr := "https://app.logz.io/#/dashboard/kibana/discover?_a=(columns:!(message),index:'logzioCustomerIndex*',interval:auto,query:(language:lucene,query:%%22%s%%22),sort:!(!('@timestamp',desc)))&_g=(filters:!(),refreshInterval:(pause:!t,value:0),time:(from:'%s',to:'%s'))&discoverTab=logz-logs-tab"
	defaultTimeLayout := "2006-01-02T15:04:05.000Z"
	t, err := time.Parse(defaultTimeLayout, l.TimeStamp)
	if err != nil {
		l.Log.Errorf("Running Enricher errored %s", err)
		return nil
	}

	fromTime := t.Add(time.Hour * time.Duration(-1))
	toTime := t.Add(time.Hour * time.Duration(1))

	formattedUrl := fmt.Sprintf(logzStr, l.SearchField, fromTime.Format(defaultTimeLayout), toTime.Format(defaultTimeLayout))
	accountId := l.resolveAccountId(data.EnvironmentName)
	accountQs := fmt.Sprintf("&accountIds=%s&switchToAccountId=%s", accountId, accountId)

	return &models.EnrichLinkBlock{Text: "Link to Logz.io", Link: formattedUrl + accountQs}
}

func (l *Logzio) resolveAccountId(environmentName string) string {

	for _, account := range l.Config.Accounts {
		if environmentName == account.Name {
			return account.Id
		}
	}
	return ""
}

func init() {
	pluginName := "logzio"
	enrichers.Add(pluginName, func(enricherConfig *models.EnricherConfig) models.Enricher {
		logger := zap.S().With(zap.String("pluginType", "enricher"), zap.String("pluginName", pluginName))

		c := Config{}
		err := enricherConfig.Values.Unmarshal(&c)
		if err != nil {
			logger.Error("couldn't load enricher config")
		}
		s := newLogzio(&c, logger)
		logger.Infof("loaded successfully")
		return s
	})
}
