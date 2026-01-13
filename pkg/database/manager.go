package database

import (
	"king-starter/config"
	"king-starter/pkg/logger"
	"sync"

	"go.uber.org/zap"
)

// InstanceManager 数据库实例管理器，支持多实例
type InstanceManager struct {
	instances map[string]*DB
	once      sync.Once
	mu        sync.RWMutex
	logger    *logger.Logger
}

// NewInstanceManager 创建数据库实例管理器
func NewInstanceManager(cfg map[string]*config.DatabaseConfig, log *logger.Logger) *InstanceManager {
	logger := log.Named("database")
	m := &InstanceManager{
		instances: make(map[string]*DB),
		logger:    logger,
	}
	// 必须要有default实例
	if _, ok := cfg["default"]; !ok {
		logger.Fatal("must have a 'default' database instance",
			zap.String("instance", "default"),
		)
	}
	// 初始化新实例
	for name, cfg := range cfg {
		db := New(cfg, logger)
		m.instances[name] = db
	}

	return m
}

// 获取实例（需要读锁）
func (m *InstanceManager) Get(name string) *DB {
	m.mu.RLock()
	defer m.mu.RUnlock()
	db, ok := m.instances[name]
	if !ok {
		m.logger.Panic("[database] instance not found",
			zap.String("instance", name),
		)
	}
	return db
}

// 添加实例（需要写锁）
func (m *InstanceManager) Add(name string, db *DB) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.instances[name] = db
}

func (m *InstanceManager) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, db := range m.instances {
		db.Close()
	}
}

func (m *InstanceManager) GetDefaultDB() (*DB, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	db, ok := m.instances["default"]
	if !ok {
		m.logger.Panic("[database] default database instance not found")
	}
	return db, nil
}
