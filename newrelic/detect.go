package newrelic

import (
	"fmt"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
)

type Detect struct{}

func (d Detect) Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	cr, err := libpak.NewConfigurationResolver(context.Buildpack, nil)
	if err != nil {
		return libcnb.DetectResult{}, fmt.Errorf("unable to create configuration resolver\n%w", err)
	}

	if !cr.ResolveBool("BP_NEW_RELIC_ENABLED") {
		return libcnb.DetectResult{Pass: false}, nil
	}

	return libcnb.DetectResult{
		Pass: true,
		Plans: []libcnb.BuildPlan{
			{
				Provides: []libcnb.BuildPlanProvide{
					{Name: "nr-agent-java"},
				},
				Requires: []libcnb.BuildPlanRequire{
					{Name: "nr-agent-java"},
					{Name: "jvm-application"},
				},
			},
			{
				Provides: []libcnb.BuildPlanProvide{
					{Name: "newrelic-nodejs"},
				},
				Requires: []libcnb.BuildPlanRequire{
					{Name: "newrelic-nodejs"},
					{Name: "node", Metadata: map[string]interface{}{"build": true}},
				},
			},
			{
				Provides: []libcnb.BuildPlanProvide{
					{Name: "newrelic-python"},
				},
				Requires: []libcnb.BuildPlanRequire{
					{Name: "newrelic-python"},
					{Name: "cpython", Metadata: map[string]interface{}{"build": true}},
				},
			},
		},
	}, nil
}
