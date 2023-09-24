package controller

import (
	"context"
	webappresv1 "demoresources/api/v1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (r *EntrycrReconciler) reconcileMyDB(ctx context.Context, entrycr *webappresv1.Entrycr, l logr.Logger) (webappresv1.Dbcr, error) {
	resName := entrycr.Name + "db"
	dbCR := &webappresv1.Dbcr{}
	findErr := r.Get(ctx, types.NamespacedName{Name: resName, Namespace: entrycr.Namespace}, dbCR)
	if findErr == nil {
		l.Info("DB resource found")
		return *dbCR, nil
	}
	if !errors.IsNotFound(findErr) {
		return *dbCR, findErr
	}
	labels := map[string]string{
		"app": resName,
	}
	dbCR = &webappresv1.Dbcr{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resName,
			Namespace: entrycr.Namespace,
			Labels:    labels,
		},
		//Spec: webappv1.DBSpec{
		//	Replicanum: entry.Spec.Dbreplica,
		//	Imgname:    entry.Spec.Dbimage,
		//},
	}
	err := r.Create(ctx, dbCR)
	return *dbCR, err

}

func (r *EntrycrReconciler) reconcileMyApp(ctx context.Context, entrycr *webappresv1.Entrycr, l logr.Logger) (webappresv1.Appcr, error) {
	resName := entrycr.Name + "app"
	appCR := &webappresv1.Appcr{}
	findErr := r.Get(ctx, types.NamespacedName{Name: resName, Namespace: entrycr.Namespace}, appCR)
	if findErr == nil {
		l.Info("App resource found")
		return *appCR, nil
	}
	if !errors.IsNotFound(findErr) {
		return *appCR, findErr
	}
	labels := map[string]string{
		"app": resName,
	}
	appCR = &webappresv1.Appcr{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resName,
			Namespace: entrycr.Namespace,
			Labels:    labels,
		},
		//Spec: webappv1.DBSpec{
		//	Replicanum: entry.Spec.Dbreplica,
		//	Imgname:    entry.Spec.Dbimage,
		//},
	}
	err := r.Create(ctx, appCR)
	return *appCR, err

}
