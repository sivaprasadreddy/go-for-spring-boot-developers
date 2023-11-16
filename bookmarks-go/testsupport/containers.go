package testsupport

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"time"
)

const postgresImage = "postgres:16-alpine"
const postgresPort = "5432"
const postgresUserName = "postgres"
const postgresPassword = "postgres"
const postgresDbName = "postgres"

type PostgresContainer struct {
	Container testcontainers.Container
	CloseFn   func()
	Host      string
	Port      string
	Database  string
	Username  string
	Password  string
}

func SetupPostgres(ctx context.Context) (*PostgresContainer, error) {
	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage(postgresImage),
		postgres.WithDatabase(postgresDbName),
		postgres.WithUsername(postgresUserName),
		postgres.WithPassword(postgresPassword),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	host, _ := container.Host(ctx)
	hostPort, _ := container.MappedPort(ctx, postgresPort)

	return &PostgresContainer{
		Container: container,
		CloseFn: func() {
			if err := container.Terminate(ctx); err != nil {
				log.Fatalf("error terminating postgres container: %s", err)
			}
		},
		Host:     host,
		Port:     hostPort.Port(),
		Database: postgresDbName,
		Username: postgresUserName,
		Password: postgresPassword,
	}, nil
}
