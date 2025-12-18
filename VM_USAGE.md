# Управление виртуальными машинами (Mock-режим)

Этот пакет предоставляет mock-реализацию для управления виртуальными машинами. Все операции выполняются в памяти, без создания реальных виртуальных машин и без нагрузки на систему.

## Быстрый старт

```go
// Создаем mock-менеджер (не требует libvirt или других зависимостей)
manager := NewMockVMManager()
defer manager.Close()

config := VMConfig{
    Name:     "test-vm",
    Memory:   512,  // MB
    VCPUs:    1,
    DiskPath: "/path/to/disk.qcow2", // Путь не проверяется в mock-режиме
    DiskSize: 5,    // GB
    Network:  "default",
}

err := manager.CreateVM(config)
```

## Примеры использования

### Базовый пример

```go
package main

import (
    "fmt"
    "log"
)

func main() {
    manager := NewMockVMManager()
    defer manager.Close()

    // Создание виртуальной машины
    config := VMConfig{
        Name:     "my-vm",
        Memory:   1024,  // 1GB RAM
        VCPUs:    2,
        DiskPath: "/mock/path/disk.qcow2",
        DiskSize: 10,   // 10GB
        Network:  "default",
    }

    if err := manager.CreateVM(config); err != nil {
        log.Fatal(err)
    }

    // Список всех ВМ
    vms, err := manager.ListVMs()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Available VMs:")
    for _, vm := range vms {
        vmInfo, _ := manager.GetVMInfo(vm)
        state, _ := manager.GetVMState(vm)
        fmt.Printf("  - %s (State: %s, Memory: %d MB, VCPUs: %d)\n",
            vm, state, vmInfo.Config.Memory, vmInfo.Config.VCPUs)
    }

    // Управление ВМ
    manager.StartVM("my-vm")
    manager.StopVM("my-vm")
    manager.DeleteVM("my-vm")
}
```

## API

### VMManagerInterface

```go
type VMManagerInterface interface {
    CreateVM(config VMConfig) error
    ListVMs() ([]string, error)
    StartVM(name string) error
    StopVM(name string) error
    DeleteVM(name string) error
    Close() error
}
```

### VMConfig

```go
type VMConfig struct {
    Name       string
    Memory     uint64 // в МБ
    VCPUs      uint   // количество виртуальных CPU
    DiskPath   string // путь к диску
    DiskSize   uint64 // размер диска в ГБ
    ISOImage   string // путь к ISO образу (опционально)
    Network    string // тип сети
}
```

## Дополнительные методы

Mock-реализация предоставляет дополнительные методы для получения информации:

```go
manager := NewMockVMManager()

// Получить полную информацию о ВМ
vmInfo, err := manager.GetVMInfo("test-vm")
if err == nil {
    fmt.Printf("VM: %s\n", vmInfo.Config.Name)
    fmt.Printf("Memory: %d MB\n", vmInfo.Config.Memory)
    fmt.Printf("VCPUs: %d\n", vmInfo.Config.VCPUs)
    fmt.Printf("State: %s\n", vmInfo.State)
}

// Получить только состояние ВМ
state, err := manager.GetVMState("test-vm")
// Возможные состояния: VMStateStopped, VMStateRunning, VMStatePaused
```

## Состояния виртуальных машин

Виртуальные машины могут находиться в следующих состояниях:

- `VMStateStopped` - виртуальная машина остановлена
- `VMStateRunning` - виртуальная машина запущена
- `VMStatePaused` - виртуальная машина приостановлена (зарезервировано для будущего использования)

## Примечания

- Все виртуальные машины хранятся в памяти и исчезают при завершении программы
- Путь к диску (`DiskPath`) не проверяется - можно использовать любой путь
- Создание ВМ автоматически запускает её (устанавливает состояние `VMStateRunning`)
- Все операции потокобезопасны благодаря использованию `sync.RWMutex`

