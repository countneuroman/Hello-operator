// TODO: Вынести в отдельный package controller
package main

import (
	"context"
	"errors"
	"time"

	corev1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	kubeinformers "k8s.io/client-go/informers"
	echoInformers "github.com/countneuroman/hello-operator/pkg/generated/informers/externalversions"


	echov1alpha1 "github.com/countneuroman/hello-operator/pkg/apis/hellocontroller/v1alpha1"
	echoclientset "github.com/countneuroman/hello-operator/pkg/generated/clientset/versioned"
	samplescheme "github.com/countneuroman/hello-operator/pkg/generated/clientset/versioned/scheme"
)

const controllerAgentName = "hello-controller"

type Controller struct {
	kubeclientset kubernetes.Interface

	echoclientset echoclientset.Interface

	//Используется информер, чтобы получать состояние сущностей, не обращаясь постоянно к kube-api
	echoInformer cache.SharedIndexInformer
	jobInformer cache.SharedIndexInformer

	namespace string

	//Рейт лимиты ставят в очередь все элементы для обработки, чтобы обеспечить обработку только единоразово
	//для каждого изменения, в случае если есть несколько воркеров.
	queue workqueue.TypedRateLimitingInterface[event]

	//Позволяет записывать эвенты для наших API объектов
	recorder record.EventRecorder

	//TODO: Брать логгер через klog.FromContext?
	logger klog.Logger
}

func (c *Controller) Run(ctx context.Context, numWorkers int) error {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	c.logger.Info("Starting controller")

	c.logger.Info("Starting informers")
	for _, i := range []cache.SharedIndexInformer{
		c.echoInformer,
	} {
		go i.Run(ctx.Done())
	}

	c.logger.Info("Waiting for informer caches to sync")
	if !cache.WaitForCacheSync(ctx.Done(), []cache.InformerSynced{
		c.echoInformer.HasSynced,
	}... ) {
		err := errors.New("Failed to wait for caches to sync")
		utilruntime.HandleError(err)
		return err
	}

	c.logger.Info("starting %d workers", numWorkers)
	for i := 0; i < numWorkers; i++ {
		go wait.Until(func() {
			c.runWorker(ctx)
		}, time.Second, ctx.Done())
	}
	c.logger.Info("Controller ready")

	<-ctx.Done()
	c.logger.Info("Controller stopping")

	return nil
}

func (c * Controller) addEcho(obj interface{}) {
	c.logger.Info("Adding echo")
	echo, ok := obj.(*echov1alpha1.Echo)
	if !ok {
		c.logger.Error(nil, "unexpected object")
		return
	}
	//TODO: По факту мы не юзаем кэш информера? Т.к оперируем объектом напрямую.
	c.queue.Add(event{
		eventType: addEcho,
		newObj:    echo.DeepCopy(),
	})
}

func NewController(
	ctx context.Context,
	kubeClientSet kubernetes.Interface,
	echoClientSet echoclientset.Interface,
	namespace string,
) *Controller {
	logger := klog.FromContext(ctx)

	//Добавляем схему нашего контроллера к схеме Куба по умолчанию, чтобы эвенты могли логгироваться для наших  типов.
	utilruntime.Must(samplescheme.AddToScheme(scheme.Scheme))

	eventBroadcaster := record.NewBroadcaster(record.WithContext(ctx))
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClientSet.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	kubeInfromerFactory := kubeinformers.NewSharedInformerFactory(kubeClientSet, time.Second*30)
	echoInformerFactory := echoInformers.NewSharedInformerFactory(echoClientSet, time.Second*30)

	echoInformer := echoInformerFactory.Hello().V1alpha1().Echos().Informer()
	jobInformer := kubeInfromerFactory.Batch().V1().Jobs().Informer()
	queue := workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[event]())

	controller := &Controller {
		kubeclientset: kubeClientSet,
		echoclientset: echoClientSet,

		namespace: namespace,

		echoInformer: echoInformer,
		jobInformer: jobInformer,

		queue: queue,

		recorder: recorder,

		logger: logger,
	}

	logger.Info("Set up event handlers")

	echoInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.addEcho,
	})

	return controller
}
