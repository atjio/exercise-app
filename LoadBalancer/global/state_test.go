package global

import (
	"testing"
)

func TestAppState_GetAllNodes(t *testing.T) {
	state := &AppState{
		HealthyNodes:   []string{"node1", "node2"},
		UnhealthyNodes: []string{"node3", "node4"},
	}

	allNodes := state.GetAllNodes()
	expectedNodes := []string{"node1", "node2", "node3", "node4"}

	if len(allNodes) != len(expectedNodes) {
		t.Errorf("Expected %d nodes, but got %d", len(expectedNodes), len(allNodes))
	}

	for _, node := range expectedNodes {
		if !contains(allNodes, node) {
			t.Errorf("Expected node %s to be in the list, but it was not", node)
		}
	}
}

func TestAppState_UpdateNodes(t *testing.T) {
	state := &AppState{}

	state.UpdateNodes([]string{"node1", "node2"}, []string{"node3", "node4"})

	if len(state.HealthyNodes) != 2 {
		t.Errorf("Expected 2 healthy nodes, but got %d", len(state.HealthyNodes))
	}

	if len(state.UnhealthyNodes) != 2 {
		t.Errorf("Expected 2 unhealthy nodes, but got %d", len(state.UnhealthyNodes))
	}
}

func TestAppState_AddNode(t *testing.T) {
	state := &AppState{}

	state.AddNode("node1")
	state.AddNode("node2")
	state.AddNode("node3")

	if len(state.HealthyNodes) != 3 {
		t.Errorf("Expected 1 healthy node, but got %d", len(state.HealthyNodes))
	}

	if state.HealthyNodes[0] != "node1" {
		t.Errorf("Expected node1 to be the first healthy node, but got %s", state.HealthyNodes[0])
	}

	if state.HealthyNodes[1] != "node2" {
		t.Errorf("Expected node1 to be the first healthy node, but got %s", state.HealthyNodes[1])
	}

	if state.HealthyNodes[2] != "node3" {
		t.Errorf("Expected node1 to be the first healthy node, but got %s", state.HealthyNodes[2])
	}
}

func TestAppState_GetNextHealthyNode(t *testing.T) {
	state := &AppState{
		HealthyNodes: []string{"node1", "node2", "node3"},
	}

	nextNode := state.GetNextHealthyNode()

	if nextNode != "node1" {
		t.Errorf("Expected next healthy node to be node1, but got %s", nextNode)
	}

	if len(state.HealthyNodes) != 3 {
		t.Errorf("Expected 3 healthy nodes left, but got %d", len(state.HealthyNodes))
	}

	if state.HealthyNodes[0] != "node2" {
		t.Errorf("Expected node2 to be the first healthy node, but got %s", state.HealthyNodes[0])
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
