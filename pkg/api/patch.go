package api

type Patch struct {
	Op    string      `json:"op,omitempty"`
	Path  string      `json:"path,omitempty"`
	From  string      `json:"from,omitempty"`
	Value interface{} `json:"value,omitempty"`
}
