load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

# gazelle:resolve go github.com/cigarframework/reconciled/pkg/proto @reconciled//pkg/proto:go_default_library
go_library(
    name = "go_default_library",
    srcs = [
        "common.go",
        "json.go",
        "kinds.go",
        "protobuf.go",
        "state.go",
        "statemap_gen.go",
        "storge.go",
        "test_utils.go",
        "validate.go",
    ],
    importpath = "github.com/cigarframework/reconciled/pkg/storage",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/proto:go_default_library",
        "@com_github_gogo_protobuf//types:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "common_test.go",
        "json_test.go",
        "kinds_test.go",
        "protobuf_test.go",
        "state_test.go",
        "storage_test.go",
        "validate_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//pkg/proto:go_default_library",
        "@com_github_stretchr_testify//assert:go_default_library",
        "@com_github_tidwall_gjson//:go_default_library",
    ],
)
