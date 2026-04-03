package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	serverURL     string
	serverProcess *exec.Cmd
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithUsername("penster"),
		postgres.WithPassword("placeholder"),
		postgres.WithDatabase("penster"),
	)
	if err != nil {
		log.Fatalf("Failed to start postgres container: %v", err)
	}

	postgresURL, err := container.ConnectionString(ctx)
	if err != nil {
		log.Fatalf("Failed to get connection string: %v", err)
	}

	dbHost, dbPort := mustExtractHostPort(postgresURL)

	os.Setenv("DB_USER", "penster")
	os.Setenv("DB_PASSWORD", "placeholder")
	os.Setenv("DB_NAME", "penster")
	os.Setenv("DB_HOST", dbHost)
	os.Setenv("DB_PORT", dbPort)
	os.Setenv("DB_SSLMODE", "disable")
	os.Setenv("APP_ENV", "test")
	os.Setenv("AUTO_MIGRATE", "true")
	os.Setenv("APP_PORT", "8081")

	serverURL = fmt.Sprintf("http://localhost:%s", os.Getenv("APP_PORT"))

	if err := waitForPostgres(dbHost, dbPort, 30*time.Second); err != nil {
		log.Fatalf("Postgres did not become ready: %v", err)
	}

	repoRoot, _ := os.Getwd()
	repoRoot = filepath.Dir(filepath.Dir(repoRoot))

	binaryPath := filepath.Join(repoRoot, "bin", "penster_test")
	if err := os.MkdirAll(filepath.Dir(binaryPath), 0755); err != nil {
		log.Fatalf("Failed to create bin directory: %v", err)
	}
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "github.com/dimasbaguspm/penster/cmd/server")
	buildCmd.Env = append(os.Environ(), "GOSUMDB=off")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	if err := buildCmd.Run(); err != nil {
		log.Fatalf("Failed to build penster server: %v", err)
	}

	serverProcess = exec.Command(binaryPath)
	serverProcess.Dir = repoRoot
	serverProcess.Env = os.Environ()
	serverProcess.Stdout = os.Stdout
	serverProcess.Stderr = os.Stderr

	if err := serverProcess.Start(); err != nil {
		log.Fatalf("Failed to start penster server: %v", err)
	}

	if err := waitForServer(serverURL, 30*time.Second); err != nil {
		serverProcess.Process.Signal(syscall.SIGTERM)
		log.Fatalf("Server did not become ready: %v", err)
	}

	exitCode := m.Run()

	if serverProcess != nil && serverProcess.Process != nil {
		serverProcess.Process.Signal(syscall.SIGTERM)
		serverProcess.Wait()
	}
	container.Terminate(ctx)

	os.Exit(exitCode)
}

func mustExtractHostPort(connectionString string) (string, string) {
	u, _ := url.Parse(connectionString)
	return u.Hostname(), u.Port()
}

func waitForPostgres(host, port string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), 2*time.Second)
		if err == nil {
			conn.Close()
			time.Sleep(1 * time.Second)
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("postgres did not become ready within %v", timeout)
}

func waitForServer(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url + "/health")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("server did not respond within %v", timeout)
}
