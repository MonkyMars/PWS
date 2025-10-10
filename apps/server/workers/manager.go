package workers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/MonkyMars/PWS/config"
	"github.com/MonkyMars/PWS/types"
	"github.com/gofiber/fiber/v3"
)

// WorkerManager coordinates all background workers with proper dependency injection
type WorkerManager struct {
	auditWorker   *AuditWorker
	healthWorker  *HealthWorker
	cleanupWorker *CleanupWorker
	logger        *config.Logger
	cfg           *config.Config
	mu            sync.RWMutex
	running       bool
}

// AuditWorker handles audit log processing
type AuditWorker struct {
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	auditChan chan types.AuditLog
	running   bool
	mu        sync.RWMutex
	stats     AuditStats
	logger    *config.Logger
	cfg       *config.Config
	dlq       *DeadLetterQueue
}

// HealthWorker handles health monitoring
type HealthWorker struct {
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	healthChan    chan types.HealthLog
	services      map[string]*RouteService
	running       bool
	mu            sync.RWMutex
	lastFlushTime time.Time
	logger        *config.Logger
	cfg           *config.Config
}

// CleanupWorker handles periodic cleanup tasks
type CleanupWorker struct {
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	running bool
	mu      sync.RWMutex
	logger  *config.Logger
	cfg     *config.Config
}

// AuditStats tracks audit worker statistics
type AuditStats struct {
	TotalProcessed int64
	TotalDropped   int64
	FailureCount   int
	LastFlushTime  time.Time
}

// Global manager instance (maintained for backward compatibility)
var (
	globalManager *WorkerManager
	managerOnce   sync.Once
)

// NewWorkerManager creates a new worker manager with all dependencies injected
func NewWorkerManager(cfg *config.Config, logger *config.Logger) *WorkerManager {
	return &WorkerManager{
		cfg:    cfg,
		logger: logger,
	}
}

// GetGlobalManager returns the global manager instance (for backward compatibility)
func GetGlobalManager() *WorkerManager {
	managerOnce.Do(func() {
		defer func() {
			if r := recover(); r != nil {
				// If there's a panic during initialization, create a safe default manager
				globalManager = &WorkerManager{
					running: false,
				}
			}
		}()

		cfg := config.Get()
		if cfg == nil {
			// Create a safe default manager if config is nil
			globalManager = &WorkerManager{
				running: false,
			}
			return
		}

		logger := config.SetupLogger()
		if logger == nil {
			// Create a safe default manager if logger is nil
			globalManager = &WorkerManager{
				cfg:     cfg,
				running: false,
			}
			return
		}

		globalManager = NewWorkerManager(cfg, logger)
	})
	return globalManager
}

// Start initializes and starts all workers
func (wm *WorkerManager) Start() error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	if wm.running {
		return fmt.Errorf("worker manager already running")
	}

	// Initialize workers
	wm.auditWorker = wm.newAuditWorker()
	wm.healthWorker = wm.newHealthWorker()
	wm.cleanupWorker = wm.newCleanupWorker()

	// Start workers in dependency order
	if wm.cfg.Audit.Enabled {
		if err := wm.auditWorker.Start(); err != nil {
			return fmt.Errorf("failed to start audit worker: %w", err)
		}
		wm.logger.Info("Audit worker started")
	}

	if wm.cfg.Health.Enabled {
		if err := wm.healthWorker.Start(); err != nil {
			return fmt.Errorf("failed to start health worker: %w", err)
		}
		wm.logger.Info("Health worker started")
	}

	if wm.cfg.Audit.Enabled && wm.cfg.Audit.RetentionDays > 0 {
		if err := wm.cleanupWorker.Start(); err != nil {
			return fmt.Errorf("failed to start cleanup worker: %w", err)
		}
		wm.logger.Info("Cleanup worker started")
	}

	wm.running = true
	wm.logger.Info("Worker manager started successfully")
	return nil
}

