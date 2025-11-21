package main

import (
	"fmt"
	"log/slog"
	"os"
)

func main() {
	file, err := os.Open("./service.pid")

	if err != nil {
		fmt.Println(fmt.Errorf("error on opening -- %w", err))
		return
	}

	defer file.Close()

	fmt.Printf("file name=%v\n", file.Name())

	var pid int
	if _, err = fmt.Fscanf(file, "%d", &pid); err != nil {
		fmt.Println(fmt.Errorf("error on reading -- %w", err))
		return
	}

	slog.Info("Killing", "pid", pid)
}
