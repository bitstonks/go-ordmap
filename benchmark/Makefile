TEST=go test -bench . -benchmem -count=3

generics:
	$(TEST) avl.go bench_test.go generics_builder.go

codegen:
	$(TEST) bench_test.go *map.go codegen_builder.go

static_generics:
	$(TEST) avl.go bench_static_test.go generics_builder.go

static_codegen:
	$(TEST) bench_static_test.go *map.go codegen_builder.go
