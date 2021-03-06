package sonarqube

import (
	"context"
	"fmt"
	sonarsourcev1alpha1 "github.com/parflesh/sonarqube-operator/pkg/apis/sonarsource/v1alpha1"
	"github.com/parflesh/sonarqube-operator/pkg/utils"
	"github.com/parflesh/sonarqube-operator/version"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_sonarqube")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new SonarQube Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileSonarQube{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("sonarqube-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource SonarQube
	err = c.Watch(&source.Kind{Type: &sonarsourcev1alpha1.SonarQube{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource StatefulSet and requeue the owner SonarQube
	err = c.Watch(&source.Kind{Type: &sonarsourcev1alpha1.SonarQubeServer{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &sonarsourcev1alpha1.SonarQube{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource PersistentVolumeClaim and requeue the owner SonarQube
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &sonarsourcev1alpha1.SonarQube{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource PersistentVolumeClaim and requeue the owner SonarQube
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &sonarsourcev1alpha1.SonarQube{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource PersistentVolumeClaim and requeue the owner SonarQube
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestsFromMapFunc{
		ToRequests: &utils.SecretMapper{Annotation: sonarsourcev1alpha1.SecretAnnotation},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileSonarQube implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileSonarQube{}

// ReconcileSonarQube reconciles a SonarQube object
type ReconcileSonarQube struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a SonarQube object and makes changes based on the state read
// and what is in the SonarQube.Spec
func (r *ReconcileSonarQube) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling SonarQube")

	// Fetch the SonarQube instance
	instance := &sonarsourcev1alpha1.SonarQube{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	newStatus := instance.DeepCopy()
	if newStatus.Status.Deployments == nil {
		newStatus.Status.Deployments = make(sonarsourcev1alpha1.DeploymentStatuses)
	}
	if newStatus.Status.SearchDeployments == nil {
		newStatus.Status.SearchDeployments = make(sonarsourcev1alpha1.DeploymentStatuses)
	}
	utils.UpdateStatus(r.client, newStatus, instance)

	secret, err := r.ReconcileSecret(instance)
	if err != nil {
		return utils.ParseErrorForReconcileResult(r.client, instance, err)
	}

	revisionHash, err := utils.GenVersion(instance.Spec, secret.Data["sonar.properties"])

	newStatus = instance.DeepCopy()
	if revisionHash != newStatus.Status.Revision {
		newStatus.Status.Revision = revisionHash
		utils.UpdateStatus(r.client, newStatus, instance)
		return reconcile.Result{Requeue: true}, nil
	}

	_, err = r.ReconcileServiceAccount(instance)
	if err != nil {
		return utils.ParseErrorForReconcileResult(r.client, instance, err)
	}

	_, err = r.ReconcileService(instance)
	if err != nil {
		return utils.ParseErrorForReconcileResult(r.client, instance, err)
	}

	_, err = r.ReconcileSonarQubeServers(instance)
	if err != nil {
		return utils.ParseErrorForReconcileResult(r.client, instance, err)
	}

	newStatus = instance.DeepCopy()

	newStatus.Status.Conditions = utils.ClearConditions(newStatus.Status.Conditions)

	utils.UpdateStatus(r.client, newStatus, instance)

	return reconcile.Result{}, nil
}

func (r *ReconcileSonarQube) Labels(cr *sonarsourcev1alpha1.SonarQube) map[string]string {
	labels := make(map[string]string)

	labels[sonarsourcev1alpha1.TypeLabel] = cr.Name
	labels[sonarsourcev1alpha1.KubeAppName] = "SonarQube"
	labels[sonarsourcev1alpha1.KubeAppInstance] = cr.Name
	labels[sonarsourcev1alpha1.KubeAppVersion] = cr.Status.Revision
	labels[sonarsourcev1alpha1.KubeAppManagedby] = fmt.Sprintf("sonarqube-operator.v%s", version.Version)

	for k, v := range cr.Labels {
		labels[k] = v
	}

	return labels
}
