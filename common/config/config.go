package config

import "fmt"

type MysqlConfig struct {
	DataSource string
}

type RabbitMQConfig struct {
	Host  string
	Port  int
	User  string
	Pass  string
	VHost string
}
type ElasticSearchConfig struct {
	Addr string
}

func (r *RabbitMQConfig) Dns() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%d/%s",
		r.User,
		r.Pass,
		r.Host,
		r.Port,
		r.VHost,
	)
}
