package targets

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	internalcodegen "github.com/aity-cloud/monty/internal/codegen"
	"github.com/aity-cloud/monty/internal/codegen/pathbuilder"
	"github.com/aity-cloud/monty/internal/codegen/templating"
	"github.com/kralicky/codegen/cli"
	"github.com/kralicky/protols/sdk/codegen"
	"github.com/kralicky/protols/sdk/codegen/generators/external"
	"github.com/kralicky/protols/sdk/codegen/generators/golang"
	"github.com/kralicky/protols/sdk/codegen/generators/golang/grpc"
	"github.com/kralicky/protols/sdk/codegen/generators/x/python"
	"github.com/magefile/mage/mg"
	_ "go.opentelemetry.io/proto/otlp/metrics/v1"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Generates Go protobuf code
func (Generate) ProtobufGo(ctx context.Context) error {
	mg.Deps(internalcodegen.GenCortexConfig)
	_, tr := Tracer.Start(ctx, "target.generate.protobuf.go")
	defer tr.End()

	generators := []codegen.Generator{
		templating.CommentRenderer{},
		golang.Generator,
		grpc.Generator,
		cli.NewGenerator(),
		pathbuilder.PathBuilderGenerator{
			Roots: []protoreflect.FullName{
				protoreflect.FullName("config.v1.GatewayConfigSpec"),
				protoreflect.FullName("ext.SampleConfiguration"),
			},
		},
	}

	out, err := codegen.GenerateCode(
		generators,
		[]string{
			"internal/cortex",
			"pkg",
			"plugins",
		},
		codegen.WithFixInvalidGoPackages(true),
	)
	if err != nil {
		return err
	}
	for _, file := range out {
		if err := file.WriteToDisk(); err != nil {
			return err
		}
	}

	return nil
}

// Generates Python protobuf code
func (Generate) ProtobufPython(ctx context.Context) error {
	_, tr := Tracer.Start(ctx, "target.generate.protobuf.python")
	defer tr.End()

	generators := []codegen.Generator{python.Generator}
	out, err := codegen.GenerateCode(generators, []string{"aiops"})
	if err != nil {
		return err
	}
	for _, file := range out {
		if err := file.WriteToDisk(); err != nil {
			return err
		}
	}
	return nil
}

func (Generate) ProtobufTypescript() error {
	mg.Deps(Build.TypescriptServiceGenerator)
	destDir := "web/pkg/monty/generated"

	searchDirs := []string{
		"pkg/config/v1",
		"pkg/config/v1beta1",
		"pkg/apis/capability/v1",
		"pkg/apis/management/v1",
		"pkg/validation",
		"plugins/metrics/apis/cortexadmin",
		"plugins/metrics/apis/cortexops",
		"plugins/metrics/apis/node",
		"plugins/logging/apis/loggingadmin",
		"internal/cortex",
	}

	out, err := codegen.GenerateCode([]codegen.Generator{
		external.NewGenerator("./web/service-generator/node_modules/.bin/protoc-gen-es", external.GeneratorOptions{
			Opt: "target=ts,import_extension=none,ts_nocheck=false",
		}),
		external.NewGenerator([]string{"./web/service-generator/generate"}, external.GeneratorOptions{
			Opt: "target=ts,import_extension=none,ts_nocheck=false",
		}),
	}, searchDirs, codegen.WithGenerateStrategy(codegen.AllDescriptorsExceptGoogleProtobuf))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error generating typescript code: %v\n", err)
		return err
	}

	for _, file := range out {
		file.SourceRelPath = filepath.Join(destDir, file.Package, file.Name)
		os.MkdirAll(filepath.Dir(file.SourceRelPath), 0o755)
		if err := file.WriteToDisk(); err != nil {
			return fmt.Errorf("error writing file %s: %w", file.SourceRelPath, err)
		}
	}
	return nil
}

// Generates all protobuf code
func (Generate) Protobuf(ctx context.Context) {
	ctx, tr := Tracer.Start(ctx, "target.generate.protobuf")
	defer tr.End()

	_, err := exec.LookPath("yarn")
	if err == nil {
		mg.CtxDeps(ctx, Generate.ProtobufGo, Generate.ProtobufPython)
	} else {
		mg.CtxDeps(ctx, Generate.ProtobufGo, Generate.ProtobufPython)
	}
}
