package jsgo

import (
	"log"
	"sync"
	"testing"

	"github.com/dop251/goja"
)

// EXT_LIB_SCIRPT

// SECURITY SCRIPT

// Below JS script enrich the user object with company name by extracting it from second part of email.
var gojaEnrichScript = `
	function f(user) {
		user.company = user.email.split('@')[1]; // goja pass go map as JS object as a reference, so that the original map gets updated
	}
`

func gojaScript(user map[string]interface{}, vms ...*goja.Runtime) {
	var vm *goja.Runtime
	if len(vms) > 0 {
		vm = vms[0]
	} else {
		vm = goja.New()
	}

	_, err := vm.RunString(gojaEnrichScript)
	if err != nil {
		log.Fatal(err)
	}

	var fn func(map[string]interface{}) map[string]interface{}
	err = vm.ExportTo(vm.Get("f"), &fn)
	if err != nil {
		log.Fatal(err)
	}
	fn(user)
}

type runtime struct {
	vm *goja.Runtime
	fn goja.Callable
}

var pool = &sync.Pool{
	New: func() interface{} {
		vm := goja.New()
		_, err := vm.RunString(gojaEnrichScript)
		if err != nil {
			log.Fatal(err)
		}

		var fn goja.Callable
		err = vm.ExportTo(vm.Get("f"), &fn)
		if err != nil {
			log.Fatal(err)
		}
		return &runtime{
			vm: vm,
			fn: fn,
		}
	},
}

func gojaScriptWithPool(user map[string]interface{}) {
	vm := pool.Get().(*runtime)
	defer pool.Put(vm)
	u := vm.vm.ToValue(user)
	_, err := vm.fn(goja.Undefined(), u)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGojaTransformScript(t *testing.T) {
	user := createUser()

	gojaScript(user)

	if user["company"] != "company.com" {
		t.Error("expected u.company")
	}
}

func BenchmarkNoReuseGojaVm(b *testing.B) {
	user := createUser()
	for i := 0; i < b.N; i++ {
		gojaScript(user)
	}
}

func BenchmarkReuseGojaVm(b *testing.B) {
	user := createUser()
	vm := goja.New()
	for i := 0; i < b.N; i++ {
		gojaScript(user, vm)
	}
}

// Is it goroutine-safe?
// No. An instance of goja.Runtime can only be used by a single goroutine at a time. You can create as many instances of Runtime as
// you like but it's not possible to pass object values between runtimes.

func TestGoja(t *testing.T) {
	u := createUser()
	gojaScriptWithPool(u)
	if u["company"] != "company.com" {
		t.Error("expected u.company")
	}
}

func BenchmarkParallelGojaVmPool(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			u := createUser()
			gojaScriptWithPool(u)
		}
	})
}
