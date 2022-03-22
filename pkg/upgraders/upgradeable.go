package upgraders

import (
	"context"
	"fmt"

	"github.com/blang/semver"
	"github.com/go-logr/logr"
	configv1 "github.com/openshift/api/config/v1"
	upgradev1alpha1 "github.com/openshift/managed-upgrade-operator/pkg/apis/upgrade/v1alpha1"
	cv "github.com/openshift/managed-upgrade-operator/pkg/clusterversion"
)

func (c *clusterUpgrader) IsUpgradeable(ctx context.Context, logger logr.Logger) (bool, error) {
	upgradeCommenced, err := c.cvClient.HasUpgradeCommenced(c.upgradeConfig)
	if err != nil {
		return false, err
	}
	if upgradeCommenced {
		logger.Info(fmt.Sprintf("Skipping upgrade step %s", upgradev1alpha1.IsClusterUpgradable))
		return true, nil
	}

	clusterVersion, err := c.cvClient.GetClusterVersion()
	if err != nil {
		return false, err
	}
	currentVersion, err := cv.GetCurrentVersion(clusterVersion)
	if err != nil {
		return false, err
	}
	parsedCurrentVersion, err := semver.Parse(currentVersion)
	if err != nil {
		return false, err
	}

	desiredVersion := c.upgradeConfig.Spec.Desired.Version
	parsedDesiredVersion, err := semver.Parse(desiredVersion)
	if err != nil {
		return false, err
	}

	// if the upgradeable is false then we need to check the current version with upgrade version for y-stream update
	for _, condition := range clusterVersion.Status.Conditions {
		if condition.Type == configv1.OperatorUpgradeable && condition.Status == configv1.ConditionFalse && parsedDesiredVersion.Major >= parsedCurrentVersion.Major && parsedDesiredVersion.Minor > parsedCurrentVersion.Minor {
			return false, fmt.Errorf("cCluster upgrade to version %s is canceled with the reason of %s containing message that %s Automated upgrades will be retried on their next scheduling cycle. If you have manually scheduled an upgrade instead, it must be rescheduled", desiredVersion, condition.Reason, condition.Message)
		}
	}

	return true, nil
}
