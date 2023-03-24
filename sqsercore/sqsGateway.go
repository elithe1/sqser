package sqsercore

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"go.uber.org/zap"
)

type SQSGateway struct {
	sqsService *sqs.SQS
}

func NewSqsGateway() SQSGateway {
	sess, err := getSession()
	if err != nil {
		zap.S().Error(fmt.Sprintf("Unable to getMessage AWS session for region"))
	}

	svc := sqs.New(sess)
	return SQSGateway{sqsService: svc}
}

func (sg *SQSGateway) GetQueueUrl(queueName string) (str string, err error) {
	params := &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	}
	resp, err := sg.sqsService.GetQueueUrl(params)
	if err != nil {
		return
	}

	return *resp.QueueUrl, nil
}

func (sg *SQSGateway) getItem(queueUrl string) (message *sqs.Message, err error) {
	var params = &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(queueUrl),
		VisibilityTimeout:     aws.Int64(180), //prolonged visibility timeout so we can delete the message :)
		WaitTimeSeconds:       aws.Int64(0),
		MaxNumberOfMessages:   aws.Int64(1),
		MessageAttributeNames: []*string{aws.String(sqs.QueueAttributeNameAll)},
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameMessageGroupId),
			aws.String(sqs.MessageSystemAttributeNameMessageDeduplicationId)},
	}
	resp, err := sg.sqsService.ReceiveMessage(params)
	if err != nil {
		zap.S().Error("Failed to receive messages", err)
		return
	}
	if len(resp.Messages) == 0 {
		fmt.Println()
		zap.S().Warn("No messages found in queue ", queueUrl)
		return nil, errors.New(fmt.Sprintf("No messages found in queue. %s", queueUrl))
	}
	return resp.Messages[0], nil

}

func (sg *SQSGateway) getQueues() (*sqs.ListQueuesOutput, error) {
	lo := sqs.ListQueuesInput{
		MaxResults: aws.Int64(100),
	}
	return sg.sqsService.ListQueues(&lo)
}

func (sg *SQSGateway) getQueueAttributes(queueUrl string) (*sqs.GetQueueAttributesOutput, error) {
	attr := sqs.GetQueueAttributesInput{
		QueueUrl:       aws.String(queueUrl),
		AttributeNames: []*string{aws.String("All")},
	}
	return sg.sqsService.GetQueueAttributes(&attr)
}

func (sg *SQSGateway) deleteItem(queueUrl string, receiptHandle string) error {
	var params = &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueUrl),
		ReceiptHandle: aws.String(receiptHandle),
	}
	_, err := sg.sqsService.DeleteMessage(params)
	if err != nil {
		zap.S().Error("Failed to delete messages", err)
		return err
	}
	return nil
}

func (sg *SQSGateway) ReceiveMessage(sourceQueueUrl string) (*sqs.ReceiveMessageOutput, error) {
	var params = &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(sourceQueueUrl),
		VisibilityTimeout:     aws.Int64(2),
		WaitTimeSeconds:       aws.Int64(0),
		MaxNumberOfMessages:   aws.Int64(10),
		MessageAttributeNames: []*string{aws.String(sqs.QueueAttributeNameAll)},
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameMessageGroupId),
			aws.String(sqs.MessageSystemAttributeNameMessageDeduplicationId)},
	}

	return sg.sqsService.ReceiveMessage(params)
}

func (sg *SQSGateway) SendMessageBatch(queueUrl string, entries []*sqs.SendMessageBatchRequestEntry) (*sqs.SendMessageBatchOutput, error) {
	batch := &sqs.SendMessageBatchInput{
		QueueUrl: aws.String(queueUrl),
		Entries:  entries,
	}
	return sg.sqsService.SendMessageBatch(batch)
}

func (sg *SQSGateway) DeleteMessageBatch(queueUrl string, entries []*sqs.DeleteMessageBatchRequestEntry) (*sqs.DeleteMessageBatchOutput, error) {
	deleteMessageBatch := &sqs.DeleteMessageBatchInput{
		Entries:  entries,
		QueueUrl: aws.String(queueUrl),
	}
	return sg.sqsService.DeleteMessageBatch(deleteMessageBatch)
}
