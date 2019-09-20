package deployments

import (
	"fmt"

	"github.com/mesosphere/bun"
)

// max number of Marathon deployments considered healthy
const maxDeployments = 10

func init() {
	builder := bun.CheckBuilder{
		Name:               "marathon-deployments",
		Description:        "Check for too many running Marathon app deployments",
		CollectFromMasters: collect,
		Aggregate:          bun.DefaultAggregate,
	}
	check := builder.Build()
	bun.RegisterCheck(check)
}

func collect(host bun.Host) (ok bool, details interface{}, err error) {
	deployments := []struct{}{}
	if err = host.ReadJSON("marathon-deployments", &deployments); err != nil {
		return
	}
	if len(deployments) > maxDeployments {
		details = fmt.Sprintf("Too many deployments: %v", len(deployments))
		return
	}
	ok = true
	return
}
