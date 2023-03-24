package app

import (
	"go.uber.org/zap"
	"net/http"
	"sqser/config"
	"sqser/models"
	"sqser/plugins/enrichers"
	"sqser/plugins/filters"
	"sqser/plugins/inputs"
	"sqser/plugins/outputs"
	"sqser/sqsercore"
	"strings"
)

type App struct {
	config.Config
	Inputs    []*models.RunningInput
	Enrichers []*models.RunningEnricher
	Filters   []*models.RunningFilter
	Outputs   []*models.RunningOutput
}

func (app *App) addInput(name string, async bool) error {
	inputCreator := inputs.Inputs[name]
	input := inputCreator()
	rp := models.NewRunningInput(input, &models.InputConfig{Name: name, Async: async})
	app.Inputs = append(app.Inputs, rp)
	return nil
}

func (app *App) addEnricher(name string, values models.RawMessage) error {
	enricherCreator := enrichers.Enrichers[name]
	config := &models.EnricherConfig{Name: name, Values: values}
	enricher := enricherCreator(config)
	rp := models.NewRunningEnricher(enricher, config)
	app.Enrichers = append(app.Enrichers, rp)
	return nil
}
func (app *App) addFilter(name string, values models.RawMessage) error {
	filterCreator := filters.Filters[name]
	config := &models.FilterConfig{Name: name, Values: values}
	filter := filterCreator(config)
	rf := models.NewRunningFilter(filter, config)
	app.Filters = append(app.Filters, rf)
	return nil
}
func (app *App) addOutput(name string) error {
	outputCreator := outputs.Outputs[name]
	output := outputCreator()
	ro := models.NewRunningOutput(output, &models.OutputConfig{})
	app.Outputs = append(app.Outputs, ro)
	return nil
}

func NewApp() *App {
	c := config.NewConfig()
	a := &App{
		Config:    *c,
		Inputs:    make([]*models.RunningInput, 0),
		Enrichers: make([]*models.RunningEnricher, 0),
	}
	for _, inputData := range c.Config.Inputs {
		a.addInput(inputData.Name, inputData.Async)
	}
	for _, enrichersData := range c.Config.Enrichers {
		a.addEnricher(enrichersData.Name, enrichersData.Values)
	}
	for _, filtersData := range c.Config.Filters {
		a.addFilter(filtersData.Name, filtersData.Values)
	}
	for _, outputData := range c.Config.Outputs {
		a.addOutput(outputData.Name)
	}
	return a
}

