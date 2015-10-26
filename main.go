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

// This file defines a cli framework for oci2aci

package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/appc/spec/schema"
	"github.com/codegangsta/cli"
	"github.com/huawei-openlab/oci2aci/convert"
)

func main() {
	log.SetLevel(log.InfoLevel)
	app := cli.NewApp()
	app.Name = "oci2aci"
	app.Usage = "Tool for conversion between oci and aci"
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "enables debug messages",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "convert",
			Usage: "Convert oci to aci",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "os",
					Value: "linux",
					Usage: "Target OS",
				},
			},
			Action: oci2aciProc,
		},
		{
			Name:   "reversal",
			Usage:  "Convert aci to oci",
			Action: aci2ociProc,
		},
	}

	app.Run(os.Args)
}

func oci2aciProc(c *cli.Context) {
	bDebug := c.GlobalBool("debug")

	args := c.Args()

	switch len(args) {
	case 1:
		err := convert.RunOCI2ACI(args[0], bDebug)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case 2:
		outFile := args[1]
		ext := filepath.Ext(outFile)
		if ext != schema.ACIExtension {
			fmt.Fprintf(os.Stderr, "oci2aci: Extension must be %s (given %s)", schema.ACIExtension, ext)
			os.Exit(1)
		}

		aciImgPath, err := convert.Oci2aciImage(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "oci2aci: Convert failed: %s", err)
			os.Exit(1)
		}

		if err = run(exec.Command("mv", aciImgPath, args[1])); err != nil {
			os.Exit(1)
		}

	default:
		cli.ShowCommandHelp(c, "convert")

	}

	return
}

func aci2ociProc(c *cli.Context) {
	return
}

func run(cmd *exec.Cmd) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)
	return cmd.Run()
}
