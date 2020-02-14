load("@io_bazel_rules_go//go:def.bzl", "go_library")

# gazelle:resolve go github.com/cigarframework/reconciled/pkg/proto @reconciled//pkg/proto:go_default_library
go_library(
    name = "go_default_library",
    srcs = [
        "options.go",
        "server.go",
        "utils.go",
    ],
    importpath = "github.com/cigarframework/reconciled/pkg/grpcserver",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/api:go_default_library",
        "//pkg/optional:go_default_library",
        "//pkg/proto:go_default_library",
        "//pkg/storage:go_default_library",
        "@com_github_gogo_protobuf//types:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
        "@org_uber_go_zap//:go_default_library",
    ],
)
