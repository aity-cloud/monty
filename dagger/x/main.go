// A generated module for X functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"dagger/x/internal/dagger"
	"encoding/json"
)

type X struct{}

// Downloads the testbin
func (m *X) TestBin(config string) (*dagger.Directory, error) {
	//client := ctx.Client()
	opts := TestBinOptions{}
	if err := json.Unmarshal([]byte(config), &opts); err != nil {
		return nil, err
	}
	opts.MountOnly = false
	ctr := RunTestBin(dag, dag.Container(), opts)

	return ctr.Directory("/src/testbin"), nil
}
