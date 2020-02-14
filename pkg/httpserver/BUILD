load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "common.go",
        "options.go",
        "server.go",
    ],
    importpath = "github.com/cigarframework/reconciled/pkg/httpserver",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/api:go_default_library",
        "//pkg/storage:go_default_library",
        "@org_uber_go_zap//:go_default_library",
    ],
)
