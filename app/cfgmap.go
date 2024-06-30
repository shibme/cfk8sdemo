package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	namespace *string
	token     *string
)

func getNamespace() string {
	if namespace == nil {
		ns := os.Getenv("NAMESPACE")
		if ns == "" {
			namespaceBytes, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
			if err != nil {
				panic(err)
			}
			ns = string(namespaceBytes)
		}
		namespace = &ns
	}
	return *namespace
}

func getToken() string {
	if token == nil {
		t := os.Getenv("API_AUTH_TOKEN")
		token = &t
	}
	return *token
}

func putMap(clientset *kubernetes.Clientset) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		mapName := ctx.Param("name")
		if getToken() == "" || ctx.GetHeader("Authorization") != "Bearer "+getToken() {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		var mapData map[string]string
		if err := ctx.BindJSON(&mapData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		confMap, err := clientset.CoreV1().ConfigMaps(getNamespace()).Get(context.Background(), mapName, metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				configMap := &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      mapName,
						Namespace: getNamespace(),
					},
					Data: mapData,
				}
				if _, err = clientset.CoreV1().ConfigMaps(getNamespace()).Create(context.Background(), configMap, metav1.CreateOptions{}); err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				log.Println("Created configmap:", configMap.Name)
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			updateRequired := false
			if len(confMap.Data) != len(mapData) {
				updateRequired = true
				confMap.Data = mapData
			} else {
				for k, v := range mapData {
					if confMap.Data[k] != v {
						updateRequired = true
						confMap.Data = mapData
						break
					}
				}
			}
			var msg string
			if updateRequired {
				if _, err = clientset.CoreV1().ConfigMaps(getNamespace()).Update(context.Background(), confMap, metav1.UpdateOptions{}); err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				msg = "Updated configmap: " + confMap.Name
			} else {
				msg = "No update required for configmap: " + confMap.Name
			}
			log.Println(msg)
		}
	}
}

func getMap(clientset *kubernetes.Clientset) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := ctx.Param("name")
		confMap, err := clientset.CoreV1().ConfigMaps(getNamespace()).Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, confMap.Data)
	}
}

func listMaps(clientset *kubernetes.Clientset) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		confMaps, err := clientset.CoreV1().ConfigMaps(getNamespace()).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			ctx.JSON(http.StatusMethodNotAllowed, gin.H{"error": err.Error()})
			return
		}
		var mapNames []string
		for _, confMap := range confMaps.Items {
			mapNames = append(mapNames, confMap.Name)
		}
		ctx.JSON(http.StatusOK, mapNames)
	}
}
