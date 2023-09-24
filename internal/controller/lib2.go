package controller

import (
	"context"
	webappresv1 "demoresources/api/v1"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *DbcrReconciler) reconcileSecret(ctx context.Context, parentResource *webappresv1.Dbcr, l logr.Logger) (corev1.Secret, error) {
	var sec = &corev1.Secret{}
	secName := parentResource.Name + "-sec"
	l.Info("secret ", "name:", secName, "namespace:", parentResource.Namespace)
	findErr := r.Get(ctx, types.NamespacedName{Name: secName, Namespace: parentResource.Namespace}, sec)
	if findErr == nil {
		l.Info("Secret Found")
		return *sec, nil
	}

	if !errors.IsNotFound(findErr) {
		return *sec, findErr
	}

	l.Info("Secret Not found, Creating new secret")
	sec = &corev1.Secret{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secName,
			Namespace: parentResource.Namespace,
		},
		Type: "Opaque",
		Data: map[string][]byte{
			"rootpass": []byte("cm9vdA=="),
			"dbpass":   []byte("ZGJAMTI3"),
		},
	}

	l.Info("Creating Secret...", "Secret name", sec.Name, "Secret namespace", sec.Namespace)

	if err := ctrl.SetControllerReference(parentResource, sec, r.Scheme); err != nil {
		return *sec, err
	}
	secErr := r.Create(ctx, sec)
	return *sec, secErr
	//return *sec, nil
}

func (r *DbcrReconciler) reconcileDBDeployment(ctx context.Context, parentResource *webappresv1.Dbcr, l logr.Logger) (appsv1.StatefulSet, error) {
	var replicaNum int32 = 1
	dep := &appsv1.StatefulSet{}
	dbResName := parentResource.Name + "-db"
	//secretName := parentResource.Name + "-sec"
	err := r.Get(ctx, types.NamespacedName{Name: dbResName, Namespace: parentResource.Namespace}, dep)
	if err == nil {
		l.Info("deployment resource Found")
		return *dep, nil
	}

	labels := map[string]string{
		"app": dbResName,
	}
	dep = &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dbResName,
			Namespace: parentResource.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replicaNum,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "mysql",
							Image: "mysql",
							Env: []corev1.EnvVar{
								{
									Name:  "MYSQL_USER",
									Value: "dbshi",
								},
								{
									Name:  "MYSQL_PASSWORD",
									Value: "db@127",
								},
								{
									Name:  "MYSQL_ROOT_PASSWORD",
									Value: "&45hhdgKjFG",
								},

								//{
								//	Name: "MYSQL_PASSWORD",
								//	ValueFrom: &corev1.EnvVarSource{
								//		SecretKeyRef: &corev1.SecretKeySelector{
								//			LocalObjectReference: corev1.LocalObjectReference{Name: secretName},
								//			Key:                  "db",
								//		},
								//	},
								//},
								//{
								//	Name: "MYSQL_ROOT_PASSWORD",
								//	ValueFrom: &corev1.EnvVarSource{
								//		SecretKeyRef: &corev1.SecretKeySelector{
								//			LocalObjectReference: corev1.LocalObjectReference{Name: secretName},
								//			Key:                  "rootpass",
								//		},
								//	},
								//},
								{
									Name:  "MYSQL_DATABASE",
									Value: "libdb",
								},
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 3306,
								},
							},
						},
					},
				},
			},
		},
	}
	l.Info("Creating DB...", "DB name", dep.Name, "DB namespace", dep.Namespace)
	if err := ctrl.SetControllerReference(parentResource, dep, r.Scheme); err != nil {
		return *dep, err
	}
	errDep := r.Create(ctx, dep)
	return *dep, errDep
}
