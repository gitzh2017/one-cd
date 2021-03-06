package deployer

import (
	"bytes"
	"fmt"
	"io"

	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	listerv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

// PodEvents 获取pod事件
func (d *Deployer) PodEvents(cluster, namespace, podName string) (list *coreV1.EventList, err error) {
	var client *kubernetes.Clientset
	if client, err = d.Client(cluster); err != nil {
		return
	}
	fieldSelector := fmt.Sprintf("involvedObject.kind=Pod,involvedObject.name=%s", podName)
	if list, err = client.CoreV1().Events(namespace).List(d.ctx, metav1.ListOptions{FieldSelector: fieldSelector}); err != nil {
		return
	}
	return
}

// PodList ...
func (d *Deployer) PodList(cluster, namespace, deploymentName string) (podList []*coreV1.Pod, err error) {
	var informer cache.SharedIndexInformer
	if informer, err = d.PodInformer(cluster); err != nil {
		return
	}
	parse, _ := labels.Parse(fmt.Sprintf("app=%s", deploymentName))
	if podList, err = listerv1.NewPodLister(informer.GetIndexer()).Pods(namespace).List(parse); err != nil {
		return
	}
	return
}

// PodLog 获取pod日志
func (d *Deployer) PodLog(cluster, namespace, podName, container string, sinceSeconds int64, previous bool) (log string, err error) {
	var client *kubernetes.Clientset
	if client, err = d.Client(cluster); err != nil {
		return
	}
	req := client.CoreV1().Pods(namespace).GetLogs(podName, &coreV1.PodLogOptions{
		Container:    container,
		SinceSeconds: &sinceSeconds,
		Previous:     previous})
	podLogs, err := req.Stream(d.ctx)
	if err != nil {
		return
	}
	defer podLogs.Close()
	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, podLogs); err != nil {
		return
	}
	log = buf.String()
	return
}

// PodDelete 删除pod
func (d *Deployer) PodDelete(cluster, namespace, podName string) (err error) {
	var client *kubernetes.Clientset
	if client, err = d.Client(cluster); err != nil {
		return
	}
	if err = client.CoreV1().Pods(namespace).Delete(d.ctx, podName, metav1.DeleteOptions{}); err != nil {
		return
	}
	return
}
