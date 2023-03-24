package models

type Enricher interface {
	Enrich(data *GetItemActionData) *EnrichLinkBlock
}

type EnrichLinkBlock struct {
	Text string
	Link string
}
