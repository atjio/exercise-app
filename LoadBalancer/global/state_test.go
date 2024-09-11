package global

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppState_GetAllNodes(t *testing.T) {
	state := &AppState{
		HealthyNodes:   []string{"http://localhost:1234", "http://localhost:5678", "http://localhost:9012"},
		UnhealthyNodes: []string{"http://localhost:1111"},
	}

	allNodes := state.GetAllNodes()

	assert.Len(t, allNodes, 4)
	assert.Equal(t, allNodes[0], "http://localhost:1234")
	assert.Equal(t, allNodes[1], "http://localhost:5678")
	assert.Equal(t, allNodes[2], "http://localhost:9012")
	assert.Equal(t, allNodes[3], "http://localhost:1111")
}

func TestAppState_UpdateNodes(t *testing.T) {
	state := &AppState{}

	healthyNodes := []string{"http://localhost:1234", "http://localhost:5678"}
	unhealthyNodes := []string{"http://localhost:9012"}

	state.UpdateNodes(healthyNodes, unhealthyNodes)

	assert.Len(t, state.HealthyNodes, 2)
	assert.Equal(t, state.HealthyNodes[0], "http://localhost:1234")
	assert.Equal(t, state.HealthyNodes[1], "http://localhost:5678")

	assert.Len(t, state.UnhealthyNodes, 1)
	assert.Equal(t, state.UnhealthyNodes[0], "http://localhost:9012")

}

func TestAppState_AddNode(t *testing.T) {
	state := &AppState{}

	state.AddNode("http://localhost:1234")
	state.AddNode("http://localhost:5678")
	state.AddNode("http://localhost:9012")

	assert.Len(t, state.HealthyNodes, 3)
	assert.Equal(t, state.HealthyNodes[0], "http://localhost:1234")
	assert.Equal(t, state.HealthyNodes[1], "http://localhost:5678")
	assert.Equal(t, state.HealthyNodes[2], "http://localhost:9012")
}

func TestAppState_GetNextHealthyNode(t *testing.T) {
	state := &AppState{
		HealthyNodes:   []string{"http://localhost:1234", "http://localhost:5678", "http://localhost:9012"},
		UnhealthyNodes: []string{"http://localhost:1111"},
	}

	nextNode := state.GetNextHealthyNode()
	assert.Equal(t, nextNode, "http://localhost:1234")

	nextNode = state.GetNextHealthyNode()
	assert.Equal(t, nextNode, "http://localhost:5678")

	nextNode = state.GetNextHealthyNode()
	assert.Equal(t, nextNode, "http://localhost:9012")

	nextNode = state.GetNextHealthyNode()
	assert.Equal(t, nextNode, "http://localhost:1234")
}
