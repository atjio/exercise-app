package global

import (
	"log"
	"strings"
	"sync"
)

type AppState struct {
	mutex sync.Mutex

	HealthyNodes []string 
	UnhealthyNodes []string 
}

func (s *AppState) GetAllNodes() []string {
	return append(s.HealthyNodes, s.UnhealthyNodes...)
}

func (s *AppState) UpdateNodes(healthyNodes []string, unhealthyNodes []string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.HealthyNodes = healthyNodes
	s.UnhealthyNodes = unhealthyNodes

	log.Println("Healthy Nodes: " + strings.Join(s.HealthyNodes, ","))
	log.Println("Unhealthy Nodes: " + strings.Join(s.UnhealthyNodes, ","))
}

func (s *AppState) AddNode(node string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.HealthyNodes = append(s.HealthyNodes, node)
}

// func (s *AppState) useNode() string {
// 	s.mutex.Lock()
// 	defer s.mutex.Unlock()
// }

var State = &AppState{
	HealthyNodes: make([]string, 0),
	UnhealthyNodes: make([]string, 0),
}