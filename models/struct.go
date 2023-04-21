package models

import (
	"net/http"
	"sqser/sqsercore"
)

type GetItemActionData struct {
	ResponseWriter   *http.ResponseWriter
	QueueName        string
	SlackData        *SlackData
	EnrichLinkBlocks []*EnrichLinkBlock
	SQSerItem        *sqsercore.SQSerItem
	EnvironmentName  string
	Error            error
}

type ListItemsActionData struct {
	ResponseWriter     *http.ResponseWriter
	SlackData          *SlackData
	FilteredQueueItems []*sqsercore.QueueCount
}

type DeleteItemActionData struct {
	ResponseWriter *http.ResponseWriter
	SlackData      *SlackData
	QueueName      string
	ReceiptHandle  string
}

type MoveItemsActionData struct {
	ResponseWriter        *http.ResponseWriter
	SlackData             *SlackData
	QueueName             string
	TotalItemsToMoveCount int
	MovedItemsCount       int
}

type SlackData struct {
	ResponseUrl string
	ChanId      string
	UserName    string
}
