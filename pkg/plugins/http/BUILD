load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "config.go",
        "plugin.go",
    ],
    importpath = "github.com/cigarframework/reconciled/pkg/plugins/http",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/api:go_default_library",
        "//pkg/common:go_default_library",
        "//pkg/storage:go_default_library",
    ],
)
