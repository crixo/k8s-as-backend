module github.com/crixo/k8s-as-backend/webhook-server

go 1.13

require (
	k8s.io/api v0.0.0-20200131112707-d64dbec685a4
	k8s.io/apiextensions-apiserver v0.0.0-00010101000000-000000000000 // indirect
	k8s.io/apimachinery v0.0.0-20200131112342-0c9ec93240c9
	k8s.io/klog v1.0.0
)

//https://stackoverflow.com/questions/59187781/revision-v0-0-0-unknown-for-go-get-k8s-io-kubernetes
//https://github.com/kubernetes/kubernetes/issues/79384
// go to https://github.com/kubernetes/kubernetes
// select Branch: k8s version/release/tag
// grab the same section/entries from go.mod file on the root
replace (
	k8s.io/api => k8s.io/api v0.0.0-20200131112707-d64dbec685a4
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20200208193839-84fe3c0be50e
	k8s.io/apimachinery => k8s.io/apimachinery v0.16.7-beta.0.0.20200131112342-0c9ec93240c9
	k8s.io/apiserver => k8s.io/apiserver v0.0.0-20200208192130-2d005a048922
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20200131120220-9674fbb91442
	k8s.io/client-go => k8s.io/client-go v0.0.0-20191016111102-bec269661e48
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20200131203752-f498d522efeb
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.0.0-20200131121422-fc6110069b18
	k8s.io/code-generator => k8s.io/code-generator v0.16.7-beta.0.0.20200131112027-a3045e5e55c0
	k8s.io/component-base => k8s.io/component-base v0.0.0-20200131113804-409d4deb41dd
	k8s.io/cri-api => k8s.io/cri-api v0.16.7-beta.0
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.0.0-20200131121824-f033562d74c3
	k8s.io/gengo => k8s.io/gengo v0.0.0-20190822140433-26a664648505
	k8s.io/heapster => k8s.io/heapster v1.2.0-beta.1
	k8s.io/klog => k8s.io/klog v0.4.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20200208192621-0eeb50407007
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.0.0-20200131121224-13b3f231e47d
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20190816220812-743ec37842bf
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.0.0-20200131120626-5b8ba5e54e1f
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.0.0-20200131121024-5f0ba0866863
	k8s.io/kubectl => k8s.io/kubectl v0.0.0-20200131122652-b28c9fbca10f
	k8s.io/kubelet => k8s.io/kubelet v0.0.0-20200131120825-905bd8eea4c4
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.0.0-20200208200602-3a1c7effd2b3
	k8s.io/metrics => k8s.io/metrics v0.0.0-20200131120008-5c623d74062d
	k8s.io/node-api => k8s.io/node-api v0.0.0-20200131122255-04077c800298
	k8s.io/repo-infra => k8s.io/repo-infra v0.0.0-20181204233714-00fe14e3d1a3
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.0.0-20200208192953-f8dc80bbc173
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.0.0-20200131120425-dca0863cb511
	k8s.io/sample-controller => k8s.io/sample-controller v0.0.0-20200131115407-2b45fb79af22
	k8s.io/utils => k8s.io/utils v0.0.0-20190801114015-581e00157fb1
)
