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
	"fmt"
	api "github.com/shahriarKhanpiyal/custom-resource/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CustomResourceReconciler reconciles a CustomResource object
type CustomResourceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=piyal.dev,resources=customresources,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=piyal.dev,resources=customresources/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=piyal.dev,resources=customresources/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CustomResource object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *CustomResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logs := log.FromContext(ctx)
	logs.WithValues("ReqName", req.Name, "ReqNamespace", req.Namespace)

	/*
		### 1: Load the CustomResource by name

		We'll fetch the CustomResource using our client.  All client methods take a
		context (to allow for cancellation) as their first argument, and the object
		in question as their last.  Get is a bit special, in that it takes a
		[`NamespacedName`](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/client?tab=doc#ObjectKey)
		as the middle argument (most don't have a middle argument, as we'll see
		below).

		Many client methods also take variadic options at the end.
	*/

	var crd api.CustomResource

	if err := r.Get(ctx, req.NamespacedName, &crd); err != nil {
		klog.Info(err, "unable to fetch crd")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, nil
	}
	klog.Info("Name:", crd.Name)

	// deploymentObject carry the all data of deployment in specific namespace and name
	var deploymentObject appsv1.Deployment

	// Naming of deployment
	deploymentName := crd.Name + "-" + crd.Spec.DeploymentName
	if crd.Spec.DeploymentName == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		//utilruntime.HandleError(fmt.Errorf("%s : deployment name must be specified", key))

		deploymentName = crd.Name + "-" + "dep"
		//return nil
	}

	//Creating NamespacedName for deploymentObject
	objectkey := client.ObjectKey{
		Namespace: req.Namespace,
		Name:      deploymentName,
	}

	if err := r.Get(ctx, objectkey, &deploymentObject); err != nil {
		if errors.IsNotFound(err) {
			klog.Info("Couldn't find existing Deployment for", crd.Name, " creating one ...")
			err := r.Client.Create(ctx, newDeployment(&crd, deploymentName))
			if err != nil {
				klog.Info("error while creating deployment %s", err)
				return ctrl.Result{}, err
			} else {
				klog.Info("%s Deployment Created ", crd.Name)
			}
		} else {
			klog.Info("error fetching deployment %s", err)
			return ctrl.Result{}, err
		}
	} else {
		if crd.Spec.Replicas != nil && *crd.Spec.Replicas != *deploymentObject.Spec.Replicas {
			klog.Info(*crd.Spec.Replicas, *deploymentObject.Spec.Replicas)
			klog.Info("Deployment replica don't match ... updating")
			// As the replica count didn't match the, we need to update it
			deploymentObject.Spec.Replicas = crd.Spec.Replicas
			if err := r.Update(ctx, &deploymentObject); err != nil {
				klog.Info("error updating deployment %s", err)
			}
			klog.Info("deployment updated")
		}
	}

	// service reconcile
	var serviceObject corev1.Service
	// service Name
	serviceName := crd.Name + "-" + crd.Spec.Service.ServiceName + "-service"
	if crd.Spec.Service.ServiceName == "" {
		serviceName = crd.Name + "-" + "-service"
	}

	objectkey = client.ObjectKey{
		Namespace: req.Namespace,
		Name:      serviceName,
	}

	// create or update service
	if err := r.Get(ctx, objectkey, &serviceObject); err != nil {
		if errors.IsNotFound(err) {
			klog.Info("Could not find existing Service for", crd.Name, "creating one ...")
			err := r.Create(ctx, newService(&crd, serviceName))
			if err != nil {
				klog.Info("error while creating service %s", err)
				return ctrl.Result{}, err
			} else {
				klog.Info("%s Service Created ", crd.Name)
			}
		} else {
			fmt.Printf("error fetching service %s", err)
			return ctrl.Result{}, err
		}
	} else {
		if crd.Spec.Replicas != nil && *crd.Spec.Replicas != crd.Status.AvailableReplicas {
			klog.Info("Is this problem ?")
			klog.Info("Service replica miss match .... updating ")
			crd.Status.AvailableReplicas = *crd.Spec.Replicas

			if err := r.Status().Update(ctx, &crd); err != nil {
				klog.Info("error while updating service %s", err)
				return ctrl.Result{}, err
			}
			klog.Info("service updated")

		}
	}
	klog.Info("reconcile done")

	return ctrl.Result{
		Requeue:      true,
		RequeueAfter: 1 * time.Minute,
	}, nil

}
