package query

import (
	"context"
	"encoding/json"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	// "github.com/kubecost/cost-model/pkg/costmodel"
	"github.com/kubecost/cost-model/pkg/kubecost"
	"github.com/kubecost/cost-model/pkg/util"
)

type aggCostModelResponse struct {
	Code int `json:"code"`
	// Data map[string]costmodel.Aggregation `json:"data"`
	Data map[string]Aggregation `json:"data"`
}

type AggCostModelParameters struct {
	RestConfig *rest.Config
	Ctx        context.Context

	KubecostNamespace   string
	ServiceName         string
	Window              string
	Aggregate           string
	AggregationSubfield string
	UseProxy            bool
}

// QueryAggCostModel queries /model/aggregatedCostModel by proxying a request to Kubecost
// through the Kubernetes API server if useProxy is true or, if it isn't, by
// temporarily port forwarding to a Kubecost pod.
func QueryAggCostModel(p AggCostModelParameters) (map[string]Aggregation, error) {
	requestParams := map[string]string{
		"window":      p.Window,
		"aggregation": p.Aggregate,
		"rate":        "monthly",
		"etl":         "true",
	}

	if p.AggregationSubfield != "" {
		requestParams["aggregationSubfield"] = p.AggregationSubfield
	}

	var bytes []byte
	var err error
	if p.UseProxy {
		clientset, err := kubernetes.NewForConfig(p.RestConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create clientset: %s", err)
		}

		bytes, err = clientset.CoreV1().Services(p.KubecostNamespace).ProxyGet("", p.ServiceName, "9090", "/model/aggregatedCostModel", requestParams).DoRaw(p.Ctx)

		if err != nil {
			return nil, fmt.Errorf("failed to proxy get kubecost. err: %s; data: %s", err, bytes)
		}
	} else {
		bytes, err = portForwardedQueryService(p.RestConfig, p.KubecostNamespace, p.ServiceName, "model/aggregatedCostModel", requestParams, p.Ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to port forward query: %s", err)
		}
	}

	var ar aggCostModelResponse
	err = json.Unmarshal(bytes, &ar)
	if err != nil {
		return ar.Data, fmt.Errorf("failed to unmarshal aggregatedCostModel response: %s", err)
	}

	return ar.Data, nil
}

// Hardcoded instead of imported because of dependency problems introduced when
// github.com/kubecost/cost-model/pkg/costmodel is imported. The breakage involves
// Azure's go-autorest, the azure-sdk-for-go, and k8s client-go. Basically, cost-model
// uses a very old version of client-go, etc. that causes a breakage.
type Aggregation struct {
	Aggregator                 string                         `json:"aggregation"`
	Subfields                  []string                       `json:"subfields,omitempty"`
	Environment                string                         `json:"environment"`
	Cluster                    string                         `json:"cluster,omitempty"`
	Properties                 *kubecost.AllocationProperties `json:"-"`
	CPUAllocationHourlyAverage float64                        `json:"cpuAllocationAverage"`
	CPUAllocationVectors       []*util.Vector                 `json:"-"`
	CPUAllocationTotal         float64                        `json:"-"`
	CPUCost                    float64                        `json:"cpuCost"`
	CPUCostVector              []*util.Vector                 `json:"cpuCostVector,omitempty"`
	CPUEfficiency              float64                        `json:"cpuEfficiency"`
	CPURequestedVectors        []*util.Vector                 `json:"-"`
	CPUUsedVectors             []*util.Vector                 `json:"-"`
	Efficiency                 float64                        `json:"efficiency"`
	GPUAllocationHourlyAverage float64                        `json:"gpuAllocationAverage"`
	GPUAllocationVectors       []*util.Vector                 `json:"-"`
	GPUCost                    float64                        `json:"gpuCost"`
	GPUCostVector              []*util.Vector                 `json:"gpuCostVector,omitempty"`
	GPUAllocationTotal         float64                        `json:"-"`
	RAMAllocationHourlyAverage float64                        `json:"ramAllocationAverage"`
	RAMAllocationVectors       []*util.Vector                 `json:"-"`
	RAMAllocationTotal         float64                        `json:"-"`
	RAMCost                    float64                        `json:"ramCost"`
	RAMCostVector              []*util.Vector                 `json:"ramCostVector,omitempty"`
	RAMEfficiency              float64                        `json:"ramEfficiency"`
	RAMRequestedVectors        []*util.Vector                 `json:"-"`
	RAMUsedVectors             []*util.Vector                 `json:"-"`
	PVAllocationHourlyAverage  float64                        `json:"pvAllocationAverage"`
	PVAllocationVectors        []*util.Vector                 `json:"-"`
	PVAllocationTotal          float64                        `json:"-"`
	PVCost                     float64                        `json:"pvCost"`
	PVCostVector               []*util.Vector                 `json:"pvCostVector,omitempty"`
	NetworkCost                float64                        `json:"networkCost"`
	NetworkCostVector          []*util.Vector                 `json:"networkCostVector,omitempty"`
	SharedCost                 float64                        `json:"sharedCost"`
	TotalCost                  float64                        `json:"totalCost"`
	TotalCostVector            []*util.Vector                 `json:"totalCostVector,omitempty"`
}
