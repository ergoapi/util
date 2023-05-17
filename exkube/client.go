package exkube

import (
	"bytes"
	"context"
	"fmt"
	"io"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Client struct {
	Clientset          kubernetes.Interface
	ExtensionClientset apiextensionsclientset.Interface // k8s api extension needed to retrieve CRDs
	DynamicClientset   dynamic.Interface
	Config             *rest.Config
}

func NewClient(cc *ClientConfig) (*Client, error) {
	config, err := NewRestConfig(cc)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	extensionClientset, err := apiextensionsclientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	dynamicClientset, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		Clientset:          clientset,
		ExtensionClientset: extensionClientset,
		Config:             config,
		DynamicClientset:   dynamicClientset,
	}, nil
}

func (c *Client) CreateSecret(ctx context.Context, namespace string, secret *corev1.Secret, opts metav1.CreateOptions) (*corev1.Secret, error) {
	return c.Clientset.CoreV1().Secrets(namespace).Create(ctx, secret, opts)
}

func (c *Client) UpdateSecret(ctx context.Context, namespace string, secret *corev1.Secret, opts metav1.UpdateOptions) (*corev1.Secret, error) {
	return c.Clientset.CoreV1().Secrets(namespace).Update(ctx, secret, opts)
}

func (c *Client) PatchSecret(ctx context.Context, namespace, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions) (*corev1.Secret, error) {
	return c.Clientset.CoreV1().Secrets(namespace).Patch(ctx, name, pt, data, opts)
}

func (c *Client) DeleteSecret(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.CoreV1().Secrets(namespace).Delete(ctx, name, opts)
}

func (c *Client) GetSecret(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.Secret, error) {
	return c.Clientset.CoreV1().Secrets(namespace).Get(ctx, name, opts)
}

func (c *Client) CreateServiceAccount(ctx context.Context, namespace string, account *corev1.ServiceAccount, opts metav1.CreateOptions) (*corev1.ServiceAccount, error) {
	return c.Clientset.CoreV1().ServiceAccounts(namespace).Create(ctx, account, opts)
}

func (c *Client) DeleteServiceAccount(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.CoreV1().ServiceAccounts(namespace).Delete(ctx, name, opts)
}

func (c *Client) GetClusterRole(ctx context.Context, name string, opts metav1.GetOptions) (*rbacv1.ClusterRole, error) {
	return c.Clientset.RbacV1().ClusterRoles().Get(ctx, name, opts)
}

func (c *Client) CreateClusterRole(ctx context.Context, role *rbacv1.ClusterRole, opts metav1.CreateOptions) (*rbacv1.ClusterRole, error) {
	return c.Clientset.RbacV1().ClusterRoles().Create(ctx, role, opts)
}

func (c *Client) DeleteClusterRole(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.RbacV1().ClusterRoles().Delete(ctx, name, opts)
}

func (c *Client) CreateClusterRoleBinding(ctx context.Context, role *rbacv1.ClusterRoleBinding, opts metav1.CreateOptions) (*rbacv1.ClusterRoleBinding, error) {
	return c.Clientset.RbacV1().ClusterRoleBindings().Create(ctx, role, opts)
}

func (c *Client) DeleteClusterRoleBinding(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.RbacV1().ClusterRoleBindings().Delete(ctx, name, opts)
}

func (c *Client) CreateRole(ctx context.Context, namespace string, role *rbacv1.Role, opts metav1.CreateOptions) (*rbacv1.Role, error) {
	return c.Clientset.RbacV1().Roles(namespace).Create(ctx, role, opts)
}

func (c *Client) UpdateRole(ctx context.Context, namespace string, role *rbacv1.Role, opts metav1.UpdateOptions) (*rbacv1.Role, error) {
	return c.Clientset.RbacV1().Roles(namespace).Update(ctx, role, opts)
}

func (c *Client) DeleteRole(ctx context.Context, namespace string, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.RbacV1().Roles(namespace).Delete(ctx, name, opts)
}

func (c *Client) CreateRoleBinding(ctx context.Context, namespace string, roleBinding *rbacv1.RoleBinding, opts metav1.CreateOptions) (*rbacv1.RoleBinding, error) {
	return c.Clientset.RbacV1().RoleBindings(namespace).Create(ctx, roleBinding, opts)
}

func (c *Client) UpdateRoleBinding(ctx context.Context, namespace string, roleBinding *rbacv1.RoleBinding, opts metav1.UpdateOptions) (*rbacv1.RoleBinding, error) {
	return c.Clientset.RbacV1().RoleBindings(namespace).Update(ctx, roleBinding, opts)
}

func (c *Client) DeleteRoleBinding(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.RbacV1().RoleBindings(namespace).Delete(ctx, name, opts)
}

func (c *Client) GetConfigMap(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.ConfigMap, error) {
	return c.Clientset.CoreV1().ConfigMaps(namespace).Get(ctx, name, opts)
}

func (c *Client) PatchConfigMap(ctx context.Context, namespace, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions) (*corev1.ConfigMap, error) {
	return c.Clientset.CoreV1().ConfigMaps(namespace).Patch(ctx, name, pt, data, opts)
}