func (app *App) GetItem(w http.ResponseWriter, req *http.Request) {
	zap.S().Infof("Running GET command")
	input := app.findInput(req.URL.Query().Get("input"))
	if input == nil {
		zap.S().Infof("couldn't find input by input name")
		w.Write([]byte("couldn't find input by input name"))
		return
	}
	ss := sqsercore.NewSQSService(app.Config.Config.Substrings.Dlq)

	zap.S().Infof("Running input: %s", input.Config.Name)
	ad := input.Input.InvokeGetItem(w, req)
	if ad == nil {
		return
	}
	if input.Config.Async {
		go app.processGetItem(ss, ad)
		return
	}
	// slack is the only plugin for now and it's async so didn't really care too much here.
	app.processGetItem(ss, ad)
}
func (app *App) ListItems(w http.ResponseWriter, req *http.Request) {
	zap.S().Infof("Running ListItems command")
	input := app.findInput(req.URL.Query().Get("input"))
	if input == nil {
		zap.S().Infof("couldn't find input by input name")
		w.Write([]byte("couldn't find input by input name"))
		return
	}
	ss := sqsercore.NewSQSService(app.Config.Config.Substrings.Dlq)

	zap.S().Infof("Running input: %s", input.Config.Name)
	ad := input.Input.InvokeListItems(w, req)

	if input.Config.Async {
		go app.processListItems(ss, ad)
		return
	}
	// slack is the only plugin for now and it's async so didn't really care too much here.
	app.processListItems(ss, ad)
}
func (app *App) DeleteItem(w http.ResponseWriter, req *http.Request) {
	zap.S().Infof("Running GET command")
	input := app.findInput(req.URL.Query().Get("input"))
	if input == nil {
		zap.S().Infof("couldn't find input by input name")
		w.Write([]byte("couldn't find input by input name"))
		return
	}
	ss := sqsercore.NewSQSService(app.Config.Config.Substrings.Dlq)

	zap.S().Infof("Running input: %s", input.Config.Name)
	idt := input.Input.InvokeDeleteItem(w, req)

	if input.Config.Async {
		go app.processDeleteItem(ss, idt)
		return
	}
	// slack is the only plugin for now and it's async so didn't really care too much here.
	app.processDeleteItem(ss, idt)
}
func (app *App) MoveItems(w http.ResponseWriter, req *http.Request) {
	zap.S().Infof("Running MoveItems command")
	input := app.findInput(req.URL.Query().Get("input"))
	if input == nil {
		zap.S().Infof("couldn't find input by input name")
		w.Write([]byte("couldn't find input by input name"))
		return
	}
	ss := sqsercore.NewSQSService(app.Config.Config.Substrings.Dlq)

	zap.S().Infof("Running input: %s", input.Config.Name)
	mi := input.Input.InvokeMoveItems(w, req)
	if mi == nil {
		return
	}
	if input.Config.Async {
		go app.processMoveItems(ss, mi)
		return
	}
	// slack is the only plugin for now and it's async so didn't really care too much here.
	app.processMoveItems(ss, mi)
}
func (app *App) processGetItem(ss *sqsercore.SQSService, ad *models.GetItemActionData) {
	item, err := ss.GetItem(ad.QueueName)
	if err != nil {
		zap.S().Errorf("Running input errored %s", err)
		return
	}
	ad.SQSerItem = item
	ad.EnvironmentName = app.resolveEnvironmentName(ad.QueueName)
	var enrichersOutput []*models.EnrichLinkBlock

	for _, enricher := range app.Enrichers {
		data := enricher.Enrich(ad)
		enrichersOutput = append(enrichersOutput, data)
	}
	ad.EnrichLinkBlocks = enrichersOutput

	for _, output := range app.Outputs {
		output.InvokeGetItem(ad)
	}
	return
}
func (app *App) processListItems(ss *sqsercore.SQSService, lad *models.ListItemsActionData) {
	items, err := ss.ListNonEmptyQueues()
	if err != nil {
		zap.S().Errorf("Running input errored %s", err)
		return
	}
	lad.FilteredQueueItems = items
	for _, filter := range app.Filters {
		lad.FilteredQueueItems = filter.ApplyFilter(lad)
	}

	for _, output := range app.Outputs {
		output.InvokeListItems(lad)
	}

	return
}
func (app *App) processDeleteItem(ss *sqsercore.SQSService, diad *models.DeleteItemActionData) {
	err := ss.DeleteItem(diad.QueueName, diad.ReceiptHandle)
	if err != nil {
		zap.S().Errorf("Running input errored %s", err)
		return
	}

	for _, output := range app.Outputs {
		output.InvokeDeleteItem(diad)
	}
}
func (app *App) processMoveItems(ss *sqsercore.SQSService, di *models.MoveItemsActionData) {
	destinationQueueName := strings.ReplaceAll(di.QueueName, app.Config.Config.Substrings.Dlq, "")
	totalMessages, messagesProcessed, err := ss.MoveItems(di.QueueName, destinationQueueName)
	di.TotalItemsToMoveCount = totalMessages
	di.MovedItemsCount = messagesProcessed
	if err != nil {
		zap.S().Errorf("Running input errored %s", err)
		return
	}
	for _, output := range app.Outputs {
		output.InvokeMoveItems(di)
	}
}

func (app *App) resolveEnvironmentName(queueName string) string {
	environments := app.Config.Config.Substrings.Environments
	for _, environment := range environments {
		if strings.Contains(queueName, environment) {
			return environment
		}
	}
	return ""
}
func (app *App) findInput(inputName string) *models.RunningInput {
	for _, input := range app.Inputs {
		if input.Config.Name == inputName {
			return input
		}
	}
	return nil
}
