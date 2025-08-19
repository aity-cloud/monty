package targets

import (
	"context"
	"fmt"

	"github.com/magefile/mage/mg"
)

// Builds the monty binary
func (Build) Monty(ctx context.Context) error {
	mg.CtxDeps(ctx, Build.Archives)

	_, tr := Tracer.Start(ctx, "target.build.monty")
	defer tr.End()

	return buildMainPackage(buildOpts{
		Path:   "./cmd/monty",
		Output: "bin/monty",
		Tags:   []string{"nomsgpack"},
	})
}

// Same as 'build:monty' but with debug symbols enabled
func (Build) MontyDebug(ctx context.Context) error {
	mg.CtxDeps(ctx, Build.Archives)

	_, tr := Tracer.Start(ctx, "target.build.monty")
	defer tr.End()

	return buildMainPackage(buildOpts{
		Path:   "./cmd/monty",
		Output: "bin/monty",
		Debug:  true,
		Tags:   []string{"nomsgpack"},
	})
}

// Same as 'build:monty' but with race detection enabled
func (Build) MontyRace(ctx context.Context) error {
	mg.CtxDeps(ctx, Build.Archives)

	_, tr := Tracer.Start(ctx, "target.build.monty")
	defer tr.End()

	return buildMainPackage(buildOpts{
		Path:   "./cmd/monty",
		Output: "bin/monty",
		Race:   true,
		Tags:   []string{"nomsgpack"},
	})
}

// Builds the monty-minimal binary
func (Build) MontyMinimal(ctx context.Context) error {
	_, tr := Tracer.Start(ctx, "target.build.monty-minimal")
	defer tr.End()

	return buildMainPackage(buildOpts{
		Path:   "./cmd/monty",
		Output: "bin/monty-minimal",
		Tags:   []string{"nomsgpack", "minimal"},
	})
}

// Builds the monty release CLI binary, requires version as input
func (Build) MontyReleaseCLI(ctx context.Context, fileSuffix string) error {
	mg.CtxDeps(ctx, Build.Archives)

	_, tr := Tracer.Start(ctx, "target.build.monty-cli-release")
	defer tr.End()

	return buildMainPackage(buildOpts{
		Path:     "./cmd/monty",
		Output:   fmt.Sprintf("bin/monty_%s", fileSuffix),
		Tags:     []string{"nomsgpack", "cli"},
		Compress: true,
	})
}
