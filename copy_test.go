package jsonstruct_test

import (
	"github.com/myshkin5/jsonstruct"

	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Copy", func() {
	FIt("does a deep copy with no shared mutables", func() {
		orig := jsonstruct.New()
		orig.SetInt("int-val", 737)
		orig.SetString("sub.sub.value", "neat string")

		copy := orig.DeepCopy()

		Expect(reflect.ValueOf(orig)).NotTo(Equal(reflect.ValueOf(copy)))

		Expect(orig["int-val"]).To(Equal(copy["int-val"]))

		origSub := orig["sub"].(map[string]interface{})
		copySub, ok := copy["sub"].(map[string]interface{})
		Expect(ok).To(BeTrue())
		Expect(reflect.ValueOf(origSub)).NotTo(Equal(reflect.ValueOf(copySub)))

		origSubSub := origSub["sub"]
		copySubSub := copySub["sub"]
		Expect(reflect.ValueOf(origSubSub)).NotTo(Equal(reflect.ValueOf(copySubSub)))
	})
})
