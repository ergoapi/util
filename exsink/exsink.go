// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package exsink

type SinkFactory struct {
}

type EventSink interface {
	SendEvent(any) error
}

func NewSinkFactory() *SinkFactory {
	return &SinkFactory{}
}
