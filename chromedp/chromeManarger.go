// Package chromedp provides a concurrent-safe Chrome browser instance manager
// with automatic lifecycle management and cleanup for web automation tasks.
package chromedp

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
)

type ChromeManager struct {
	instances map[string]*ChromeInstance
	mutex     sync.RWMutex
	maximum   int
	ttl       time.Duration
	stopChan  chan struct{} 
	once      sync.Once
	executeTimeout	  time.Duration
}

type ChromeInstance struct {
	Context   context.Context
	Cancel    context.CancelFunc
	TTL       time.Duration
	LastUsed  time.Time
}

var Manager ChromeManager

func InitManager(maximum int, ttl time.Duration, timeout time.Duration){
	Manager = *NewChromeManager(maximum, ttl, timeout)
}

func NewChromeManager(maximum int, ttl time.Duration, timeout time.Duration) *ChromeManager {
	cm := &ChromeManager{
		instances: make(map[string]*ChromeInstance, maximum),
		maximum:   maximum,
		ttl:       ttl,
		stopChan:  make(chan struct{}),
		executeTimeout: timeout,
	}
	
	// Start background goroutine for automatic cleanup
	go cm.startAutoCleanup()
	
	return cm
}

func (cm *ChromeManager) CreateVisibleInstance(allocOpts []chromedp.ExecAllocatorOption, opts ...chromedp.ContextOption) (string, *ChromeInstance, error) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	if len(cm.instances) >= cm.maximum {
		return "", nil, fmt.Errorf("exceeded maximum number of Chrome instances (%d)", cm.maximum)
	}
	
	
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), allocOpts...)
	
	instance := NewChromeInstanceWithAllocator(allocCtx, cm.ttl, opts...)
	
	id := uuid.New().String()
	cm.instances[id] = instance
	
	return id, instance, nil
}

func NewChromeInstanceWithAllocator(allocCtx context.Context, ttl time.Duration, opts ...chromedp.ContextOption) *ChromeInstance {
	ctx, cancel := chromedp.NewContext(allocCtx, opts...)
	now := time.Now()
	
	return &ChromeInstance{
		Context:   ctx,
		Cancel:    cancel,
		TTL:       ttl,
		LastUsed:  now,
	}
}

func NewChromeInstance(ttl time.Duration, opts ...chromedp.ContextOption) *ChromeInstance {
	ctx, cancel := chromedp.NewContext(context.Background(), opts...)
	now := time.Now()
	
	return &ChromeInstance{
		Context:   ctx,
		Cancel:    cancel,
		TTL:       ttl,
		LastUsed:  now,
	}
}

func (cm *ChromeManager) CreateNewInstance(opts ...chromedp.ContextOption) (string, *ChromeInstance, error) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	// Check if reached the maximum number of instances
	if len(cm.instances) >= cm.maximum {
		return "", nil, fmt.Errorf("exceeded maximum number of Chrome instances (%d); stop some instances before creating a new one", cm.maximum)
	}
	
	// Create new Chrome instance
	instance := NewChromeInstance(cm.ttl, opts...)
	
	// Generate unique ID for this instance
	id := uuid.New().String()
	cm.instances[id] = instance
	
	return id, instance, nil
}

// GetInstance retrieves an existing Chrome instance by ID
func (cm *ChromeManager) GetInstance(id string) (*ChromeInstance, error) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	instance, exists := cm.instances[id]
	if !exists {
		return nil, errors.New("chrome instance not found; it may have timed out or been closed, create a new instance first")
	}
	
	// Check if instance has been idle for too long
	if time.Since(instance.LastUsed) > instance.TTL {
		// Clean up idle instance
		instance.Cancel()
		delete(cm.instances, id)
		return nil, errors.New("chrome instance has been idle too long and expired; create a new instance")
	}
	
	// Update last used timestamp
	instance.LastUsed = time.Now()
	
	return instance, nil
}

func (cm *ChromeManager) CloseInstance(id string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	instance, exists := cm.instances[id]
	if !exists {
		return
	}
	
	instance.Cancel()
	delete(cm.instances, id)
}

// CloseAll closes all Chrome instances
func (cm *ChromeManager) CloseAll() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	for id, instance := range cm.instances {
		instance.Cancel()
		delete(cm.instances, id)
	}
}

// Shutdown gracefully shuts down the Chrome manager and all instances
func (cm *ChromeManager) Shutdown() {
	cm.once.Do(func() {
		close(cm.stopChan)
		cm.CloseAll()
	})
}

// InstanceInfo provides information about a Chrome instance
type InstanceInfo struct {
	ID        string
	LastUsed  time.Time
	TTL       time.Duration
	IsExpired bool
}

// GetInstancesInfo returns information about all managed instances
func (cm *ChromeManager) GetInstancesInfo() map[string]InstanceInfo {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	info := make(map[string]InstanceInfo)
	for id, instance := range cm.instances {
		info[id] = InstanceInfo{
			ID:        id,
			LastUsed:  instance.LastUsed,
			TTL:       instance.TTL,
			IsExpired: time.Since(instance.LastUsed) > instance.TTL, // Check based on LastUsed
		}
	}
	
	return info
}

// GetInstanceCount returns the current number of active instances
func (cm *ChromeManager) GetInstanceCount() int {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return len(cm.instances)
}

// startAutoCleanup runs in a separate goroutine to automatically clean up expired instances
func (cm *ChromeManager) startAutoCleanup() {
	ticker := time.NewTicker(time.Minute) // Check every minute
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			cm.cleanupExpiredInstances()
		case <-cm.stopChan:
			return
		}
	}
}

// cleanupExpiredInstances removes all idle Chrome instances
func (cm *ChromeManager) cleanupExpiredInstances() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	now := time.Now()
	idleIDs := make([]string, 0)
	
	// Find all idle instances (based on LastUsed)
	for id, instance := range cm.instances {
		if now.Sub(instance.LastUsed) > instance.TTL {
			idleIDs = append(idleIDs, id)
		}
	}
	
	// Clean up idle instances
	for _, id := range idleIDs {
		if instance, exists := cm.instances[id]; exists {
			instance.Cancel()
			delete(cm.instances, id)
		}
	}
}

// Execute runs ChromeDP actions on a specific instance
func (cm *ChromeManager) Execute(id string, actions ...chromedp.Action) error {
	instance, err := cm.GetInstance(id)
	if err != nil {
		return err
	}
	
	done := make(chan error, 1)
	
	go func() {
		done <- chromedp.Run(instance.Context, actions...)
	}()
	
	select {
	case err := <-done:
		return err
	case <-time.After(cm.executeTimeout):
		return fmt.Errorf("chromedp execute timeout: %v", cm.executeTimeout)
	}
}

// ExecuteWithTimeout runs ChromeDP actions with a timeout on a specific instance
func (cm *ChromeManager) ExecuteWithTimeout(id string, timeout time.Duration, actions ...chromedp.Action) error {
	instance, err := cm.GetInstance(id)
	if err != nil {
		return err
	}
	
	ctx, cancel := context.WithTimeout(instance.Context, timeout)
	defer cancel()
	
	return chromedp.Run(ctx, actions...)
}
