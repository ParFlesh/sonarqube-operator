module github.com/parflesh/sonarqube-operator

go 1.13

require (
	github.com/magiconair/properties v1.8.0
	github.com/operator-framework/operator-sdk v0.17.0
	github.com/spf13/pflag v1.0.5
	github.com/thanhpk/randstr v1.0.4
	golang.org/x/mod v0.2.0
	golang.org/x/net v0.0.0-20200226121028-0de0cce0169b
	k8s.io/api v0.17.4
	k8s.io/apimachinery v0.17.4
	k8s.io/client-go v12.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.5.2
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible // Required by OLM
	k8s.io/client-go => k8s.io/client-go v0.17.4 // Required by prometheus-operator
)
