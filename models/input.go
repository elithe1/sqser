package models

import (
	"net/http"
)

type Input interface {
	InvokeGetItem(w http.ResponseWriter, req *http.Request) *GetItemActionData
	InvokeListItems(w http.ResponseWriter, req *http.Request) *ListItemsActionData
	InvokeDeleteItem(w http.ResponseWriter, req *http.Request) *DeleteItemActionData
	InvokeMoveItems(w http.ResponseWriter, req *http.Request) *MoveItemsActionData
}
