package inertia

import (
	"context"
	"fmt"
	"sync"

	gonertia "github.com/romsar/gonertia/v2"
)

type InmemFlashProvider struct {
	mu           sync.RWMutex
	flash        map[string]Props
	errors       map[string]gonertia.ValidationErrors
	clearHistory map[string]bool
	getSessionID func(ctx context.Context) (string, error)
}

// Constructor
func NewInmemFlashProvider() *InmemFlashProvider {
	return &InmemFlashProvider{
		flash:        make(map[string]Props),
		errors:       make(map[string]gonertia.ValidationErrors),
		clearHistory: make(map[string]bool),
	}
}

// Setter for session ID callback
func (p *InmemFlashProvider) SetSessionIDFunc(f func(ctx context.Context) (string, error)) {
	p.getSessionID = f
}

// Internal helper to get session ID
func (p *InmemFlashProvider) sessionID(ctx context.Context) (string, error) {
	if p.getSessionID == nil {
		return "", fmt.Errorf("session ID function not set")
	}
	return p.getSessionID(ctx)
}

// FlashErrors stores validation errors
func (p *InmemFlashProvider) FlashErrors(ctx context.Context, errors gonertia.ValidationErrors) error {
	sessionID, err := p.sessionID(ctx)
	if err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.errors[sessionID] = errors
	return nil
}

// Flash stores props
func (p *InmemFlashProvider) Flash(ctx context.Context, props Props) error {
	sessionID, err := p.sessionID(ctx)
	if err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.flash[sessionID] = props
	return nil
}

// GetFlash retrieves and deletes flash
func (p *InmemFlashProvider) GetFlash(ctx context.Context) (Props, error) {
	sessionID, err := p.sessionID(ctx)
	if err != nil {
		return Props{}, err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	props := p.flash[sessionID]
	delete(p.flash, sessionID)
	return props, nil
}

// GetErrors retrieves and deletes validation errors
func (p *InmemFlashProvider) GetErrors(ctx context.Context) (gonertia.ValidationErrors, error) {
	sessionID, err := p.sessionID(ctx)
	if err != nil {
		return nil, err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	errors := p.errors[sessionID]
	delete(p.errors, sessionID)
	return errors, nil
}

// FlashClearHistory marks session history for clearing
func (p *InmemFlashProvider) FlashClearHistory(ctx context.Context) error {
	sessionID, err := p.sessionID(ctx)
	if err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.clearHistory[sessionID] = true
	return nil
}

// ShouldClearHistory checks and clears history flag
func (p *InmemFlashProvider) ShouldClearHistory(ctx context.Context) (bool, error) {
	sessionID, err := p.sessionID(ctx)
	if err != nil {
		return false, err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	clearHistory := p.clearHistory[sessionID]
	delete(p.clearHistory, sessionID)
	return clearHistory, nil
}
