package vm

import (
	"fmt"
	"log"
	"sync"
)

// VMManagerInterface определяет интерфейс для управления виртуальными машинами
type VMManagerInterface interface {
	CreateVM(config VMConfig) error
	ListVMs() ([]string, error)
	StartVM(name string) error
	StopVM(name string) error
	DeleteVM(name string) error
	Close() error
}

type VMConfig struct {
	Name string
	Memory uint64
	VCPUs uint
	DiskPath string
	DiskSize uint64
	ISOImage string
	Network string
}

// VMState представляет состояние виртуальной машины
type VMState string

const (
	VMStateStopped VMState = "stopped"
	VMStateRunning VMState = "running"
	VMStatePaused  VMState = "paused"
)

// MockVM представляет виртуальную машину в mock-режиме
type MockVM struct {
	Config VMConfig
	State  VMState
}

// MockVMManager - mock-реализация менеджера виртуальных машин
// Хранит все данные в памяти, не создает реальные виртуальные машины
type MockVMManager struct {
	vms  map[string]*MockVM
	mu   sync.RWMutex
	next int // для генерации уникальных ID
}

// NewMockVMManager создает новый mock-менеджер виртуальных машин
func NewMockVMManager() *MockVMManager {
	return &MockVMManager{
		vms:  make(map[string]*MockVM),
		next: 1,
	}
}

// Close закрывает mock-менеджер (ничего не делает, но реализует интерфейс)
func (m *MockVMManager) Close() error {
	log.Println("[MOCK] Closing VM manager (no-op in mock mode)")
	return nil
}

// CreateVM создает новую виртуальную машину в памяти
func (m *MockVMManager) CreateVM(config VMConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Проверяем, не существует ли уже ВМ с таким именем
	if _, exists := m.vms[config.Name]; exists {
		return fmt.Errorf("virtual machine with name '%s' already exists", config.Name)
	}

	// Валидация конфигурации
	if config.Name == "" {
		return fmt.Errorf("VM name cannot be empty")
	}
	if config.Memory == 0 {
		return fmt.Errorf("VM memory cannot be zero")
	}
	if config.VCPUs == 0 {
		return fmt.Errorf("VM VCPUs cannot be zero")
	}

	// Создаем mock-виртуальную машину
	mockVM := &MockVM{
		Config: config,
		State:  VMStateStopped,
	}

	m.vms[config.Name] = mockVM

	log.Printf("[MOCK] Virtual machine '%s' created successfully (Memory: %d MB, VCPUs: %d, Disk: %s)",
		config.Name, config.Memory, config.VCPUs, config.DiskPath)

	// Автоматически запускаем ВМ (в mock-режиме это просто изменение состояния)
	mockVM.State = VMStateRunning
	log.Printf("[MOCK] Virtual machine '%s' started successfully", config.Name)

	return nil
}

// ListVMs возвращает список всех виртуальных машин
func (m *MockVMManager) ListVMs() ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	vmNames := make([]string, 0, len(m.vms))
	for name := range m.vms {
		vmNames = append(vmNames, name)
	}

	log.Printf("[MOCK] Listed %d virtual machine(s)", len(vmNames))
	return vmNames, nil
}

// StartVM запускает виртуальную машину по имени
func (m *MockVMManager) StartVM(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	vm, exists := m.vms[name]
	if !exists {
		return fmt.Errorf("virtual machine '%s' not found", name)
	}

	if vm.State == VMStateRunning {
		log.Printf("[MOCK] Virtual machine '%s' is already running", name)
		return nil
	}

	vm.State = VMStateRunning
	log.Printf("[MOCK] Virtual machine '%s' started", name)
	return nil
}

// StopVM останавливает виртуальную машину
func (m *MockVMManager) StopVM(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	vm, exists := m.vms[name]
	if !exists {
		return fmt.Errorf("virtual machine '%s' not found", name)
	}

	if vm.State == VMStateStopped {
		log.Printf("[MOCK] Virtual machine '%s' is already stopped", name)
		return nil
	}

	vm.State = VMStateStopped
	log.Printf("[MOCK] Virtual machine '%s' stopped", name)
	return nil
}

// DeleteVM удаляет виртуальную машину
func (m *MockVMManager) DeleteVM(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	vm, exists := m.vms[name]
	if !exists {
		return fmt.Errorf("virtual machine '%s' not found", name)
	}

	// Останавливаем, если запущена
	if vm.State == VMStateRunning {
		vm.State = VMStateStopped
		log.Printf("[MOCK] Stopped virtual machine '%s' before deletion", name)
	}

	// Удаляем из хранилища
	delete(m.vms, name)
	log.Printf("[MOCK] Virtual machine '%s' deleted", name)
	return nil
}

// GetVMInfo возвращает информацию о виртуальной машине (дополнительный метод для mock)
func (m *MockVMManager) GetVMInfo(name string) (*MockVM, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	vm, exists := m.vms[name]
	if !exists {
		return nil, fmt.Errorf("virtual machine '%s' not found", name)
	}

	return vm, nil
}

// GetVMState возвращает состояние виртуальной машины (дополнительный метод для mock)
func (m *MockVMManager) GetVMState(name string) (VMState, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	vm, exists := m.vms[name]
	if !exists {
		return "", fmt.Errorf("virtual machine '%s' not found", name)
	}

	return vm.State, nil
}

