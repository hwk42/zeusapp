/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	nativeaidevv1 "github.com/hwk42/zeusapp/api/v1"
)

// ZeusappReconciler reconciles a Zeusapp object
type ZeusappReconciler struct {
	client.Client
	Log      logr.Logger
	Recorder record.EventRecorder
	Scheme   *runtime.Scheme
}

//+kubebuilder:rbac:groups=nativeai.dev,resources=zeusapps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=nativeai.dev,resources=zeusapps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=nativeai.dev,resources=zeusapps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Zeusapp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *ZeusappReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	logger := r.Log.WithValues("zeusapp", req.NamespacedName)

	instance := &nativeaidevv1.Zeusapp{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Reconcile k8s deployment.
	deployment, err := generateDeployment(instance, logger, r)
	if err != nil {
		return ctrl.Result{}, err
	}
	if err := ctrl.SetControllerReference(instance, deployment, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	if err := ReconcileDeployment(ctx, r.Client, deployment, logger); err != nil {
		return ctrl.Result{}, err
	}

	// Reconcile k8s service.
	service := generateService(instance)
	if err := ctrl.SetControllerReference(instance, service, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	if err := ReconcileService(ctx, r.Client, service, logger); err != nil {
		return ctrl.Result{}, err
	}

	foundDeployment := &appsv1.Deployment{}
	//Update the instance.Status.ReadyReplicas if the foundDeployment.Status.ReadyReplicas
	//has changed.
	_err := r.Get(ctx, types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, foundDeployment)

	if _err != nil {
		if errors.IsNotFound(_err) {
			logger.Info("Deployment not found...", "deployment", deployment.Name)
		} else {
			return ctrl.Result{}, _err
		}
	} else {
		if foundDeployment.Status.ReadyReplicas != instance.Status.ReadyReplicas {
			logger.Info("Updating Status", "namespace", instance.Namespace, "name", instance.Name)
			instance.Status.ReadyReplicas = foundDeployment.Status.ReadyReplicas
		}

		_err = r.Status().Update(ctx, instance)
		if _err != nil {
			return ctrl.Result{}, _err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ZeusappReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&nativeaidevv1.Zeusapp{}).
		Complete(r)
}
