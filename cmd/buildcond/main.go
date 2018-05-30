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

	templText := `// +build {{if .Not}}!{{end}}{{.Tag}}

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
	templParams := struct {
		Not          bool
		Tag, PkgLeaf string
	}{
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
	templParams.Not = false
	filename := filepath.Join(g.Output, g.Tag+".go")
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create %v: %v", filename, err)
	}
	err = templ.Execute(file, templParams)
	if err != nil {
		file.Close()
		return fmt.Errorf("write to %v: %v", filename, err)
	}

	// gen notag.go
	templParams.Not = true
	filename = filepath.Join(g.Output, "no"+g.Tag+".go")
	file, err = os.Create(filename)
	if err != nil {
		return fmt.Errorf("create %v: %v", filename, err)
	}
	err = templ.Execute(file, templParams)
	if err != nil {
		file.Close()
		return fmt.Errorf("write to %v: %v", filename, err)
	}

	return nil
}

func main() {
	app := gli.New(&globalCmd{})
	app.Name = "buildcond"
	app.Desc = "A go-generate command to generate a {build tag}.Do function (generates {tag}.go and no{tag}.go) "
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
