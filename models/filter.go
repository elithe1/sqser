package models

import "sqser/sqsercore"

type Filter interface {
	ApplyFilter(queueCounts *ListItemsActionData) []*sqsercore.QueueCount
}
