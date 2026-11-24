package ai

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type PythonSupervisor struct {
	cmd         *exec.Cmd
	port        int
	authToken   string
	baseURL     string
	httpClient  *http.Client
	mu          sync.RWMutex
	stopped     bool
	geminiKey   string
	geminiModel string
}

func NewPythonSupervisor(geminiKey, geminiModel string) *PythonSupervisor {
	return &PythonSupervisor{
		httpClient: &http.Client{
			Timeout: 25 * time.Second,
		},
		geminiKey:   geminiKey,
		geminiModel: geminiModel,
	}
}

func (p *PythonSupervisor) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.cmd != nil && p.cmd.Process != nil {
		return nil
	}

	tokenBytes := make([]byte, 24)
	if _, err := rand.Read(tokenBytes); err != nil {
		return fmt.Errorf("failed to generate random auth token: %w", err)
	}
	p.authToken = hex.EncodeToString(tokenBytes)

	pyBin, err := findPythonBinary()
	if err != nil {
		return err
	}

	scriptPath, policyPath, err := resolveScriptAndPolicyPaths()
	if err != nil {
		return err
	}

	args := []string{
		scriptPath,
		"--port", "0",
		"--auth-token", p.authToken,
		"--policy", policyPath,
	}
	if p.geminiModel != "" {
		args = append(args, "--gemini-model", p.geminiModel)
	}

	cmd := exec.Command(pyBin, args...)
	env := os.Environ()
	if p.geminiKey != "" {
		env = append(env, "GEMINI_API_KEY="+p.geminiKey)
	}
	if p.geminiModel != "" {
		env = append(env, "GEMINI_MODEL="+p.geminiModel)
	}
	cmd.Env = env

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe for python: %w", err)
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start python ai process: %w", err)
	}
	p.cmd = cmd
	p.stopped = false

	portChan := make(chan int, 1)
	errChan := make(chan error, 1)

	go func() {
		scanner := bufio.NewScanner(stdout)
		foundPort := false
		for scanner.Scan() {
			line := scanner.Text()
			if !foundPort && strings.HasPrefix(line, "HUNTIQ_AI_PORT=") {
				pStr := strings.TrimPrefix(line, "HUNTIQ_AI_PORT=")
				port, err := strconv.Atoi(strings.TrimSpace(pStr))
				if err == nil && port > 0 {
					foundPort = true
					portChan <- port
				}
			}
		}
		if !foundPort {
			if err := scanner.Err(); err != nil {
				errChan <- err
			} else {
				errChan <- fmt.Errorf("python process exited without publishing port")
			}
		}
	}()

	select {
	case port := <-portChan:
		p.port = port
		p.baseURL = fmt.Sprintf("http://127.0.0.1:%d", port)
		log.Printf("[HUNTIQ AI] Python AI process successfully initialized on %s (PID: %d)", p.baseURL, cmd.Process.Pid)
	case err := <-errChan:
		_ = cmd.Process.Kill()
		return fmt.Errorf("python ai startup failed: %w", err)
	case <-time.After(10 * time.Second):
		_ = cmd.Process.Kill()
		return fmt.Errorf("timeout waiting for python ai startup (10s)")
	}

	// Verify health check
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := p.healthCheck(ctx); err != nil {
		log.Printf("[HUNTIQ AI] Initial health check warning: %v", err)
	}

	return nil
}

func (p *PythonSupervisor) healthCheck(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.baseURL+"/health", nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Internal-Token", p.authToken)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check returned status %d", resp.StatusCode)
	}
	return nil
}

func (p *PythonSupervisor) Classify(ctx context.Context, req ClassifyRequest) (*ClassifyResponse, error) {
	p.mu.RLock()
	baseURL := p.baseURL
	token := p.authToken
	p.mu.RUnlock()

	if baseURL == "" {
		return nil, fmt.Errorf("python ai engine is not running")
	}

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/classify", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Internal-Token", token)

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error contacting python ai engine: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("python ai engine error (%d): %s", resp.StatusCode, string(b))
	}

	var res ClassifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("invalid json from python ai engine: %w", err)
	}
	return &res, nil
}

func (p *PythonSupervisor) Chat(ctx context.Context, req PythonChatRequest) (*ChatResponse, error) {
	p.mu.RLock()
	baseURL := p.baseURL
	token := p.authToken
	p.mu.RUnlock()

	if baseURL == "" {
		return nil, fmt.Errorf("python ai engine is not running")
	}

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/chat", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Internal-Token", token)

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error contacting python ai engine: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("python ai engine error (%d): %s", resp.StatusCode, string(b))
	}

	var res ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("invalid json from python ai engine: %w", err)
	}
	return &res, nil
}

func (p *PythonSupervisor) Translate(ctx context.Context, req TranslateRequest) (*TranslateResponse, error) {
	p.mu.RLock()
	baseURL := p.baseURL
	token := p.authToken
	p.mu.RUnlock()

	if baseURL == "" {
		return nil, fmt.Errorf("python ai engine is not running")
	}

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/translate", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Internal-Token", token)

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error contacting python ai engine: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("python ai engine error (%d): %s", resp.StatusCode, string(b))
	}

	var res TranslateResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("invalid json from python ai engine: %w", err)
	}
	return &res, nil
}

func (p *PythonSupervisor) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.stopped || p.cmd == nil || p.cmd.Process == nil {
		return
	}
	p.stopped = true

	log.Println("[HUNTIQ AI] Stopping Python AI process gracefully...")

	// Attempt graceful shutdown endpoint
	if p.baseURL != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/shutdown", nil)
		if err == nil {
			req.Header.Set("X-Internal-Token", p.authToken)
			_, _ = p.httpClient.Do(req)
		}
		cancel()
	}

	done := make(chan error, 1)
	go func() {
		done <- p.cmd.Wait()
	}()

	select {
	case <-done:
		log.Println("[HUNTIQ AI] Python AI process stopped cleanly.")
	case <-time.After(3 * time.Second):
		log.Println("[HUNTIQ AI] Python process did not exit in time; killing process...")
		_ = p.cmd.Process.Kill()
	}
}

func findPythonBinary() (string, error) {
	for _, bin := range []string{"python3", "python"} {
		p, err := exec.LookPath(bin)
		if err == nil && p != "" {
			return p, nil
		}
	}
	return "", fmt.Errorf("no python3 or python binary found in PATH")
}

func resolveScriptAndPolicyPaths() (string, string, error) {
	dir, err := os.Getwd()
	if err != nil {
		dir = "."
	}

	var scriptPath string
	for i := 0; i < 6; i++ {
		try1 := filepath.Join(dir, "backend", "ai", "script.py")
		if _, err := os.Stat(try1); err == nil {
			scriptPath = try1
			break
		}
		try2 := filepath.Join(dir, "ai", "script.py")
		if _, err := os.Stat(try2); err == nil {
			scriptPath = try2
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	if scriptPath == "" {
		return "", "", fmt.Errorf("could not locate backend/ai/script.py")
	}

	policyPath := filepath.Join(filepath.Dir(scriptPath), "policy.json")
	if _, err := os.Stat(policyPath); err != nil {
		return "", "", fmt.Errorf("could not locate backend/ai/policy.json at %s", policyPath)
	}

	return scriptPath, policyPath, nil
}
