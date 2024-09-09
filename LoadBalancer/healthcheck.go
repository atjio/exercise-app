package main

import (
	"net/http"
	"log"
	"strings"
	"sync"
)

type nodeStatus struct {
	url string
	healthy bool
}

func healthcheck(healthcheckUrl string) (healthyNodes, unhealthyNodes []string) {
	nodeList := append(state.healthyNodes, state.unhealthyNodes...)
	httpClient := http.Client {
		Timeout: nodeMaxTimeout,
	}

	nodes := make([]nodeStatus,len(nodeList))

	var wg sync.WaitGroup

	for i, node := range nodeList {
		wg.Add(1)

		go func(i int, node string) {
			defer wg.Done()
			resp, err := httpClient.Get(node + healthcheckUrl)
			nodes[i] = nodeStatus{url: node, healthy: (err == nil && resp.StatusCode == 200)}
		} (i, node)
	}

	wg.Wait()
	for _, node := range nodes {
		if (node.healthy) {
			healthyNodes = append(healthyNodes, node.url)
		} else {
			unhealthyNodes = append(unhealthyNodes, node.url)
		}
	}

	log.Println("Healthy Nodes: " + strings.Join(healthyNodes, ","))
	log.Println("Unhealthy Nodes: " + strings.Join(unhealthyNodes, ","))

	return 
}