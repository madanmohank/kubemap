package kubemap

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	apps_v1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/tools/cache"
)

func (m *Mapper) map_v1_DeploymentObj(obj ResourceEvent, store cache.Store) (MapResult, error) {
	var deployment apps_v1.Deployment
	var namespaceKeys []string

	if obj.Event != nil {
		deployment = *obj.Event.(*apps_v1.Deployment).DeepCopy()

		keys := store.ListKeys()
		for _, b64Key := range keys {
			encodedKey, _ := base64.StdEncoding.DecodeString(b64Key)
			key := fmt.Sprintf("%s", encodedKey)
			if len(strings.Split(key, "$")) > 0 {
				if strings.Split(key, "$")[0] == obj.Namespace {
					namespaceKeys = append(namespaceKeys, key)
				}
			}
		}

		for _, namespaceKey := range namespaceKeys {
			metaIdentifierString := strings.Split(namespaceKey, "$")[1]
			metaIdentifier := MetaIdentifier{}

			json.Unmarshal([]byte(metaIdentifierString), &metaIdentifier)

			//Try matching with Service
			for _, svcID := range metaIdentifier.ServicesIdentifier.MatchLabels {
				if reflect.DeepEqual(deployment.Spec.Selector.MatchLabels, svcID) {
					//Service and deployment matches. Add service to this mapped resource
					// mappedResource, _ := getObjectFromStore(namespaceKey, store)
					mappedResource, _ := getObjectFromStore(base64.StdEncoding.EncodeToString([]byte(namespaceKey)), store)

					for i, mappedDeployment := range mappedResource.Kube.DeploymentsV1 {
						if mappedDeployment.Name == deployment.Name {
							mappedResource.Kube.DeploymentsV1[i] = deployment

							return MapResult{
								Action:         "Updated",
								Key:            namespaceKey,
								IsMapped:       true,
								MappedResource: mappedResource,
								Message:        fmt.Sprintf("Deployment %s is updated in Common Label %s after matching with service", deployment.Name, mappedResource.CommonLabel),
							}, nil
						}
					}

					mappedResource.Kube.DeploymentsV1 = append(mappedResource.Kube.DeploymentsV1, deployment)
					return MapResult{
						Action:         "Updated",
						Key:            namespaceKey,
						IsMapped:       true,
						MappedResource: mappedResource,
						Message:        fmt.Sprintf("Deployment %s is added to Common Label %s after matching with service", deployment.Name, mappedResource.CommonLabel),
					}, nil
				}
			}

			//Try matching with Deployment
			for _, depID := range metaIdentifier.DeploymentsIdentifier.MatchLabels {
				if reflect.DeepEqual(deployment.Spec.Selector.MatchLabels, depID) {
					//Service and deployment matches. Add service to this mapped resource
					// mappedResource, _ := getObjectFromStore(namespaceKey, store)
					mappedResource, _ := getObjectFromStore(base64.StdEncoding.EncodeToString([]byte(namespaceKey)), store)

					for i, mappedDeployment := range mappedResource.Kube.DeploymentsV1 {
						if mappedDeployment.Name == deployment.Name {
							mappedResource.Kube.DeploymentsV1[i] = deployment

							return MapResult{
								Action:         "Updated",
								Key:            namespaceKey,
								IsMapped:       true,
								MappedResource: mappedResource,
								Message:        fmt.Sprintf("Deployment %s is updated io Common Label %s after matching with deployment", deployment.Name, mappedResource.CommonLabel),
							}, nil
						}
					}
				}
			}

			//Try matching with Replica set
			for _, rsID := range metaIdentifier.ReplicaSetsIdentifier {
				for _, ownerReference := range rsID.OwnerReferences {
					if ownerReference == deployment.Name {
						//Deployment and RS matches. Add deployment to this mapped resource
						// mappedResource, _ := getObjectFromStore(namespaceKey, store)
						mappedResource, _ := getObjectFromStore(base64.StdEncoding.EncodeToString([]byte(namespaceKey)), store)

						for i, mappedDeployment := range mappedResource.Kube.DeploymentsV1 {
							if mappedDeployment.Name == deployment.Name {
								mappedResource.Kube.DeploymentsV1[i] = deployment

								return MapResult{
									Action:         "Updated",
									Key:            namespaceKey,
									IsMapped:       true,
									MappedResource: mappedResource,
									Message:        fmt.Sprintf("Deployment %s is updated io Common Label %s after matching with replica set", deployment.Name, mappedResource.CommonLabel),
								}, nil
							}
						}

						mappedResource.Kube.DeploymentsV1 = append(mappedResource.Kube.DeploymentsV1, deployment)
						if len(mappedResource.Kube.DeploymentsV1) < 2 { //Set Common Label to deployment name.
							mappedResource.CommonLabel = deployment.Name
						}
						return MapResult{
							Action:         "Updated",
							Key:            namespaceKey,
							IsMapped:       true,
							MappedResource: mappedResource,
							Message:        fmt.Sprintf("Deployment %s is added to Common Label %s after matching with replica set", deployment.Name, mappedResource.CommonLabel),
						}, nil
					}
				}
			}

			//Try matching with Pod
			for _, podID := range metaIdentifier.PodsIdentifier {
				podMatchedLabels := make(map[string]string)
				for podKey, podValue := range podID.MatchLabels {
					if val, ok := deployment.Spec.Selector.MatchLabels[podKey]; ok {
						if val == podValue {
							podMatchedLabels[podKey] = podValue
						}
					}
				}

				if reflect.DeepEqual(deployment.Spec.Selector.MatchLabels, podMatchedLabels) {
					//Deployment and RS matches. Add deployment to this mapped resource
					// mappedResource, _ := getObjectFromStore(namespaceKey, store)
					mappedResource, _ := getObjectFromStore(base64.StdEncoding.EncodeToString([]byte(namespaceKey)), store)

					for i, mappedDeployment := range mappedResource.Kube.DeploymentsV1 {
						if mappedDeployment.Name == deployment.Name {
							mappedResource.Kube.DeploymentsV1[i] = deployment

							return MapResult{
								Action:         "Updated",
								Key:            namespaceKey,
								IsMapped:       true,
								MappedResource: mappedResource,
								Message:        fmt.Sprintf("Deployment %s is updated io Common Label %s after matching with pod", deployment.Name, mappedResource.CommonLabel),
							}, nil
						}
					}

					mappedResource.Kube.DeploymentsV1 = append(mappedResource.Kube.DeploymentsV1, deployment)
					if len(mappedResource.Kube.DeploymentsV1) < 2 { //Set Common Label to deployment name.
						mappedResource.CommonLabel = deployment.Name
					}
					return MapResult{
						Action:         "Updated",
						Key:            namespaceKey,
						IsMapped:       true,
						MappedResource: mappedResource,
						Message:        fmt.Sprintf("Deployment %s is added to Common Label %s after matching with pod", deployment.Name, mappedResource.CommonLabel),
					}, nil
				}
			}
		}

		//Didn't find any match. Create new resource
		newMappedService := MappedResource{}
		newMappedService.CommonLabel = deployment.Name
		newMappedService.CurrentType = "deployment"
		newMappedService.Namespace = deployment.Namespace
		newMappedService.Kube.DeploymentsV1 = append(newMappedService.Kube.DeploymentsV1, deployment)

		return MapResult{
			Action:         "Added",
			IsMapped:       true,
			MappedResource: newMappedService,
			Message:        fmt.Sprintf("New deployment %s is created with Common Label %s", deployment.Name, newMappedService.CommonLabel),
		}, nil
	}

	//Handle Delete
	if obj.EventType == "DELETED" {
		m.info(fmt.Sprintf("DELETE received. - K8s Type - %s Name - %s Namespace - %s", obj.ResourceType, obj.Name, obj.Namespace))

		keys := store.ListKeys()
		for _, b64Key := range keys {
			encodedKey, _ := base64.StdEncoding.DecodeString(b64Key)
			key := fmt.Sprintf("%s", encodedKey)
			if len(strings.Split(key, "$")) > 0 {
				if strings.Split(key, "$")[0] == obj.Namespace {
					namespaceKeys = append(namespaceKeys, key)
				}
			}
		}

		var newDepSet []apps_v1.Deployment
		for _, namespaceKey := range namespaceKeys {
			metaIdentifierString := strings.Split(namespaceKey, "$")[1]
			metaIdentifier := MetaIdentifier{}

			json.Unmarshal([]byte(metaIdentifierString), &metaIdentifier)

			for _, mappedDepName := range metaIdentifier.DeploymentsIdentifier.Names {
				if mappedDepName == obj.Name {
					//Pod is being deleted.
					// mappedResource, _ := getObjectFromStore(namespaceKey, store)
					mappedResource, _ := getObjectFromStore(base64.StdEncoding.EncodeToString([]byte(namespaceKey)), store)

					newDepSet = nil
					for _, mappedDeployment := range mappedResource.Kube.DeploymentsV1 {
						if mappedDeployment.Name != obj.Name {
							newDepSet = append(newDepSet, mappedDeployment)
						}
					}

					if len(mappedResource.Kube.IngressesNW) > 0 || len(mappedResource.Kube.Services) > 0 || len(mappedResource.Kube.ReplicaSetsV1) > 0 || len(mappedResource.Kube.Pods) > 0 || len(mappedResource.Kube.DeploymentsV1) > 1 {
						//It has another resources.
						mappedResource.Kube.DeploymentsV1 = nil
						mappedResource.Kube.DeploymentsV1 = newDepSet

						m.info(fmt.Sprintf("DELETE Completed. - K8s Type - %s Name - %s Namespace - %s CL %s updated.", obj.ResourceType, obj.Name, obj.Namespace, mappedResource.CommonLabel))
						return MapResult{
							Action:         "Updated",
							Key:            namespaceKey,
							IsMapped:       true,
							MappedResource: mappedResource,
							Message:        fmt.Sprintf("Deployment %s is deleted from Common Label %s", deployment.Name, mappedResource.CommonLabel),
						}, nil
					}

					m.info(fmt.Sprintf("DELETE Completed. - K8s Type - %s Name - %s Namespace - %s CL %s deleted.", obj.ResourceType, obj.Name, obj.Namespace, mappedResource.CommonLabel))
					return MapResult{
						Action:         "Deleted",
						Key:            namespaceKey,
						IsMapped:       true,
						CommonLabel:    mappedResource.CommonLabel,
						MappedResource: mappedResource,
						Message:        fmt.Sprintf("Deployment %s is deleted from Common Label %s", deployment.Name, mappedResource.CommonLabel),
					}, nil

				}
			}
		}
	}
	return MapResult{}, nil
}

