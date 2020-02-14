load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "config.go",
        "plugin.go",
    ],
    importpath = "github.com/cigarframework/reconciled/pkg/plugins/authz",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/api:go_default_library",
        "//pkg/optional:go_default_library",
        "//pkg/storage:go_default_library",
        "@com_github_ory_ladon//:go_default_library",
        "@com_github_ory_ladon//manager/memory:go_default_library",
        "@org_uber_go_zap//:go_default_library",
    ],
)
