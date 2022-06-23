package newrelic

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/effect"
	"github.com/paketo-buildpacks/libpak/sherpa"
)

type PythonAgent struct {
	buildpackPath    string
	ApplicationPath  string
	Executor         effect.Executor
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger
}

func NewPythonAgent(applicationPath string, buildpackPath string, dependency libpak.BuildpackDependency, cache libpak.DependencyCache, logger bard.Logger) PythonAgent {
	contributor, _ := libpak.NewDependencyLayer(dependency, cache, libcnb.LayerTypes{Launch: true})
	return PythonAgent{
		ApplicationPath:  applicationPath,
		buildpackPath:    buildpackPath,
		Executor:         effect.NewExecutor(),
		LayerContributor: contributor,
		Logger:           logger,
	}
}

func (p PythonAgent) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	p.LayerContributor.Logger = p.Logger

	layer, err := p.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		p.Logger.Bodyf("Installing to %s", layer.Path)

		file := filepath.Join(layer.Path, filepath.Base(p.LayerContributor.Dependency.URI))
		if err := sherpa.CopyFile(artifact, file); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to copy artifact to %s\n%w", file, err)
		}

		fmt.Println("Installing New Relic Python Agent using pip3")
		if err := p.Executor.Execute(effect.Execution{
			Command: "pip3",
			Args:    []string{"install", "newrelic==7.10.0.175"},
			Dir:     p.ApplicationPath,
			Stdout:  p.Logger.InfoWriter(),
			Stderr:  p.Logger.InfoWriter(),
		}); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to run pip install\n%w", err)
		}

		fmt.Println("Checking for New Relic Config file...")
		file = filepath.Join(p.buildpackPath, "resources", "newrelic.ini")
		in, err := os.Open(file)
		if err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to open %s\n%w", file, err)
		}
		defer in.Close()

		fmt.Println("Copying New Relic Config file...")
		file = filepath.Join(p.ApplicationPath, "newrelic.ini")
		if err := sherpa.CopyFile(in, file); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to copy %s to %s\n%w", in.Name(), file, err)
		}

		return layer, nil
	})
	if err != nil {
		return libcnb.Layer{}, fmt.Errorf("unable to install python agent\n%w", err)
	}
	return layer, nil
}

func (p PythonAgent) Name() string {
	return p.LayerContributor.LayerName()
}