func (m *Mapper) map_v1_ReplicaSetObj(obj ResourceEvent, store cache.Store) (MapResult, error) {
	var replicaSet apps_v1.ReplicaSet
	var namespaceKeys []string

	if obj.Event != nil {
		replicaSet = *obj.Event.(*apps_v1.ReplicaSet).DeepCopy()

		keys := store.ListKeys()
		for _, b64Key := range keys {
			encodedKey, _ := base64.StdEncoding.DecodeString(b64Key)
			key := fmt.Sprintf("%s", encodedKey)
			if len(strings.Split(key, "$")) > 0 {
				if strings.Split(key, "$")[0] == obj.Namespace {
					namespaceKeys = append(namespaceKeys, key)
				}
			}
		}

		for _, namespaceKey := range namespaceKeys {
			metaIdentifierString := strings.Split(namespaceKey, "$")[1]
			metaIdentifier := MetaIdentifier{}

			json.Unmarshal([]byte(metaIdentifierString), &metaIdentifier)

			//Try matching with Service
			if metaIdentifier.ServicesIdentifier.MatchLabels != nil {
				for _, svcID := range metaIdentifier.ServicesIdentifier.MatchLabels {
					rsMatchedLabels := make(map[string]string)
					if svcID != nil && replicaSet.Spec.Selector.MatchLabels != nil {
						for svcKey, svcValue := range svcID {
							if val, ok := replicaSet.Spec.Selector.MatchLabels[svcKey]; ok {
								if val == svcValue {
									rsMatchedLabels[svcKey] = svcValue
								}
							}
						}
					}
					if reflect.DeepEqual(rsMatchedLabels, svcID) {
						//Service and pod matches. Add pod to this mapped resource
						// mappedResource, _ := getObjectFromStore(namespaceKey, store)
						mappedResource, _ := getObjectFromStore(base64.StdEncoding.EncodeToString([]byte(namespaceKey)), store)

						for i, mappedReplicaSet := range mappedResource.Kube.ReplicaSetsV1 {
							if mappedReplicaSet.Name == replicaSet.Name {
								mappedResource.Kube.ReplicaSetsV1[i] = replicaSet

								return MapResult{
									Action:         "Updated",
									Key:            namespaceKey,
									IsMapped:       true,
									MappedResource: mappedResource,
									Message:        fmt.Sprintf("Replica set %s is updated in Common Label %s after matching with service", replicaSet.Name, mappedResource.CommonLabel),
								}, nil
							}
						}

						mappedResource.Kube.ReplicaSetsV1 = append(mappedResource.Kube.ReplicaSetsV1, replicaSet)

						return MapResult{
							Action:         "Updated",
							Key:            namespaceKey,
							IsMapped:       true,
							MappedResource: mappedResource,
							Message:        fmt.Sprintf("Replica set %s is added to Common Label %s after matching with service", replicaSet.Name, mappedResource.CommonLabel),
						}, nil
					}
				}
			}

			//Try matching with Deployment
			for _, depID := range metaIdentifier.DeploymentsIdentifier.MatchLabels {
				rsMatchedLabels := make(map[string]string)
				for depKey, depValue := range depID {
					if val, ok := replicaSet.Spec.Selector.MatchLabels[depKey]; ok {
						if val == depValue {
							rsMatchedLabels[depKey] = depValue
						}
					}
				}
				if reflect.DeepEqual(rsMatchedLabels, depID) {
					//Service and deployment matches. Add service to this mapped resource
					// mappedResource, _ := getObjectFromStore(namespaceKey, store)
					mappedResource, _ := getObjectFromStore(base64.StdEncoding.EncodeToString([]byte(namespaceKey)), store)

					for i, mappedReplicaSet := range mappedResource.Kube.ReplicaSetsV1 {
						if mappedReplicaSet.Name == replicaSet.Name {
							mappedResource.Kube.ReplicaSetsV1[i] = replicaSet

							return MapResult{
								Action:         "Updated",
								Key:            namespaceKey,
								IsMapped:       true,
								MappedResource: mappedResource,
								Message:        fmt.Sprintf("Replica set %s is updated in Common Label %s after matching with deployment", replicaSet.Name, mappedResource.CommonLabel),
							}, nil
						}
					}

					mappedResource.Kube.ReplicaSetsV1 = append(mappedResource.Kube.ReplicaSetsV1, replicaSet)
					return MapResult{
						Action:         "Updated",
						Key:            namespaceKey,
						IsMapped:       true,
						MappedResource: mappedResource,
						Message:        fmt.Sprintf("Replica set %s is added to Common Label %s after matching with deployment", replicaSet.Name, mappedResource.CommonLabel),
					}, nil
				}
			}

			//Try matching with Replica set
			for _, rsID := range metaIdentifier.ReplicaSetsIdentifier {
				if reflect.DeepEqual(replicaSet.Spec.Selector.MatchLabels, rsID.MatchLabels) {
					//Service and deployment matches. Add service to this mapped resource
					// mappedResource, _ := getObjectFromStore(namespaceKey, store)
					mappedResource, _ := getObjectFromStore(base64.StdEncoding.EncodeToString([]byte(namespaceKey)), store)

					for i, mappedReplicaSet := range mappedResource.Kube.ReplicaSetsV1 {
						if mappedReplicaSet.Name == replicaSet.Name {
							mappedResource.Kube.ReplicaSetsV1[i] = replicaSet

							return MapResult{
								Action:         "Updated",
								Key:            namespaceKey,
								IsMapped:       true,
								MappedResource: mappedResource,
								Message:        fmt.Sprintf("Replica set %s is updated io Common Label %s after matching with replica set", replicaSet.Name, mappedResource.CommonLabel),
							}, nil
						}
					}
				}
			}

			//Try matching with Pod
			for _, podID := range metaIdentifier.PodsIdentifier {
				rsMatchedLabels := make(map[string]string)
				for podKey, podValue := range podID.MatchLabels {
					if val, ok := replicaSet.Spec.Selector.MatchLabels[podKey]; ok {
						if val == podValue {
							rsMatchedLabels[podKey] = podValue
						}
					}
				}
				if reflect.DeepEqual(rsMatchedLabels, replicaSet.Spec.Selector.MatchLabels) {
					//Service and deployment matches. Add service to this mapped resource
					// mappedResource, _ := getObjectFromStore(namespaceKey, store)
					mappedResource, _ := getObjectFromStore(base64.StdEncoding.EncodeToString([]byte(namespaceKey)), store)

					for i, mappedReplicaSet := range mappedResource.Kube.ReplicaSetsV1 {
						if mappedReplicaSet.Name == replicaSet.Name {
							mappedResource.Kube.ReplicaSetsV1[i] = replicaSet

							return MapResult{
								Action:         "Updated",
								Key:            namespaceKey,
								IsMapped:       true,
								MappedResource: mappedResource,
								Message:        fmt.Sprintf("Replica set %s is updated in Common Label %s after matching with pod", replicaSet.Name, mappedResource.CommonLabel),
							}, nil
						}
					}

					mappedResource.Kube.ReplicaSetsV1 = append(mappedResource.Kube.ReplicaSetsV1, replicaSet)
					if len(mappedResource.Kube.ReplicaSetsV1) < 2 { //Set Common Label to service name.
						mappedResource.CommonLabel = replicaSet.Name
					}
					return MapResult{
						Action:         "Updated",
						Key:            namespaceKey,
						IsMapped:       true,
						MappedResource: mappedResource,
						Message:        fmt.Sprintf("Replica set %s is added to Common Label %s after matching with pod", replicaSet.Name, mappedResource.CommonLabel),
					}, nil
				}
			}
		}

		//Didn't find any match. Create new resource
		newMappedService := MappedResource{}
		newMappedService.CommonLabel = replicaSet.Name
		newMappedService.CurrentType = "replicaset"
		newMappedService.Namespace = replicaSet.Namespace
		newMappedService.Kube.ReplicaSetsV1 = append(newMappedService.Kube.ReplicaSetsV1, replicaSet)

		return MapResult{
			Action:         "Added",
			IsMapped:       true,
			MappedResource: newMappedService,
			Message:        fmt.Sprintf("New replica set %s is added with Common Label %s", replicaSet.Name, newMappedService.CommonLabel),
		}, nil

	}

	//Handle Delete
	if obj.EventType == "DELETED" {
		m.info(fmt.Sprintf("DELETE received. - K8s Type - %s Name - %s Namespace - %s", obj.ResourceType, obj.Name, obj.Namespace))

		keys := store.ListKeys()
		for _, b64Key := range keys {
			encodedKey, _ := base64.StdEncoding.DecodeString(b64Key)
			key := fmt.Sprintf("%s", encodedKey)
			if len(strings.Split(key, "$")) > 0 {
				if strings.Split(key, "$")[0] == obj.Namespace {
					namespaceKeys = append(namespaceKeys, key)
				}
			}
		}

		var newRsSet []apps_v1.ReplicaSet
		for _, namespaceKey := range namespaceKeys {
			metaIdentifierString := strings.Split(namespaceKey, "$")[1]
			metaIdentifier := MetaIdentifier{}

			json.Unmarshal([]byte(metaIdentifierString), &metaIdentifier)

			for _, rsChileSet := range metaIdentifier.ReplicaSetsIdentifier {
				if rsChileSet.Name == obj.Name {
					//Pod is being deleted.
					// mappedResource, _ := getObjectFromStore(namespaceKey, store)
					mappedResource, _ := getObjectFromStore(base64.StdEncoding.EncodeToString([]byte(namespaceKey)), store)

					newRsSet = nil
					for _, mappedRs := range mappedResource.Kube.ReplicaSetsV1 {
						if mappedRs.Name != obj.Name {
							newRsSet = append(newRsSet, mappedRs)
						}
					}

					if len(mappedResource.Kube.IngressesNW) > 0 || len(mappedResource.Kube.Services) > 0 || len(mappedResource.Kube.DeploymentsV1) > 0 || len(mappedResource.Kube.Pods) > 0 || len(mappedResource.Kube.ReplicaSetsV1) > 1 {
						//It has another resources.
						mappedResource.Kube.ReplicaSetsV1 = nil
						mappedResource.Kube.ReplicaSetsV1 = newRsSet

						m.info(fmt.Sprintf("DELETE Completed. - K8s Type - %s Name - %s Namespace - %s CL %s updated.", obj.ResourceType, obj.Name, obj.Namespace, mappedResource.CommonLabel))
						return MapResult{
							Action:         "Updated",
							Key:            namespaceKey,
							IsMapped:       true,
							MappedResource: mappedResource,
							Message:        fmt.Sprintf("Replica set %s is deleted from Common Label %s", replicaSet.Name, mappedResource.CommonLabel),
						}, nil
					}

					m.info(fmt.Sprintf("DELETE Completed. - K8s Type - %s Name - %s Namespace - %s CL %s deleted.", obj.ResourceType, obj.Name, obj.Namespace, mappedResource.CommonLabel))
					return MapResult{
						Action:         "Deleted",
						Key:            namespaceKey,
						IsMapped:       true,
						CommonLabel:    mappedResource.CommonLabel,
						MappedResource: mappedResource,
						Message:        fmt.Sprintf("Replica set %s is deleted from Common Label %s", replicaSet.Name, mappedResource.CommonLabel),
					}, nil

				}
			}
		}
	}

	return MapResult{}, nil
}
