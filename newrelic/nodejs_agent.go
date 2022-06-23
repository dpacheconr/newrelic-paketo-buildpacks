package newrelic

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/effect"
	"github.com/paketo-buildpacks/libpak/sherpa"
)

type NodeJSAgent struct {
	buildpackPath    string
	ApplicationPath  string
	Executor         effect.Executor
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger
}

func NewNodeJSAgent(applicationPath string, buildpackPath string, dependency libpak.BuildpackDependency, cache libpak.DependencyCache, logger bard.Logger) NodeJSAgent {
	contributor, _ := libpak.NewDependencyLayer(dependency, cache, libcnb.LayerTypes{Launch: true})
	return NodeJSAgent{
		ApplicationPath:  applicationPath,
		buildpackPath:    buildpackPath,
		Executor:         effect.NewExecutor(),
		LayerContributor: contributor,
		Logger:           logger,
	}
}

func (n NodeJSAgent) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	n.LayerContributor.Logger = n.Logger

	layer, err := n.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		n.Logger.Bodyf("Installing to %s", layer.Path)

		fmt.Println("Installing New Relic module version " + filepath.Base(n.LayerContributor.Dependency.Version))
		if err := n.Executor.Execute(effect.Execution{
			Command: "npm",
			Args:    []string{"install", "newrelic@" + filepath.Base(n.LayerContributor.Dependency.Version), "--save"},
			Dir:     layer.Path,
			Stdout:  n.Logger.InfoWriter(),
			Stderr:  n.Logger.InfoWriter(),
		}); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to run npm install\n%w", err)
		}

		layer.LaunchEnvironment.Prepend("NODE_PATH", string(os.PathListSeparator), filepath.Join(layer.Path, "node_modules"))

		file := filepath.Join(layer.Path, filepath.Base(n.LayerContributor.Dependency.URI))
		if err := sherpa.CopyFile(artifact, file); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to copy artifact to %s\n%w", file, err)
		}

		fmt.Println("Checking for New Relic Config file...")
		file = filepath.Join(n.buildpackPath, "resources", "newrelic.js")
		in, err := os.Open(file)
		if err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to open %s\n%w", file, err)
		}
		defer in.Close()

		fmt.Println("Copying New Relic Config file...")
		file = filepath.Join(n.ApplicationPath, "newrelic.js")
		if err := sherpa.CopyFile(in, file); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to copy %s to %s\n%w", in.Name(), file, err)
		}

		return layer, nil
	})
	if err != nil {
		return libcnb.Layer{}, fmt.Errorf("unable to install node module\n%w", err)
	}

	m, err := sherpa.NodeJSMainModule(n.ApplicationPath)
	if err != nil {
		return libcnb.Layer{}, fmt.Errorf("unable to find main module in %s\n%w", n.ApplicationPath, err)
	}

	file := filepath.Join(n.ApplicationPath, m)
	c, err := ioutil.ReadFile(file)
	if err != nil {
		return libcnb.Layer{}, fmt.Errorf("unable to read contents of %s\n%w", file, err)
	}

	if !regexp.MustCompile(`require\(['"]newrelic['"]\);`).Match(c) {
		n.Logger.Header("Requiring 'newrelic' module")

		if err := ioutil.WriteFile(file, append([]byte("require('newrelic');\n"), c...), 0644); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to write main module %s\n%w", file, err)
		}
	}

	return layer, nil
}

func (n NodeJSAgent) Name() string {
	return n.LayerContributor.LayerName()
}
