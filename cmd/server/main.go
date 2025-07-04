package main

import (
    "log"
    "github.com/eager/cyberpunkmp/internal/server"
)

func main() {
    if err := server.Start(); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
