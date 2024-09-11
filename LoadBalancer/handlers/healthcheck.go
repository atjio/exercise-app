package handlers

import (
	"net/http"
	"sync"
	"time"

	"LoadBalancer/global"
)

type nodeStatus struct {
	url     string
	healthy bool
}

func InitiateHealthcheck() {
	for {
		RefreshNodeHealthStatus()
		time.Sleep(global.HEALTHCHECK_DELAY_IN_MS)
	}
}

func RefreshNodeHealthStatus() {
	nodeList := global.State.GetAllNodes()
	var (
		healthyNodes   []string
		unhealthyNodes []string
		wg             sync.WaitGroup
		nodes          = make([]nodeStatus, len(nodeList))
	)

	for i, node := range nodeList {
		wg.Add(1)
		go func(i int, node string) {
			defer wg.Done()
			nodes[i] = nodeStatus{url: node, healthy: healthcheck(node)}
		}(i, node)
	}

	wg.Wait()
	for _, node := range nodes {
		if node.healthy {
			healthyNodes = append(healthyNodes, node.url)
		} else {
			unhealthyNodes = append(unhealthyNodes, node.url)
		}
	}

	global.State.UpdateNodes(healthyNodes, unhealthyNodes)
}

var healthcheck = func(url string) bool {
	httpClient := http.Client{
		Timeout: global.NODE_MAX_TIMEOUT,
	}

	resp, err := httpClient.Get(url + global.HEALTHCHECK_URL)
	return (err == nil && resp.StatusCode == 200)
}
