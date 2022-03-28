package upgraders

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/openshift/managed-upgrade-operator/config"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	fioNamespace     string = "openshift-file-integrity"
	fioObject        string = "osd-fileintegrity"
	reinitAnnotation        = map[string]string{"file-integrity.openshift.io/re-init": ""}
)

type configManagerSpec struct {
	ConfigManager struct {
		Source     string `yaml:"source"`
		OcmBaseURL string `yaml:"ocmBaseUrl"`
	} `yaml:"configManager"`
}

// PostUpgradeProcedures are any misc tasks that are needed to be completed after an upgrade has finished to ensure healthy state
// Currently the only task is to reinit file integrity operator due to changes that come from upgrades
func (c *clusterUpgrader) PostUpgradeProcedures(ctx context.Context, logger logr.Logger) (bool, error) {

	frCluster, err := c.frClusterCheck(ctx)
	if err != nil {
		return false, err
	}
	if !frCluster {
		logger.Info("Non-FedRAMP environment...skipping PostUpgradeFIOReInit ")
		return true, nil
	}
	err = c.postUpgradeFIOReInit(ctx, logger)
	if err != nil {
		return false, err
	}
	return true, nil
}

// frClusterCheck checks to see if the upgrading cluster is a FedRAMP cluster to determine if we need to re-init the File Integrity Operator
func (c *clusterUpgrader) frClusterCheck(ctx context.Context) (bool, error) {
	ocmConfig := &corev1.ConfigMap{}
	err := c.client.Get(context.TODO(), client.ObjectKey{Namespace: config.OperatorNamespace, Name: config.ConfigMapName}, ocmConfig)
	if err != nil {
		return false, fmt.Errorf("failed to fetch %s config map to parse: %v", config.ConfigMapName, err)
	}

	var cm configManagerSpec
	err = yaml.Unmarshal([]byte(ocmConfig.Data["config.yaml"]), &cm)
	if err != nil {
		return false, fmt.Errorf("failed to parse %s config map for OCM URL: %v", config.ConfigMapName, err)
	}

	if cm.ConfigManager.Source == "OCM" {
		ocmBaseUrl := strings.TrimPrefix(cm.ConfigManager.OcmBaseURL, "https://")
		if ocmBaseUrl != "TENTATIVE-FEDRAMP-OCM-URL" {
			return false, nil
		}
	}
	return true, nil
}

// postUpgradeFIOReInit reinitializes the AIDE DB in file integrity operator to track file changes due to upgrades
func (c *clusterUpgrader) postUpgradeFIOReInit(ctx context.Context, logger logr.Logger) error {
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "fileintegrity.openshift.io",
		Kind:    "FileIntegrity",
		Version: "v1alpha1",
	})

	logger.Info("FedRAMP Environment...Fetching File Integrity for re-initialization")
	err := c.client.Get(context.TODO(), client.ObjectKey{Namespace: fioNamespace, Name: fioObject}, u)
	if err != nil {
		return fmt.Errorf("failed to fetch file integrity %s in %s namespace: %v", fioObject, fioNamespace, err)
	}

	logger.Info("Setting re-init annotation")
	u.SetAnnotations(reinitAnnotation)
	err = c.client.Update(context.TODO(), u)
	if err != nil {
		logger.Error(err, "Failed to annotate File Integrity object")
		return err
	}
	logger.Info("File Integrity Operator AIDE Datbase reinitialized")
	return nil
}
