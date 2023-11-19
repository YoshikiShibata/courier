// Copyright © 2023 Yoshiki Shibata. All rights reserved.

package courier

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type Config struct {
	MakeDir           string
	MakeBuildTarget   string
	ServiceBinaryPath string
	ServiceName       string
	GRPCPort          int
	Envs              []string
	Verbose           bool
	CoverageDir       string
}

func InvokeService(config *Config) (func() error, error) {
	ctx, cancel := context.WithCancel(context.Background())

	if err := os.Chdir(config.MakeDir); err != nil {
		cancel()
		return nil, err
	}

	// Build service
	cmd := exec.CommandContext(ctx, "make", config.MakeBuildTarget)

	if config.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Start(); err != nil {
		cancel()
		return nil, fmt.Errorf("cmd.Start() failed: %w", err)
	}
	if err := cmd.Wait(); err != nil {
		cancel()
		return nil, fmt.Errorf("cmd.Wait() failed: %w", err)
	}

	// InvokeService
	log.Printf("Invoking %s", config.ServiceBinaryPath)
	cmd = exec.CommandContext(ctx, config.ServiceBinaryPath)
	cmd.Env = append(cmd.Env, config.Envs...)
	cmd.Env = append(cmd.Env, "GOCOVERDIR="+config.CoverageDir)

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

	ready := asyncWaitForServiceReady(config)

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

		if err == nil {
			showHowToSeeCoverage(config)
		}
		return err
	}, nil
}

func asyncWaitForServiceReady(config *Config) (ready chan struct{}) {
	ready = make(chan struct{})

	go func() {
		clientOptions := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}
		cc, err := grpc.Dial(fmt.Sprintf("localhost:%d", config.GRPCPort), clientOptions...)
		if err != nil {
			log.Fatalf("grpc.Dial failed: %v", err)
		}
		defer func() { _ = cc.Close() }()

		client := healthpb.NewHealthClient(cc)

		for {
			res, err := client.Check(context.Background(), &healthpb.HealthCheckRequest{
				Service: config.ServiceName,
			})
			if status.Code(err) != codes.OK {
				log.Printf("Check failed: %v", err)
				time.Sleep(time.Second)
				continue
			}

			if res.Status != healthpb.HealthCheckResponse_SERVING {
				log.Printf("Status is %v", res.Status)
				time.Sleep(time.Second)
				continue
			}

			log.Printf("The service is now SERVING")
			close(ready)
			return
		}
	}()

	return ready
}

func terminateService(cmd *exec.Cmd) {
	cmd.Process.Signal(syscall.SIGTERM)
}

func showHowToSeeCoverage(config *Config) {
	log.Printf("To see coverage, run the follwoing command")
	log.Printf("  $ (cd %s; go tool covdata percent -i=%s)", config.MakeDir, config.CoverageDir)
	log.Printf("Or to see which lines have been executed, run the following command")
	log.Printf("  $ (cd %s; go tool covdata textfmt -i=%s -o profile.txt; go tool cover -html=profile.txt)",
		config.MakeDir, config.CoverageDir)
}
