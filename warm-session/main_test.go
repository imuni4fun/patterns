package main

import (
	"testing"
	"time"
)

func testMisc(t *testing.T) {
	initSynchronizer()

	go func() {
		pullImage("nginx:latest", "Always", time.Duration(2*time.Minute))
	}()
	go func() {
		pullImage("nginx:A", "Always", time.Duration(2*time.Minute))
	}()
	go func() {
		pullImage("nginx:B", "Always", time.Duration(2*time.Minute))
	}()
	go func() {
		pullImage("nginx:C", "Always", time.Duration(2*time.Minute))
	}()

	killAfter(10)
	// assert.Equal(t, 1, 2, "oops")
}

func TestDuplicates(t *testing.T) {
	initSynchronizer()

	go func() {
		pullImage("nginx:latest", "Always", time.Duration(2*time.Minute))
	}()
	go func() {
		pullImage("nginx:A", "Always", time.Duration(2*time.Minute))
	}()
	go func() {
		pullImage("nginx:B", "Always", time.Duration(2*time.Minute))
	}()
	go func() {
		pullImage("nginx:B", "Always", time.Duration(2*time.Minute))
	}()
	go func() {
		pullImage("nginx:B", "Always", time.Duration(2*time.Minute))
	}()
	go func() {
		pullImage("nginx:B", "Always", time.Duration(2*time.Minute))
	}()
	go func() {
		pullImage("nginx:C", "Always", time.Duration(2*time.Minute))
	}()

	killAfter(10)
	// assert.Equal(t, 1, 2, "oops")
}
