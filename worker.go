// TODO: Вынести в отдельный package controller
package main

import (
	"context"
	"fmt"

	echov1alpha1 "github.com/countneuroman/hello-operator/pkg/apis/hellocontroller/v1alpha1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Controller) runWorker(ctx context.Context) {
	for c.processNextItem(ctx) {
	}
}

func (c *Controller) processNextItem(ctx context.Context) bool {
	obj, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	defer c.queue.Done(obj)

	err := c.processEvent(ctx, obj)
	if err == nil {
		c.logger.Info("Processed item")
		c.queue.Forget(obj)
	} else {
		c.logger.Error(err, "error processing event")
		c.queue.Forget(obj)
		utilruntime.HandleError(err)
	}

	return true
}

func (c *Controller) processEvent(ctx context.Context, event event) error {
	switch event.eventType {
	case addEcho:
		return c.processEcho(ctx, event.newObj.(*echov1alpha1.Echo))
	}
	return nil
}

func (c *Controller) processEcho(ctx context.Context, echo *echov1alpha1.Echo) error {
	job := createJob(echo, c.namespace)
	exists, err := resurceExists(job, c.echoInformer.GetIndexer())
	if err != nil {
		return fmt.Errorf("error checking job existence %v", err)
	}
	if exists {
		c.logger.Info("job aledray exists, skipping")
		return nil
	}

	_, err = c.kubeclientset.BatchV1().
	Jobs(c.namespace).
	Create(ctx, job, metav1.CreateOptions{})
	return err
}

func resurceExists(obj interface{}, indexer cache.Indexer) (bool, error) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		return false, fmt.Errorf("error getting key %v", err)
	}
	_, exists, err := indexer.GetByKey(key)
	return exists, err
}