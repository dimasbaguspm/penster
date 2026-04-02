package config

type KafkaConfig struct {
	Brokers string
}

func LoadKafkaConfig() KafkaConfig {
	return KafkaConfig{
		Brokers: getEnv("KAFKA_BROKERS", "localhost:9092"),
	}
}
