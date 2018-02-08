package main

import (
	"flag"
	"fmt"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
	"strings"
)

func in(item string, items []string) bool {
	for _, i := range items {
		if strings.Contains(i, item) {
			return true
		}
	}
	return false
}

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig",
			filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	user := flag.String("user", "", "k8ts user")
	group := flag.String("group", "", "k8ts user group")
	namespace := flag.String("namespace", "default", "k8ts namespace")

	flag.Parse()

	// both unset - bad!
	if *user == "" && *group == "" {
		fmt.Fprintln(os.Stderr, "Invalid configuration. Both options unset.")
		os.Exit(1)
	}

	// both set - bad!
	if *user != "" && *group != "" {
		fmt.Fprintln(os.Stderr, "Invalid configuration. Both options were set")
		os.Exit(1)
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to build kubeconfig. Reason: %v\n", err.Error())
		os.Exit(1)
	}

	emptyGetOpts := metaV1.GetOptions{}
	emptyListOpts := metaV1.ListOptions{}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to initialize kubeconfig. Reason: %v\n", err.Error())
		os.Exit(1)
	}

	v1RBAC := clientset.RbacV1()
	rolesBindings, err := v1RBAC.RoleBindings(*namespace).List(emptyListOpts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to list role bindings. Reason: %v\n", err.Error())
		os.Exit(1)
	}

	var configmaps []string
	// we need to loop over all role bindings within given namespace
	for _, r := range rolesBindings.Items {
		// each role binding can have more than one subject
		for _, s := range r.Subjects {
			// we seek for particular objects: User or Group
			if (s.Kind == "User" && s.Name == *user) || (s.Kind == "Group" && s.Name == *group) {
				// we have to do revers thing: get role by binding of it
				role, err := v1RBAC.Roles(*namespace).Get(r.RoleRef.Name, emptyGetOpts)
				// missing role is a bad case, unfortunately kubectl nor kube API
				// doesn't valida such cases along with the other edge cases, lame!
				if err != nil {
					fmt.Fprintf(os.Stderr, "Unable to get particular role '%v'. Reason: %v\n\n",
						r.RoleRef.Name, err.Error())
					os.Exit(1)
				}
				// each role has number of rules, so we need to check each rule,
				// because kubectl nor kube API doesn't validate rules overlapping
				for _, rule := range role.Rules {
					// that's why we check each rule for the configmap (see comments at templates/roles.yml)
					if in("configmaps", rule.Resources) {
						// this one is tricky, we can let users basically access
						// all configmaps by specifying resourceNames as '[*]',
						// so this code will show you current configmaps instead of * (all)
						configmaps = append(configmaps, rule.ResourceNames...)
						if len(configmaps) == 1 && configmaps[0] == "*" {
							cMaps, err := clientset.CoreV1().ConfigMaps(*namespace).List(emptyListOpts)
							if err != nil {
								fmt.Fprintf(os.Stderr, "Unable to list configmaps. Reason: %v\n\n", err.Error())
								os.Exit(1)
							}
							if len(cMaps.Items) != 0 {
								configmaps = []string{}
								for _, m := range cMaps.Items {
									configmaps = append(configmaps, m.Name)
								}
							}
						}
					}
				}
				fmt.Printf("%v '%v' granted with access according "+
					"to its role '%v' to the following configmaps: %v\n",
					s.Kind, s.Name, role.Name, configmaps)
				os.Exit(0)
			}
		}
	}
}
