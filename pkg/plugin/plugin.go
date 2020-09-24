package plugin

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubectl/pkg/cmd/util"
	resourcehelper "k8s.io/kubectl/pkg/util/resource"
)

func RunPlugin(configFlags *genericclioptions.ConfigFlags, cmd *cobra.Command) error {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return errors.Wrap(err, "failed to read kubeconfig")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "failed to create clientset")
	}
	return getNodeResource(clientset, cmd)
}

//get Resource
func getNodeResource(cli *kubernetes.Clientset, cmd *cobra.Command) error {
	selector := util.GetFlagString(cmd, "selector")
	fieldSelector := util.GetFlagString(cmd, "field-selector")
	nodeName := util.GetFlagString(cmd, "node")
	fmt.Printf("Node:\t\t\t\tCPU\tMemory\t\tCPURequests\tCPULimits\tMemoryRequests\t\tMemoryLimits\n")
	if len(nodeName) > 0 {
		node, err := cli.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
		if err != nil {
			return errors.Wrap(err, "failed to get node")
		}
		namespace := ""
		fieldSelector, err := fields.ParseSelector("spec.nodeName=" + nodeName + ",status.phase!=" + string(corev1.PodSucceeded) + ",status.phase!=" + string(corev1.PodFailed))
		if err != nil {
			return err
		}
		podList, err := cli.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{FieldSelector: fieldSelector.String()})
		describeNodeResource(podList, node)
		return nil
	}
	nodeList, err := cli.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{
		LabelSelector: selector,
		FieldSelector: fieldSelector,
	})
	if err != nil {
		return errors.Wrap(err, "failed to list  nodes ")
	}
	for _, n := range nodeList.Items {
		namespace := ""
		fieldSelector, err := fields.ParseSelector("spec.nodeName=" + n.Name + ",status.phase!=" + string(corev1.PodSucceeded) + ",status.phase!=" + string(corev1.PodFailed))
		if err != nil {
			return err
		}
		podList, err := cli.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{FieldSelector: fieldSelector.String()})
		describeNodeResource(podList, &n)
	}
	return nil

}

func getPodsTotalRequestsAndLimits(podList *corev1.PodList) (reqs map[corev1.ResourceName]resource.Quantity, limits map[corev1.ResourceName]resource.Quantity) {
	reqs, limits = map[corev1.ResourceName]resource.Quantity{}, map[corev1.ResourceName]resource.Quantity{}
	for _, pod := range podList.Items {
		podReqs, podLimits := resourcehelper.PodRequestsAndLimits(&pod)
		for podReqName, podReqValue := range podReqs {
			if value, ok := reqs[podReqName]; !ok {
				reqs[podReqName] = podReqValue.DeepCopy()
			} else {
				value.Add(podReqValue)
				reqs[podReqName] = value
			}
		}
		for podLimitName, podLimitValue := range podLimits {
			if value, ok := limits[podLimitName]; !ok {
				limits[podLimitName] = podLimitValue.DeepCopy()
			} else {
				value.Add(podLimitValue)
				limits[podLimitName] = value
			}
		}
	}
	return
}

func describeNodeResource(nodeNonTerminatedPodsList *corev1.PodList, node *corev1.Node) {

	allocatable := node.Status.Capacity
	if len(node.Status.Allocatable) > 0 {
		allocatable = node.Status.Allocatable
	}
	allocatablememory := allocatable[corev1.ResourceMemory]
	allocatablecpu := allocatable[corev1.ResourceCPU]
	reqs, limits := getPodsTotalRequestsAndLimits(nodeNonTerminatedPodsList)
	cpuReqs, cpuLimits, memoryReqs, memoryLimits, _, _ :=
		reqs[corev1.ResourceCPU], limits[corev1.ResourceCPU], reqs[corev1.ResourceMemory], limits[corev1.ResourceMemory], reqs[corev1.ResourceEphemeralStorage], limits[corev1.ResourceEphemeralStorage]
	fractionCpuReqs := float64(0)
	fractionCpuLimits := float64(0)
	if allocatable.Cpu().MilliValue() != 0 {
		fractionCpuReqs = float64(cpuReqs.MilliValue()) / float64(allocatable.Cpu().MilliValue()) * 100
		fractionCpuLimits = float64(cpuLimits.MilliValue()) / float64(allocatable.Cpu().MilliValue()) * 100
	}
	fractionMemoryReqs := float64(0)
	fractionMemoryLimits := float64(0)
	if allocatable.Memory().Value() != 0 {
		fractionMemoryReqs = float64(memoryReqs.Value()) / float64(allocatable.Memory().Value()) * 100
		fractionMemoryLimits = float64(memoryLimits.Value()) / float64(allocatable.Memory().Value()) * 100
	}
	//resourcePrinter:=printers.NewTablePrinter(printers.PrintOptions{})
	//tableWriter:=printers.GetNewTabWriter(os.Stdout)
	//tableWriter.RememberedWidths()
	//TODO
	fmt.Printf("%s\t%s\t%fG", node.Name, allocatablecpu.String(), float64(allocatablememory.Value()/1024/1024/1024))
	fmt.Printf("\t%.2fcore (%d%%)\t%.2fcore(%d%%)", float64(cpuReqs.MilliValue())/float64(1000), int64(fractionCpuReqs), float64(cpuLimits.MilliValue())/float64(1000), int64(fractionCpuLimits))
	//fmt.Printf("%s\t%s (%d%%)\t%s (%d%%)\n", corev1.ResourceCPU, cpuReqs.String(), int64(fractionCpuReqs), cpuLimits.String(), int64(fractionCpuLimits))
	fmt.Printf("\t%s (%d%%)\t%s (%d%%)\n", memoryReqs.String(), int64(fractionMemoryReqs), memoryLimits.String(), int64(fractionMemoryLimits))
}
