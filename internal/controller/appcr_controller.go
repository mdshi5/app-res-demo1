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

package controller

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	webappresv1 "demoresources/api/v1"
)

// AppcrReconciler reconciles a Appcr object
type AppcrReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=webappres.shi.io,resources=appcrs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=webappres.shi.io,resources=appcrs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=webappres.shi.io,resources=appcrs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Appcr object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *AppcrReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	l.Info("app reconciliation")
	appCR := &webappresv1.Appcr{}
	r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, appCR)

	appDep, appDepErr := r.reconcileappDeployment(ctx, appCR, l)
	if appDepErr != nil {
		l.Info("error creating app deployment")
		return ctrl.Result{}, appDepErr
	}
	l.Info("created app deployment", "name:", appDep.Name, "namespace", appDep.Namespace)

	appSvc, appSvcErr := r.reconcileAppSvc(ctx, appCR, l)
	if appSvcErr != nil {
		l.Info("error in creating app svc")
		return ctrl.Result{}, appSvcErr
	}
	l.Info("created app deployment", "name:", appSvc.Name, "namespace", appSvc.Namespace)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AppcrReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappresv1.Appcr{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
