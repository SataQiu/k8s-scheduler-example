package nodememoryusagelimit

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	listerv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/klog/v2"
	resourcehelper "k8s.io/kubectl/pkg/util/resource"
	"k8s.io/kubernetes/pkg/scheduler/framework"

	customschedulerinternal "github.com/SataQiu/k8s-scheduler-example/api"
)

const (
	// Name is the name of the plugin used in Registry and configurations.
	Name = "NodeMemoryUsageLimit"

	// defaultNodeMemoryUsageLimit defines the default node memory usage limit
	defaultNodeMemoryUsageLimit = 70
)

// NodeMemoryUsageLimit is a scheduler plugin that only permit node with enough memory to run Pod.
type NodeMemoryUsageLimit struct {
	handle    framework.Handle
	podLister listerv1.PodLister

	nodeMemoryUsageLimit int
}

var _ framework.FilterPlugin = &NodeMemoryUsageLimit{}

// New initializes and returns a new NodeMemoryUsageLimit plugin.
func New(obj runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	args, ok := obj.(*customschedulerinternal.NodeMemoryUsageLimitArgs)
	if !ok {
		return nil, fmt.Errorf("want args to be of type NodeMemoryUsageLimit, got %T", obj)
	}

	nodeMemoryUsageLimit := defaultNodeMemoryUsageLimit
	if args.NodeMemoryUsageLimit > 0 {
		nodeMemoryUsageLimit = args.NodeMemoryUsageLimit
	}

	return &NodeMemoryUsageLimit{
		handle:               handle,
		podLister:            handle.SharedInformerFactory().Core().V1().Pods().Lister(),
		nodeMemoryUsageLimit: nodeMemoryUsageLimit,
	}, nil
}

// Name returns name of the plugin. It is used in logs, etc.
func (cs *NodeMemoryUsageLimit) Name() string {
	return Name
}

func (cs *NodeMemoryUsageLimit) Filter(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	if nodeInfo.Node() == nil || nodeInfo.Allocatable == nil {
		return framework.NewStatus(framework.Error, "node not found")
	}

	pods, err := cs.podLister.List(labels.Everything())
	if err != nil {
		return framework.NewStatus(framework.Error, "can not list Pods in cluster")
	}

	req := getPodsTotalRequestsOnNode(pods, nodeInfo.Node().Name)
	memoryReq := req[corev1.ResourceMemory]

	node := nodeInfo.Node()
	allocatable := node.Status.Capacity
	if len(node.Status.Allocatable) > 0 {
		allocatable = node.Status.Allocatable
	}

	fractionMemoryReq := float64(memoryReq.Value()) / float64(allocatable.Memory().Value()) * 100

	klog.Infof("Current Node memory usage %v, limit %v", fractionMemoryReq, cs.nodeMemoryUsageLimit)

	if fractionMemoryReq > float64(cs.nodeMemoryUsageLimit) {
		return framework.NewStatus(framework.Unschedulable,
			fmt.Sprintf("node memory usage reach the limit %v", cs.nodeMemoryUsageLimit))
	}

	return nil
}

func getPodsTotalRequestsOnNode(podList []*corev1.Pod, nodeName string) (reqs map[corev1.ResourceName]resource.Quantity) {
	reqs = map[corev1.ResourceName]resource.Quantity{}
	for i := range podList {
		pod := podList[i]
		if pod.Spec.NodeName != nodeName {
			continue
		}
		podReqs, _ := resourcehelper.PodRequestsAndLimits(pod)
		for podReqName, podReqValue := range podReqs {
			if value, ok := reqs[podReqName]; !ok {
				reqs[podReqName] = podReqValue.DeepCopy()
			} else {
				value.Add(podReqValue)
				reqs[podReqName] = value
			}
		}
	}
	return
}
