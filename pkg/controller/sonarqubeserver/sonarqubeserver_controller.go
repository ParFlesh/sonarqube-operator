package sonarqubeserver

import (
	"context"
	"github.com/operator-framework/operator-sdk/pkg/status"
	appsv1 "k8s.io/api/apps/v1"
	"strings"

	sonarsourcev1alpha1 "github.com/parflesh/sonarqube-operator/pkg/apis/sonarsource/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_sonarqubeserver")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new SonarQubeServer Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileSonarQubeServer{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("sonarqubeserver-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource SonarQubeServer
	err = c.Watch(&source.Kind{Type: &sonarsourcev1alpha1.SonarQubeServer{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Deployment and requeue the owner SonarQube
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &sonarsourcev1alpha1.SonarQube{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Service and requeue the owner SonarQube
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &sonarsourcev1alpha1.SonarQube{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Secret and requeue the owner SonarQube
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &sonarsourcev1alpha1.SonarQube{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource PersistentVolumeClaim and requeue the owner SonarQube
	err = c.Watch(&source.Kind{Type: &corev1.PersistentVolumeClaim{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &sonarsourcev1alpha1.SonarQube{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Secret and requeue the watcher
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestsFromMapFunc{
		ToRequests: secretHandlerFunc,
	})
	if err != nil {
		return err
	}

	return nil
}

type secretMapper func(handler.MapObject) []reconcile.Request

func (r secretMapper) Map(o handler.MapObject) []reconcile.Request {
	return r(o)
}

var secretHandlerFunc secretMapper = func(o handler.MapObject) []reconcile.Request {
	var output []reconcile.Request
	for k, v := range o.Meta.GetAnnotations() {
		if k == sonarsourcev1alpha1.ServerSecretAnnotation {
			for _, e := range strings.Split(v, ",") {
				output = append(output, reconcile.Request{
					NamespacedName: types.NamespacedName{
						Namespace: o.Meta.GetNamespace(),
						Name:      e,
					},
				})
			}
		}
	}
	return output
}

// blank assignment to verify that ReconcileSonarQubeServer implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileSonarQubeServer{}

// ReconcileSonarQubeServer reconciles a SonarQubeServer object
type ReconcileSonarQubeServer struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a SonarQubeServer object and makes changes based on the state read
// and what is in the SonarQubeServer.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileSonarQubeServer) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling SonarQubeServer")

	// Fetch the SonarQubeServer instance
	instance := &sonarsourcev1alpha1.SonarQubeServer{}
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

	newStatus := &sonarsourcev1alpha1.SonarQubeServerStatus{}
	*newStatus = instance.Status
	if newStatus.Deployment == nil {
		newStatus.Deployment = make(sonarsourcev1alpha1.DeploymentStatus)
	}
	r.updateStatus(newStatus, instance)

	_, err = r.ReconcileSecret(instance)
	if err != nil {
		return r.ParseErrorForReconcileResult(instance, err)
	}

	_, err = r.ReconcileServiceAccount(instance)
	if err != nil {
		return r.ParseErrorForReconcileResult(instance, err)
	}

	_, err = r.ReconcileService(instance)
	if err != nil {
		return r.ParseErrorForReconcileResult(instance, err)
	}

	_, err = r.ReconcileDeployment(instance)
	if err != nil {
		return r.ParseErrorForReconcileResult(instance, err)
	}

	*newStatus = instance.Status

	if newStatus.Conditions.IsTrueFor(sonarsourcev1alpha1.ConditionInvalid) {
		newStatus.Conditions.SetCondition(status.Condition{
			Type:   sonarsourcev1alpha1.ConditionInvalid,
			Status: corev1.ConditionFalse,
		})
	}

	if newStatus.Conditions.IsTrueFor(sonarsourcev1alpha1.ConditionProgressing) {
		newStatus.Conditions.SetCondition(status.Condition{
			Type:   sonarsourcev1alpha1.ConditionProgressing,
			Status: corev1.ConditionFalse,
		})
	}

	if newStatus.Conditions.IsTrueFor(sonarsourcev1alpha1.ConditionPending) {
		newStatus.Conditions.SetCondition(status.Condition{
			Type:   sonarsourcev1alpha1.ConditionPending,
			Status: corev1.ConditionFalse,
		})
	}

	r.updateStatus(newStatus, instance)

	return reconcile.Result{}, nil
}