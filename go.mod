module github.com/myonlyzzy/kubectl-resource-view

go 1.15

require (
	github.com/docker/spdystream v0.0.0-20181023171402-6480d4af844c // indirect
	github.com/fatih/color v1.7.0
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/mattn/go-colorable v0.1.1 // indirect
	github.com/mattn/go-isatty v0.0.7 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.4.0
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/cli-runtime v0.19.2
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/kube-openapi v0.0.0-20200831175022-64514a1d5d59 // indirect
	k8s.io/kubectl v0.19.2
	k8s.io/utils v0.0.0-20200912215256-4140de9c8800 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20200914174313-52bf62410745
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20200916235632-714f1137f89b
	k8s.io/cli-runtime v0.19.2 => k8s.io/cli-runtime v0.0.0-20200915100420-3cc3835b3ec2
	k8s.io/client-go v11.0.0+incompatible => k8s.io/client-go v0.0.0-20200917000235-cba7285b7f29
	k8s.io/kubectl v0.19.2 => k8s.io/kubectl v0.0.0-20200921122246-67718c957b49

)
