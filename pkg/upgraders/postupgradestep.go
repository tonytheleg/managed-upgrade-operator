package upgraders

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/openshift/managed-upgrade-operator/pkg/ocm"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

var (
	fioNamespace     string = "openshift-file-integrity"
	fioObject        string = "osd-fileintegrity"
	reinitAnnotation        = map[string]string{"file-integrity.openshift.io/re-init": ""}
)

/*
	ToDo:
	- Need to find a way to determine if the cluster being upgraded is a FedRAMP cluster from within cluster as this only applies there
		- information is not available in the UpgradeConfig
		- is there something in cluster we can get to determine if its FR, since clusterdeployments are on hive only?
		- something added/in OCM?

	Solution: We can trigger of the OCM base URL for API calls available in the managed-upgrade-operator-config deployed by MCC.
			  This tells MUO what OCM to poll, and since the OCM API is region specific, we can check for the FedRAMP OCM API address first
*/

// PostUpgradeProcedures are any misc tasks that are needed to be completed after an upgrade has finished to ensure healthy state
func (c *clusterUpgrader) PostUpgradeProcedures(ctx context.Context, logger logr.Logger) (bool, error) {

	// FIO is a FedRAMP specific operator, PostUpgradeFIOReInit is only for FedRAMP clusters
	// Check if this is an FR environment via OCM first
	ocmConfig := ocm.OcmClientConfig{}
	ocmBaseUrl := ocmConfig.GetOCMBaseURL()
	fmt.Println("OCM BASE URL is ", ocmBaseUrl.Host)
	if ocmBaseUrl.Host != "api.stage.openshift.com" {
		logger.Info("Non-FedRAMP environment...skipping PostUpgradeFIOReInit ")
	}
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
		return fmt.Errorf("failed to fetch config object for in cluster authentication: %v", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %v", err)
	}

	gvr := schema.GroupVersionResource{
		Group:    "fileintegrity.openshift.io",
		Version:  "v1alpha1",
		Resource: "fileintegrities",
	}

	logger.Info("Fetching File Integrity for re-initialization")
	fio, err := dynamicClient.Resource(gvr).Namespace(fioNamespace).Get(context.Background(), fioObject, v1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to fetch file integrity %s in %s namespace: %v", fioObject, fioNamespace, err)
	}

	logger.Info("Setting re-init annotation")
	fio.SetAnnotations(reinitAnnotation)
	_, err = dynamicClient.Resource(gvr).Namespace(fioNamespace).Update(context.Background(), fio, v1.UpdateOptions{})

	if err != nil {
		logger.Error(err, "Failed to annotate File Integrity object")
		return err
	}
	logger.Info("File Integrity Operator AIDE Datbase reinitialized")
	return nil
}
