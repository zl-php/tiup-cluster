// Copyright 2020 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package scripts

import (
	"bytes"
	"errors"
	"io/ioutil"
	"path"
	"text/template"

	"github.com/pingcap-incubator/tiup-cluster/pkg/embed"
)

// DMMasterScript represent the data to generate TiDB config
type DMMasterScript struct {
	Name      string
	Scheme    string
	IP        string
	Port      int
	PeerPort  int
	DeployDir string
	DataDir   string
	LogDir    string
	NumaNode  string
	Endpoints []*DMMasterScript
}

// NewDMMasterScript returns a DMMasterScript with given arguments
func NewDMMasterScript(name, ip, deployDir, dataDir, logDir string) *DMMasterScript {
	return &DMMasterScript{
		Name:      name,
		Scheme:    "http",
		IP:        ip,
		Port:      8261,
		PeerPort:  8291,
		DeployDir: deployDir,
		DataDir:   dataDir,
		LogDir:    logDir,
	}
}

// WithScheme set Scheme field of PDScript
func (c *DMMasterScript) WithScheme(scheme string) *DMMasterScript {
	c.Scheme = scheme
	return c
}

// WithPort set Port field of DMMasterScript
func (c *DMMasterScript) WithPort(port int) *DMMasterScript {
	c.Port = port
	return c
}

// WithNumaNode set NumaNode field of DMMasterScript
func (c *DMMasterScript) WithNumaNode(numa string) *DMMasterScript {
	c.NumaNode = numa
	return c
}

// WithPeerPort set PeerPort field of DMMasterScript
func (c *DMMasterScript) WithPeerPort(port int) *DMMasterScript {
	c.PeerPort = port
	return c
}

// AppendEndpoints add new DMMasterScript to Endpoints field
func (c *DMMasterScript) AppendEndpoints(ends ...*DMMasterScript) *DMMasterScript {
	c.Endpoints = append(c.Endpoints, ends...)
	return c
}

// Config generate the config file data.
func (c *DMMasterScript) Config() ([]byte, error) {
	fp := path.Join("/templates", "scripts", "run_dm_master.sh.tpl")
	tpl, err := embed.ReadFile(fp)
	if err != nil {
		return nil, err
	}
	return c.ConfigWithTemplate(string(tpl))
}

// ConfigToFile write config content to specific path
func (c *DMMasterScript) ConfigToFile(file string) error {
	config, err := c.Config()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, config, 0755)
}

// ConfigWithTemplate generate the TiDB config content by tpl
func (c *DMMasterScript) ConfigWithTemplate(tpl string) ([]byte, error) {
	tmpl, err := template.New("dm-master").Parse(tpl)
	if err != nil {
		return nil, err
	}

	if c.Name == "" {
		return nil, errors.New("empty name")
	}
	for _, s := range c.Endpoints {
		if s.Name == "" {
			return nil, errors.New("empty name")
		}
	}

	content := bytes.NewBufferString("")
	if err := tmpl.Execute(content, c); err != nil {
		return nil, err
	}

	return content.Bytes(), nil
}
