/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

var verbose = flag.Bool("verbose", false, "Print verbose output")
var delay = flag.Duration("delay", 0, "A duration to delay for")
var wait = flag.Duration("wait", 0, "A duration to wait between restarting the sub process on exit")

func main() {
	flag.Parse()

	subCommand := flag.Args()
	if len(subCommand) == 0 {
		fatalf("please provide a sub command to run")
	}

	pwd, err := os.Getwd()
	if err != nil {
		fatal(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	if err := runLoop(quit, pwd, subCommand...); err != nil {
		fatal(err)
	}

	os.Exit(0)
}

// resolveBin splits a slice of command arguments into the binary (i.e. the
// first argument) and the arguments. It also will resolve the first argument
// to a binary on the PATH. E.g. `ls` gets replaced with `/bin/ls`.
func resolveBin(subCommand ...string) (string, []string, error) {
	bin := subCommand[0]
	binPath, err := exec.LookPath(bin)
	if err != nil {
		return "", nil, err
	}

	verbosef("%q resolved to %q", bin, binPath)
	return binPath, subCommand[1:], nil
}

func createSub(pwd, binPath string, args ...string) *exec.Cmd {
	sub := exec.Command(binPath, args...)
	sub.Env = os.Environ()
	sub.Dir = pwd
	sub.Stdout = os.Stdout
	sub.Stderr = os.Stderr
	return sub
}

func runLoop(quit chan os.Signal, pwd string, subCommand ...string) error {
	if delay != nil && *delay > 0 {
		verbosef("delaying %v before starting", *delay)
		alarm := time.After(*delay)
		select {
		case <-alarm:
			break
		case s := <-quit:
			verbosef("received SIGINT (%s) during delay, exiting", s)
			return nil
		}
	}

	var sub *exec.Cmd
	var err error
	var didQuit bool
	var abort chan struct{}
	var aborted chan struct{}

	binPath, args, err := resolveBin(subCommand...)
	if err != nil {
		return err
	}

	for {
		abort = make(chan struct{})
		aborted = make(chan struct{})

		sub = createSub(pwd, binPath, args...)
		if err := sub.Start(); err != nil {
			return err
		}

		// kick off monitor
		go func() {
			select {
			case s := <-quit:
				verbosef("received SIGINT (%s) while sub process is running, killing sub process", s)
				didQuit = true
				_ = sub.Process.Kill()
				return
			case <-abort:
				close(aborted)
				return
			}
		}()

		// wait for the sub process to exit
		if err := sub.Wait(); err != nil {
			verbosef("sub process exit: %v", err)
		}

		if didQuit {
			return nil
		}

		// abort the monitor
		close(abort)
		// wait for monitor to return
		<-aborted

		if wait != nil && *wait > 0 {
			verbosef("waiting %v before restart", *wait)
			alarm := time.After(*wait)
			select {
			case <-alarm:
				continue
			case s := <-quit:
				verbosef("received SIGINT (%s) during wait, exiting", s)
				return nil
			}
		}
	}
}

func verbosef(format string, args ...interface{}) {
	if verbose != nil && *verbose {
		fmt.Fprintf(os.Stdout, "recover: "+format+"\n", args...)
	}
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "recover: "+format+"\n", args...)
	os.Exit(1)
}

func fatal(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "recover: %v\n", err)
		os.Exit(1)
	}
}
