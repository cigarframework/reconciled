load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["plugins.go"],
    importpath = "github.com/cigarframework/reconciled/pkg/plugins",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/api:go_default_library",
        "//pkg/plugins/authn:go_default_library",
        "//pkg/plugins/authz:go_default_library",
        "//pkg/plugins/grpc:go_default_library",
        "//pkg/plugins/http:go_default_library",
        "//pkg/plugins/webhook:go_default_library",
        "@org_uber_go_zap//:go_default_library",
    ],
)
