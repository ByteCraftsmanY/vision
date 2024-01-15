package dtos

type HealthDTO struct {
	Status      int    `json:"status"`
	Message     string `json:"message"`
	DBStatus    bool   `json:"db_status"`
	KafkaStatus bool   `json:"kafka_status"`
}
