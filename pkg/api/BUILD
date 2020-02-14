load("@io_bazel_rules_go//go:def.bzl", "go_library")

# gazelle:resolve go github.com/cigarframework/reconciled/pkg/proto @reconciled//pkg/proto:go_default_library
go_library(
    name = "go_default_library",
    srcs = [
        "auth.go",
        "context.go",
        "errors.go",
        "patch.go",
        "plugin.go",
        "server.go",
        "user.go",
    ],
    importpath = "github.com/cigarframework/reconciled/pkg/api",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/optional:go_default_library",
        "//pkg/storage:go_default_library",
    ],
)
