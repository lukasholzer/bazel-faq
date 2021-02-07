load("@bazel_gazelle//:deps.bzl", "go_repository")

# Download all dependencies with `go get -u all`
# Then automatically add it with:
# npx bazelisk run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%install_dependencies
def install_dependencies():
    go_repository(
        name = "com_github_evanw_esbuild",
        importpath = "github.com/evanw/esbuild",
        sum = "h1:4rZ1OJWV+yi5hFvLJ1dq9ZACrQZxzpa4u+NKoSWv6XM=",
        version = "v0.8.42",
    )
    go_repository(
        name = "com_github_gorilla_websocket",
        importpath = "github.com/gorilla/websocket",
        sum = "h1:+/TMaTYc4QFitKJxsQ7Yye35DkWvkdLcvGKqM+x0Ufc=",
        version = "v1.4.2",
    )

    go_repository(
        name = "org_golang_x_sys",
        importpath = "golang.org/x/sys",
        sum = "h1:VwygUrnw9jn88c4u8GD3rZQbqrP/tgas88tPUbBxQrk=",
        version = "v0.0.0-20210124154548-22da62e12c0c",
    )
