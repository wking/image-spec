// Copyright 2016 The Linux Foundation
//
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

package image

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/opencontainers/image-spec/image/cas"
	"github.com/opencontainers/image-spec/schema"
	imageSpecs "github.com/opencontainers/image-spec/specs-go"
	runtimeSpecs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type cfg struct {
	User         string
	Memory       int64
	MemorySwap   int64
	CPUShares    int64 `json:"CpuShares"`
	ExposedPorts map[string]struct{}
	Env          []string
	Entrypoint   []string
	Cmd          []string
	Volumes      map[string]struct{}
	WorkingDir   string
}

type config struct {
	Architecture string `json:"architecture"`
	OS           string `json:"os"`
	Config       cfg    `json:"config"`
}

func findConfig(ctx context.Context, engine cas.Engine, descriptor *imageSpecs.Descriptor) (*config, error) {
	reader, err := engine.Get(ctx, descriptor.Digest)
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errors.Wrapf(err, "%s: error reading manifest", descriptor.Digest)
	}

	if err := schema.MediaTypeImageSerializationConfig.Validate(bytes.NewReader(buf)); err != nil {
		return nil, errors.Wrapf(err, "%s: config validation failed", descriptor.Digest)
	}

	var c config
	if err := json.Unmarshal(buf, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *config) runtimeSpec(rootfs string) (*runtimeSpecs.Spec, error) {
	if c.OS != "linux" {
		return nil, fmt.Errorf("%s: unsupported OS", c.OS)
	}

	var s runtimeSpecs.Spec
	s.Version = "0.5.0"
	s.Root.Path = rootfs
	s.Process.Cwd = c.Config.WorkingDir
	s.Process.Env = append([]string(nil), c.Config.Env...)
	s.Process.Args = append([]string(nil), c.Config.Entrypoint...)
	s.Process.Args = append(s.Process.Args, c.Config.Cmd...)

	if uid, err := strconv.Atoi(c.Config.User); err == nil {
		s.Process.User.UID = uint32(uid)
	} else if ug := strings.Split(c.Config.User, ":"); len(ug) == 2 {
		uid, err := strconv.Atoi(ug[0])
		if err != nil {
			return nil, errors.New("config.User: unsupported uid format")
		}

		gid, err := strconv.Atoi(ug[1])
		if err != nil {
			return nil, errors.New("config.User: unsupported gid format")
		}

		s.Process.User.UID = uint32(uid)
		s.Process.User.GID = uint32(gid)
	} else {
		return nil, errors.New("config.User: unsupported format")
	}

	s.Platform.OS = c.OS
	s.Platform.Arch = c.Architecture

	mem := uint64(c.Config.Memory)
	swap := uint64(c.Config.MemorySwap)
	shares := uint64(c.Config.CPUShares)

	s.Linux.Resources = &runtimeSpecs.Resources{
		CPU: &runtimeSpecs.CPU{
			Shares: &shares,
		},

		Memory: &runtimeSpecs.Memory{
			Limit:       &mem,
			Reservation: &mem,
			Swap:        &swap,
		},
	}

	for vol := range c.Config.Volumes {
		s.Mounts = append(
			s.Mounts,
			runtimeSpecs.Mount{
				Destination: vol,
				Type:        "bind",
				Options:     []string{"rbind"},
			},
		)
	}

	return &s, nil
}
