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
func (c *clusterUpgrader) PostUpgradeProcedures(ctx context.Context, logger logr.Logger) (bool, error) {

	// FIO is a FedRAMP specific operator, PostUpgradeFIOReInit is only for FedRAMP clusters
	// Check if this is an FR environment by looking at the MUO Operator config
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
			logger.Info("Non-FedRAMP environment...skipping PostUpgradeFIOReInit ")
		}
		err = c.PostUpgradeFIOReInit(ctx, logger)
		if err != nil {
			return false, err
		}
	}
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
