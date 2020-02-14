load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

# gazelle:resolve proto github.com/gogo/protobuf/gogoproto/gogo.proto @gogo_special_proto//github.com/gogo/protobuf/gogoproto
# gazelle:resolve proto go github.com/gogo/protobuf/gogoproto/gogo.proto @com_github_gogo_protobuf//gogoproto:go_default_library
proto_library(
    name = "proto_proto",
    srcs = [
        "plugin.proto",
        "state.proto",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "@com_google_protobuf//:empty_proto",
        "@com_google_protobuf//:struct_proto",
        "@com_google_protobuf//:timestamp_proto",
        "@com_google_protobuf//:wrappers_proto",
        "@gogo_special_proto//github.com/gogo/protobuf/gogoproto",
    ],
)

# gazelle:go_grpc_compilers @io_bazel_rules_go//proto:gogoslick_grpc
go_proto_library(
    name = "proto_go_proto",
    compilers = ["@io_bazel_rules_go//proto:gogoslick_grpc"],
    importpath = "github.com/cigarframework/reconciled/pkg/proto",
    proto = ":proto_proto",
    visibility = ["//visibility:public"],
    deps = ["@com_github_gogo_protobuf//gogoproto:go_default_library"],
)

go_library(
    name = "go_default_library",
    embed = [":proto_go_proto"],
    importpath = "github.com/cigarframework/reconciled/pkg/proto",
    visibility = ["//visibility:public"],
)
