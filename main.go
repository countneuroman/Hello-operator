package main

import (
	"flag"
	"time"

	signals "github.com/countneuroman/hello-operator/pkg/signals"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	clientset "github.com/countneuroman/hello-operator/pkg/generated/clientset/versioned"
	echoInformers "github.com/countneuroman/hello-operator/pkg/generated/informers/externalversions"
)

var (
	masterURL      string
	kubeconfigPath string
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	ctx := signals.SetupSignalHandler()
	logger := klog.FromContext(ctx)

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
	if err != nil {
		logger.Error(err, "Error building kubeconfig.")
		klog.FlushAndExit(klog.ExitFlushTimeout, 1)
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		logger.Error(err, "Error building kubernetes client")
		klog.FlushAndExit(klog.ExitFlushTimeout, 1)
	}

	echoClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		logger.Error(err, "Error building example client")
		klog.FlushAndExit(klog.ExitFlushTimeout, 1)
	}

	kubeInfromerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	echoInformerFactory := echoInformers.NewSharedInformerFactory(echoClient, time.Second*30)

	//TODO: Получать неймспейс из конфига
	controller := NewController(ctx, kubeClient, echoClient, "default")

	kubeInfromerFactory.Start(ctx.Done())
	echoInformerFactory.Start(ctx.Done())

	if err = controller.Run(ctx, 2); err != nil {
		logger.Error(err, "Error running controller")
		klog.FlushAndExit(klog.ExitFlushTimeout, 1)
	}
}

func init() {
	flag.StringVar(&kubeconfigPath, "kubeconfig", "", "Path to kubeconfig.")
	flag.StringVar(&masterURL, "master", "", "The adress of kubernetes Api server. Rewrite kubeconfig value. Only requied if out of cluster.")
}
