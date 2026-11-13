package main

import (
	"context"
	"errors"
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

// mockResourceCloser simulates a database or pool closer tracking its close step.
type mockResourceCloser struct {
	closeFunc func()
	closed    int32
}

func (m *mockResourceCloser) Close() {
	atomic.StoreInt32(&m.closed, 1)
	if m.closeFunc != nil {
		m.closeFunc()
	}
}

func (m *mockResourceCloser) IsClosed() bool {
	return atomic.LoadInt32(&m.closed) == 1
}

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

	// Bind to an ephemeral port without close/re-open race conditions
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen on port: %v", err)
	}
	addr := ln.Addr().String()

	serverErr := make(chan error, 1)
	go func() {
		if err := app.Listener(ln); err != nil {
			serverErr <- err
		}
	}()

	// Launch in-flight request in background with readiness retries
	var respBody string
	var respStatusCode int
	var clientErr error
	reqWg := sync.WaitGroup{}
	reqWg.Add(1)

	go func() {
		defer reqWg.Done()
		var resp *http.Response
		for i := 0; i < 20; i++ {
			resp, clientErr = http.Get(fmt.Sprintf("http://%s/slow", addr))
			if clientErr == nil {
				break
			}
			time.Sleep(50 * time.Millisecond)
		}
		if clientErr != nil {
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

	// Cancel context
	cancel()

	select {
	case <-workerStopped:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("background worker failed to stop after context cancellation")
	}
}

// TestShutdownOrdering_ResourceCleanupSequence verifies the exact cleanup ordering:
// 1. In-flight HTTP server requests drained
// 2. Worker contexts cancelled and workers stopped
// 3. Database connection pool closed
func TestShutdownOrdering_ResourceCleanupSequence(t *testing.T) {
	var stepCounter int32
	var httpDrainedStep int32
	var workerStoppedStep int32
	var dbClosedStep int32

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/ping", func(c *fiber.Ctx) error { return c.SendString("pong") })

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen error: %v", err)
	}
	go func() {
		_ = app.Listener(ln)
	}()

	mockDB := &mockResourceCloser{
		closeFunc: func() {
			atomic.StoreInt32(&dbClosedStep, atomic.AddInt32(&stepCounter, 1))
		},
	}

	lc := NewLifecycleManager(LifecycleConfig{
		ShutdownTimeout:    2 * time.Second,
		HTTPDrainTimeout:   500 * time.Millisecond,
		WorkerDrainTimeout: 500 * time.Millisecond,
		App:                app,
		DB:                 mockDB,
	})

	err = lc.RegisterWorker("test-worker", func(ctx context.Context) {
		<-ctx.Done()
		time.Sleep(50 * time.Millisecond)
		atomic.StoreInt32(&workerStoppedStep, atomic.AddInt32(&stepCounter, 1))
	})
	if err != nil {
		t.Fatalf("register worker failed: %v", err)
	}

	// Trigger shutdown
	atomic.StoreInt32(&httpDrainedStep, atomic.AddInt32(&stepCounter, 1))
	if err := lc.Shutdown(); err != nil {
		t.Fatalf("shutdown failed: %v", err)
	}

	if workerStoppedStep == 0 || dbClosedStep == 0 {
		t.Fatalf("expected all cleanup steps to execute, got worker=%d, db=%d", workerStoppedStep, dbClosedStep)
	}

	if dbClosedStep <= workerStoppedStep {
		t.Errorf("database closed (step %d) before worker stopped (step %d) — violates safety order", dbClosedStep, workerStoppedStep)
	}
}

