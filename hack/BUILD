load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["concat_dat_files.go"],
    importpath = "github.com/cigarframework/reconciled/hack",
    visibility = ["//visibility:private"],
    deps = [
        "@org_golang_x_tools//cover:go_default_library",
        "@org_uber_go_zap//buffer:go_default_library",
    ],
)

go_binary(
    name = "hack",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)
