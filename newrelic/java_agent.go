package newrelic

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/sherpa"
)

type JavaAgent struct {
	buildpackPath    string
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger
}

func NewJavaAgent(buildpackPath string, dependency libpak.BuildpackDependency, cache libpak.DependencyCache, logger bard.Logger) JavaAgent {
	contrib, _ := libpak.NewDependencyLayer(dependency, cache, libcnb.LayerTypes{
		Launch: true,
	})
	return JavaAgent{
		buildpackPath:  buildpackPath,
		LayerContributor: contrib,
		Logger: logger, 
	}
}

func (j JavaAgent) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	j.LayerContributor.Logger = j.Logger

	return j.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		j.Logger.Bodyf("Copying to %s", layer.Path)

		file := filepath.Join(layer.Path, filepath.Base(j.LayerContributor.Dependency.URI))
		if err := sherpa.CopyFile(artifact, file); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to copy artifact to %s\n%w", file, err)
		}

		layer.LaunchEnvironment.Appendf("JAVA_TOOL_OPTIONS", " ", "-javaagent:%s", file)
        
		fmt.Println("Checking for New Relic Config file...")
		file = filepath.Join(j.buildpackPath, "resources", "newrelic.yml")
		in, err := os.Open(file)
		if err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to open %s\n%w", file, err)
		}
		defer in.Close()
        
		fmt.Println("Copying New Relic Config file...")
		file = filepath.Join(layer.Path, "newrelic.yml")
		if err := sherpa.CopyFile(in, file); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to copy %s to %s\n%w", in.Name(), file, err)
		}

		return layer, nil
	})
}

func (j JavaAgent) Name() string {
	return j.LayerContributor.LayerName()
}
