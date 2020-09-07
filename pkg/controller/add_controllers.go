package controller

import (
	"gitlab.cee.redhat.com/service/osd-adoption-metrics-exporter/pkg/controller/clusterrole"
	"gitlab.cee.redhat.com/service/osd-adoption-metrics-exporter/pkg/controller/oauth"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, oauth.Add, clusterrole.Add)
}
