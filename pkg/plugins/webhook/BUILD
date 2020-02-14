load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "config.go",
        "plugin.go",
        "pluginmap_gen.go",
        "spec.go",
    ],
    importpath = "github.com/cigarframework/reconciled/pkg/plugins/webhook",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/api:go_default_library",
        "//pkg/common:go_default_library",
        "//pkg/optional:go_default_library",
        "//pkg/plugins/grpc:go_default_library",
        "//pkg/plugins/http:go_default_library",
        "//pkg/storage:go_default_library",
        "@org_uber_go_zap//:go_default_library",
    ],
)
