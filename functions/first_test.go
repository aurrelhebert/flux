package functions_test

import (
	"testing"

	"github.com/influxdata/flux"
	"github.com/influxdata/flux/execute"
	"github.com/influxdata/flux/execute/executetest"
	"github.com/influxdata/flux/functions"
	"github.com/influxdata/flux/plan"
	"github.com/influxdata/flux/plan/plantest"
	"github.com/influxdata/flux/querytest"
)

func TestFirstOperation_Marshaling(t *testing.T) {
	data := []byte(`{"id":"first","kind":"first","spec":{"column":"foo"}}`)
	op := &flux.Operation{
		ID: "first",
		Spec: &functions.FirstOpSpec{
			SelectorConfig: execute.SelectorConfig{
				Column: "foo",
			},
		},
	}

	querytest.OperationMarshalingTestHelper(t, data, op)
}

func TestFirst_Process(t *testing.T) {
	testCases := []struct {
		name string
		data *executetest.Table
		want [][]int
	}{
		{
			name: "first",
			data: &executetest.Table{
				ColMeta: []flux.ColMeta{
					{Label: "_time", Type: flux.TTime},
					{Label: "_value", Type: flux.TFloat},
					{Label: "t1", Type: flux.TString},
					{Label: "t2", Type: flux.TString},
				},
				Data: [][]interface{}{
					{execute.Time(0), 0.0, "a", "y"},
					{execute.Time(10), 5.0, "a", "x"},
					{execute.Time(20), 9.0, "a", "y"},
					{execute.Time(30), 4.0, "a", "x"},
					{execute.Time(40), 6.0, "a", "y"},
					{execute.Time(50), 8.0, "a", "x"},
					{execute.Time(60), 1.0, "a", "y"},
					{execute.Time(70), 2.0, "a", "x"},
					{execute.Time(80), 3.0, "a", "y"},
					{execute.Time(90), 7.0, "a", "x"},
				},
			},
			want: [][]int{{0}, nil, nil, nil, nil, nil, nil, nil, nil, nil},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			executetest.IndexSelectorFuncTestHelper(
				t,
				new(functions.FirstSelector),
				tc.data,
				tc.want,
			)
		})
	}
}

func BenchmarkFirst(b *testing.B) {
	executetest.IndexSelectorFuncBenchmarkHelper(b, new(functions.FirstSelector), NormalTable)
}

func TestFirst_PushDown_Match(t *testing.T) {
	spec := new(functions.FirstProcedureSpec)
	from := new(functions.FromProcedureSpec)

	// Should not match when an aggregate is set
	from.AggregateSet = true
	plantest.PhysicalPlan_PushDown_Match_TestHelper(t, spec, from, []bool{false})

	// Should match when no aggregate is set
	from.AggregateSet = false
	plantest.PhysicalPlan_PushDown_Match_TestHelper(t, spec, from, []bool{true})
}

func TestFirst_PushDown(t *testing.T) {
	spec := new(functions.FirstProcedureSpec)
	root := &plan.Procedure{
		Spec: new(functions.FromProcedureSpec),
	}
	want := &plan.Procedure{
		Spec: &functions.FromProcedureSpec{
			BoundsSet: true,
			Bounds: flux.Bounds{
				Start: flux.MinTime,
				Stop:  flux.Now,
			},
			LimitSet:      true,
			PointsLimit:   1,
			DescendingSet: true,
			Descending:    false,
		},
	}

	plantest.PhysicalPlan_PushDown_TestHelper(t, spec, root, false, want)
}
func TestFirst_PushDown_Duplicate(t *testing.T) {
	spec := new(functions.FirstProcedureSpec)
	root := &plan.Procedure{
		Spec: &functions.FromProcedureSpec{
			BoundsSet: true,
			Bounds: flux.Bounds{
				Start: flux.MinTime,
				Stop:  flux.Now,
			},
			LimitSet:      true,
			PointsLimit:   1,
			DescendingSet: true,
			Descending:    false,
		},
	}
	want := &plan.Procedure{
		// Expect the duplicate has been reset to zero values
		Spec: new(functions.FromProcedureSpec),
	}

	plantest.PhysicalPlan_PushDown_TestHelper(t, spec, root, true, want)
}
