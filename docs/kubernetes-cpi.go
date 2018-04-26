package main

import (
	"fmt"
	"os"

	"flag"
	"path/filepath"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"github.com/cppforlife/bosh-cpi-go/rpc"
	//	"k8s.io/apimachinery/pkg/api/errors"

	// "k8s.io/client-go/1.5/kubernetes"
	// "k8s.io/client-go/1.5/pkg/api/v1"
	// "k8s.io/client-go/1.5/tools/clientcmd"
	"github.com/satori/go.uuid"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type CPIFactory struct{}

type CPI struct{}

var _ apiv1.CPIFactory = CPIFactory{}
var _ apiv1.CPI = CPI{}

var k8sClient *kubernetes.Clientset
var namespace = "default"

func main() {
	u2 := uuid.NewV4()
	fmt.Printf("UUIDv4: %s\n", u2)

	var err error
	k8sConfigPath := filepath.Join(".", "kubeconfig")
	k8sClient, err = initK8s(k8sConfigPath)
	if err != nil {
		panic(err.Error())
	}

	logger := boshlog.NewLogger(boshlog.LevelNone)

	cli := rpc.NewFactory(logger).NewCLI(CPIFactory{})

	err = cli.ServeOnce()
	if err != nil {
		logger.Error("main", "Serving once: %s", err)
		os.Exit(1)
	}

}

func initK8s(k8sConfigPath string) (*kubernetes.Clientset, error) {
	kubeconfig := flag.String("kubeconfig", k8sConfigPath, "path to the kubeconfig file")
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	// create the clientset
	k8sClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return k8sClient, nil
}

// Empty CPI implementation

func (f CPIFactory) New(_ apiv1.CallContext) (apiv1.CPI, error) {
	return CPI{}, nil
}

func (c CPI) Info() (apiv1.Info, error) {
	return apiv1.Info{}, nil
}

func (c CPI) CreateStemcell(imagePath string, _ apiv1.StemcellCloudProps) (apiv1.StemcellCID, error) {
	return apiv1.NewStemcellCID("stemcell-cid"), nil
}

func (c CPI) DeleteStemcell(cid apiv1.StemcellCID) error {
	return nil
}

func (c CPI) CreateVM(
	agentID apiv1.AgentID, stemcellCID apiv1.StemcellCID,
	cloudProps apiv1.VMCloudProps, networks apiv1.Networks,
	associatedDiskCIDs []apiv1.DiskCID, env apiv1.VMEnv) (apiv1.VMCID, error) {

	//read the config to create a pod instead of a VM.
	//check pods (shouldn't work yet)
	podsClient := k8sClient.CoreV1().Pods(corev1.NamespaceDefault)
	pod, err := podsClient.Create(&corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-pod",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "jupyter-notebook",
					Image: "jupyter/minimal-notebook",
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: 8888,
						},
					},
				},
			},
		},
	})

	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Pod %s info namespace: %s clustername: %s \n", pod, pod.Namespace, pod.ClusterName)

	return apiv1.NewVMCID("vm-cid"), nil
}

func (c CPI) DeleteVM(cid apiv1.VMCID) error {

	podsClient := k8sClient.CoreV1().Pods(corev1.NamespaceDefault)
	pods, err := podsClient.Get(cid.AsString(), metav1.GetOptions{})

	fmt.Printf("We have found %s number of pods ", pods.Size)

	if err != nil {
		panic(err.Error())
	}

	// err = podsClient.Delete("my-pod", metav1.DeleteOptions{
	// 	GracePeriodSeconds: 10,
	// 	Preconditions :
	// 	OrphanDependents : false,
	// 	PropagationPolicy : metav1.Prop
	// })

	// fmt.Printf("Pod %s info namespace: %s clustername: %s \n", pod, pod.Namespace, pod.ClusterName)
	return nil
}

func (c CPI) CalculateVMCloudProperties(res apiv1.VMResources) (apiv1.VMCloudProps, error) {
	return apiv1.NewVMCloudPropsFromMap(map[string]interface{}{}), nil
}

func (c CPI) SetVMMetadata(cid apiv1.VMCID, metadata apiv1.VMMeta) error {
	return nil
}

func (c CPI) HasVM(cid apiv1.VMCID) (bool, error) {
	return false, nil
}

func (c CPI) RebootVM(cid apiv1.VMCID) error {
	return nil
}

func (c CPI) GetDisks(cid apiv1.VMCID) ([]apiv1.DiskCID, error) {
	return []apiv1.DiskCID{}, nil
}

func (c CPI) CreateDisk(size int,
	cloudProps apiv1.DiskCloudProps, associatedVMCID *apiv1.VMCID) (apiv1.DiskCID, error) {

	return apiv1.NewDiskCID("disk-cid"), nil
}

func (c CPI) DeleteDisk(cid apiv1.DiskCID) error {
	return nil
}

func (c CPI) AttachDisk(vmCID apiv1.VMCID, diskCID apiv1.DiskCID) error {
	return nil
}

func (c CPI) DetachDisk(vmCID apiv1.VMCID, diskCID apiv1.DiskCID) error {
	return nil
}

func (c CPI) HasDisk(cid apiv1.DiskCID) (bool, error) {
	return false, nil
}
