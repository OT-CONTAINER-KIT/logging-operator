package elasticsearch

import (
	loggingv1alpha1 "logging-operator/api/v1alpha1"

	"github.com/go-logr/logr"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("elasticsearch")

func ssetLogger(cr *loggingv1alpha1.Elasticsearch) logr.Logger {
	return log.WithValues(
		"namespace", cr.Namespace,
		"elasticsearch_name", cr.ObjectMeta.Name,
	)
}
