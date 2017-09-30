/*******************************************************************************
*
* Copyright 2017 Stefan Majewsky <majewsky@gmx.net>
*
* This program is free software: you can redistribute it and/or modify it under
* the terms of the GNU General Public License as published by the Free Software
* Foundation, either version 3 of the License, or (at your option) any later
* version.
*
* This program is distributed in the hope that it will be useful, but WITHOUT ANY
* WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR
* A PARTICULAR PURPOSE. See the GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License along with
* this program. If not, see <http://www.gnu.org/licenses/>.
*
*******************************************************************************/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/MasterMinds/sprig"
)

const defaultDirectory = `/etc/replicator.d`
const configGlob = `/etc/replicator.d/*.toml`

func failIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func readConfig(dir string) map[string]interface{} {
	paths, err := filepath.Glob(filepath.Join(dir, "*.toml"))
	failIf(err)

	result := map[string]interface{}{}
	for _, path := range paths {
		bytes, err := ioutil.ReadFile(path)
		failIf(err)
		var next map[string]interface{}
		failIf(toml.Unmarshal(bytes, &next))
		result, err = MergeTables(result, next)
		failIf(err)
	}

	return result
}

func main() {
	dir := os.Getenv("REPLICATOR_VARIABLES_DIR")
	if dir == "" {
		dir = defaultDirectory
	}

	var locals struct {
		Vars map[string]interface{}
	}
	locals.Vars = readConfig(dir)

	tmplText, err := ioutil.ReadAll(os.Stdin)
	failIf(err)
	tmpl, err := template.New("stdin").Funcs(sprig.TxtFuncMap()).Parse(string(tmplText))
	failIf(err)
	failIf(tmpl.Execute(os.Stdout, &locals))
}
