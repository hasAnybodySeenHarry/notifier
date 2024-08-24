package consumer

type Event struct {
	Metadata metadata               `json:"metadata"`
	Data     map[string]interface{} `json:"data"`
}

type metadata struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
	Source    string `json:"source"`
	Version   string `json:"version"`
}
