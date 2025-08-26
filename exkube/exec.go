// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package exkube

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/remotecommand"
)

const (
	defaultReadFromByteCmd = "tail -c+%d %s"
)

type ExecResult struct {
	Stdout bytes.Buffer
	Stderr bytes.Buffer
}

type ExecParameters struct {
	Namespace string
	Pod       string
	Container string
	Command   []string
	TTY       bool // fuses stderr into stdout if 'true', needed for Ctrl-C support
}

func (c *Client) execInPodWithWriters(connCtx, killCmdCtx context.Context, p ExecParameters, stdout, stderr io.Writer) error {
	req := c.Clientset.CoreV1().RESTClient().Post().Resource("pods").Name(p.Pod).Namespace(p.Namespace).SubResource("exec")

	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		return fmt.Errorf("error adding to scheme: %w", err)
	}

	parameterCodec := runtime.NewParameterCodec(scheme)

	execOpts := &corev1.PodExecOptions{
		Command:   p.Command,
		Container: p.Container,
		Stdin:     p.TTY,
		Stdout:    true,
		Stderr:    true,
		TTY:       p.TTY,
	}
	req.VersionedParams(execOpts, parameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(c.Config, "POST", req.URL())
	if err != nil {
		return fmt.Errorf("error while creating executor: %w", err)
	}

	var stdin io.ReadCloser
	if p.TTY {
		// CtrlCReader sends Ctrl-C/D sequence if context is cancelled
		stdin = NewCtrlCReader(killCmdCtx)
		// Graceful close of stdin once we are done, no Ctrl-C is sent
		// if execution finishes before the context expires.
		defer stdin.Close()
	}

	return exec.StreamWithContext(connCtx, remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
		Tty:    p.TTY,
	})
}

func (c *Client) execInPod(ctx context.Context, p ExecParameters) (*ExecResult, error) {
	result := &ExecResult{}
	err := c.execInPodWithWriters(ctx, nil, p, &result.Stdout, &result.Stderr)

	return result, err
}

// CopyFromPod is to copy srcFile in a given pod to local destFile with defaultMaxTries.
func (c *Client) CopyFromPod(ctx context.Context, namespace, pod, container, fromFile, destFile string, retryLimit int) error {
	pipe := newPipe(&CopyOptions{
		MaxTries: retryLimit,
		ReadFunc: readFromPod(ctx, c, namespace, pod, container, fromFile),
	})

	outFile, err := os.OpenFile(destFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if _, err = io.Copy(outFile, pipe); err != nil {
		return err
	}
	return nil
}

func readFromPod(ctx context.Context, client *Client, namespace, pod, container, srcFile string) ReadFunc {
	return func(offset uint64, writer io.Writer) error {
		command := []string{"sh", "-c", fmt.Sprintf(defaultReadFromByteCmd, offset, srcFile)}
		return client.execInPodWithWriters(ctx, nil, ExecParameters{
			Namespace: namespace,
			Pod:       pod,
			Container: container,
			Command:   command,
		}, writer, writer)
	}
}
