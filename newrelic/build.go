package newrelic

import (
	"fmt"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

type Build struct {
	Logger bard.Logger
}

func (b Build) Build(context libcnb.BuildContext) (libcnb.BuildResult, error) {
	b.Logger.Title(context.Buildpack)
	result := libcnb.NewBuildResult()

	pr := libpak.PlanEntryResolver{Plan: context.Plan}

	dr, err := libpak.NewDependencyResolver(context)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create dependency resolver\n%w", err)
	}

	dc, err := libpak.NewDependencyCache(context)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create dependency cache\n%w", err)
	}
	dc.Logger = b.Logger

	if _, ok, err := pr.Resolve("nr-agent-java"); err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve nr-agent-java plan entry\n%w", err)
	} else if ok {
		agentDependency, err := dr.Resolve("nr-agent-java", "")
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
		}

		result.Layers = append(result.Layers, NewJavaAgent(context.Buildpack.Path, agentDependency, dc, b.Logger))
	}

	if _, ok, err := pr.Resolve("newrelic-nodejs"); err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve newrelic-nodejs plan entry\n%w", err)
	} else if ok {
		dep, err := dr.Resolve("newrelic-nodejs", "")
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
		}

		result.Layers = append(result.Layers, NewNodeJSAgent(context.Application.Path, context.Buildpack.Path, dep, dc, b.Logger))
	}
	if _, ok, err := pr.Resolve("newrelic-python"); err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve newrelic-python plan entry\n%w", err)
	} else if ok {
		dep, err := dr.Resolve("newrelic-python", "")
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
		}

		result.Layers = append(result.Layers, NewPythonAgent(context.Application.Path, context.Buildpack.Path, dep, dc, b.Logger))
	}

	return result, nil
}
