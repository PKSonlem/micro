package cmd

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
)

func NewClickHouseConn(cfg *clickHouseConfiguration) (clickhouse.Conn, error) {
	var (
		ctx       = context.Background()
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{cfg.addr},
			Auth: clickhouse.Auth{
				Database: cfg.db,
				Username: cfg.user,
				Password: cfg.password,
			},
		})
	)

	if err != nil {
		return nil, fmt.Errorf("failed conn to clickhouse: %w", err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, fmt.Errorf("failed ping clickhouse: %w", err)
	}

	return conn, nil
}
