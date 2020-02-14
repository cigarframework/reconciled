load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "github.com/cigarframework/reconciled/cmd/rd-server",
    visibility = ["//visibility:private"],
    deps = [
        "//cmd/rd-server/app:go_default_library",
        "@in_gopkg_yaml_v2//:go_default_library",
        "@org_uber_go_zap//:go_default_library",
    ],
)

go_binary(
    name = "rd-server",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)
