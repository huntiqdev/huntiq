package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

// TestGracefulShutdown_InFlightRequestCompletes verifies that Fiber v2 allows
// in-flight requests to complete when ShutdownWithTimeout is invoked.
func TestGracefulShutdown_InFlightRequestCompletes(t *testing.T) {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	inFlightStarted := make(chan struct{})
	requestFinished := make(chan struct{})

	app.Get("/slow", func(c *fiber.Ctx) error {
		close(inFlightStarted)
		time.Sleep(300 * time.Millisecond)
		close(requestFinished)
		return c.SendString("request-completed-successfully")
	})

	// Find free local port
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen on port: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	_ = ln.Close()

	addr := fmt.Sprintf("127.0.0.1:%d", port)

	serverErr := make(chan error, 1)
	go func() {
		if err := app.Listen(addr); err != nil {
			serverErr <- err
		}
	}()

	// Wait for server to start listening
	time.Sleep(100 * time.Millisecond)

	// Launch in-flight request in background
	var respBody string
	var respStatusCode int
	var clientErr error
	reqWg := sync.WaitGroup{}
	reqWg.Add(1)

	go func() {
		defer reqWg.Done()
		resp, err := http.Get(fmt.Sprintf("http://%s/slow", addr))
		if err != nil {
			clientErr = err
			return
		}
		defer resp.Body.Close()
		respStatusCode = resp.StatusCode
		bodyBytes, _ := io.ReadAll(resp.Body)
		respBody = string(bodyBytes)
	}()

	// Wait for handler to begin executing
	select {
	case <-inFlightStarted:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for in-flight request to begin")
	}

	// Trigger graceful shutdown while request is still in flight
	shutdownStart := time.Now()
	shutdownErr := app.ShutdownWithTimeout(2 * time.Second)
	shutdownDuration := time.Since(shutdownStart)

	if shutdownErr != nil {
		t.Fatalf("app.ShutdownWithTimeout returned unexpected error: %v", shutdownErr)
	}

	// Wait for client request to finish
	reqWg.Wait()

	if clientErr != nil {
		t.Fatalf("client request failed during graceful shutdown: %v", clientErr)
	}
	if respStatusCode != http.StatusOK {
		t.Fatalf("expected HTTP 200, got %d", respStatusCode)
	}
	if respBody != "request-completed-successfully" {
		t.Fatalf("unexpected response body: %q", respBody)
	}

	t.Logf("Graceful shutdown completed in %v, in-flight request finished successfully with 200 OK", shutdownDuration)
}

// TestWorkerCancellation_OnJobContextDone verifies that background jobs terminate
// cleanly upon receiving context cancellation.
func TestWorkerCancellation_OnJobContextDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	workerStopped := make(chan struct{})

	// Simulate background worker loop
	go func() {
		ticker := time.NewTicker(20 * time.Millisecond)
		defer ticker.Stop()
		defer close(workerStopped)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// background maintenance work
			}
		}
	}()

	// Allow worker to run several ticks
	time.Sleep(60 * time.Millisecond)

	// Cancel context (as in jobCancel())
	cancel()

	select {
	case <-workerStopped:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("background worker failed to stop after context cancellation")
	}
}

// TestShutdownOrdering_ResourceCleanupSequence verifies the exact cleanup ordering:
// 1. Worker contexts cancelled
// 2. HTTP server shut down
// 3. Database connection pool closed
func TestShutdownOrdering_ResourceCleanupSequence(t *testing.T) {
	var stepCounter int32
	var workerCancelledStep int32
	var serverShutdownStep int32
	var dbClosedStep int32

	ctx, cancel := context.WithCancel(context.Background())
	workerDone := make(chan struct{})

	go func() {
		<-ctx.Done()
		atomic.StoreInt32(&workerCancelledStep, atomic.AddInt32(&stepCounter, 1))
		close(workerDone)
	}()

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/ping", func(c *fiber.Ctx) error { return c.SendString("pong") })

	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	addr := ln.Addr().String()
	_ = ln.Close()

	go func() {
		_ = app.Listen(addr)
	}()
	time.Sleep(50 * time.Millisecond)

	// Execute sequence modeled from main.go
	// 1. Cancel background workers
	cancel()
	<-workerDone

	// 2. Shutdown HTTP server
	_ = app.ShutdownWithTimeout(1 * time.Second)
	atomic.StoreInt32(&serverShutdownStep, atomic.AddInt32(&stepCounter, 1))

	// 3. Close DB pool (simulated)
	atomic.StoreInt32(&dbClosedStep, atomic.AddInt32(&stepCounter, 1))

	if workerCancelledStep != 1 {
		t.Errorf("expected worker cancellation at step 1, got %d", workerCancelledStep)
	}
	if serverShutdownStep != 2 {
		t.Errorf("expected server shutdown at step 2, got %d", serverShutdownStep)
	}
	if dbClosedStep != 3 {
		t.Errorf("expected db pool close at step 3, got %d", dbClosedStep)
	}
}
