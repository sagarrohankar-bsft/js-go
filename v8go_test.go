package jsgo

import (
	"fmt"
	"log"
	"testing"

	"rogchap.com/v8go"
)

var v8goUserEnrichScirpt = `
	const f = (u) => {
		u.company = u.email.split('@')[1];
		return u
	}
	result = f(%+v);
	result
`

var script = fmt.Sprintf(v8goUserEnrichScirpt, userjson)

func v8goScript(isolates ...*v8go.Isolate) string {
	var isolate *v8go.Isolate
	if len(isolates) > 0 {
		isolate = isolates[0]
	}
	ctx := v8go.NewContext(isolate)
	if isolate == nil {
		// When we pass isolate as nil, then NewContext function create the new one internally that need to be explicitly disposed of
		defer ctx.Isolate().Dispose()
	}
	defer ctx.Close()
	output, err := ctx.RunScript(script, "function.js")
	if err != nil {
		log.Fatal(err)
	}
	json, err := output.MarshalJSON()
	if err != nil {
		log.Fatal(err)
	}
	return string(json)
}

func TestV8GoTransformScript(t *testing.T) {
	v8goScript()
}

func BenchmarkNoReuseV8GoVm(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = v8goScript()
	}
}

func BenchmarkReuseV8GoVm(b *testing.B) {
	vm := v8go.NewIsolate()
	defer vm.Dispose()
	for i := 0; i < b.N; i++ {
		_ = v8goScript(vm)
	}
}

func BenchmarkParallelNoReuseV8GoVm(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_ = v8goScript()
		}
	})
}

func BenchmarkParallelReuseV8GoVm(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		vm := v8go.NewIsolate()
		defer vm.Dispose()
		for p.Next() {
			_ = v8goScript(vm)
		}
	})
}