func (c *Client) UpdateConfigMap(ctx context.Context, configMap *corev1.ConfigMap, opts metav1.UpdateOptions) (*corev1.ConfigMap, error) {
	return c.Clientset.CoreV1().ConfigMaps(configMap.Namespace).Update(ctx, configMap, opts)
}

func (c *Client) CreateService(ctx context.Context, namespace string, service *corev1.Service, opts metav1.CreateOptions) (*corev1.Service, error) {
	return c.Clientset.CoreV1().Services(namespace).Create(ctx, service, opts)
}

func (c *Client) DeleteService(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.CoreV1().Services(namespace).Delete(ctx, name, opts)
}

func (c *Client) GetService(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.Service, error) {
	return c.Clientset.CoreV1().Services(namespace).Get(ctx, name, opts)
}

func (c *Client) CreateEndpoints(ctx context.Context, namespace string, ep *corev1.Endpoints, opts metav1.CreateOptions) (*corev1.Endpoints, error) {
	return c.Clientset.CoreV1().Endpoints(namespace).Create(ctx, ep, opts)
}

func (c *Client) GetEndpoints(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.Endpoints, error) {
	return c.Clientset.CoreV1().Endpoints(namespace).Get(ctx, name, opts)
}

func (c *Client) DeleteEndpoints(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.CoreV1().Endpoints(namespace).Delete(ctx, name, opts)
}

func (c *Client) CreateDeployment(ctx context.Context, namespace string, deployment *appsv1.Deployment, opts metav1.CreateOptions) (*appsv1.Deployment, error) {
	return c.Clientset.AppsV1().Deployments(namespace).Create(ctx, deployment, opts)
}

func (c *Client) GetDeployment(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*appsv1.Deployment, error) {
	return c.Clientset.AppsV1().Deployments(namespace).Get(ctx, name, opts)
}

func (c *Client) DeleteDeployment(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.AppsV1().Deployments(namespace).Delete(ctx, name, opts)
}

func (c *Client) PatchDeployment(ctx context.Context, namespace, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions) (*appsv1.Deployment, error) {
	return c.Clientset.AppsV1().Deployments(namespace).Patch(ctx, name, pt, data, opts)
}

func (c *Client) CheckDeploymentStatus(ctx context.Context, namespace, deployment string) error {
	d, err := c.GetDeployment(ctx, namespace, deployment, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if d == nil {
		return fmt.Errorf("deployment is not available")
	}

	if d.Status.ObservedGeneration != d.Generation {
		return fmt.Errorf("observed generation (%d) is older than generation of the desired state (%d)",
			d.Status.ObservedGeneration, d.Generation)
	}

	if d.Status.Replicas == 0 {
		return fmt.Errorf("replicas count is zero")
	}

	if d.Status.AvailableReplicas != d.Status.Replicas {
		return fmt.Errorf("only %d of %d replicas are available", d.Status.AvailableReplicas, d.Status.Replicas)
	}

	if d.Status.ReadyReplicas != d.Status.Replicas {
		return fmt.Errorf("only %d of %d replicas are ready", d.Status.ReadyReplicas, d.Status.Replicas)
	}

	if d.Status.UpdatedReplicas != d.Status.Replicas {
		return fmt.Errorf("only %d of %d replicas are up-to-date", d.Status.UpdatedReplicas, d.Status.Replicas)
	}

	return nil
}

func (c *Client) ListDeployment(ctx context.Context, namespace string, o metav1.ListOptions) (*appsv1.DeploymentList, error) {
	return c.Clientset.AppsV1().Deployments(namespace).List(ctx, o)
}

func (c *Client) CreateNamespace(ctx context.Context, namespace string, opts metav1.CreateOptions) (*corev1.Namespace, error) {
	return c.Clientset.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}, opts)
}

func (c *Client) GetNamespace(ctx context.Context, namespace string, options metav1.GetOptions) (*corev1.Namespace, error) {
	return c.Clientset.CoreV1().Namespaces().Get(ctx, namespace, options)
}

func (c *Client) DeleteNamespace(ctx context.Context, namespace string, opts metav1.DeleteOptions) error {
	return c.Clientset.CoreV1().Namespaces().Delete(ctx, namespace, opts)
}

func (c *Client) GetPod(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.Pod, error) {
	return c.Clientset.CoreV1().Pods(namespace).Get(ctx, name, opts)
}

func (c *Client) CreatePod(ctx context.Context, namespace string, pod *corev1.Pod, opts metav1.CreateOptions) (*corev1.Pod, error) {
	return c.Clientset.CoreV1().Pods(namespace).Create(ctx, pod, opts)
}

func (c *Client) DeletePod(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.CoreV1().Pods(namespace).Delete(ctx, name, opts)
}

func (c *Client) DeletePodCollection(ctx context.Context, namespace string, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return c.Clientset.CoreV1().Pods(namespace).DeleteCollection(ctx, opts, listOpts)
}

