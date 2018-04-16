/*
Kubernetes labels and annotations used in Conduit's control plane and data plane
Kubernetes configs.
*/

package k8s

import (
	"fmt"
	"strings"

	"github.com/runconduit/conduit/pkg/version"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	/*
	 * Labels
	 */

	// ControllerComponentLabel identifies this object as a component of Conduit's
	// control plane (e.g. web, controller).
	ControllerComponentLabel = "conduit.io/control-plane-component"

	// ControllerNSLabel is injected into mesh-enabled apps, identifying the
	// namespace of the Conduit control plane.
	ControllerNSLabel = "conduit.io/control-plane-ns"

	// ProxyDeploymentLabel is injected into mesh-enabled apps, identifying the
	// deployment that this proxy belongs to.
	ProxyDeploymentLabel = "conduit.io/proxy-deployment"

	// ProxyReplicationControllerLabel is injected into mesh-enabled apps,
	// identifying the ReplicationController that this proxy belongs to.
	ProxyReplicationControllerLabel = "conduit.io/proxy-replication-controller"

	// ProxyReplicaSetLabel is injected into mesh-enabled apps, identifying the
	// ReplicaSet that this proxy belongs to.
	ProxyReplicaSetLabel = "conduit.io/proxy-replica-set"

	// ProxyJobLabel is injected into mesh-enabled apps, identifying the Job that
	// this proxy belongs to.
	ProxyJobLabel = "conduit.io/proxy-job"

	// ProxyDaemonSetLabel is injected into mesh-enabled apps, identifying the
	// DaemonSet that this proxy belongs to.
	ProxyDaemonSetLabel = "conduit.io/proxy-daemon-set"

	/*
	 * Annotations
	 */

	// CreatedByAnnotation indicates the source of the injected data plane
	// (e.g. conduit/cli v0.1.3).
	CreatedByAnnotation = "conduit.io/created-by"

	// ProxyVersionAnnotation indicates the version of the injected data plane
	// (e.g. v0.1.3).
	ProxyVersionAnnotation = "conduit.io/proxy-version"

	/*
	 * Component Names
	 */

	// CertificateBundleName is the name of the ConfigMap that holds the root
	// certificate
	CertificateBundleName = "conduit-ca-bundle"
)

var proxyLabels = []string{
	ProxyDeploymentLabel,
	ProxyReplicationControllerLabel,
	ProxyReplicaSetLabel,
	ProxyJobLabel,
	ProxyDaemonSetLabel,
}

// CreatedByAnnotationValue returns the value associated with
// CreatedByAnnotation.
func CreatedByAnnotationValue() string {
	return fmt.Sprintf("conduit/cli %s", version.Version)
}

// GetOwnerLabels returns the set of prometheus owner labels that can be
// extracted from the proxy labels that have been added to an injected
// kubernetes resource
func GetOwnerLabels(objectMeta meta.ObjectMeta) map[string]string {
	labels := make(map[string]string)
	for _, label := range proxyLabels {
		if labelValue, ok := objectMeta.Labels[label]; ok {
			labels[toOwnerLabel(label)] = labelValue
		}
	}
	return labels
}

// toOwnerLabel converts a proxy label to a prometheus label, following the
// relabel conventions from the prometheus scrape config file
func toOwnerLabel(proxyLabel string) string {
	stripped := strings.TrimPrefix(proxyLabel, "conduit.io/proxy-")
	if stripped == "job" {
		return "k8s_job"
	}
	return strings.Replace(stripped, "-", "_", -1)
}
