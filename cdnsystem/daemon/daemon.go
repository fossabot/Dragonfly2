/*
 *     Copyright 2020 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package daemon

import (
	_ "d7y.io/dragonfly/v2/cdnsystem/storedriver/local"
	"github.com/pkg/errors"
)

import (
	"d7y.io/dragonfly/v2/cdnsystem/config"
	"d7y.io/dragonfly/v2/cdnsystem/plugins"
	"d7y.io/dragonfly/v2/cdnsystem/server"
	logger "d7y.io/dragonfly/v2/pkg/dflog"
)

// Daemon is a struct to identify main instance of cdn.
type Daemon struct {
	config *config.Config
	server *server.Server
}

// New creates a new Daemon.
func New(cfg *config.Config) (*Daemon, error) {
	if err := plugins.Initialize(cfg); err != nil {
		return nil, err
	}
	s, err := server.New(cfg)
	if err != nil {
		return nil, err
	}
	return &Daemon{
		config: cfg,
		server: s,
	}, nil
}

// Serve runs the daemon.
func (d *Daemon) Serve() error {
	if err := d.server.Start(); err != nil {
		return errors.Wrapf(err,"failed to start cdn system")
	}
	logger.Info("successfully start cdn system")
	return nil
}