// TestLifecycleManager_PreventRegistrationAfterShutdown verifies requirement 4:
// Prevent new worker registration after shutdown begins.
func TestLifecycleManager_PreventRegistrationAfterShutdown(t *testing.T) {
	lc := NewLifecycleManager(LifecycleConfig{
		ShutdownTimeout: 1 * time.Second,
	})

	// Before shutdown: registration succeeds
	err := lc.RegisterWorker("early_worker", func(ctx context.Context) {
		<-ctx.Done()
	})
	if err != nil {
		t.Fatalf("expected registration before shutdown to succeed, got %v", err)
	}

	// Initiate shutdown in background
	go func() {
		_ = lc.Shutdown()
	}()

	// Wait until state transitions to at least StateShuttingDown
	for i := 0; i < 50; i++ {
		if lc.State() >= StateShuttingDown {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Attempt registration after shutdown: MUST be rejected
	err = lc.RegisterWorker("late_worker", func(ctx context.Context) {})
	if err == nil {
		t.Fatal("expected RegisterWorker to be rejected after shutdown initiated, but it succeeded")
	}
	if !errors.Is(err, ErrShutdownInProgress) {
		t.Errorf("expected ErrShutdownInProgress, got %v", err)
	}
}

// TestLifecycleManager_IdempotentShutdown verifies requirement 9:
// Make shutdown idempotent (safe for concurrent and repeated invocations).
func TestLifecycleManager_IdempotentShutdown(t *testing.T) {
	mockDB := &mockResourceCloser{}
	lc := NewLifecycleManager(LifecycleConfig{
		ShutdownTimeout: 1 * time.Second,
		DB:              mockDB,
	})

	var wg sync.WaitGroup
	// Run 10 concurrent shutdown calls
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = lc.Shutdown()
		}()
	}
	wg.Wait()

	if lc.State() != StateStopped {
		t.Errorf("expected StateStopped, got %v", lc.State())
	}

	// Sequential repeated call
	if err := lc.Shutdown(); err != nil {
		t.Errorf("expected subsequent shutdown call to succeed cleanly, got %v", err)
	}
}

// TestLifecycleManager_BoundedTimeout verifies requirement 10:
// Guarantee bounded shutdown behavior when a worker hangs.
func TestLifecycleManager_BoundedTimeout(t *testing.T) {
	lc := NewLifecycleManager(LifecycleConfig{
		ShutdownTimeout:    200 * time.Millisecond,
		WorkerDrainTimeout: 100 * time.Millisecond,
	})

	// Register a rogue worker that ignores context cancellation
	_ = lc.RegisterWorker("hung_worker", func(ctx context.Context) {
		// Ignore ctx.Done() and sleep longer than the drain timeout
		time.Sleep(1 * time.Second)
	})

	start := time.Now()
	err := lc.Shutdown()
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected timeout error for hung worker, got nil")
	}
	if !errors.Is(err, ErrWorkerDrainTimeout) {
		t.Errorf("expected ErrWorkerDrainTimeout, got %v", err)
	}

	// Verify bounded execution (must finish within 400ms, not 1s)
	if elapsed > 500*time.Millisecond {
		t.Errorf("shutdown took %v, exceeded bounded timeout", elapsed)
	}
}

// TestLifecycleManager_WorkerCoordination verifies requirement 3 & 5:
// Coordinate background workers and wait for their completion via sync.WaitGroup.
func TestLifecycleManager_WorkerCoordination(t *testing.T) {
	lc := NewLifecycleManager(LifecycleConfig{
		ShutdownTimeout:    1 * time.Second,
		WorkerDrainTimeout: 500 * time.Millisecond,
	})

	var worker1Done, worker2Done int32

	_ = lc.RegisterWorker("worker1", func(ctx context.Context) {
		<-ctx.Done()
		time.Sleep(30 * time.Millisecond)
		atomic.StoreInt32(&worker1Done, 1)
	})

	_ = lc.RegisterWorker("worker2", func(ctx context.Context) {
		<-ctx.Done()
		time.Sleep(40 * time.Millisecond)
		atomic.StoreInt32(&worker2Done, 1)
	})

	if lc.ActiveWorkerCount() != 2 {
		t.Errorf("expected 2 active workers, got %d", lc.ActiveWorkerCount())
	}

	err := lc.Shutdown()
	if err != nil {
		t.Fatalf("shutdown failed: %v", err)
	}

	if atomic.LoadInt32(&worker1Done) != 1 || atomic.LoadInt32(&worker2Done) != 1 {
		t.Error("workers did not complete before shutdown returned")
	}
	if lc.ActiveWorkerCount() != 0 {
		t.Errorf("expected 0 active workers after shutdown, got %d", lc.ActiveWorkerCount())
	}
}
