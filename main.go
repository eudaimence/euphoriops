package main

import (
    "bufio"
    "embed"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

//go:embed resources/bnet.exe
var embeddedFiles embed.FS

type Game struct {
    name    string
    prodID  string
    display string
}

var games = []Game{
    {"Call of Duty Black Ops 4", "viper", "Call of Duty: Black Ops 4"},
    {"Call of Duty Black Ops Cold War", "zeus", "Call of Duty: Black Ops Cold War"},
    {"Call of Duty Vanguard", "fore", "Call of Duty: Vanguard"},
}

func cleanPath(path string) string {
    path = strings.Trim(path, "'\"() ")
    path = strings.ReplaceAll(path, "\\", "/")
    path = strings.ReplaceAll(path, "/", "\\")
    return path
}

func main() {
    tempFile, err := os.CreateTemp("", "bnet-*.exe")
    if err != nil {
        fmt.Printf("Error creating temp file: %v\n", err)
        return
    }
    defer os.Remove(tempFile.Name())

    executableBytes, err := embeddedFiles.ReadFile("resources/bnet.exe")
    if err != nil {
        fmt.Printf("Error reading embedded file: %v\n", err)
        return
    }

    if _, err := tempFile.Write(executableBytes); err != nil {
        fmt.Printf("Error writing to temp file: %v\n", err)
        return
    }
    tempFile.Close()

    for {
        fmt.Println("\nAvailable Games:")
        for i, game := range games {
            fmt.Printf("%d. %s\n", i+1, game.display)
        }
        fmt.Println("0. Exit")

        var choice int
        fmt.Print("\nSelect a game (0-", len(games), "): ")
        _, err := fmt.Scan(&choice)
        if err != nil {
            fmt.Println("Invalid input. Please try again.")
            continue
        }

        if choice == 0 {
            break
        }

        if choice < 1 || choice > len(games) {
            fmt.Println("Invalid choice. Please try again.")
            continue
        }

        bufio.NewReader(os.Stdin).ReadString('\n')

        reader := bufio.NewReader(os.Stdin)
        fmt.Printf("\nEnter the installation directory for %s\n(e.g., D:\\Games\\%s): ", 
            games[choice-1].display, 
            games[choice-1].name)

        installDir, err := reader.ReadString('\n')
        if err != nil {
            fmt.Printf("Error reading input: %v\n", err)
            continue
        }

        installDir = cleanPath(strings.TrimSpace(installDir))

        if installDir == "" {
            fmt.Println("Invalid path: Path cannot be empty")
            continue
        }

        if !strings.Contains(installDir, ":\\") {
            fmt.Println("Invalid path: Must include drive letter (e.g., D:\\Games\\...)")
            continue
        }

        if err := os.MkdirAll(installDir, 0755); err != nil {
            fmt.Printf("Error creating directory: %v\n", err)
            continue
        }

        cmd := exec.Command(tempFile.Name(),
            "--prod", games[choice-1].prodID,
            "--lang", "enUS",
            "--dir", installDir)

        fmt.Printf("\nStarting installation for %s...\n", games[choice-1].display)
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr

        if err := cmd.Run(); err != nil {
            fmt.Printf("Error running installer: %v\n", err)
            continue
        }
    }
}