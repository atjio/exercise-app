package main

import (
	"net/http"
	"log"
	"strings"
)

func healthcheck(healthcheckUrl string) (healthyNodes, unhealthyNodes []string) {
	nodeList := append(state.healthyNodes, state.unhealthyNodes...)

	for _, node := range nodeList {
		resp, err := http.Get(node + healthcheckUrl)
		if (err != nil || resp.StatusCode != 200) {
			unhealthyNodes = append(unhealthyNodes, node)
		} else {
			healthyNodes = append(healthyNodes, node)
		}
	}

	log.Println("Healthy Nodes: " + strings.Join(healthyNodes, ","))
	log.Println("Unhealthy Nodes: " + strings.Join(unhealthyNodes, ","))

	return 
}