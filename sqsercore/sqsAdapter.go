package sqsercore

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/umpc/go-sortedmap"
	"github.com/umpc/go-sortedmap/desc"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type SQSerItem struct {
	MessageGroupId string           `json:"messageGroupId"`
	MessageBody    *json.RawMessage `json:"body"`
	ReceiptHandle  *string          `json:"receiptHandle"`
}

type QueueCount struct {
	QueueName string `json:"queueName"`
	Count     string `json:"count"`
}

type SQSAdapter struct {
	sqsGateway *SQSGateway
}

func NewSqsAdapter() SQSAdapter {
	sg := NewSqsGateway()
	return SQSAdapter{sqsGateway: &sg}
}

func (sa *SQSAdapter) GetItem(queueName string) (item *SQSerItem, err error) {
	queueUrl, err := sa.sqsGateway.GetQueueUrl(queueName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Couldnt resolve queue URL %s", queueName))
	}

	m, err := sa.sqsGateway.getItem(queueUrl)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Couldnt get item %s", queueUrl))
	}

	messageBody := m.Body
	var messageGroupId string
	if mgi, ok := m.Attributes[sqs.MessageSystemAttributeNameMessageGroupId]; ok {
		messageGroupId = *mgi
	}

	dec := json.NewDecoder(bytes.NewReader([]byte(*messageBody)))
	var body *json.RawMessage
	if err := dec.Decode(&body); err != nil {
		fmt.Printf("error decoding messageBody: %s\n", err)
		return nil, errors.New("error decoding messageBody")
	}

	return &SQSerItem{
		MessageGroupId: messageGroupId,
		MessageBody:    body,
		ReceiptHandle:  m.ReceiptHandle,
	}, nil

}

func (sa *SQSAdapter) ListItems() ([]*QueueCount, error) {
	queues, err := sa.sqsGateway.getQueues()
	if err != nil {
		zap.S().Error(fmt.Sprintf("Unable to List queues"))
		return nil, err
	}

	sm := sortedmap.New(len(queues.QueueUrls), desc.Int)
	for _, queueUrl := range queues.QueueUrls {
		queueAttributes, err := sa.sqsGateway.getQueueAttributes(*queueUrl)
		if err != nil {
			zap.S().Error(fmt.Sprintf("Unable to get attrebutes of queueUrl"))
			return nil, err
		}
		totalMessages, _ := strconv.Atoi(*queueAttributes.Attributes["ApproximateNumberOfMessages"])
		queueName := (*queueUrl)[strings.LastIndex(*queueUrl, "/")+1:]
		sm.Insert(queueName, totalMessages)
		zap.S().Infof("Number of messages in %s is %d", queueName, totalMessages)

	}
	iterCh, err := sm.IterCh()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer iterCh.Close()
	sortedQueueCounts := make([]*QueueCount, len(queues.QueueUrls))
	i := 0
	for rec := range iterCh.Records() {
		strKey := fmt.Sprintf("%v", rec.Key)
		strValue := fmt.Sprintf("%v", rec.Val)
		sortedQueueCounts[i] = &QueueCount{
			QueueName: strKey,
			Count:     strValue,
		}
		i++
	}

	return sortedQueueCounts, nil
}

func (sa *SQSAdapter) DeleteItem(queueName string, receiptHandle string) (err error) {
	queueUrl, err := sa.sqsGateway.GetQueueUrl(queueName)
	if err != nil {
		return errors.New(fmt.Sprintf("Couldnt resolve queue URL %s", queueName))
	}
	err = sa.sqsGateway.deleteItem(queueUrl, receiptHandle)
	return err
}

func (sa *SQSAdapter) MoveItems(sourceQueueName string, destinationQueueName string) (totalMessages int, messagesProcessed int, err error) {
	err = nil
	totalMessages = 0
	messagesProcessed = 0
	sourceQueueUrl, err := sa.sqsGateway.GetQueueUrl(sourceQueueName)
	if err != nil {
		return
	}
	destinationQueueUrl, err := sa.sqsGateway.GetQueueUrl(destinationQueueName)
	if err != nil {
		return
	}
	queueAttributes, err := sa.sqsGateway.getQueueAttributes(sourceQueueUrl)
	if err != nil {
		zap.S().Error("Failed to resolve queue attributes", err)
		return
	}
	totalMessages, _ = strconv.Atoi(*queueAttributes.Attributes["ApproximateNumberOfMessages"])
	zap.S().Infof("Number of mwssages in DLQ is %d", totalMessages)
	for {
		resp, e := sa.sqsGateway.ReceiveMessage(sourceQueueUrl)
		if e != nil {
			zap.S().Error("Failed to receive messages", e)
			return totalMessages, messagesProcessed, e
		}

		if len(resp.Messages) == 0 || messagesProcessed == totalMessages {
			fmt.Println()
			zap.S().Infof("Done. Moved %s messages", strconv.Itoa(totalMessages))
			break
		}

		messagesToCopy := resp.Messages

		if len(resp.Messages)+messagesProcessed > totalMessages {
			messagesToCopy = resp.Messages[0 : totalMessages-messagesProcessed]
		}

		sendResp, e := sa.sqsGateway.SendMessageBatch(destinationQueueUrl, convertToEntries(messagesToCopy))

		if e != nil {
			zap.S().Error("Failed to un-queue messages to the destination", e)
			return totalMessages, messagesProcessed, e
		}

		if len(sendResp.Failed) > 0 {
			zap.S().Error("%d messages failed to enqueue, see details below", len(sendResp.Failed))
			for index, failed := range sendResp.Failed {
				zap.S().Error("%d - (%s) %s", index, *failed.Code, *failed.Message)
			}
			return
		}

		if len(sendResp.Successful) == len(messagesToCopy) {

			deleteResp, err := sa.sqsGateway.DeleteMessageBatch(sourceQueueUrl, convertSuccessfulMessageToBatchRequestEntry(messagesToCopy))

			if err != nil {
				zap.S().Error("Failed to delete messages from source queue", err)
				return totalMessages, messagesProcessed, err
			}

			if len(deleteResp.Failed) > 0 {
				zap.S().Error("Error deleting messages, the following were not deleted\n %s", deleteResp.Failed)
				return totalMessages, messagesProcessed, err
			}

			messagesProcessed += len(messagesToCopy)
		}

	}
	return
}

func convertToEntries(messages []*sqs.Message) []*sqs.SendMessageBatchRequestEntry {
	result := make([]*sqs.SendMessageBatchRequestEntry, len(messages))
	for i, message := range messages {
		requestEntry := &sqs.SendMessageBatchRequestEntry{
			MessageBody:       message.Body,
			Id:                message.MessageId,
			MessageAttributes: message.MessageAttributes,
		}

		if messageGroupId, ok := message.Attributes[sqs.MessageSystemAttributeNameMessageGroupId]; ok {
			requestEntry.MessageGroupId = messageGroupId
		}

		if messageDeduplicationId, ok := message.Attributes[sqs.MessageSystemAttributeNameMessageDeduplicationId]; ok {
			requestEntry.MessageDeduplicationId = messageDeduplicationId
		}

		result[i] = requestEntry
	}

	return result
}

func convertSuccessfulMessageToBatchRequestEntry(messages []*sqs.Message) []*sqs.DeleteMessageBatchRequestEntry {
	result := make([]*sqs.DeleteMessageBatchRequestEntry, len(messages))
	for i, message := range messages {
		result[i] = &sqs.DeleteMessageBatchRequestEntry{
			ReceiptHandle: message.ReceiptHandle,
			Id:            message.MessageId,
		}
	}

	return result
}
