package kubernetes

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// ## Defines the table structure, columns, and the function to list the data.
// ## Fungsi ini tidak perlu diubah.
func tableKubernetesCloudekaGuard(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "kubernetes_cloudeka_guard",
		Description: "CloudekaGuard is a custom resource for network policy management (Dekaguard CRD).",
		List: &plugin.ListConfig{
			Hydrate: listKubernetesCloudekaGuards,
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the CloudekaGuard."},
			{Name: "namespace", Type: proto.ColumnType_STRING, Description: "Namespace of the CloudekaGuard."},
			{Name: "uid", Type: proto.ColumnType_STRING, Description: "UID of the CloudekaGuard."},
			{Name: "labels", Type: proto.ColumnType_JSON, Description: "Labels of the CloudekaGuard."},
			{Name: "annotations", Type: proto.ColumnType_JSON, Description: "Annotations of the CloudekaGuard."},
			{Name: "endpoint_selector", Type: proto.ColumnType_JSON, Description: "EndpointSelector defines which pods this policy applies to."},
			{Name: "ingress", Type: proto.ColumnType_JSON, Description: "Ingress is a list of ingress rules."},
			{Name: "egress", Type: proto.ColumnType_JSON, Description: "Egress is a list of egress rules."},
		},
	}
}

type CloudekaGuard struct {
	// inherit from metav1.ObjectMeta
	Name             string                 `json:"name"`
	Namespace        string                 `json:"namespace"`
	UID              string                 `json:"uid"`
	Labels           map[string]string      `json:"labels"`
	Annotations      map[string]string      `json:"annotations"`
	EndpointSelector map[string]interface{} `json:"endpoint_selector"`
	Ingress          []interface{}          `json:"ingress"`
	Egress           []interface{}          `json:"egress"`
}

// ## Lists all CloudekaGuard resources from the cluster.
// ## INI ADALAH FUNGSI YANG DIPERBAIKI DAN DISERDERHANAKAN.
func listKubernetesCloudekaGuards(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Info("listKubernetesCloudekaGuards", "listing cloudekaguards")

	// Get a dynamic client to query CRDs
	dynamicClient, err := GetNewClientDynamic(ctx, d)
	if err != nil {
		logger.Error("listKubernetesCloudekaGuards", "GetNewClientDynamic error", err)
		return nil, err
	}

	// Define the Group, Version, and Resource (GVR) for CloudekaGuard
	gvr := schema.GroupVersionResource{
		Group:    "tenants.cloudeka.ai",
		Version:  "v1alpha2",
		Resource: "cloudekaguards",
	}

	// Directly list the resources across all namespaces.
	// This is simpler and more reliable since you have cluster-wide permissions.
	list, err := dynamicClient.Resource(gvr).List(ctx, metav1.ListOptions{})
	if err != nil {
		logger.Error("listKubernetesCloudekaGuards", "failed to list resources", err)
		return nil, err
	}

	// Loop through the found items and stream them to Steampipe
	for _, item := range list.Items {
		// Extract spec data
		spec := item.Object["spec"]
		var endpointSelector map[string]interface{}
		var ingress []interface{}
		var egress []interface{}

		if spec != nil {
			if specMap, ok := spec.(map[string]interface{}); ok {
				if es, exists := specMap["endpointSelector"]; exists {
					endpointSelector, _ = es.(map[string]interface{})
				}
				if ing, exists := specMap["ingress"]; exists {
					ingress, _ = ing.([]interface{})
				}
				if eg, exists := specMap["egress"]; exists {
					egress, _ = eg.([]interface{})
				}
			}
		}

		d.StreamListItem(ctx, CloudekaGuard{
			Name:             item.GetName(),
			Namespace:        item.GetNamespace(),
			UID:              string(item.GetUID()),
			Labels:           item.GetLabels(),
			Annotations:      item.GetAnnotations(),
			EndpointSelector: endpointSelector,
			Ingress:          ingress,
			Egress:           egress,
		})

		// Stop processing if the query has been cancelled or the limit has been reached
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

// NOTE: Fungsi 'listK8sNamespacesForCloudekaGuard' sudah tidak diperlukan lagi dan telah dihapus.
