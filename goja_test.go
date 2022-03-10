package jsgo

import (
	"fmt"
	"log"
	"testing"

	"github.com/dop251/goja"
)

// EXT_LIB_SCIRPT

// SECURITY SCRIPT

// Below JS script enrich the user object with company name by extracting it from second part of email.
var gojaEnrichScript = `
	function f(user) {
		user.company = user.email.split('@')[1];
		return user;
	}
`

func gojaScript(user map[string]interface{}, vms ...*goja.Runtime) map[string]interface{} {
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
	return fn(user)
}

func TestGojaTransformScript(t *testing.T) {
	user := createUser()

	user = gojaScript(user)

	fmt.Println(user)
}

func BenchmarkNoReuseGojaVm(b *testing.B) {
	user := createUser()
	for i := 0; i < b.N; i++ {
		_ = gojaScript(user)
	}
}

func BenchmarkReuseGojaVm(b *testing.B) {
	user := createUser()
	vm := goja.New()
	for i := 0; i < b.N; i++ {
		_ = gojaScript(user, vm)
	}
}

// Is it goroutine-safe?
// No. An instance of goja.Runtime can only be used by a single goroutine at a time. You can create as many instances of Runtime as
// you like but it's not possible to pass object values between runtimes.

func BenchmarkParallelNoReuseGojaVm(b *testing.B) {
	user := createUser()
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_ = gojaScript(user)
		}
	})
}

func BenchmarkParallelReuseGojaVm(b *testing.B) {
	user := createUser()
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		vm := goja.New()
		for p.Next() {
			_ = gojaScript(user, vm)
		}
	})
}
