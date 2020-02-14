load("@io_bazel_rules_go//go:def.bzl", "go_library")

package(default_visibility = ["//visibility:public"])

load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:exclude portal

#gazelle:exclude npm
gazelle(
    prefix = "github.com/cigarframework/reconciled",
    name = "gazelle",
    command = "fix",
    args = [
        "-build_file_name",
        "BUILD,BUILD.bazel",
    ],
)

# gazelle:resolve go github.com/cigarframework/reconciled/pkg/proto @reconciled//pkg/proto:go_default_library
go_library(
    name = "com_github_cigarframework_reconciled",
    importpath = "github.com/cigarframework/reconciled",
    deps = [
        "//pkg/api:go_default_library",
        "//pkg/client:go_default_library",
        "//pkg/common:go_default_library",
        "//pkg/grpcserver:go_default_library",
        "//pkg/httpserver:go_default_library",
        "//pkg/optional:go_default_library",
        "//pkg/plugins:go_default_library",
        "//pkg/server:go_default_library",
        "//pkg/storage:go_default_library",
        "//pkg/proto:go_default_library",
    ],
)
