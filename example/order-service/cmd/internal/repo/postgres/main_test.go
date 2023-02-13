package postgres

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

var (
	mu   sync.Mutex
	conn *sql.DB
	url  string
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	terminate := initDB(ctx)

	code := m.Run()
	if err := terminate(ctx); err != nil {
		log.Fatalf("could not connect to the database: %s", err)
	}

	os.Exit(code)
}

func initDB(ctx context.Context) func(context.Context) error {
	container, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				FromDockerfile: testcontainers.FromDockerfile{
					Context:    "./../../../../",
					Dockerfile: "config/database/Dockerfile",
				},
				ExposedPorts: []string{"5432/tcp"},
				WaitingFor: (&wait.LogStrategy{
					Log:          "database system is ready to accept connections",
					Occurrence:   2,
					PollInterval: 100 * time.Millisecond,
				}).WithStartupTimeout(2 * time.Minute),
				Env: map[string]string{
					"POSTGRES_DB":       "order-service",
					"POSTGRES_PASSWORD": "order-service",
					"POSTGRES_USER":     "order-service",
				},
			},
			Started: true,
			Logger:  log.Default(),
		},
	)
	if err != nil {
		log.Fatalf("could not connect to the database: %s", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("could not get container host: %s", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("could not get container port: %s", err)
	}

	url = fmt.Sprintf("postgres://order-service:order-service@%v:%v/order-service?sslmode=disable", host, port.Port())
	return container.Terminate
}

func getDB(t *testing.T) *sql.DB {
	mu.Lock()
	defer mu.Unlock()
	if conn != nil {
		return conn
	}

	var err error
	conn, err = sql.Open("postgres", url)
	if err != nil {
		t.Fatalf("could not open sql connection: %s", err)
	}

	return conn
}
