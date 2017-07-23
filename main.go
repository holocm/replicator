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
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/MasterMinds/sprig"
)

const configGlob = `/etc/replicator.d/*.toml`

func failIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func readConfig() map[string]interface{} {
	paths, err := filepath.Glob(configGlob)
	failIf(err)

	text := ""
	for _, path := range paths {
		bytes, err := ioutil.ReadFile(path)
		failIf(err)
		text += string(bytes) + "\n"
	}

	result := map[string]interface{}{}
	failIf(toml.Unmarshal([]byte(text), &result))
	return result
}

func main() {
	var locals struct {
		Vars map[string]interface{}
	}
	locals.Vars = readConfig()

	tmplText, err := ioutil.ReadAll(os.Stdin)
	failIf(err)
	tmpl, err := template.New("stdin").Funcs(sprig.FuncMap()).Parse(string(tmplText))
	failIf(err)
	failIf(tmpl.Execute(os.Stdout, &locals))
}