// Stop gracefully shuts down all workers
func (wm *WorkerManager) Stop(ctx context.Context) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	if !wm.running {
		return nil
	}

	wm.logger.Info("Stopping worker manager...")

	// Create a channel to collect errors
	errChan := make(chan error, 3)
	var wg sync.WaitGroup

	// Stop workers concurrently with timeout
	if wm.auditWorker != nil {
		wg.Go(func() {
			if err := wm.auditWorker.Stop(ctx); err != nil {
				errChan <- fmt.Errorf("audit worker stop error: %w", err)
			}
		})
	}

	if wm.healthWorker != nil {
		wg.Go(func() {
			if err := wm.healthWorker.Stop(ctx); err != nil {
				errChan <- fmt.Errorf("health worker stop error: %w", err)
			}
		})
	}

	if wm.cleanupWorker != nil {
		wg.Go(func() {
			if err := wm.cleanupWorker.Stop(ctx); err != nil {
				errChan <- fmt.Errorf("cleanup worker stop error: %w", err)
			}
		})
	}

	// Wait for all workers to stop or timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		wm.logger.Info("All workers stopped successfully")
	case <-ctx.Done():
		wm.logger.Warn("Worker shutdown timed out")
		return ctx.Err()
	}

	// Collect any errors
	close(errChan)
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	wm.running = false

	if len(errors) > 0 {
		return fmt.Errorf("worker shutdown errors: %v", errors)
	}

	return nil
}

// DiscoverRoutes auto-discovers routes for health monitoring
func (wm *WorkerManager) DiscoverRoutes(app *fiber.App) {
	if wm.healthWorker != nil {
		wm.healthWorker.DiscoverRoutes(app)
	}
}

// AddAuditLog adds an audit log entry (backward compatibility)
func (wm *WorkerManager) AddAuditLog(entry types.AuditLog) {
	if wm.auditWorker != nil {
		wm.auditWorker.AddLog(entry)
	}
}

// RecordHealthMetric records a health metric (backward compatibility)
func (wm *WorkerManager) RecordHealthMetric(serviceName string, statusCode int, latency time.Duration) {
	if wm.healthWorker != nil {
		wm.healthWorker.RecordRequest(serviceName, statusCode, latency)
	}
}

// HealthStatus returns the overall health status of all workers
func (wm *WorkerManager) HealthStatus() map[string]any {
	if wm == nil {
		return map[string]any{
			"manager_running": false,
			"timestamp":       time.Now(),
			"is_healthy":      false,
			"error":           "worker manager not initialized",
		}
	}

	if wm.cfg == nil {
		return map[string]any{
			"manager_running": false,
			"timestamp":       time.Now(),
			"is_healthy":      false,
			"error":           "worker manager configuration is nil",
		}
	}

	wm.mu.RLock()
	defer wm.mu.RUnlock()

	status := map[string]any{
		"manager_running": wm.running,
		"timestamp":       time.Now(),
	}

	if wm.auditWorker != nil {
		status["audit"] = wm.auditWorker.HealthStatus()
	} else {
		status["audit"] = map[string]any{
			"enabled":        false,
			"worker_running": false,
			"is_healthy":     false,
		}
	}

	if wm.healthWorker != nil {
		status["health"] = wm.healthWorker.HealthStatus()
	} else {
		status["health"] = map[string]any{
			"enabled":        false,
			"worker_running": false,
			"is_healthy":     false,
		}
	}

	if wm.cleanupWorker != nil {
		status["cleanup"] = wm.cleanupWorker.HealthStatus()
	} else {
		status["cleanup"] = map[string]any{
			"enabled":        false,
			"worker_running": false,
			"is_healthy":     false,
		}
	}

	// Overall health calculation
	isHealthy := wm.running
	if wm.cfg != nil && wm.cfg.Audit.Enabled && wm.auditWorker != nil {
		auditHealth := wm.auditWorker.HealthStatus()
		if auditHealth != nil {
			if healthy, ok := auditHealth["is_healthy"].(bool); ok {
				isHealthy = isHealthy && healthy
			} else {
				isHealthy = false
			}
		} else {
			isHealthy = false
		}
	}

	if wm.cfg != nil && wm.cfg.Health.Enabled && wm.healthWorker != nil {
		healthHealth := wm.healthWorker.HealthStatus()
		if healthHealth != nil {
			if healthy, ok := healthHealth["is_healthy"].(bool); ok {
				isHealthy = isHealthy && healthy
			} else {
				isHealthy = false
			}
		} else {
			isHealthy = false
		}
	}

	status["is_healthy"] = isHealthy
	return status
}

// TriggerCleanup manually triggers cleanup operations
func (wm *WorkerManager) TriggerCleanup() error {
	if wm.cleanupWorker != nil {
		return wm.cleanupWorker.TriggerCleanup()
	}
	return fmt.Errorf("cleanup worker not available")
}

