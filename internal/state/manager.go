package state

import (
	"sync"
	"time"

	"x.localhost/rvabot/internal/logger"
	"x.localhost/rvabot/internal/states"
)

// UserStateEntry содержит состояние пользователя с метаданными
type UserStateEntry struct {
	State     states.State
	LastSeen  time.Time
	CreatedAt time.Time
}

// Manager управляет состояниями пользователей с TTL
type Manager struct {
	states   map[int]*UserStateEntry
	mutex    sync.RWMutex
	ttl      time.Duration
	cleanup  time.Duration
	stopChan chan struct{}
}

// NewManager создает новый менеджер состояний
func NewManager(ttl, cleanupInterval time.Duration) *Manager {
	manager := &Manager{
		states:   make(map[int]*UserStateEntry),
		ttl:      ttl,
		cleanup:  cleanupInterval,
		stopChan: make(chan struct{}),
	}

	// Запускаем cleanup goroutine
	go manager.cleanupLoop()

	return manager
}

// GetState возвращает состояние пользователя
func (m *Manager) GetState(userID int) (states.State, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	entry, exists := m.states[userID]
	if !exists {
		return states.State{}, false
	}

	// Обновляем время последнего обращения
	entry.LastSeen = time.Now()

	return entry.State, true
}

// SetState устанавливает состояние пользователя
func (m *Manager) SetState(userID int, state states.State) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	now := time.Now()
	entry, exists := m.states[userID]

	if exists {
		entry.State = state
		entry.LastSeen = now
	} else {
		m.states[userID] = &UserStateEntry{
			State:     state,
			LastSeen:  now,
			CreatedAt: now,
		}
		logger.UserInfo(userID, "Новый пользователь")
	}
}

// GetOrCreateState возвращает существующее состояние или создает новое
func (m *Manager) GetOrCreateState(userID int) states.State {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	entry, exists := m.states[userID]
	if !exists {
		now := time.Now()
		entry = &UserStateEntry{
			State:     states.SetStart(),
			LastSeen:  now,
			CreatedAt: now,
		}
		m.states[userID] = entry
		logger.UserInfo(userID, "Новый пользователь")
	} else {
		entry.LastSeen = time.Now()
	}

	return entry.State
}

// DeleteState удаляет состояние пользователя
func (m *Manager) DeleteState(userID int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.states, userID)
	logger.UserInfo(userID, "Состояние пользователя удалено")
}

// cleanupLoop периодически очищает устаревшие состояния
func (m *Manager) cleanupLoop() {
	ticker := time.NewTicker(m.cleanup)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.cleanupExpiredStates()
		case <-m.stopChan:
			return
		}
	}
}

// cleanupExpiredStates удаляет состояния, которые не использовались дольше TTL
func (m *Manager) cleanupExpiredStates() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	now := time.Now()
	expiredCount := 0

	for userID, entry := range m.states {
		if now.Sub(entry.LastSeen) > m.ttl {
			delete(m.states, userID)
			expiredCount++
		}
	}

	if expiredCount > 0 {
		logger.BotInfo("Очищено %d устаревших состояний пользователей", expiredCount)
	}
}

// GetStats возвращает статистику менеджера состояний
func (m *Manager) GetStats() map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	now := time.Now()
	activeCount := 0
	expiredCount := 0

	for _, entry := range m.states {
		if now.Sub(entry.LastSeen) <= m.ttl {
			activeCount++
		} else {
			expiredCount++
		}
	}

	return map[string]interface{}{
		"total_states":   len(m.states),
		"active_states":  activeCount,
		"expired_states": expiredCount,
		"ttl_seconds":    int(m.ttl.Seconds()),
	}
}

// Shutdown корректно завершает работу менеджера
func (m *Manager) Shutdown() {
	close(m.stopChan)
	logger.BotInfo("State manager shutdown completed")
}
