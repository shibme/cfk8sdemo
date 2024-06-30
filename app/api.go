package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	r.POST("/api/map/:name", putMap(clientset))
	r.PUT("/api/map/:name", putMap(clientset))
	r.GET("/api/map/:name", getMap(clientset))
	r.GET("/api/map", listMaps(clientset))
	r.GET("/api/maps", listMaps(clientset))
	if err = r.Run(fmt.Sprintf(":%d", 8888)); err != nil {
		panic(err)
	}
}
