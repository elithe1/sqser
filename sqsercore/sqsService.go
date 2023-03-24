package sqsercore

import (
	"errors"
	"strings"
)

type SQSService struct {
	sqsAdapter   *SQSAdapter
	dlqSubstring string
}

func NewSQSService(dlqSubstring string) *SQSService {
	sa := NewSqsAdapter()
	return &SQSService{sqsAdapter: &sa, dlqSubstring: dlqSubstring}
}

func (ss *SQSService) GetItem(queueName string) (item *SQSerItem, err error) {
	if !strings.Contains(queueName, ss.dlqSubstring) {
		return nil, errors.New("reading items from non-dlq queues are forbidden")
	}
	return ss.sqsAdapter.GetItem(queueName)
}

func (ss *SQSService) ListNonEmptyQueues() ([]*QueueCount, error) {
	items, err := ss.sqsAdapter.ListItems()
	if err != nil {
		return nil, err
	}
	for i, item := range items {
		if item.Count == "0" {
			//relying on items to be ordered so whe we see the fist "0" there will be only 0 left.
			items = items[:i]
			break
		}
	}
	return items, nil
}

func (ss *SQSService) DeleteItem(queueName string, receiptHandle string) (err error) {
	if !strings.Contains(queueName, ss.dlqSubstring) {
		return errors.New("deleting items from non-dlq queues is forbidden")
	}
	return ss.sqsAdapter.DeleteItem(queueName, receiptHandle)
}

func (ss *SQSService) MoveItems(sourceQueueName string, destinationQueueName string) (totalMessages int, messagesProcessed int, err error) {
	if !strings.Contains(sourceQueueName, ss.dlqSubstring) {
		return 0, 0, errors.New("moving items from non-dlq queues is forbidden")
	}

	return ss.sqsAdapter.MoveItems(sourceQueueName, destinationQueueName)
}
