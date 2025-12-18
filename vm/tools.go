package vm

import (
	"fmt"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

// Структуры аргументов и результатов каждого действия

// CreateVMArgs - аргументы для создания ВМ
type CreateVMArgs struct {
	Name string `json:"name"`
	Memory uint64 `json:"memory"` // в МБ
	VCPUs uint `json:"vcpus"`
	DiskPath string `json:"disk_path,omitempty"`
	DiskSize uint64 `json:"disk_size,omitempty"` // в ГБ
	ISOImage string `json:"iso_image,omitempty"`
	Network string `json:"network,omitempty"`
}

// CreateVMResult - результат создания ВМ
type CreateVMResult struct {
	Message string `json:"message"`
	VMName string `json:"vm_name"`
}

// StartVMArgs - аргументы для запуска ВМ
type StartVMArgs struct {
	Name string `json:"name"`
}

// StartVMResult - результат запуска ВМ
type StartVMResult struct {
	Message string `json:"message"`
}

// StopVMArgs - аргументы для остановки ВМ
type StopVMArgs struct {
	Name string `json:"name"`
}

// StopVMResult - результат остановки ВМ
type StopVMResult struct {
	Message string `json:"message"`
}

// ListVMsResult - результат списка ВМ
type ListVMsResult struct {
	VMs []string `json:"vms"`
}

// DeleteVMArgs - аргументы для удаления ВМ
type DeleteVMArgs struct {
	Name string `json:"name"`
}

// DeleteVMResult - результат удаления ВМ
type DeleteVMResult struct {
	Message string `json:"message"`
}

// NewVMTools создает набор инструментов для управления ВМ
func NewVMTools(manager VMManagerInterface) ([]tool.Tool, error) {
	var tools []tool.Tool

	// Инструмент для создания ВМ
	createVMTool, err := functiontool.New(
		functiontool.Config{
			Name: "create_vm",
			Description: "Creates a new virtual machine with the specified configuration.",
		},
		func(ctx tool.Context, args CreateVMArgs) (CreateVMResult, error) {
			config := VMConfig{
				Name: args.Name,
				Memory:   args.Memory,
				VCPUs:    args.VCPUs,
				DiskPath: args.DiskPath,
				DiskSize: args.DiskSize,
				ISOImage: args.ISOImage,
				Network:  args.Network,
			}

			if err := manager.CreateVM(config); err != nil {
				return CreateVMResult{}, fmt.Errorf("failed to create a VM: %w", err)
			}

			return CreateVMResult{
				Message: fmt.Sprintf("VM '%s' has created successfully!", args.Name),
				VMName: args.Name,
			}, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create VM tool: %w", err)
	}
	tools = append(tools, createVMTool)

	// Инструмент для запуска ВМ
	startVMTool, err := functiontool.New(
		functiontool.Config{
			Name: "start_vm",
			Description: "Starts a specific virtual machine.",
		},
		func(ctx tool.Context, args StartVMArgs) (StartVMResult, error) {
			if err := manager.StartVM(args.Name); err != nil {
				return StartVMResult{}, fmt.Errorf("failed to start '%s' VM; err: %w", args.Name, err)
			}
			return StartVMResult{
				Message: "Virtual machine '%s' has started successfully!",
			}, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create start_vm tool: %w", err)
	}
	tools = append(tools, startVMTool)

	// Инструмент для остановки ВМ
	stopVMTool, err := functiontool.New(
		functiontool.Config{
			Name:        "stop_vm",
			Description: "Stops a virtual machine by name",
		},
		func(ctx tool.Context, args StopVMArgs) (StopVMResult, error) {
			if err := manager.StopVM(args.Name); err != nil {
				return StopVMResult{}, fmt.Errorf("failed to stop VM: %w", err)
			}
			return StopVMResult{
				Message: fmt.Sprintf("Virtual machine '%s' stopped successfully", args.Name),
			}, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create stop_vm tool: %w", err)
	}
	tools = append(tools, stopVMTool)

	// Инструмент для списка ВМ
	listVMsTool, err := functiontool.New(
		functiontool.Config{
			Name:        "list_vms",
			Description: "Lists all available virtual machines",
		},
		func(ctx tool.Context, args struct{}) (ListVMsResult, error) {
			vms, err := manager.ListVMs()
			if err != nil {
				return ListVMsResult{}, fmt.Errorf("failed to list VMs: %w", err)
			}
			return ListVMsResult{
				VMs: vms,
			}, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create list_vms tool: %w", err)
	}
	tools = append(tools, listVMsTool)

	// Инструмент для удаления ВМ
	deleteVMTool, err := functiontool.New(
		functiontool.Config{
			Name:        "delete_vm",
			Description: "Deletes a virtual machine by name",
		},
		func(ctx tool.Context, args DeleteVMArgs) (DeleteVMResult, error) {
			if err := manager.DeleteVM(args.Name); err != nil {
				return DeleteVMResult{}, fmt.Errorf("failed to delete VM: %w", err)
			}
			return DeleteVMResult{
				Message: fmt.Sprintf("Virtual machine '%s' deleted successfully", args.Name),
			}, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create delete_vm tool: %w", err)
	}
	tools = append(tools, deleteVMTool)

	return tools, nil
}