package space

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"

	"github.com/pspaces/gospace/container"
	"github.com/pspaces/gospace/function"
	"github.com/pspaces/gospace/policy"
)

// Coin picks a tuple randomly by a coin toss or creates an empty tuple.
func Coin(ts ...container.Intertuple) (r container.Intertuple) {
	upper := big.NewInt(int64(len(ts)))

	if len(ts) > 0 {
		coin, _ := rand.Int(rand.Reader, upper)
		r = ts[int(coin.Int64())]
	} else {
		t := container.NewTuple()
		r = &t
	}

	return r
}

// TemplateIdentity preservers the template.
func TemplateIdentity(i interface{}) (tp container.Template) {
	tpf := i.([]interface{})

	tp = container.NewTemplate(tpf...)

	return tp
}

// TupleIdentity preserves the tuple.
func TupleIdentity(i interface{}) (it container.Intertuple) {
	tf := i.([]interface{})

	t := container.NewTuple(tf...)
	it = &t

	return it
}

// ABACTemplates generates an policy by producing the Cartesian product between a set of functions and a multiset of templates.
// ABACTemplates adds aggregation policies to an existing composable policy and will not overwrite existing policies.
// Functions are operations related to a tuple space, but in principle could be more complicated.
// In ABACTemplates, templates, matched tuples and resulting aggregated tuple are all preserved.
// ABACTemplates returns all the generated label sets containing a single label and if the addition of all policies were successful.
func ABACTemplates(pol *policy.Composable, funs []interface{}, tps []container.Template) (slbls []container.Labels, b bool) {
	b = pol != nil && len(funs) > 0 && len(tps) > 0

	if b {
		polName := function.Name(ABACTemplates)
		slbls = make([]container.Labels, 0, len(funs)*len(tps))
		for _, fun := range funs {
			funName := function.Name(fun)
			for i, tp := range tps {
				sl := fmt.Sprintf("%s-%s-%d", polName, funName, i)
				l := container.NewLabel(sl)
				lm := container.NewLabels(l)
				slbls = append(slbls, lm)
				ltf := tp.Fields()
				//ltp := container.NewTemplate(ltf...)

				a := policy.NewAction(fun, ltf...)
				templateTrans := policy.NewTransformation(TemplateIdentity)
				tupleTrans := policy.NewTransformation(TupleIdentity)
				resultTrans := policy.NewTransformation(TupleIdentity)
				transformations := policy.NewTransformations(&templateTrans, &tupleTrans, &resultTrans)
				rule := policy.NewAggregationRule(*a, *transformations)
				aggPol := policy.NewAggregation(l, rule)
				b = b && pol.Add(aggPol)
			}
		}
	}

	return slbls, b
}

// ABACLabelledTemplates generates an policy by producing the Cartesian product between a set of functions and a multiset of templates.
// ABACLabelledTemplates adds aggregation policies to an existing composable policy and will not overwrite existing policies.
// Functions are operations related to a tuple space, but in principle could be more complicated.
// In ABACLabelledTemplates, templates, matched tuples and resulting aggregated tuple are preserved.
// ABACLabelledTemplates injects labels into the templates such that the templates become labelled.
// ABACLabelledTemplates returns all the generated label sets containing a single label and if the addition of all policies were successful.
func ABACLabelledTemplates(pol *policy.Composable, funs []interface{}, tps []container.Template) (slbls []container.Labels, b bool) {
	b = pol != nil && len(funs) > 0 && len(tps) > 0

	if b {
		spc := new(Space)
		polName := function.Name(ABACLabelledTemplates)
		slbls = make([]container.Labels, 0, len(funs)*len(tps))
		for _, fun := range funs {
			funName := function.Name(fun)
			for i, tp := range tps {
				sl := fmt.Sprintf("%s-%s-%d", polName, funName, i)
				l := container.NewLabel(sl)
				lm := container.NewLabels(l)
				slbls = append(slbls, lm)

				tf := tp.Fields()
				var ltf []interface{}
				if container.NewSignature(1, fun) == container.NewSignature(1, spc.PutAgg) ||
					container.NewSignature(1, fun) == container.NewSignature(1, spc.GetAgg) ||
					container.NewSignature(1, fun) == container.NewSignature(1, spc.QueryAgg) {
					ltf = make([]interface{}, len(tf)+1)
					copy(ltf[:2], []interface{}{tf[0], lm})
					copy(ltf[2:], tf[1:])
				} else {
					ltf = make([]interface{}, len(tf)+1)
					copy(ltf[:1], []interface{}{l})
					copy(ltf[1:], tf)
				}

				//ltp := container.NewTemplate(ltf...)
				a := policy.NewAction(fun, ltf...)
				templateTrans := policy.NewTransformation(TemplateIdentity)
				tupleTrans := policy.NewTransformation(TupleIdentity)
				resultTrans := policy.NewTransformation(TupleIdentity)
				transformations := policy.NewTransformations(&templateTrans, &tupleTrans, &resultTrans)
				rule := policy.NewAggregationRule(*a, *transformations)
				aggPol := policy.NewAggregation(l, rule)
				b = b && pol.Add(aggPol)
			}
		}
	}

	return slbls, b
}

