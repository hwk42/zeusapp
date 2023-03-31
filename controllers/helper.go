package controllers

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/go-logr/logr"
	nativeaidevv1 "github.com/hwk42/zeusapp/api/v1"
	"google.golang.org/protobuf/proto"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	zeusappFinalizer  = "zeusapp.nativeai.dev"
	Ascend310Resource = "huawei.com/Ascend310"
	NvidiaResource    = "nvidia.com/gpu"
)

func generateDeployment(app *nativeaidevv1.Zeusapp, log logr.Logger, r *ZeusappReconciler) (*appsv1.Deployment, error) {
	//var volumeMounts []corev1.VolumeMount
	//var volumes []corev1.Volume
	//var mountpath, subpath string = " ", ""
	var affinity = &corev1.Affinity{}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: proto.Int32(app.Spec.MinReplicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"zeusapp": app.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"zeusapp": app.Name},
				},
				Spec: corev1.PodSpec{
					//HostNetwork:   true,
					//DNSPolicy:     corev1.DNSClusterFirstWithHostNet,
					Affinity:      affinity,
					RestartPolicy: corev1.RestartPolicyAlways,
					Containers: []corev1.Container{
						{
							Name:            "zeusapp",
							Image:           app.Spec.Image,
							ImagePullPolicy: corev1.PullAlways,
							//Command:         app.Spec.Command,
							//WorkingDir: "/",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: app.Spec.ContainerPort,
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:      "VIE_POD_IP",
									ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "status.podIP"}},
								},
								{
									Name:      "NODE_IP",
									ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "status.hostIP"}},
								},

								{
									Name:  "HTTP_WEB_PORT",
									Value: strconv.FormatInt(int64(app.Spec.ContainerPort), 10),
								},
							},
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									NvidiaResource: resource.MustParse("1"),
								},
							},
							//VolumeMounts: volumeMounts,
						},
					},
					//Volumes: volumes,
				},
			},
		},
	}, nil
}

func generateDeploymentForAscend(app *nativeaidevv1.Zeusapp, log logr.Logger, r *ZeusappReconciler) (*appsv1.Deployment, error) {
	var volumeMounts []corev1.VolumeMount
	var volumes []corev1.Volume
	var affinity = &corev1.Affinity{}

	volumeMounts = append(volumeMounts,
		corev1.VolumeMount{
			Name:      "npu",
			MountPath: "/usr/local/bin/npu-smi"},
		corev1.VolumeMount{
			Name:      "dcmi",
			MountPath: "/usr/local/dcmi"},
		corev1.VolumeMount{
			Name:      "Ascend",
			MountPath: "/usr/local/Ascend"},
	)
	volumes = append(volumes,
		corev1.Volume{
			Name: "npu",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/usr/local/bin/npu-smi",
				},
			},
		},
		corev1.Volume{
			Name: "dcmi",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/usr/local/dcmi",
				},
			},
		},
		corev1.Volume{
			Name: "Ascend",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/usr/local/Ascend",
				},
			},
		},
	)

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: proto.Int32(app.Spec.MinReplicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"zeusapp": app.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"zeusapp": app.Name},
				},
				Spec: corev1.PodSpec{
					//HostNetwork:   true,
					//DNSPolicy:     corev1.DNSClusterFirstWithHostNet,
					Affinity:      affinity,
					RestartPolicy: corev1.RestartPolicyAlways,
					Containers: []corev1.Container{
						{
							Name:            "zeusapp",
							Image:           app.Spec.Image,
							ImagePullPolicy: corev1.PullAlways,
							//Command:         app.Spec.Command,
							//WorkingDir: "/",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: app.Spec.ContainerPort,
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:      "VIE_POD_IP",
									ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "status.podIP"}},
								},
								{
									Name:      "NODE_IP",
									ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "status.hostIP"}},
								},
								{
									Name:  "HTTP_WEB_PORT",
									Value: strconv.FormatInt(int64(app.Spec.ContainerPort), 10),
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									Ascend310Resource: resource.MustParse("1"),
								},
								Limits: corev1.ResourceList{
									Ascend310Resource: resource.MustParse("1"),
								},
							},
							VolumeMounts: volumeMounts,
						},
					},
					Volumes: volumes,
				},
			},
		},
	}, nil
}

func generateService(app *nativeaidevv1.Zeusapp) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Type:     "ClusterIP",
			Selector: map[string]string{"zeusapp": app.Name},
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Name:       "http-" + app.Name,
					Port:       80,
					TargetPort: intstr.FromInt(int(app.Spec.ContainerPort)),
				},
			},
		},
	}
}

func CopyDeploymentSetFields(from, to *appsv1.Deployment) bool {
	requireUpdate := false
	for k, v := range to.Labels {
		if from.Labels[k] != v {
			requireUpdate = true
		}
	}
	to.Labels = from.Labels

	for k, v := range to.Annotations {
		if from.Annotations[k] != v {
			requireUpdate = true
		}
	}
	to.Annotations = from.Annotations

	if from.Spec.Replicas != to.Spec.Replicas {
		to.Spec.Replicas = from.Spec.Replicas
		requireUpdate = true
	}

	if !reflect.DeepEqual(to.Spec.Template.Spec, from.Spec.Template.Spec) {
		requireUpdate = true
	}
	to.Spec.Template.Spec = from.Spec.Template.Spec

	return requireUpdate
}

// CopyServiceFields copies the owned fields from one Service to another
func CopyServiceFields(from, to *corev1.Service) bool {
	requireUpdate := false
	for k, v := range to.Labels {
		if from.Labels[k] != v {
			requireUpdate = true
		}
	}
	to.Labels = from.Labels

	for k, v := range to.Annotations {
		if from.Annotations[k] != v {
			requireUpdate = true
		}
	}
	to.Annotations = from.Annotations

	// Don't copy the entire Spec, because we can't overwrite the clusterIp field

	if !reflect.DeepEqual(to.Spec.Selector, from.Spec.Selector) {
		requireUpdate = true
	}
	to.Spec.Selector = from.Spec.Selector

	if !reflect.DeepEqual(to.Spec.Ports, from.Spec.Ports) {
		requireUpdate = true
	}
	to.Spec.Ports = from.Spec.Ports

	return requireUpdate
}

func getEnvVariable(envVar string) (string, error) {
	if lookupEnv, exists := os.LookupEnv(envVar); exists {
		return lookupEnv, nil
	} else {
		return "", fmt.Errorf("environment variable %v is not set", envVar)
	}
}
