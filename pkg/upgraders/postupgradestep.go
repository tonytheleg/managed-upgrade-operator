package upgraders

import (
	"context"

	"github.com/go-logr/logr"
	fileintegrityv1alpha1 "github.com/openshift/file-integrity-operator/pkg/apis/fileintegrity/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

/*
	ToDo:
	- Need to find a way to determine if the cluster being upgraded is a FedRAMP cluster fro within cluster as this only applies there
		- information is not available in the UpgradeConfig
		- is there something in cluster we can get to determine if its FR, since clusterdeployments are on hive only?
		- something added/in OCM?
*/

// PostUpgradeProcedures are any misc tasks that are needed to be completed after an upgrade has finished to ensure healthy state
func (c *clusterUpgrader) PostUpgradeProcedures(ctx context.Context, logger logr.Logger) (bool, error) {
	err := c.PostUpgradeFIOReInit(ctx, logger)
	if err != nil {
		return false, err
	}
	return true, err
}

// PostUpgradeFIOReInit reinitializes the AIDE DB in file integrity operator to track file changes due to upgrades
func (c *clusterUpgrader) PostUpgradeFIOReInit(ctx context.Context, logger logr.Logger) error {
	logger.Info("Fetching File Integrity for re-initialization")
	instance := &fileintegrityv1alpha1.FileIntegrity{}
	var osd_file_integrity = types.NamespacedName{Namespace: "openshift-file-integrity", Name: "osd-file-integrity"}

	// Create a client to connect to cluster since we are not in a reconcile loop here
	kubeConfig := controllerruntime.GetConfigOrDie()
	kclient, err := client.New(kubeConfig, client.Options{})
	if err != nil {
		return err
	}
	// Get the FileIntegrity object
	err = kclient.Get(ctx, osd_file_integrity, instance)
	if err != nil {
		return err
	}
	// Add the re-init annotation
	instance.Annotations["file-integrity.openshift.io/re-init"] = ""
	err = kclient.Update(ctx, instance)
	if err != nil {
		return err
	}
	return nil
}
