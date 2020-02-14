load("@io_bazel_rules_go//go:def.bzl", "go_library")

# gazelle:resolve go github.com/cigarframework/reconciled/pkg/proto @reconciled//pkg/proto:go_default_library
go_library(
    name = "go_default_library",
    srcs = [
        "config.go",
        "server.go",
    ],
    importpath = "github.com/cigarframework/reconciled/cmd/rd-server/app",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/api:go_default_library",
        "//pkg/common:go_default_library",
        "//pkg/grpcserver:go_default_library",
        "//pkg/httpserver:go_default_library",
        "//pkg/plugins:go_default_library",
        "//pkg/proto:go_default_library",
        "//pkg/server:go_default_library",
        "@io_etcd_go_etcd//clientv3:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//credentials:go_default_library",
        "@org_golang_x_net//context:go_default_library",
        "@org_uber_go_zap//:go_default_library",
    ],
)
