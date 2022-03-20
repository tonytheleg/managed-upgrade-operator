package upgraders

import (
	"context"

	"github.com/go-logr/logr"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
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
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}

	gvr := schema.GroupVersionResource{
		Group:    "fileintegrity.openshift.io",
		Version:  "v1alpha1",
		Resource: "fileintegrities",
	}

	logger.Info("Fetching File Integrity for re-initialization")
	fio, err := dynamicClient.Resource(gvr).Namespace("openshift-file-integrity").Get(context.Background(), "example-fileintegrity", v1.GetOptions{})
	if err != nil {
		logger.Error(err, "Failed to get File Integrity object")
		return err
	}

	logger.Info("Setting re-init annotation")
	reinit := map[string]string{"file-integrity.openshift.io/re-init": ""}
	fio.SetAnnotations(reinit)
	_, err = dynamicClient.Resource(gvr).Namespace("openshift-file-integrity").Update(context.Background(), fio, v1.UpdateOptions{})

	if err != nil {
		logger.Error(err, "Failed to annotate File Integrity object")
		return err
	}
	return nil
}