// Worker factory methods
func (wm *WorkerManager) newAuditWorker() *AuditWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &AuditWorker{
		ctx:       ctx,
		cancel:    cancel,
		auditChan: make(chan types.AuditLog, wm.cfg.Audit.ChannelSize),
		logger:    wm.logger,
		cfg:       wm.cfg,
		dlq:       NewDeadLetterQueue(wm.cfg, wm.logger),
		stats: AuditStats{
			LastFlushTime: time.Now(),
		},
	}
}

func (wm *WorkerManager) newHealthWorker() *HealthWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &HealthWorker{
		ctx:           ctx,
		cancel:        cancel,
		healthChan:    make(chan types.HealthLog, wm.cfg.Health.ChannelSize),
		services:      make(map[string]*RouteService),
		logger:        wm.logger,
		cfg:           wm.cfg,
		lastFlushTime: time.Now(),
	}
}

func (wm *WorkerManager) newCleanupWorker() *CleanupWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &CleanupWorker{
		ctx:    ctx,
		cancel: cancel,
		logger: wm.logger,
		cfg:    wm.cfg,
	}
}

// Backward compatibility functions
func StartAuditWorker() {
	manager := GetGlobalManager()
	if manager.auditWorker == nil && manager.cfg.Audit.Enabled {
		manager.auditWorker = manager.newAuditWorker()
		manager.auditWorker.Start()
	}
}

func StopAuditWorker() {
	manager := GetGlobalManager()
	if manager.auditWorker != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		manager.auditWorker.Stop(ctx)
	}
}

func StartHealthLogWorker() {
	manager := GetGlobalManager()
	if manager.healthWorker == nil && manager.cfg.Health.Enabled {
		manager.healthWorker = manager.newHealthWorker()
		manager.healthWorker.Start()
	}
}

func StopHealthLogWorker() {
	manager := GetGlobalManager()
	if manager.healthWorker != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		manager.healthWorker.Stop(ctx)
	}
}

func StartCleanupScheduler() {
	manager := GetGlobalManager()
	if manager.cleanupWorker == nil && manager.cfg.Audit.Enabled {
		manager.cleanupWorker = manager.newCleanupWorker()
		manager.cleanupWorker.Start()
	}
}

func StopAuditCleanupScheduler() {
	manager := GetGlobalManager()
	if manager.cleanupWorker != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		manager.cleanupWorker.Stop(ctx)
	}
}

func AddAuditLog(entry types.AuditLog) {
	manager := GetGlobalManager()
	manager.AddAuditLog(entry)
}

func DiscoverRoutes(app *fiber.App) {
	manager := GetGlobalManager()
	manager.DiscoverRoutes(app)
}

func TriggerCleanupNow() error {
	manager := GetGlobalManager()
	return manager.TriggerCleanup()
}

// Health status functions for backward compatibility
func AuditHealthStatus() map[string]any {
	manager := GetGlobalManager()
	if manager == nil {
		return map[string]any{
			"enabled":        false,
			"worker_running": false,
			"is_healthy":     false,
			"error":          "worker manager not initialized",
		}
	}

	if manager.auditWorker != nil {
		status := manager.auditWorker.HealthStatus()
		if status != nil {
			return status
		}
	}
	return map[string]any{
		"enabled":        false,
		"worker_running": false,
		"is_healthy":     false,
		"error":          "audit worker not initialized",
	}
}

func ServiceHealthStatus() map[string]any {
	manager := GetGlobalManager()
	if manager == nil {
		return map[string]any{
			"enabled":        false,
			"worker_running": false,
			"is_healthy":     false,
			"error":          "worker manager not initialized",
		}
	}

	if manager.healthWorker != nil {
		status := manager.healthWorker.HealthStatus()
		if status != nil {
			return status
		}
	}
	return map[string]any{
		"enabled":        false,
		"worker_running": false,
		"is_healthy":     false,
		"error":          "health worker not initialized",
	}
}

type WorkerManagerInterface interface {
	Start() error
	Stop(ctx context.Context) error
	DiscoverRoutes(app *fiber.App)
	AddAuditLog(entry types.AuditLog)
	RecordHealthMetric(serviceName string, statusCode int, latency time.Duration)
	HealthStatus() map[string]any
	TriggerCleanup() error
}
