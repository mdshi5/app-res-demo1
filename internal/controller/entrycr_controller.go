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
	"k8s.io/apimachinery/pkg/types"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	webappresv1 "demoresources/api/v1"
)

// EntrycrReconciler reconciles a Entrycr object
type EntrycrReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=webappres.shi.io,resources=entrycrs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=webappres.shi.io,resources=entrycrs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=webappres.shi.io,resources=entrycrs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Entrycr object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *EntrycrReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	entrycr := &webappresv1.Entrycr{}
	errfind := r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, entrycr)
	if errfind == nil {
		l.Info("resource EntryCR found")
	}
	l.Info("Created ")
	dbCR, errDB := r.reconcileMyDB(ctx, entrycr, l)
	if errDB != nil {
		l.Info("error in db creation")
		return ctrl.Result{}, errDB
	}
	l.Info("created dbCR", "name:", dbCR.Name, "namespace", dbCR.Namespace)

	appCR, errApp := r.reconcileMyApp(ctx, entrycr, l)
	if errApp != nil {
		l.Info("error in app creation")
		return ctrl.Result{}, errApp
	}
	l.Info("created app", "name:", appCR.Name, "namespace", appCR.Namespace)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EntrycrReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappresv1.Entrycr{}).
		Owns(&webappresv1.Dbcr{}).
		Complete(r)
}
