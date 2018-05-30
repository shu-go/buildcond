package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"bitbucket.org/shu/gli"
)

type globalCmd struct {
	Pkg    string `cli:"pkg=path/to/pkg"  help:"path to the pkg"`
	Output string `cli:"output=path/to/pkg"  help:"path to the output directory (default: --pkg)"`
	Tag    string `cli:"tag=BUILD_TAG" help:"buildtag"`
	Doc    string `help:"a document for the function pkg.Do (default: Do executes function f if the buildtag '{tag}' is enabled.)"`
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
	if g.Doc == "" {
		g.Doc = "Do executes function f if the buildtag '" + g.Tag + "' is enabled."
	}

	templText := `// +build {{.Not}}{{.Tag}}

package {{.PkgLeaf}}

// {{.Doc}}
func Do(f func()) {
	{{.MayBeAComment}}f()
}
`
	templParams := struct {
		Not, Tag, PkgLeaf, Doc, MayBeAComment string
	}{
		Tag:     g.Tag,
		PkgLeaf: pkgleaf(g.Pkg),
		Doc:     g.Doc,
	}

	templ, err := template.New("tag").Parse(templText)
	if err != nil {
		return fmt.Errorf("parsing template text: %v", err)
	}

	err = os.MkdirAll(g.Output, os.ModePerm)
	if err != nil {
		return fmt.Errorf("mkdir %v: %v", g.Output, err)
	}

	filename := filepath.Join(g.Output, g.Tag+".go")
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create %v: %v", filename, err)
	}
	templParams.Not = ""
	templParams.MayBeAComment = ""
	err = templ.Execute(file, templParams)
	if err != nil {
		file.Close()
		return fmt.Errorf("write to %v: %v", filename, err)
	}

	filename = filepath.Join(g.Output, "no"+g.Tag+".go")
	file, err = os.Create(filename)
	if err != nil {
		return fmt.Errorf("create %v: %v", filename, err)
	}
	templParams.Not = "!"
	templParams.MayBeAComment = "// "
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
	app.Desc = "A go-generate command to generate a {buildtag}.Do function (generates {tag}.go and no{tag}.go) "
	app.Version = "0.1.0"
	app.Usage = `buildcond --tag=mytag
buildcond --tag=mytag --pkg=my --output=tag/my`
	app.Copyright = "(C) 2018 Shuhei Kubota"

	err := app.Run(os.Args)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func pkgleaf(pkg string) string {
	_, leaf := path.Split(pkg)
	return leaf
}
