package model

type LinkStatus string

const (
	StatusAvailable    LinkStatus = "available"
	StatusNotAvailable LinkStatus = "not available"
	StatusPending      LinkStatus = "pending"
)

type Batch struct {
	ID     int64                 `json:"id"`
	Links  []string              `json:"links"`
	Status map[string]LinkStatus `json:"status"`
}

type State struct {
	NextID  int64           `json:"next_id"`
	Batches map[int64]Batch `json:"batches"`
}
