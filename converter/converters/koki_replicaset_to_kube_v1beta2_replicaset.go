package converters

import (
	apps "k8s.io/api/apps/v1beta2"

	labelselector "github.com/koki/short/converter/converters/affinity"
	"github.com/koki/short/types"
	"github.com/koki/short/util"
)

func Convert_Koki_ReplicaSet_to_Kube_v1_ReplicaSet(rc *types.ReplicaSetWrapper) (*apps.ReplicaSet, error) {
	var err error
	kubeRC := &apps.ReplicaSet{}
	kokiRC := rc.ReplicaSet

	kubeRC.Name = kokiRC.Name
	kubeRC.Namespace = kokiRC.Namespace
	kubeRC.APIVersion = kokiRC.Version
	kubeRC.Kind = "ReplicaSet"
	kubeRC.ClusterName = kokiRC.Cluster
	kubeRC.Labels = kokiRC.Labels
	kubeRC.Annotations = kokiRC.Annotations

	kubeRC.Spec.Replicas = kokiRC.Replicas
	kubeRC.Spec.MinReadySeconds = kokiRC.MinReadySeconds
	kubeRC.Spec.Selector, err = labelselector.ParseLabelSelector(kokiRC.PodSelector)
	if err != nil {
		return nil, err
	}

	kubeTemplate, err := revertTemplate(kokiRC.Template)
	if err != nil {
		return nil, err
	}

	if kubeTemplate == nil {
		return nil, util.TypeValueErrorf(kokiRC, "missing pod template")
	}

	kubeRC.Spec.Template = *kubeTemplate

	if kokiRC.Status != nil {
		kubeRC.Status = *kokiRC.Status
	}

	return kubeRC, nil
}