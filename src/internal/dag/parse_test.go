// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package dag

import (
	"slices"
	"strings"
	"testing"
)

const diamond = `
NONE < a < b, c < d;
`

func mustParse(t *testing.T, dag string) *Graph {
	t.Helper()
	g, err := Parse(dag)
	if err != nil {
		t.Fatal(err)
	}
	return g
}

func wantEdges(t *testing.T, g *Graph, edges string) {
	t.Helper()

	wantEdges := strings.Fields(edges)
	wantEdgeMap := make(map[string]bool)
	for _, e := range wantEdges {
		wantEdgeMap[e] = true
	}

	for _, n1 := range g.Nodes {
		for _, n2 := range g.Nodes {
			golangt := g.HasEdge(n1, n2)
			want := wantEdgeMap[n1+"->"+n2]
			if golangt && want {
				t.Logf("%s->%s", n1, n2)
			} else if golangt && !want {
				t.Errorf("%s->%s present but not expected", n1, n2)
			} else if want && !golangt {
				t.Errorf("%s->%s missing but expected", n1, n2)
			}
		}
	}
}

func TestParse(t *testing.T) {
	// Basic smoke test for graph parsing.
	g := mustParse(t, diamond)

	wantNodes := strings.Fields("a b c d")
	if !slices.Equal(wantNodes, g.Nodes) {
		t.Fatalf("want nodes %v, golangt %v", wantNodes, g.Nodes)
	}

	// Parse returns the transitive closure, so it adds d->a.
	wantEdges(t, g, "b->a c->a d->a d->b d->c")
}
