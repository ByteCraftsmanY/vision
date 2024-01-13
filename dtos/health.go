package dtos

type HealthDTO struct {
	Status      int    `json:"status,omitempty"`
	Message     string `json:"message,omitempty"`
	DBStatus    bool   `json:"db_status,omitempty"`
	KafkaStatus bool   `json:"kafka_status,omitempty"`
}
