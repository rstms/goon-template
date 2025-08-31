/*
Copyright © 2025 Matt Krueger <mkrueger@rstms.net>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

 1. Redistributions of source code must retain the above copyright notice,
    this list of conditions and the following disclaimer.

 2. Redistributions in binary form must reproduce the above copyright notice,
    this list of conditions and the following disclaimer in the documentation
    and/or other materials provided with the distribution.

 3. Neither the name of the copyright holder nor the names of its contributors
    may be used to endorse or promote products derived from this software
    without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/
package cmd

import (
	"embed"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed template/*
var sourceTemplate embed.FS

//go:embed scripts/*
var scriptTemplate embed.FS

var initCmd = &cobra.Command{
	Use:   "init PROGRAM_NAME",
	Short: "Create a new CLI project from this template",
	Long: `
The goon-template project is both the template and the tempate installer.
goon-template init PROGRAM_NAME will create a new project and initialize
it using the goon_init script.
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		content, err := scriptTemplate.ReadFile("scripts/goon_init")
		cobra.CheckErr(err)
		initScript, err := os.CreateTemp("", "goon_init")
		if !ViperGetBool("debug") {
			defer os.Remove(initScript.Name())
		}
		_, err = initScript.Write(content)
		cobra.CheckErr(err)
		initScript.Close()
		tempDir, err := os.MkdirTemp("", "goon-init-*")
		cobra.CheckErr(err)
		if !ViperGetBool("debug") {
			defer os.RemoveAll(tempDir)
		}
		templateFiles, err := fs.Glob(sourceTemplate, "template/*")
		cobra.CheckErr(err)
		for _, pathname := range templateFiles {
			err = copyTemplateFile(tempDir, pathname)
			cobra.CheckErr(err)
		}
		env := []string{}
		for k, v := range os.Environ() {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		for k, v := range ViperGetStringMapString("env") {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		if ViperGetBool("debug") {
			env = append(env, "DEBUG=1")
		}
		if ViperGetBool("verbose") {
			env = append(env, "VERBOSE=1")
		}
		cmdArgs := strings.Fields(ViperGetString("shell_args"))
		cmdArgs = append(cmdArgs, initScript.Name())
		cmdArgs = append(cmdArgs, args[0])
		cmdArgs = append(cmdArgs, tempDir)
		command := exec.Command(ViperGetString("shell"), cmdArgs...)
		command.Env = env
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		if ViperGetBool("debug") {
			log.Printf("command: %v\n", command)
		}
		err = command.Run()
		_, isExit := err.(*exec.ExitError)
		if isExit {
			os.Exit(command.ProcessState.ExitCode())
		}
		cobra.CheckErr(err)
	},
}

func copyTemplateFile(tempDir, srcPathname string) error {
	_, filename := filepath.Split(srcPathname)
	srcFile, err := sourceTemplate.Open(srcPathname)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(filepath.Join(tempDir, filename))
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	CobraAddCommand(rootCmd, rootCmd, initCmd)
}
