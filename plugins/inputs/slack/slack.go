package inputs

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"sqser/models"
	"sqser/plugins/inputs"
	"strings"
)

type Slack struct {
	Log *zap.SugaredLogger
}

func newSlack(logger *zap.SugaredLogger) *Slack {
	return &Slack{
		Log: logger,
	}
}

func (s *Slack) InvokeGetItem(w http.ResponseWriter, req *http.Request) *models.GetItemActionData {
	s.Log.Infof("Running InvokeGetItem")
	req.ParseForm()
	s.Log.Infof(req.Form.Encode())

	queueName := req.FormValue("text")
	responseUrl := req.FormValue("response_url")
	chanId := req.FormValue("channel_id")
	userName := req.FormValue("user_name")
	if queueName == "" {
		w.Write([]byte("Please provide a valid queue name. For example production-sqs-integrations-dlq.fifo"))
		return nil
	}
	w.Write([]byte(fmt.Sprintf("Let me fetch you a message from *%s* dlq real quick...", queueName)))

	return &models.GetItemActionData{
		ResponseWriter: &w,
		QueueName:      queueName,
		SlackData: &models.SlackData{
			ResponseUrl: responseUrl,
			ChanId:      chanId,
			UserName:    userName,
		},
	}
}

func (s *Slack) InvokeListItems(w http.ResponseWriter, req *http.Request) *models.ListItemsActionData {
	s.Log.Infof("Running InvokeListItems")
	req.ParseForm()
	s.Log.Infof(req.Form.Encode())

	responseUrl := req.FormValue("response_url")
	chanId := req.FormValue("channel_id")
	userName := req.FormValue("user_name")
	w.Write([]byte(fmt.Sprintf("Let me fetch you all non empty queues real quick...")))

	return &models.ListItemsActionData{
		ResponseWriter: &w,
		SlackData: &models.SlackData{
			ResponseUrl: responseUrl,
			ChanId:      chanId,
			UserName:    userName,
		},
	}
}

func (s *Slack) InvokeDeleteItem(w http.ResponseWriter, req *http.Request) *models.DeleteItemActionData {
	s.Log.Infof("Running InvokeDeleteItem")
	req.ParseForm()
	s.Log.Infof(req.Form.Encode())
	text := req.FormValue("text")
	responseUrl := req.FormValue("response_url")
	chanId := req.FormValue("channel_id")
	userName := req.FormValue("user_name")

	args := strings.Split(text, " ")
	if len(args) != 2 {
		w.Write([]byte(fmt.Sprintf("Should have exactly two arguments but got %d", len(args))))
		return nil
	}
	queueName := args[0]
	if queueName == "" {
		w.Write([]byte("Please provide a valid queue name. For example sqs-integrations"))
		return nil
	}

	w.Write([]byte(fmt.Sprintf("Let me delete the message for you from *%s* real quick...", queueName)))

	return &models.DeleteItemActionData{
		ResponseWriter: &w,
		SlackData: &models.SlackData{
			ResponseUrl: responseUrl,
			ChanId:      chanId,
			UserName:    userName,
		},
		QueueName:     queueName,
		ReceiptHandle: args[1],
	}
}

func (s *Slack) InvokeMoveItems(w http.ResponseWriter, req *http.Request) *models.MoveItemsActionData {
	s.Log.Infof("Running InvokeMoveItems")
	req.ParseForm()
	s.Log.Infof(req.Form.Encode())
	queueName := req.FormValue("text")
	responseUrl := req.FormValue("response_url")
	chanId := req.FormValue("channel_id")
	userName := req.FormValue("user_name")

	if queueName == "" {
		w.Write([]byte("Please provide a valid queue name. For example sqs-integrations"))
		return nil
	}
	w.Write([]byte(fmt.Sprintf("Let me move all the messages from *%s* dlq real quick...", queueName)))
	return &models.MoveItemsActionData{
		ResponseWriter: &w,
		SlackData: &models.SlackData{
			ResponseUrl: responseUrl,
			ChanId:      chanId,
			UserName:    userName,
		},
		QueueName: queueName,
	}
}
func init() {
	pluginName := "slack"
	inputs.Add(pluginName, func() models.Input {
		logger := zap.S().With(zap.String("pluginType", "input"), zap.String("pluginName", pluginName))
		s := newSlack(logger)
		s.Log.Infof("loaded successfully")
		return s
	})
}
