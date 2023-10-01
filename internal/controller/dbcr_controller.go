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

// DbcrReconciler reconciles a Dbcr object
type DbcrReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=webappres.shi.io,resources=dbcrs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=webappres.shi.io,resources=dbcrs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=webappres.shi.io,resources=dbcrs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Dbcr object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *DbcrReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	l.Info("db reconciliation")
	dbCR := &webappresv1.Dbcr{}
	r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, dbCR)
	secDB, secDBErr := r.reconcileSecret(ctx, dbCR, l)
	if secDBErr != nil {
		l.Info("error in db secret creation")
		return ctrl.Result{}, secDBErr
	}
	l.Info("created db secret", "name:", secDB.Name, "namespace", secDB.Namespace)

	depDB, depDBErr := r.reconcileDBDeployment(ctx, dbCR, l)
	if depDBErr != nil {
		l.Info("error in db deployment creation")
		return ctrl.Result{}, depDBErr
	} else {
		dbCR.Status.IsDBReady = "init"
		if err := r.Client.Status().Update(context.TODO(), dbCR); err != nil {
			return ctrl.Result{}, err
		}

	}
	l.Info("created db deployment", "name:", depDB.Name, "namespace", depDB.Namespace)

	svcDB, svcDBErr := r.reconcileDBSvc(ctx, dbCR, l)
	if svcDBErr != nil {
		l.Info("error in db svc creation")
		return ctrl.Result{}, svcDBErr
	}
	l.Info("created db deployment", "name:", svcDB.Name, "namespace", svcDB.Namespace)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DbcrReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappresv1.Dbcr{}).
		Owns(&corev1.Secret{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
