package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"bitbucket.org/shu/gli"
)

type globalCmd struct {
	Tag    string `cli:"tag=BUILD_TAG"       help:"a build tag"`
	Pkg    string `cli:"pkg=path/to/pkg"     help:"path to the pkg"`
	Output string `cli:"output=path/to/pkg"  help:"path to the output directory (default: --pkg)"`
}

const templText = `// +build {{if .Not}}!{{end}}{{.Tag}}

package {{.PkgLeaf}}

// If{{capitalize .Tag}} executes function f if the build tag '{{.Tag}}' is enabled.
func If{{capitalize .Tag}}(f func()) {
	{{if .Not}}// {{end}}f()
}

// Unless{{capitalize .Tag}} executes function f if the build tag '{{.Tag}}' is disabled.
func Unless{{capitalize .Tag}}(f func()) {
	{{if .Not}}{{else}}// {{end}}f()
}

// Is{{capitalize .Tag}} returns true if the build tag '{{.Tag}}' is enabled.
func Is{{capitalize .Tag}}() bool {
	return {{if .Not}}false{{else}}true{{end}}
}
`

type templParams struct {
	Not          bool
	Tag, PkgLeaf string
}

func (g globalCmd) Run(args []string) error {
	if g.Tag == "" {
		if len(args) == 0 {
			return errors.New("tag not provided")
		}
		g.Tag = args[0]
	}
	if g.Pkg == "" {
		g.Pkg = g.Tag
	}
	if g.Output == "" {
		g.Output = g.Pkg
	}

	params := templParams{
		Tag:     g.Tag,
		PkgLeaf: pkgleaf(g.Pkg),
	}

	funcMap := template.FuncMap{
		"capitalize": strings.Title,
	}
	templ, err := template.New("tag").Funcs(funcMap).Parse(templText)
	if err != nil {
		return fmt.Errorf("parsing template text: %v", err)
	}

	// mkdir
	err = os.MkdirAll(g.Output, os.ModePerm)
	if err != nil {
		return fmt.Errorf("mkdir %v: %v", g.Output, err)
	}

	// gen tag.go
	params.Not = false
	err = generate(templ, params, g.Output, g.Tag+".go")
	if err != nil {
		return fmt.Errorf("failed to generate %v.go: %v", g.Tag, err)
	}

	// gen notag.go
	params.Not = true
	err = generate(templ, params, g.Output, "no"+g.Tag+".go")
	if err != nil {
		return fmt.Errorf("failed to generate no%v.go: %v", g.Tag, err)
	}

	return nil
}

func generate(templ *template.Template, params templParams, output, filename string) error {
	path := filepath.Join(output, filename)
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %v: %v", path, err)
	}
	defer file.Close()

	err = templ.Execute(file, params)
	if err != nil {
		return fmt.Errorf("write to %v: %v", path, err)
	}

	return nil
}

func main() {
	app := gli.New(&globalCmd{})
	app.Name = "buildcond"
	app.Desc = "A go-generate command to generate a If{tag}/Unless{tag}/Is{tag} function (generates {tag}.go and no{tag}.go) "
	app.Version = "0.1.0"
	app.Usage = `buildcond --tag=mytag
buildcond --tag=mytag --pkg=my --output=cond/my`
	app.Copyright = "(C) 2018 Shuhei Kubota"

	//app.SuppressErrorOutput = true

	err := app.Run(os.Args)

	if err != nil {
		//fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func pkgleaf(pkg string) string {
	_, leaf := path.Split(pkg)
	return leaf
}
