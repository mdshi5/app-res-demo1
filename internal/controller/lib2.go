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
	"k8s.io/apimachinery/pkg/util/intstr"
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
									Value: "root",
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

func (r *DbcrReconciler) reconcileDBSvc(ctx context.Context, parentResource *webappresv1.Dbcr, l logr.Logger) (corev1.Service, error) {
	resName := "mysql-db-service"
	dbDep := parentResource.Name + "-db"
	//resName :=parentResource.Name + "-dbsvc"
	svc := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{Name: resName, Namespace: parentResource.Namespace}, svc)
	if err == nil {
		l.Info("db svc Found")
		return *svc, err
	}

	l.Info("db svc not found, Creating new db svc")

	svc = &corev1.Service{

		ObjectMeta: metav1.ObjectMeta{
			Name:      resName,
			Namespace: parentResource.Namespace,
		},

		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": dbDep},

			Ports: []corev1.ServicePort{
				{
					Port:       3306,
					Name:       "svcport",
					Protocol:   "TCP",
					TargetPort: intstr.FromInt(3306),
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
	l.Info("Creating DB svc...", "DB svc name", svc.Name, "DB svc  namespace", svc.Namespace)
	if err := ctrl.SetControllerReference(parentResource, svc, r.Scheme); err != nil {
		return *svc, err
	}
	errSvc := r.Create(ctx, svc)
	return *svc, errSvc

}

func (r *AppcrReconciler) reconcileappDeployment(ctx context.Context, parentResource *webappresv1.Appcr, l logr.Logger) (appsv1.Deployment, error) {
	var replicaNum int32 = 1
	resName := parentResource.Name + "-appdep"
	dep := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: resName, Namespace: parentResource.Namespace}, dep)
	if err == nil {
		l.Info("app deployment resource Found")
		return *dep, err
	}

	if !errors.IsNotFound(err) {
		return *dep, err
	}

	labels := map[string]string{
		"app": resName,
	}

	dep = &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resName,
			Namespace: parentResource.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
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
							Name:  "app-dep",
							Image: "mdshihabuddin/locallibrary:2.0",
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 8080,
								},
							},
						},
					},
				},
			},
		},
	}
	l.Info("Creating app...", "app name", dep.Name, "app namespace", dep.Namespace)
	errAppcr := r.Create(ctx, dep)
	return *dep, errAppcr
}

func (r *AppcrReconciler) reconcileAppSvc(ctx context.Context, parentResource *webappresv1.Appcr, l logr.Logger) (corev1.Service, error) {
	resName := parentResource.Name + "appsvc"
	appSelected := parentResource.Name + "-appdep"
	svc := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{Name: resName, Namespace: parentResource.Namespace}, svc)
	if err == nil {
		l.Info("App svc Found")

	}

	l.Info("svc Not found, Creating new svc")

	svc = &corev1.Service{

		ObjectMeta: metav1.ObjectMeta{
			Name:      resName,
			Namespace: parentResource.Namespace,
		},

		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": appSelected},

			Ports: []corev1.ServicePort{
				{
					Protocol:   "TCP",
					Port:       8080,
					TargetPort: intstr.FromInt(8080),
					NodePort:   30100,
				},
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}
	l.Info("Creating app service...", "App service name", svc.Name, "App service namespace", svc.Namespace)

	if err := ctrl.SetControllerReference(parentResource, svc, r.Scheme); err != nil {
		return *svc, err
	}
	errDep := r.Create(ctx, svc)
	return *svc, errDep

}
