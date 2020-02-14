load("@io_bazel_rules_go//go:def.bzl", "go_library")

# gazelle:resolve go github.com/cigarframework/reconciled/pkg/proto @reconciled//pkg/proto:go_default_library
go_library(
    name = "go_default_library",
    srcs = [
        "common.go",
        "create.go",
        "get.go",
        "list.go",
        "options.go",
        "patch.go",
        "remove.go",
        "server.go",
        "subscription.go",
        "update.go",
    ],
    importpath = "github.com/cigarframework/reconciled/pkg/server",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/api:go_default_library",
        "//pkg/common:go_default_library",
        "//pkg/optional:go_default_library",
        "//pkg/storage:go_default_library",
        "@com_github_antonmedv_expr//:go_default_library",
        "@com_github_antonmedv_expr//vm:go_default_library",
        "@com_github_coreos_etcd//mvcc/mvccpb:go_default_library",
        "@com_github_evanphx_json_patch//:go_default_library",
        "@com_github_satori_go_uuid//:go_default_library",
        "@io_etcd_go_etcd//clientv3:go_default_library",
        "@org_uber_go_zap//:go_default_library",
    ],
)
