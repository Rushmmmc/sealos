// Copyright 2019 Intel Corporation. All Rights Reserved.
// Copyright 2022 Sealer Corporation
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shim

import (
	"fmt"
	"os"
	"sync"

	"github.com/labring/image-cri-shim/pkg/server"
	"github.com/labring/image-cri-shim/pkg/shim/client"

	"github.com/labring/sealos/pkg/utils/logger"
)

const (
	// DisableService is used to mark a socket/service to not be connected.
	DisableService = client.DontConnect
)

// Options contains the configurable options of our CRI shim.
type Options struct {
	// ShimSocket is the socket path for the CRI shim services.
	ShimSocket string
	// ImageSocket is the socket path for the (real) CRI image services.
	ImageSocket string
}

// Shim is the interface we expose for controlling our CRI shim.
type Shim interface {
	// Setup prepares the shim to start processing CRI requests.
	Setup() error
	// Start starts the shim.
	Start() error
	// Stop stops the shim.
	Stop()
	// Client returns the shim client interface.
	Client() client.Client
	// Server returns the shim server interface.
	Server() server.Server
}

// shim is the implementation of Shim.
type shim struct {
	sync.Mutex               // hmm... do *we* need to be lockable, or the upper layer(s) ?
	options    Options       // shim options
	client     client.Client // shim CRI client
	server     server.Server // shim CRI server
}

// NewShim creates a new shim instance.
func NewShim(options Options) (Shim, error) {
	var err error

	r := &shim{
		options: options,
	}

	cltopts := client.Options{
		ImageSocket: r.options.ImageSocket,
		DialNotify:  r.dialNotify,
	}
	if r.client, err = client.NewClient(cltopts); err != nil {
		return nil, shimError("failed to create shim client: %v", err)
	}

	srvopts := server.Options{
		Socket: r.options.ShimSocket,
		User:   -1,
		Group:  -1,
		Mode:   0660,
	}
	if r.server, err = server.NewServer(srvopts); err != nil {
		return nil, shimError("failed to create shim server: %v", err)
	}

	return r, nil
}

// Setup prepares the shim to start processing requests.
func (r *shim) Setup() error {
	if err := r.client.Connect(client.ConnectOptions{Wait: true}); err != nil {
		return shimError("client connection failed: %v", err)
	}

	if r.options.ImageSocket != DisableService {
		if err := r.server.RegisterImageService(r); err != nil {
			return shimError("failed to register image service: %v", err)
		}
	}

	return nil
}

// Start starts the shim request processing goroutine.
func (r *shim) Start() error {
	if err := r.server.Start(); err != nil {
		return shimError("failed to start shim: %v", err)
	}

	return nil
}

// Stop stops the shim.
func (r *shim) Stop() {
	r.client.Close()
	r.server.Stop()
}

// Client returns the shim Client interface.
func (r *shim) Client() client.Client {
	return r.client
}

// Server returns the shim Server interface.
func (r *shim) Server() server.Server {
	return r.server
}

func (r *shim) dialNotify(socket string, uid int, gid int, mode os.FileMode, err error) {
	if err != nil {
		logger.Error("failed to determine permissions/ownership of client socket %q: %v",
			socket, err)
		return
	}

	if err = r.server.Chown(uid, gid); err != nil {
		logger.Error("server socket ownership change request failed: %v", err)
	} else {
		if err := r.server.Chmod(mode); err != nil {
			logger.Error("server socket permissions change request failed: %v", err)
		}
	}
}

// shimError creates a formatted shim-specific error.
func shimError(format string, args ...interface{}) error {
	return fmt.Errorf("cri/shim: "+format, args...)
}
