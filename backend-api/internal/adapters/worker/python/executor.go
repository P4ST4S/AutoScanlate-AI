package python

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/P4ST4S/manga-translator/backend-api/internal/infrastructure/config"
	"github.com/P4ST4S/manga-translator/backend-api/internal/ports"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type pythonExecutor struct {
	pythonPath string
	workerPath string
	timeout    int
	logger     *zap.Logger
}

// NewPythonExecutor creates a new Python worker executor
func NewPythonExecutor(cfg *config.WorkerConfig, logger *zap.Logger) ports.WorkerExecutor {
	return &pythonExecutor{
		pythonPath: cfg.PythonPath,
		workerPath: cfg.WorkerPath,
		timeout:    int(cfg.Timeout.Seconds()),
		logger:     logger,
	}
}

func (e *pythonExecutor) Translate(
	ctx context.Context,
	inputPath string,
	onProgress ports.ProgressCallback,
) (*ports.TranslationOutput, error) {
	e.logger.Info("starting translation",
		zap.String("input_path", inputPath),
	)

	// Create job-specific temp directory
	jobID := uuid.New().String()
	tempDir := filepath.Join("storage", "temp", jobID)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer func() {
		// Cleanup temp directory
		os.RemoveAll(tempDir)
	}()

	// Build absolute path to main.py
	mainPyPath := filepath.Join(e.workerPath, "main.py")
	if _, err := os.Stat(mainPyPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Python worker not found at: %s", mainPyPath)
	}

	// Convert to absolute path for Python
	absInputPath, err := filepath.Abs(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Build command
	cmd := exec.CommandContext(ctx, e.pythonPath, mainPyPath, absInputPath)
	cmd.Dir = e.workerPath // Set working directory to ai-worker

	// Set environment variables
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("TEMP_DIR=%s", tempDir),
	)

	// Capture stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the process
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start worker: %w", err)
	}

	// Parse output in goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	go e.parseStdout(stdout, onProgress, &wg)
	go e.parseStderr(stderr, &wg)

	// Wait for process to complete
	done := make(chan error, 1)
	go func() {
		wg.Wait()
		done <- cmd.Wait()
	}()

	// Wait for completion or context cancellation
	var processErr error
	select {
	case <-ctx.Done():
		// Context cancelled, kill process
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return nil, ctx.Err()
	case processErr = <-done:
		if processErr != nil {
			return nil, fmt.Errorf("worker process failed: %w", processErr)
		}
	}

	// Find output files
	output, err := e.findOutputFiles(absInputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to find output files: %w", err)
	}

	e.logger.Info("translation completed",
		zap.String("input_path", inputPath),
		zap.Int("pages", len(output.Pages)),
	)

	return output, nil
}

func (e *pythonExecutor) parseStdout(reader io.Reader, onProgress ports.ProgressCallback, wg *sync.WaitGroup) {
	defer wg.Done()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		e.logger.Debug("worker output", zap.String("line", line))

		// Parse progress from output
		progress, message := parseProgressLine(line)
		if progress >= 0 {
			if onProgress != nil {
				onProgress(progress, message)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		e.logger.Error("error reading stdout", zap.Error(err))
	}
}

func (e *pythonExecutor) parseStderr(reader io.Reader, wg *sync.WaitGroup) {
	defer wg.Done()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		e.logger.Warn("worker stderr", zap.String("line", line))
	}

	if err := scanner.Err(); err != nil {
		e.logger.Error("error reading stderr", zap.Error(err))
	}
}

func (e *pythonExecutor) findOutputFiles(inputPath string) (*ports.TranslationOutput, error) {
	output := &ports.TranslationOutput{
		Pages: []ports.PageOutput{},
	}

	// The Python worker creates output files in its working directory (ai-worker/)
	// not in the directory of the input file
	inputBase := filepath.Base(inputPath)
	inputExt := filepath.Ext(inputBase)
	inputName := strings.TrimSuffix(inputBase, inputExt)

	// Check if input was a ZIP file
	if strings.HasSuffix(strings.ToLower(inputPath), ".zip") {
		// Look for {name}_translated.zip in the worker directory
		translatedZip := filepath.Join(e.workerPath, inputName+"_translated.zip")

		if _, err := os.Stat(translatedZip); os.IsNotExist(err) {
			return nil, fmt.Errorf("translated zip not found: %s", translatedZip)
		}

		output.OutputPath = translatedZip

		// Note: For ZIP files, we'll need to extract and catalog pages
		// This will be handled by the translation service
		e.logger.Info("found translated zip", zap.String("path", translatedZip))

		return output, nil
	}

	// Single image file - Python worker creates translated_{name}.jpg in worker directory
	translatedPath := filepath.Join(e.workerPath, "translated_"+inputName+".jpg")

	if _, err := os.Stat(translatedPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("translated image not found: %s", translatedPath)
	}

	output.OutputPath = translatedPath
	output.Pages = []ports.PageOutput{
		{
			PageNumber:     1,
			OriginalPath:   inputPath,
			TranslatedPath: translatedPath,
		},
	}

	e.logger.Info("found translated image", zap.String("path", translatedPath))

	return output, nil
}
