// This file is kept for backward compatibility
// The main application entry point is now at cmd/server/main.go
package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("Starting Simple Auth with Roles server...")
	fmt.Println("Main application moved to cmd/server/main.go")

	// Set environment variable to run migrations automatically
	os.Setenv("MIGRATE", "true")

	// Execute the main server
	cmd := exec.Command("go", "run", "./cmd/server")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ() // Pass all environment variables

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running server: %v\n", err)
		os.Exit(1)
	}
}
