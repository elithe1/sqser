package outputs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net/http"
	"sqser/models"
	"sqser/plugins/outputs"
	"strings"
)

type Slack struct {
	Log *zap.SugaredLogger
}

func newSlack(logger *zap.SugaredLogger) *Slack {
	return &Slack{Log: logger}
}

type SimpleSlackMessageResponse struct {
	Text     string `json:"text"`
	RespType string `json:"response_type"`
}

func (s *Slack) sendMessage(toSend interface{}, responseUrl string) {
	if responseUrl == "" {
		s.Log.Error("Response url was empty. Not sending ")
		return
	}
	json_data, err := json.Marshal(toSend)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(responseUrl, "application/json",
		bytes.NewBuffer(json_data))

	if err != nil {
		s.Log.Error("err  sending message to slack ", err)
	}

	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)

	fmt.Println(res["json"])

}

func (s *Slack) InvokeGetItem(actionData *models.GetItemActionData) {
	s.Log.Infof("Running InvokeGetItem")

	var prettyJSON bytes.Buffer
	item := actionData.SQSerItem
	error := json.Indent(&prettyJSON, *item.MessageBody, "", "\t")
	if error != nil {
		log.Println("JSON parse error: ", error)
		return
	}

	json.MarshalIndent(item.MessageBody, "", "\t")
	messageBody := &SimpleSlackMessageResponse{
		RespType: "peripheral",
		Text:     fmt.Sprintf("Here is your requested message *%s*\n%s", strings.ToTitle(actionData.SlackData.UserName), string(prettyJSON.Bytes())),
	}
	s.sendMessage(messageBody, actionData.SlackData.ResponseUrl)

	receiptHandle := &SimpleSlackMessageResponse{
		RespType: "peripheral",
		Text:     fmt.Sprintf("*messageGroupId:* %s\n*receiptHandle for deletion:* %s\n", item.MessageGroupId, *item.ReceiptHandle),
	}
	s.sendMessage(receiptHandle, actionData.SlackData.ResponseUrl)
	for _, block := range actionData.EnrichLinkBlocks {
		logLink := &SimpleSlackMessageResponse{
			RespType: "peripheral",
			Text:     fmt.Sprintf("<%s|%s>", block.Link, block.Text),
		}
		s.sendMessage(logLink, actionData.SlackData.ResponseUrl)
	}

}

func (s *Slack) InvokeListItems(actionData *models.ListItemsActionData) {
	s.Log.Infof("Running InvokeListItems")
	respType := "peripheral"
	// send private message if moved nothing.
	if actionData.SlackData.UserName == "" {
		actionData.SlackData.UserName = "Mate"
		respType = "in_channel"
	}
	if len(actionData.FilteredQueueItems) == 0 {
		messageBody := &SimpleSlackMessageResponse{
			RespType: respType,
			Text:     fmt.Sprintf("Seems like all queues on are on 0. Nicely done!! Time to take a break and grab a cookie"),
		}
		s.sendMessage(messageBody, actionData.SlackData.ResponseUrl)
		return
	}
	m, _ := json.MarshalIndent(actionData.FilteredQueueItems, "", "\t")
	messageBody := &SimpleSlackMessageResponse{
		RespType: respType,
		Text:     fmt.Sprintf("Hey *%s* here found some non empty dlqs:\n%s", strings.Title(actionData.SlackData.UserName), string(m)),
	}
	s.sendMessage(messageBody, actionData.SlackData.ResponseUrl)

}
func (s *Slack) InvokeDeleteItem(actionData *models.DeleteItemActionData) {
	s.Log.Infof("Running InvokeDeleteItem")
	messageBody := &SimpleSlackMessageResponse{
		RespType: "peripheral",
		Text:     fmt.Sprintf("Message deleted successfully"),
	}
	s.sendMessage(messageBody, actionData.SlackData.ResponseUrl)
}
func (s *Slack) InvokeMoveItems(actionData *models.MoveItemsActionData) {
	s.Log.Infof("Running InvokeMoveItems")
	respType := "in_channel"
	// send private message if moved nothing.
	if actionData.MovedItemsCount == 0 {
		respType = "peripheral"
	}
	resp := &SimpleSlackMessageResponse{
		RespType: respType,
		Text:     fmt.Sprintf("Hey *%s*, found %d messages in *%s* dlq. Moved %d messages.", strings.Title(actionData.SlackData.UserName), actionData.TotalItemsToMoveCount, actionData.QueueName, actionData.MovedItemsCount),
	}
	s.sendMessage(resp, actionData.SlackData.ResponseUrl)
}
func init() {
	pluginName := "slack"
	outputs.Add(pluginName, func() models.Output {
		logger := zap.S().With(zap.String("pluginType", "output"), zap.String("pluginName", pluginName))
		s := newSlack(logger)
		logger.Infof("loaded successfully")
		return s
	})
}