// BenchmarkRandomDisclosurePolicy randomly picks a tuple from the matched tuples returns it according to a policy.
// Depending on operation, multiple tuples may be consumed and aggregated.
// Execute with: go test -bench=RandomDisclosurePolicy -run=none -gcflags "-N -l" -cpuprofile=cprof
func BenchmarkRandomDisclosurePolicy(b *testing.B) {
	var i int
	var f float64
	var s string

	fspc := new(Space)

	disclose := policy.NewComposable()
	funs := []interface{}{fspc.PutAgg, fspc.GetAgg, fspc.QueryAgg}
	temps := []container.Template{container.NewTemplate(Coin, &s, &i, &f)}
	slabels, _ := ABACLabelledTemplates(disclose, funs, temps)

	// Parameters for considering tuple space size from 1 to 1250 tuples.
	m := 5
	isize := []int{1, 250, 500, 750, 1000, 1250, 1500, 1750, 2000}

	spc := NewSpace("tcp4://localhost:31415/spc?CONN")
	aspc := NewSpace("tcp4://localhost:31416/aspc?CONN", disclose)

	b.ResetTimer()

	b.Run("Constructor-ComposablePolicy", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			disclose = policy.NewComposable()
		}
	})

	funs = []interface{}{fspc.PutAgg, fspc.GetAgg, fspc.QueryAgg}
	temps = []container.Template{container.NewTemplate(Coin, &s, &i, &f)}

	b.Run("ABACLabelledTemplates-RandomDisclosure", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			b.StopTimer()
			disclose = policy.NewComposable()
			b.StartTimer()
			slabels, _ = ABACLabelledTemplates(disclose, funs, temps)
			b.StopTimer()
		}
	})

	OperationRandomDisclosureLabelled := func(n int, s *Space, op func(interface{}, ...interface{}) (container.Tuple, error), aggf func(...container.Intertuple) container.Intertuple, plbls container.Labels, b *testing.B) func(*testing.B) {
		return func(b *testing.B) {
			var i int
			var f float64
			var str string

			b.ResetTimer()
			b.StopTimer()
			for k := 0; k < b.N; k++ {
				// Initialize the tuple space (excluded from timing).
				for l := 0; l < n; l++ {
					s.Put(plbls, "1", l, 3.14)
				}
				b.StartTimer()

				// Benchmark the operation.
				op(aggf, plbls, &str, &i, &f)

				// Clean-up the tuple space (excluded from timing).
				b.StopTimer()
				s.GetAll(plbls, &str, &i, &f)
			}
		}
	}

	// For all possible aggregate operators and all possible input sizes, run the test.
	b.Run("ExcludingPolicy", func(b *testing.B) {
		for _, sz := range isize[0:m] {
			// Labels allowing GetAgg.
			lbls := slabels[0]
			benchOperation := OperationRandomDisclosureLabelled(sz, &spc, spc.PutAgg, Coin, lbls, b)
			b.ResetTimer()
			b.Run("PutAgg-"+fmt.Sprintf("%d-Tuple", sz), benchOperation)

			// Labels allowing GetAgg.
			lbls = slabels[1]
			benchOperation = OperationRandomDisclosureLabelled(sz, &spc, spc.GetAgg, Coin, lbls, b)
			b.ResetTimer()
			b.Run("GetAgg-"+fmt.Sprintf("%d-Tuple", sz), benchOperation)

			// Labels allowing QueryAgg.
			lbls = slabels[2]
			benchOperation = OperationRandomDisclosureLabelled(sz, &spc, spc.QueryAgg, Coin, lbls, b)
			b.ResetTimer()
			b.Run("QueryAgg-"+fmt.Sprintf("%d-Tuple", sz), benchOperation)
		}
	})

	// For all possible aggregate operators and all possible input sizes, run the test.
	b.Run("IncludingPolicy", func(b *testing.B) {
		for _, sz := range isize[0:m] {
			// Labels allowing GetAgg.
			lbls := slabels[0]
			benchOperation := OperationRandomDisclosureLabelled(sz, &aspc, aspc.PutAgg, Coin, lbls, b)
			b.ResetTimer()
			b.Run("PutAgg-"+fmt.Sprintf("%d-Tuple", sz), benchOperation)

			// Labels allowing GetAgg.
			lbls = slabels[1]
			benchOperation = OperationRandomDisclosureLabelled(sz, &aspc, aspc.GetAgg, Coin, lbls, b)
			b.ResetTimer()
			b.Run("GetAgg-"+fmt.Sprintf("%d-Tuple", sz), benchOperation)

			// Labels allowing QueryAgg.
			lbls = slabels[2]
			benchOperation = OperationRandomDisclosureLabelled(sz, &aspc, aspc.QueryAgg, Coin, lbls, b)
			b.ResetTimer()
			b.Run("QueryAgg-"+fmt.Sprintf("%d-Tuple", sz), benchOperation)
		}
	})

	// Both spaces should not have anything left in them.
	sz, _ := spc.Size()
	fmt.Printf("Size of spc: %d\n", sz)

	asz, _ := aspc.Size()
	fmt.Printf("Size of aspc: %d\n", asz)

	return
}
