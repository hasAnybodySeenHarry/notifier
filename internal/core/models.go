package core

type Entity struct {
	ID   int64  `json:"id,omitempty" bson:"id,omitempty"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
}

func EntityToMap(e *Entity) map[string]interface{} {
	return map[string]interface{}{
		"id":   e.ID,
		"name": e.Name,
	}
}
