// TODO: Вынести в отдельный package controller
package main

import echov1alpha1 (
	echo "github.com/countneuroman/hello-operator/pkg/apis/hellocontroller/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
)
func createJob(newEcho *echov1alpha1.Echo, namespace string) *batchv1.Job {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: newEcho.ObjectMeta.Name,
			Namespace: namespace,
			labels: make(map[string]string),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(
					newEcho,
					echov1alpha1.SchemeGroupVersion.WithKind(echo.EchoKind)
				)
			}
		},
		Spec: createJobSpec(echo.Name, namespace, echo.Spec.Message)
	}
}

func createJobSpec(name, namespace, msg string) batchv1.Spec {
	return batchv1.JobSpec{
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: name + "-",
				Namespace: namespace,
				Labels: make(map[string]string)
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name: 	name,
						Image:	"bussybox:1.37.0",
						Command: []string{"echo", msg},
						ImagePullPolicy: "IfNotPresent"
					}
				},
				RestartPolicy: corev1.RestartPolicyNever
			}
		}
	}
}