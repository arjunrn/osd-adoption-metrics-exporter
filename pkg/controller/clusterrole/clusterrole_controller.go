package clusterrole

import (
	"context"

	"gitlab.cee.redhat.com/service/osd-adoption-metrics-exporter/pkg/controller/utils"
	"gitlab.cee.redhat.com/service/osd-adoption-metrics-exporter/pkg/metrics"
	rbacv1 "k8s.io/api/rbac/v1"
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

var log = logf.Log.WithName("controller_cluster_role")

const (
	finalizer = "adoption-exporter/finalizer"
)

// Add creates a new Cluster Role Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileClusterRole{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("cluster-role-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ClusterRole
	err = c.Watch(&source.Kind{Type: &rbacv1.ClusterRole{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileClusterRole implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileClusterRole{}

// ReconcileClusterRole reconciles a ClusterRole object
type ReconcileClusterRole struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a ClusterRole object and makes changes based on the state read
// and what is in the ClusterRole.Spec
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileClusterRole) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ClusterRole")

	// Fetch the ClusterRole instance
	instance := &rbacv1.ClusterRole{}
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

	if instance.Name != "cluster-admin" {
		return reconcile.Result{}, nil
	}

	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		if !utils.ContainsString(instance.ObjectMeta.Finalizers, finalizer) {
			instance.ObjectMeta.Finalizers = append(instance.ObjectMeta.Finalizers, finalizer)
			if err = r.client.Update(context.Background(), instance); err != nil {
				return reconcile.Result{}, err
			}
		}
		metrics.Aggregator.SetClusterAdmin(true)
	} else {
		if utils.ContainsString(instance.ObjectMeta.Finalizers, finalizer) {
			instance.ObjectMeta.Finalizers = utils.RemoveString(instance.ObjectMeta.Finalizers, finalizer)
			if err = r.client.Update(context.Background(), instance); err != nil {
				return reconcile.Result{}, err
			}
		}
		metrics.Aggregator.DeleteAuthIDP(instance.Name, instance.Namespace)
		metrics.Aggregator.SetClusterAdmin(false)
	}

	return reconcile.Result{}, nil
}
