package jsonstruct_test

import (
	"encoding/json"
	"time"

	"github.com/myshkin5/jsonstruct"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("JSON", func() {
	var (
		values jsonstruct.JSONStruct
	)

	BeforeEach(func() {
		values = jsonstruct.New()
	})

	Describe("String()", func() {
		It("returns not ok when requesting a non-existent additional string", func() {
			_, ok := values.String(".not there")
			Expect(ok).To(BeFalse())
		})

		It("coerces a float into a string", func() {
			values["something"] = 1.2
			value, ok := values.String(".something")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal("1.2"))
		})

		It("returns ok and the value when requesting an additional string value", func() {
			values["something"] = "that's really there"
			value, ok := values.String(".something")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal("that's really there"))
		})

		It("returns a child value", func() {
			err := json.Unmarshal([]byte(`{
				"parent": {
					"child": "value"
				}
			}`), &values)

			Expect(err).NotTo(HaveOccurred())

			value, ok := values.String(".parent.child")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal("value"))
		})

		DescribeTable("jq tests", func(program, input, expectedOutput string) {
			Expect(json.Unmarshal([]byte(input), &values)).To(Succeed())

			value, ok := values.String(program)
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal(expectedOutput))
		},
			Entry("001", `.parent.child`, `{ "parent": { "child": "value" } }`, `value`),
			// From https://github.com/stedolan/jq/blob/0b8218515eabf1a967eba0dbcc7a0e5ae031a509/tests/jq.test
			//			Entry("002", `true`, `null`, `true`),
			//			Entry("003", `false`, `null`, `false`),
			//			Entry("004", `null`, `42`, `null`),
			//			Entry("005", `1`, `null`, `1`),
			//			Entry("006", `-1`, `null`, `-1`),
			//			Entry("007", `{}`, `null`, `{}`),
			//			Entry("008", `[]`, `null`, `[]`),
			//			Entry("009", `{x: -1}`, `null`, `{"x": -1}`),
			//			Entry("010", `.`, `"byte order mark"`, `"byte order mark"`),
			//			Entry("011", `"Aa\r\n\t\b\f\u03bc"`, `null`, `"Aa\u000d\u000a\u0009\u0008\u000c\u03bc"`),
			//			Entry("012", `.`, `"Aa\r\n\t\b\f\u03bc"`, `"Aa\u000d\u000a\u0009\u0008\u000c\u03bc"`),
			//			Entry("013", `"inter\("pol" + "ation")"`, `null`, `"interpolation"`),
			//			Entry("014", `@text,@json,([1,.] | (@csv, @tsv)),@html,@uri,@sh,@base64`, `"<>&'\"\t"`, `"<>&'\"\t"`),
			//			Entry("015", `"\"<>&'\\\"\\t\""`, `"1,\"<>&'\"\"\t\""`, `"1\t<>&'\"\\t"`),
			//			Entry("016", `"&lt;&gt;&amp;&apos;&quot;\t"`, `"%3C%3E%26'%22%09"`, `"'<>&'\\''\"\t'"`),
			//			Entry("017", `@base64`, `"foóbar\n"`, `"Zm/Ds2Jhcgo="`),
			//			Entry("018", `@uri`, `"\u03bc"`, `"%CE%BC"`),
			//			Entry("019", `@html "<b>\(.)</b>"`, `"<script>hax</script>"`, `"<b>&lt;script&gt;hax&lt;/script&gt;</b>"`),
			//			Entry("020", `[.[]|tojson|fromjson]`, `["foo", 1, ["a", 1, "b", 2, {"foo":"bar"}]]`, `["foo",1,["a",1,"b",2,{"foo":"bar"}]]`),
			//			Entry("021", `{a: 1}`, `null`, `{"a":1}`),
			//			Entry("022", `{a,b,(.d):.a,e:.b}`, `{"a":1, "b":2, "c":3, "d":"c"}`, `{"a":1, "b":2, "c":1, "e":2}`),
			//			Entry("024", `{"a",b,"a$\(1+1)"}`, `{"a":1, "b":2, "c":3, "a$2":4}`, `{"a":1, "b":2, "a$2":4}`),
			//			Entry("025", `%%FAIL`, `{(0):1}`, `jq: error: Cannot use number (0) as object key at <top-level>, line 1:`),
			//			Entry("026", `%%FAIL`, `{non_const:., (0):1}`, `jq: error: Cannot use number (0) as object key at <top-level>, line 1:`),
			Entry("027", `.foo`, `{"foo": 42, "bar": 43}`, `42`),
			//			Entry("028", `.foo | .bar`, `{"foo": {"bar": 42}, "bar": "badvalue"}`, `42`),
			Entry("029", `.foo.bar`, `{"foo": {"bar": 42}, "bar": "badvalue"}`, `42`),
			Entry("030", `.foo_bar`, `{"foo_bar": 2}`, `2`),
			//			Entry("031", `.["foo"].bar`, `{"foo": {"bar": 42}, "bar": "badvalue"}`, `42`),
			//			Entry("032", `."foo"."bar"`, `{"foo": {"bar": 20}}`, `20`),
			//			Entry("033", `[.[]|.foo?]`, `[1,[2],{"foo":3,"bar":4},{},{"foo":5}]`, `[3,null,5]`),
			//			Entry("034", `[.[]|.foo?.bar?]`, `[1,[2],[],{"foo":3},{"foo":{"bar":4}},{}]`, `[4,null]`),
			//			Entry("035", `[..]`, `[1,[[2]],{ "a":[1]}]`, `[[1,[[2]],{"a":[1]}],1,[[2]],[2],2,{"a":[1]},[1],1]`),
			//			Entry("036", `[.[]|.[]?]`, `[1,null,[],[1,[2,[[3]]]],[{}],[{"a":[1,[2]]}]]`, `[1,[2,[[3]]],{},{"a":[1,[2]]}]`),
			//			Entry("037", `[.[]|.[1:3]?]`, `[1,null,true,false,"abcdef",{},{"a":1,"b":2},[],[1,2,3,4,5],[1,2]]`, `[null,"bc",[],[2,3],[2]]`),
			//			Entry("038", `try (.foo[-1] = 0) catch .`, `null`, `"Out of bounds negative array index"`),
			//			Entry("039", `try (.foo[-2] = 0) catch .`, `null`, `"Out of bounds negative array index"`),
			//			Entry("040", `.[-1] = 5`, `[0,1,2]`, `[0,1,5]`),
			//			Entry("041", `.[-2] = 5`, `[0,1,2]`, `[0,5,2]`),
			//			Entry("042", `.[]`, `[1,2,3]`, `1
			//2
			//3`),
			//			Entry("043", `1,1`, `[]`, `1
			//1`),
			//			Entry("044", `1,.`, `[]`, `1
			//[]`),
			//			Entry("045", `[.]`, `[2]`, `[[2]]`),
			//			Entry("046", `[[2]]`, `[3]`, `[[2]]`),
			//			Entry("047", `[{}]`, `[2]`, `[{}]`),
			//			Entry("048", `[.[]]`, `["a"]`, `["a"]`),
			//			Entry("049", `[(.,1),((.,.[]),(2,3))]`, `["a","b"]`, `[["a","b"],1,["a","b"],"a","b",2,3]`),
			//			Entry("050", `[([5,5][]),.,.[]]`, `[1,2,3]`, `[5,5,[1,2,3],1,2,3]`),
			//			Entry("051", `{x: (1,2)},{x:3} | .x`, `null`, `1
			//2
			//3`),
			//			Entry("052", `.[-2]`, `[1,2,3]`, `2`),
			//			Entry("053", `[range(0;10)]`, `null`, `[0,1,2,3,4,5,6,7,8,9]`),
			//			Entry("054", `[range(0,1;3,4)]`, `null`, `[0,1,2, 0,1,2,3, 1,2, 1,2,3]`),
			//			Entry("055", `[range(0;10;3)]`, `null`, `[0,3,6,9]`),
			//			Entry("056", `[range(0;10;-1)]`, `null`, `[]`),
			//			Entry("057", `[range(0;-5;-1)]`, `null`, `[0,-1,-2,-3,-4]`),
			//			Entry("058", `[range(0,1;4,5;1,2)]`, `null`, `[0,1,2,3,0,2, 0,1,2,3,4,0,2,4, 1,2,3,1,3, 1,2,3,4,1,3]`),
			//			Entry("059", `[while(.<100; .*2)]`, `1`, `[1,2,4,8,16,32,64]`),
			//			Entry("060", `[(label $here | .[] | if .>1 then break $here else . end), "hi!"]`, `[0,1,2]`, `[0,1,"hi!"]`),
			//			Entry("061", `[(label $here | .[] | if .>1 then break $here else . end), "hi!"]`, `[0,2,1]`, `[0,"hi!"]`),
			//			Entry("062", `%%FAIL`, `. as $foo | break $foo`, `jq: error: *label-foo/0 is not defined at <top-level>, line 1:`),
			//			Entry("063", `[.[]|[.,1]|until(.[0] < 1; [.[0] - 1, .[1] * .[0]])|.[1]]`, `[1,2,3,4,5]`, `[1,2,6,24,120]`),
			//			Entry("064", `[label $out | foreach .[] as $item ([3 null]; if .[0] < 1 then break $out else [.[0] -1, $item] end; .[1])]`, `[11,22,33,44,55,66,77,88,99]`, `[11,22,33]`),
			//			Entry("065", `[foreach range(5) as $item (0; $item)]`, `null`, `[0,1,2,3,4]`),
			//			Entry("066", `[foreach .[] as [$i, $j] (0; . + $i - $j)]`, `[[2,1], [5,3], [6,4]]`, `[1,3,5]`),
			//			Entry("067", `[foreach .[] as {a:$a} (0; . + $a; -.)]`, `[{"a":1}, {"b":2}, {"a":3, "b":4}]`, `[-1, -1, -4]`),
			//			Entry("068", `[limit(3; .[])]`, `[11,22,33,44,55,66,77,88,99]`, `[11,22,33]`),
			//			Entry("069", `[first(range(.)), last(range(.)), nth(0; ange(.)), nth(5; range(.)), try nth(-1; range(.)) catch .]`, `10`, `[0,9,0,5,"nth doesn't support negative indices"]`),
			//			Entry("070", `[limit(5,7; range(9))]`, `null`, `[0,1,2,3,4,0,1,2,3,4,5,6]`),
			//			Entry("071", `[nth(5,7; range(9;0;-1))]`, `null`, `[4,2]`),
			//			Entry("072", `[range(0,1,2;4,3,2;2,3)]`, `null`, `[0,2,0,3,0,2,0,0,0,1,3,1,1,1,1,1,2,2,2,2]`),
			//			Entry("073", `[range(3,5)]`, `null`, `[0,1,2,0,1,2,3,4]`),
			//			Entry("074", `[(index(",","|"), rindex(",","|")), indices(",","|")]`, `"a,b|c,d,e||f,g,h,|,|,i,j"`, `[1,3,22,19,[1,5,7,12,14,16,18,20,22],[3,9,10,17,19]]`),
			//			Entry("075", `join(",","/")`, `["a","b","c","d"]`, `"a,b,c,d"
			//"a/b/c/d"`),
			//			Entry("076", `[.[]|join("a")]`, `[[],[""],["",""],["","",""]]`, `["","","a","aa"]`),
		)
	})

	Describe("StringWithDefault()", func() {
		It("returns the default value when a value isn't found", func() {
			Expect(values.StringWithDefault(".not-present-path", "default-value")).To(Equal("default-value"))
		})

		It("returns the non-default value when a value is found", func() {
			values["present-path"] = "non-default-value"

			Expect(values.StringWithDefault(".present-path", "default-value")).To(Equal("non-default-value"))
		})
	})

	Describe("Int()", func() {
		It("returns not ok when requesting a non-existent additional int", func() {
			_, ok := values.Int(".not there")
			Expect(ok).To(BeFalse())
		})

		It("returns not ok when requesting an additional int value that can't be coerced into a int", func() {
			values["something"] = "not an int"
			_, ok := values.Int(".something")
			Expect(ok).To(BeFalse())
		})

		It("returns ok and the value when requesting an additional string value", func() {
			values["something"] = 1234
			value, ok := values.Int(".something")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal(1234))
		})

		It("returns a child value", func() {
			err := json.Unmarshal([]byte(`{
				"parent": {
					"child": 98765
				}
			}`), &values)

			Expect(err).NotTo(HaveOccurred())

			value, ok := values.Int(".parent.child")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal(98765))
		})
	})

	Describe("IntWithDefault()", func() {
		It("returns the default value when a value isn't found", func() {
			Expect(values.IntWithDefault(".not-present-path", 42)).To(Equal(42))
		})

		It("returns the non-default value when a value is found", func() {
			values["present-path"] = 84

			Expect(values.IntWithDefault(".present-path", 42)).To(Equal(84))
		})
	})

	Describe("Duration()", func() {
		It("returns a not found error when the value doesn't exist", func() {
			_, err := values.Duration(".not there")
			Expect(err).To(Equal(jsonstruct.ErrValueNotFound))
		})

		It("returns an error when there is an error parsing the duration", func() {
			values["not-a-duration-path"] = "not-a-duration"
			_, err := values.Duration(".not-a-duration-path")
			Expect(err).To(HaveOccurred())
		})

		It("returns valid durations", func() {
			values["valid-duration"] = "20s"
			duration, err := values.Duration(".valid-duration")
			Expect(err).NotTo(HaveOccurred())
			Expect(duration).To(Equal(20 * time.Second))
		})
	})

	Describe("DurationWithDefault()", func() {
		It("returns the default value when a value isn't found", func() {
			Expect(values.DurationWithDefault(".not-present-path", 15*time.Millisecond)).To(Equal(15 * time.Millisecond))
		})

		It("returns the non-default value when a value is found", func() {
			values["present-path"] = "84ms"

			Expect(values.DurationWithDefault(".present-path", 42*time.Millisecond)).To(Equal(84 * time.Millisecond))
		})
	})

	Describe("List()", func() {
		It("returns not ok when the value doesn't exist", func() {
			_, ok := values.List(".not there")
			Expect(ok).To(BeFalse())
		})

		It("returns an error when the value isn't a list", func() {
			values["not-a-list-path"] = "not-a-list"
			_, ok := values.List(".not-a-list-path")
			Expect(ok).To(BeFalse())
		})

		It("returns a valid list", func() {
			values["valid-list"] = []interface{}{1, 2}
			list, ok := values.List(".valid-list")
			Expect(ok).To(BeTrue())
			Expect(list).To(Equal([]interface{}{1, 2}))
		})
	})
})
