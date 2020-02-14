load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "config.go",
        "plugin.go",
    ],
    importpath = "github.com/cigarframework/reconciled/pkg/plugins/authn",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/api:go_default_library",
        "@com_github_dgrijalva_jwt_go//:go_default_library",
        "@com_github_tidwall_buntdb//:go_default_library",
        "@org_uber_go_zap//:go_default_library",
    ],
)
