// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gardener/replica-reloader/cmd/replica-reloader/cmd"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 2)
	signal.Notify(c, []os.Signal{os.Interrupt, syscall.SIGTERM}...)

	go func() {
		<-c
		cancel()
		<-c
		os.Exit(1)
	}()

	reloader := cmd.NewReplicaReloader(ctx)
	if err := reloader.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
