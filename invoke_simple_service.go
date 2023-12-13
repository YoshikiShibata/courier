// Copyright © 2023 Yoshiki Shibata. All rights reserved.

package courier

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// SimpleConfig provides a simple configuration without any steps
// to build service binary.
type SimpleConfig struct {
	// Directory to change.
	Chdir string
	// Command to invoke service.
	ServiceCmd string
	// Args to pass when invoking the service.
	ServiceArgs []string
	// Environments variables
	Envs []string
	// Redirect the output from the service to standard out/err.
	Verbose bool
	// Function to probe if the service is ready.
	// This function must retry probing until the service is ready.
	ReadyProbe func()
}

func InvokeSimpleService(config *SimpleConfig) (func() error, error) {
	ctx, cancel := context.WithCancel(context.Background())

	if err := os.Chdir(config.Chdir); err != nil {
		cancel()
		return nil, err
	}

	// Build service
	log.Printf("Running %s %s", config.ServiceCmd, strings.Join(config.ServiceArgs, " "))
	cmd := exec.CommandContext(ctx, config.ServiceCmd, config.ServiceArgs...)
	cmd.Env = append(cmd.Env, config.Envs...)

	if config.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Start(); err != nil {
		cancel()
		return nil, fmt.Errorf("cmd.Start() failed: %w", err)
	}

	exitedErr := make(chan error)
	go func() {
		exitedErr <- cmd.Wait()
		close(exitedErr)
	}()

	ready := asyncWaitForSimpleServiceReady(config)

	select {
	// Wait for the service being ready.
	case <-ready:

	case err := <-exitedErr:
		cancel()
		return nil, err
	}

	return func() error {
		// The service may have already terminated.
		select {
		case err := <-exitedErr:
			return fmt.Errorf("The service has already terminated: %w", err)
		default:
		}

		// The service is still running
		terminateService(cmd)
		err := <-exitedErr

		cancel()

		// The invoked service may not handle TERM signal gracefully and
		// just terminated with the sinal. So it should be treated as ok.
		if err != nil && strings.Contains(err.Error(), "signal: terminated") {
			return nil
		}
		return err
	}, nil
}

func asyncWaitForSimpleServiceReady(config *SimpleConfig) (ready chan struct{}) {
	if config.ReadyProbe == nil {
		log.Fatal("ReadyProbe is not specified")
	}

	ready = make(chan struct{})

	go func() {
		config.ReadyProbe()
		close(ready)
	}()

	return ready
}
