package resourcepool

import (
	"github.com/stretchr/testify/require"
	"testing"

	k8sCore "k8s.io/api/core/v1"

	machineTypeV1 "github.com/Netflix/titus-controllers-api/api/machinetype/v1"
	"github.com/Netflix/titus-resource-pool/machine"
	"github.com/Netflix/titus-resource-pool/node"
)

func TestKubeletNodesAreExcluded(t *testing.T) {
	pool := EmptyResourcePool()
	nodes := []*k8sCore.Node{
		node.NewNode("node1", pool.Name, machine.R5Metal()),
		node.ButNodeLabel(node.NewNode("node2", pool.Name, machine.R5Metal()),
			node.NodeLabelBackend, node.NodeBackendKubelet),
	}

	// With kubelet
	snapshot := NewStaticResourceSnapshot(pool, []*machineTypeV1.MachineTypeConfig{}, nodes, []*k8sCore.Pod{},
		0, true)
	require.Equal(t, 2, len(snapshot.Nodes))
	require.Equal(t, 0, len(snapshot.ExcludedNodes))

	// Without kubelet
	snapshot = NewStaticResourceSnapshot(pool, []*machineTypeV1.MachineTypeConfig{}, nodes, []*k8sCore.Pod{},
		0, false)
	require.Equal(t, 1, len(snapshot.Nodes))
	require.Equal(t, 1, len(snapshot.ExcludedNodes))
}
