load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["resolver.go"],
    importpath = "github.com/cigarframework/reconciled/pkg/grpclb",
    visibility = ["//visibility:public"],
    deps = [
        "@org_golang_google_grpc//resolver:go_default_library",
        "@org_uber_go_atomic//:go_default_library",
    ],
)
