module miga

go 1.15

replace github.com/pressly/goose v2.7.0+incompatible => github.com/chapsuk/goose v2.1.1-0.20210207132231-8dfe4480e4dd+incompatible

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/ClickHouse/clickhouse-go v1.4.3
	github.com/go-pg/pg v8.0.3+incompatible
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang-migrate/migrate/v4 v4.2.5
	github.com/im-kulikov/migrate v0.1.0
	github.com/jinzhu/inflection v0.0.0-20180308033659-04140366298a // indirect
	github.com/lib/pq v1.0.0
	github.com/onsi/ginkgo v1.15.0 // indirect
	github.com/onsi/gomega v1.10.5 // indirect
	github.com/pkg/errors v0.8.1
	github.com/pressly/goose v2.7.0+incompatible
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/viper v1.3.2
	go.uber.org/atomic v1.3.2 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.9.1
	gopkg.in/urfave/cli.v2 v2.0.0-20180128182452-d3ae77c26ac8
	mellium.im/sasl v0.2.1 // indirect
)
