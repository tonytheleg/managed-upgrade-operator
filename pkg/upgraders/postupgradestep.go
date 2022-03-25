package upgraders

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/openshift/managed-upgrade-operator/pkg/configmanager"
	"github.com/openshift/managed-upgrade-operator/pkg/specprovider"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	fioNamespace     string = "openshift-file-integrity"
	fioObject        string = "osd-fileintegrity"
	reinitAnnotation        = map[string]string{"file-integrity.openshift.io/re-init": ""}
)

// PostUpgradeProcedures are any misc tasks that are needed to be completed after an upgrade has finished to ensure healthy state
func (c *clusterUpgrader) PostUpgradeProcedures(ctx context.Context, logger logr.Logger) (bool, error) {

	// FIO is a FedRAMP specific operator, PostUpgradeFIOReInit is only for FedRAMP clusters
	// Check if this is an FR environment via OCM first

	/*
		the upgradeconfigmanager uses a specprovider to pull upgrade config specifications
		the specprovider interface has multiple implementations depending on the configuration in the configmap
		which defines where to pull upgrade configs from.
		if pulling from ocm, an ocmprovider will be used, which is what will then use the
		base URL for communicating to ocm for the purpose of pulling upgrade policies.
	*/
	cmb := configmanager.NewBuilder()
	spec, err := specprovider.NewBuilder().New(c.client, cmb)
	if err != nil {
		return false, err
	}
	fmt.Println(spec)
	//ocmConfig := &corev1.ConfigMap{}
	//_ = c.client.Get(context.Background(), client.ObjectKey{Namespace: fioNamespace, Name: fioObject}, ocmConfig)

	//if ocmBaseUrl.Host != "TENTATIVE-FEDRAMP-OCM-URL" {
	//	logger.Info("Non-FedRAMP environment...skipping PostUpgradeFIOReInit ")
	//}
	// err = c.PostUpgradeFIOReInit(ctx, logger)
	// if err != nil {
	// return false, err
	// }
	return true, err
}

// PostUpgradeFIOReInit reinitializes the AIDE DB in file integrity operator to track file changes due to upgrades
func (c *clusterUpgrader) PostUpgradeFIOReInit(ctx context.Context, logger logr.Logger) error {
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "fileintegrity.openshift.io",
		Kind:    "FileIntegrity",
		Version: "v1alpha1",
	})

	logger.Info("Fetching File Integrity for re-initialization")
	err := c.client.Get(context.Background(), client.ObjectKey{Namespace: fioNamespace, Name: fioObject}, u)
	if err != nil {
		return fmt.Errorf("failed to fetch file integrity %s in %s namespace: %v", fioObject, fioNamespace, err)
	}

	logger.Info("Setting re-init annotation")
	u.SetAnnotations(reinitAnnotation)
	err = c.client.Update(context.Background(), u)
	if err != nil {
		logger.Error(err, "Failed to annotate File Integrity object")
		return err
	}
	logger.Info("File Integrity Operator AIDE Datbase reinitialized")
	return nil
}
