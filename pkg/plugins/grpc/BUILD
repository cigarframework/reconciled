load("@io_bazel_rules_go//go:def.bzl", "go_library")

# gazelle:resolve go github.com/cigarframework/reconciled/pkg/proto @reconciled//pkg/proto:go_default_library
go_library(
    name = "go_default_library",
    srcs = [
        "config.go",
        "plugin.go",
    ],
    importpath = "github.com/cigarframework/reconciled/pkg/plugins/grpc",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/api:go_default_library",
        "//pkg/common:go_default_library",
        "//pkg/grpclb:go_default_library",
        "//pkg/proto:go_default_library",
        "//pkg/storage:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//balancer/roundrobin:go_default_library",
        "@org_golang_google_grpc//credentials:go_default_library",
        "@org_golang_google_grpc//resolver:go_default_library",
    ],
)
