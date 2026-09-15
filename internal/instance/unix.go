//go:build !windows

// ai-generated
package instance

import (
	"clipsync/gui"
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/allan-simon/go-singleinstance"
)

var (
	lockFile *os.File
	lockPath = filepath.Join(os.TempDir(), "clipsync.lock")
)

func checkInstance() bool {
	f, err := singleinstance.CreateLockFile(lockPath)
	if err != nil {
		log.Println("[Instance] Another instance is already running. Notifying it...")
		notifyExistingProcess()
		return true
	}

	lockFile = f
	return false
}

func notifyExistingProcess() {
	pid, err := singleinstance.GetLockFilePid(lockPath)
	if err != nil {
		log.Printf("[Instance] Could not read PID from lock file: %v\n", err)
		return
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		log.Printf("[Instance] Could not find process with PID %d: %v\n", pid, err)
		return
	}

	err = process.Signal(syscall.SIGUSR1)
	if err != nil {
		log.Printf("[Instance] Failed to send SIGUSR1 to PID %d: %v\n", pid, err)
	}
}

func startListener(ctx context.Context, showGUI func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGUSR1)

	go func() {
		defer signal.Stop(sigChan)
		for {
			select {
			case <-sigChan:
				log.Println("[Instance] Received wake-up signal (SIGUSR1), opening GUI...")
				if gui.IsWindowOpen == false{
					showGUI()
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func cleanup() {
	if lockFile != nil {
		_ = lockFile.Close()
		_ = os.Remove(lockPath)
		lockFile = nil
	}
}