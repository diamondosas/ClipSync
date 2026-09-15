//go:build windows

// ai-generated
package instance

import (
	"context"
	"errors"
	"log"
	"syscall"

	"golang.org/x/sys/windows"
)

var (
	instanceMutex windows.Handle
	wakeupEvent   windows.Handle
)

const (
	mutexName = "Local\\ClipSync_Instance_Mutex"
	eventName = "Local\\ClipSync_Wakeup_Event"
)

func checkInstance() bool {
	mName, err := windows.UTF16PtrFromString(mutexName)
	if err != nil {
		log.Printf("[Instance] Failed to encode mutex name: %v\n", err)
		return false
	}

	hMutex, err := windows.CreateMutex(nil, true, mName)
	// If the mutex already exists, Windows returns ERROR_ALREADY_EXISTS (183)
	if errors.Is(err, windows.ERROR_ALREADY_EXISTS) || errors.Is(err, syscall.Errno(183)) {
		log.Println("[Instance] Another instance is already running. Notifying it...")
		notifyExistingProcess()
		if hMutex != 0 {
			_ = windows.CloseHandle(hMutex)
		}
		return true
	}

	if err != nil && err != syscall.Errno(0) {
		log.Printf("[Instance] Warning: CreateMutex returned unexpected error: %v\n", err)
	}

	instanceMutex = hMutex

	// Create Named Event for waking up the primary instance
	eName, err := windows.UTF16PtrFromString(eventName)
	if err != nil {
		log.Printf("[Instance] Failed to encode event name: %v\n", err)
		return false
	}

	// manualReset = 0 (auto-reset), initialState = 0 (nonsignaled)
	hEvent, err := windows.CreateEvent(nil, 0, 0, eName)
	if err != nil && err != syscall.Errno(0) && !errors.Is(err, windows.ERROR_ALREADY_EXISTS) {
		log.Printf("[Instance] Warning: CreateEvent returned error: %v\n", err)
	}
	wakeupEvent = hEvent

	return false
}

func notifyExistingProcess() {
	eName, err := windows.UTF16PtrFromString(eventName)
	if err != nil {
		log.Printf("[Instance] Failed to encode event name: %v\n", err)
		return
	}

	hEvent, err := windows.OpenEvent(windows.EVENT_MODIFY_STATE, false, eName)
	if err != nil {
		log.Printf("[Instance] Could not open wake-up event: %v\n", err)
		return
	}
	defer windows.CloseHandle(hEvent)

	err = windows.SetEvent(hEvent)
	if err != nil {
		log.Printf("[Instance] Failed to signal wake-up event: %v\n", err)
	}
}

func startListener(ctx context.Context, showGUI func()) {
	if wakeupEvent == 0 {
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// Wait for secondary instance to signal event with 500ms timeout to periodically check ctx.Done
				res, _ := windows.WaitForSingleObject(wakeupEvent, 500)
				if res == windows.WAIT_OBJECT_0 {
					log.Println("[Instance] Received wake-up event, opening GUI...")
					showGUI()
				}
			}
		}
	}()
}

func cleanup() {
	if wakeupEvent != 0 {
		_ = windows.CloseHandle(wakeupEvent)
		wakeupEvent = 0
	}
	if instanceMutex != 0 {
		_ = windows.CloseHandle(instanceMutex)
		instanceMutex = 0
	}
}