func (c *Client) ListPods(ctx context.Context, namespace string, options metav1.ListOptions) (*corev1.PodList, error) {
	return c.Clientset.CoreV1().Pods(namespace).List(ctx, options)
}

func (c *Client) PodLogs(namespace, name string, opts *corev1.PodLogOptions) *rest.Request {
	return c.Clientset.CoreV1().Pods(namespace).GetLogs(name, opts)
}

func (c *Client) ListServices(ctx context.Context, namespace string, options metav1.ListOptions) (*corev1.ServiceList, error) {
	return c.Clientset.CoreV1().Services(namespace).List(ctx, options)
}

func (c *Client) ExecInPodWithStderr(ctx context.Context, namespace, pod, container string, command []string) (bytes.Buffer, bytes.Buffer, error) {
	result, err := c.execInPod(ctx, ExecParameters{
		Namespace: namespace,
		Pod:       pod,
		Container: container,
		Command:   command,
	})
	return result.Stdout, result.Stderr, err
}

func (c *Client) ExecInPod(ctx context.Context, namespace, pod, container string, command []string) (bytes.Buffer, error) {
	result, err := c.execInPod(ctx, ExecParameters{
		Namespace: namespace,
		Pod:       pod,
		Container: container,
		Command:   command,
	})
	if err != nil {
		return bytes.Buffer{}, err
	}

	if errString := result.Stderr.String(); errString != "" {
		return bytes.Buffer{}, fmt.Errorf("command failed: %s", errString)
	}

	return result.Stdout, nil
}

func (c *Client) ExecInPodWithWriters(connCtx, killCmdCtx context.Context, namespace, pod, container string, command []string, stdout, stderr io.Writer) error {
	execParams := ExecParameters{
		Namespace: namespace,
		Pod:       pod,
		Container: container,
		Command:   command,
	}
	if killCmdCtx != nil {
		execParams.TTY = true
	}
	err := c.execInPodWithWriters(connCtx, killCmdCtx, execParams, stdout, stderr)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) CreateConfigMap(ctx context.Context, namespace string, config *corev1.ConfigMap, opts metav1.CreateOptions) (*corev1.ConfigMap, error) {
	return c.Clientset.CoreV1().ConfigMaps(namespace).Create(ctx, config, opts)
}

func (c *Client) DeleteConfigMap(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.CoreV1().ConfigMaps(namespace).Delete(ctx, name, opts)
}

func (c *Client) CreateDaemonSet(ctx context.Context, namespace string, ds *appsv1.DaemonSet, opts metav1.CreateOptions) (*appsv1.DaemonSet, error) {
	return c.Clientset.AppsV1().DaemonSets(namespace).Create(ctx, ds, opts)
}

func (c *Client) PatchDaemonSet(ctx context.Context, namespace, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions) (*appsv1.DaemonSet, error) {
	return c.Clientset.AppsV1().DaemonSets(namespace).Patch(ctx, name, pt, data, opts)
}

func (c *Client) GetDaemonSet(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*appsv1.DaemonSet, error) {
	return c.Clientset.AppsV1().DaemonSets(namespace).Get(ctx, name, opts)
}

func (c *Client) ListDaemonSet(ctx context.Context, namespace string, o metav1.ListOptions) (*appsv1.DaemonSetList, error) {
	return c.Clientset.AppsV1().DaemonSets(namespace).List(ctx, o)
}

func (c *Client) DeleteDaemonSet(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.AppsV1().DaemonSets(namespace).Delete(ctx, name, opts)
}

func (c *Client) GetCRD(ctx context.Context, name string, opts metav1.GetOptions) (*apiextensions.CustomResourceDefinition, error) {
	return c.ExtensionClientset.ApiextensionsV1().CustomResourceDefinitions().Get(ctx, name, opts)
}

// Kubernetes Network Policies specific commands

func (c *Client) ListKubernetesNetworkPolicies(ctx context.Context, namespace string, opts metav1.ListOptions) (*networkingv1.NetworkPolicyList, error) {
	return c.Clientset.NetworkingV1().NetworkPolicies(namespace).List(ctx, opts)
}

func (c *Client) GetKubernetesNetworkPolicy(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*networkingv1.NetworkPolicy, error) {
	return c.Clientset.NetworkingV1().NetworkPolicies(namespace).Get(ctx, name, opts)
}

func (c *Client) CreateKubernetesNetworkPolicy(ctx context.Context, policy *networkingv1.NetworkPolicy, opts metav1.CreateOptions) (*networkingv1.NetworkPolicy, error) {
	return c.Clientset.NetworkingV1().NetworkPolicies(policy.Namespace).Create(ctx, policy, opts)
}

func (c *Client) UpdateKubernetesNetworkPolicy(ctx context.Context, policy *networkingv1.NetworkPolicy, opts metav1.UpdateOptions) (*networkingv1.NetworkPolicy, error) {
	return c.Clientset.NetworkingV1().NetworkPolicies(policy.Namespace).Update(ctx, policy, opts)
}

func (c *Client) DeleteKubernetesNetworkPolicy(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	return c.Clientset.NetworkingV1().NetworkPolicies(namespace).Delete(ctx, name, opts)
}
