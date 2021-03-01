def _esbuild_impl(ctx):
    inputs = []
    outdir = ctx.actions.declare_directory("out")

    for dep in ctx.attr.srcs:
        if hasattr(dep, "files"):
            inputs.extend(dep.files.to_list())

    args = ctx.actions.args()

    args.add_joined(['--entry', 'bazel-out/darwin-fastbuild/bin/packages/app/src/main.js'], join_with = "=")
    args.add('--root=packages/app')
    args.add('--index=packages/app/src/index.html')

    print(args)

    ctx.actions.run(
        arguments = [args],
        executable = ctx.executable._esbuild,
        inputs = inputs,
        outputs = [outdir],
        mnemonic = "Build",
    )

    return DefaultInfo(files = depset([outdir]))

esbuild = rule(
    implementation = _esbuild_impl,
    attrs = {
        "srcs": attr.label_list(
            allow_files = True,
            default = [],
        ),
        "_esbuild": attr.label(
            default = Label("//tools/build"),
            executable = True,
            cfg = "host",
        ),
    },
)
