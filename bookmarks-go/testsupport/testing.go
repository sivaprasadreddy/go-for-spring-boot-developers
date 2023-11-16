package testsupport

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"

	log "github.com/sirupsen/logrus"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func InitPostgresContainer() *PostgresContainer {
	ctx := context.Background()
	pgContainer, err := SetupPostgres(ctx)
	if err != nil {
		log.Fatalf("failed to setup Postgres container: %v", err)
		return nil
	}
	overrideEnv(pgContainer)
	return pgContainer
}

func overrideEnv(pgC *PostgresContainer) {
	os.Setenv("DB_HOST", pgC.Host)
	os.Setenv("DB_PORT", fmt.Sprint(pgC.Port))
	os.Setenv("DB_USERNAME", pgC.Username)
	os.Setenv("DB_PASSWORD", pgC.Password)
	os.Setenv("DB_NAME", pgC.Database)
}
