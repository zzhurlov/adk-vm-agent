package main

import (
	"context"
	"log"
	"os"
	"test/vm"

	"github.com/joho/godotenv"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"
)

func main() {
    // Загружаем переменные из .env файла
    if err := godotenv.Load(".env"); err != nil {
        log.Println("Warning: .env file not found, using environment variables")
    }

    ctx := context.Background()

    model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
        APIKey: os.Getenv("GOOGLE_API_KEY"),
    })
    if err != nil {
        log.Fatalf("Failed to create model: %v", err)
    }

    VMTools := getVMTools()

    VMAgent, err := llmagent.New(llmagent.Config{
        Name:        "vm_agent",
        Model:       model,
        Description: "Manage some virtual machines using common interface",
        Instruction: "You are a manager of virtual machines, you can creating, starting, stopping, deleting virtual machines, get some information about them.",
        Tools: VMTools,
    })
    if err != nil {
        log.Fatalf("Failed to create agent: %v", err)
    }

    config := &launcher.Config{
        AgentLoader: agent.NewSingleLoader(VMAgent),
    }

    l := full.NewLauncher()
    if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
        log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
    }
}


func getVMTools() []tool.Tool {
    manager := vm.NewMockVMManager()
    VMTools, err := vm.NewVMTools(manager)
    if err != nil {
        log.Fatalf("Failed to create VM tools: %w", err)
    }
    
    return VMTools
}