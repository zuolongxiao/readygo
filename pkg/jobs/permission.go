package jobs

import (
	"sync"
	"time"

	"readygo/models"
	"readygo/services"
)

// Permission permissions
type permission struct {
	mu          sync.RWMutex
	permissions []string
}

var ps = permission{}

func (p *permission) set() {
	s := services.New(&models.Permission{})
	var list []models.PermissionView
	conds := map[string]interface{}{
		"is_enabled": "Y",
	}
	if err := s.GetRows(&list, conds); err != nil {
		return
	}

	p.mu.Lock()
	p.permissions = p.permissions[:0]
	for _, v := range list {
		p.permissions = append(p.permissions, v.Name)
	}
	p.mu.Unlock()
}

func (p *permission) get() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.permissions
}

// SetPermissions update permissions periodically
func SetPermissions() {
	for {
		ps.set()
		time.Sleep(time.Second * 600)
	}
}

// GetPermissions get permissions
func GetPermissions() []string {
	return ps.get()
}
