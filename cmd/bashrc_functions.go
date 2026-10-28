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
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

var bashrcFunctionsCmd = &cobra.Command{
	Use:   "functions",
	Short: "output bash function source",
	Long: `
Write function definitions to stdout.  The function definitions may be 
written to a file and sourced from ~/.bashrc

These functions are included:

    goon	    display or activate a project
    ungoon	    undo the effects of goon 
    cdg		    cd to the goon project dir
    versions	    display installed go versions 
    is_function	    (internal use)
`,
	Run: func(cmd *cobra.Command, args []string) {
		content, err := scriptTemplate.ReadFile("scripts/bash_functions")
		cobra.CheckErr(err)
		fmt.Println(strings.TrimSpace(string(content)))
	},
}

func init() {
	CobraAddCommand(rootCmd, bashrcCmd, bashrcFunctionsCmd)
}
