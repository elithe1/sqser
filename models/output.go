package models

type Output interface {
	InvokeGetItem(actionData *GetItemActionData)
	InvokeListItems(actionData *ListItemsActionData)
	InvokeDeleteItem(actionData *DeleteItemActionData)
	InvokeMoveItems(actionData *MoveItemsActionData)
}
