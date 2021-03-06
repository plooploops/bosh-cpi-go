package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"github.com/cppforlife/bosh-cpi-go/rpc"
	uuid "github.com/satori/go.uuid"
	//	"k8s.io/apimachinery/pkg/api/errors"
	// "k8s.io/client-go/1.5/kubernetes"
	// "k8s.io/client-go/1.5/pkg/api/v1"
	// "k8s.io/client-go/1.5/tools/clientcmd"
	corev1 "k8s.io/api/core/v1"
	apimachv1 "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type CPIFactory struct{}

type CPI struct{}

// var _ apiv1.CPIFactory = CPIFactory{}
// var _ apiv1.CPI = CPI{}

var k8sClient *kubernetes.Clientset
var namespace = "default"

func main() {
	logger := boshlog.NewLogger(boshlog.LevelNone)

	var err error
	k8sConfigPath := filepath.Join(".", "kubeconfig")
	k8sClient, err = initK8s(k8sConfigPath)
	if err != nil {
		logger.Error("main", "Serving once: %s", err)
		os.Exit(1)
	}

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
	// Handwave here; we need an image from a container registry
	// Ideally we'd take the stemcell tarball and build an image out of, upload that to the registry, and point to that
	// But the most widely used real stemcells are Ubuntu Trusty, so we'll cheat and use that here
	// (though it can't work "for real" because they don't have bosh agents)
	return apiv1.NewStemcellCID("phusion/baseimage"), nil
}

func (c CPI) DeleteStemcell(cid apiv1.StemcellCID) error {
	return nil
}

func (c CPI) CreateVM(
	agentID apiv1.AgentID, stemcellCID apiv1.StemcellCID,
	cloudProps apiv1.VMCloudProps, networks apiv1.Networks,
	associatedDiskCIDs []apiv1.DiskCID, env apiv1.VMEnv) (apiv1.VMCID, error) {

	vmcid := uuid.NewV4().String()

	//read the config to create a pod instead of a VM.
	//check pods (shouldn't work yet)
	podsClient := k8sClient.CoreV1().Pods(namespace)
	_, err := podsClient.Create(&corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: vmcid,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  vmcid,
					Image: stemcellCID.AsString(),
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: 22, //ssh
						},
					},
				},
			},
		},
	})

	if err != nil {
		return apiv1.NewVMCID(""), err
	}

	return apiv1.NewVMCID(vmcid), nil
}

func (c CPI) DeleteVM(cid apiv1.VMCID) error {
	podsClient := k8sClient.CoreV1().Pods(namespace)
	deletePolicy := metav1.DeletePropagationForeground
	err := podsClient.Delete(cid.AsString(), &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})

	return err
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

func (c CPI) CreateDisk(size int, cloudProps apiv1.DiskCloudProps, associatedVMCID *apiv1.VMCID) (apiv1.DiskCID, error) {
	scn := "azurefile"
	diskCID := uuid.NewV4().String()
	quantity, err := resource.ParseQuantity(fmt.Sprintf("%dGi", size/1024))
	if err != nil {
		return apiv1.NewDiskCID(""), err
	}

	pvcClient := k8sClient.CoreV1().PersistentVolumeClaims(namespace)
	_, err = pvcClient.Create(&corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: diskCID,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			StorageClassName: &scn,
			Resources: corev1.ResourceRequirements{
				Requests: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceStorage: quantity}},
		},
	})
	if err != nil {
		return apiv1.NewDiskCID(""), err
	}

	return apiv1.NewDiskCID(diskCID), nil
}

func (c CPI) DeleteDisk(cid apiv1.DiskCID) error {
	return nil
}

func (c CPI) AttachDisk(vmCID apiv1.VMCID, diskCID apiv1.DiskCID) error {
	podsClient := k8sClient.CoreV1().Pods(namespace)

	// Retrieve the current state of the pod
	pod, err := podsClient.Get(vmCID.AsString(), metav1.GetOptions{})
	if err != nil {
		return err
	}

	deletePolicy := metav1.DeletePropagationForeground
	err = podsClient.Delete(vmCID.AsString(), &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})

	//spin wait?
	//pod still in memory
	// Update the pod with the PVC
	pod.Spec.Volumes = []corev1.Volume{
		{
			Name: diskCID.AsString(),
			HostPath: nil // Do like https://github.com/MicrosoftDX/MB-ForensicWatermark/blob/master/k8s/submit-job-tempdisk.yaml#L35-L38
		},
	}
	pod.Spec.Containers[0].VolumeMounts = []corev1.VolumeMount{
		{
			Name:      diskCID.AsString(),
			MountPath: strings.Join([]string{"/mnt", diskCID.AsString()}, "/"),
		},
	}
	pod.ObjectMeta.ResourceVersion = ""

	deleteFinished := true
	for deleteFinished {
		_, err = podsClient.Create(pod)
		deleteFinished = apimachv1.IsAlreadyExists(err) //check if the pod is already exists (pending delete) then try to create it again
		time.Sleep(1 * time.Second)                     //sleep to loosen calls
	}

	return nil
}

func (c CPI) DetachDisk(vmCID apiv1.VMCID, diskCID apiv1.DiskCID) error {
	return nil
}

func (c CPI) HasDisk(cid apiv1.DiskCID) (bool, error) {
	return false, nil
}
