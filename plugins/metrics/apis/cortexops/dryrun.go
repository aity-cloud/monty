package cortexops

import (
	"context"
	"fmt"
	strings "strings"

	cliutil "github.com/aity-cloud/monty/pkg/monty/cliutil"
	"github.com/aity-cloud/monty/pkg/plugins/driverutil"
	"github.com/nsf/jsondiff"
	"github.com/samber/lo"
	cobra "github.com/spf13/cobra"
	"github.com/ttacon/chalk"
	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func BuildDryRunCmd() *cobra.Command {
	var diffFull bool
	var diffFormat string
	dryRunCmd := &cobra.Command{
		Use: "config dry-run",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cliutil.BasePreRunE(cmd, args); err != nil {
				return err
			}
			// inject the dry-run client into the context
			client, ok := CortexOpsClientFromContext(cmd.Context())
			if ok {
				cmd.SetContext(ContextWithCortexOpsClient(cmd.Context(), &DryRunClient{
					Client: client,
				}))
			}
			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if err := cliutil.BasePostRunE(cmd, args); err != nil {
				return err
			}
			// print the dry-run response
			client, ok := CortexOpsClientFromContext(cmd.Context())
			if !ok {
				return nil
			}
			dryRunClient, ok := client.(*DryRunClient)
			if !ok {
				return nil
			}
			response := dryRunClient.Response
			if errs := response.GetValidationErrors(); len(errs) > 0 {
				cmd.Println(fmt.Sprintf("validation errors occurred (%d):", len(errs)))
				for _, e := range errs {
					switch e.GetSeverity() {
					case driverutil.ValidationError_Warning:
						cmd.Print("[" + chalk.Yellow.Color("WARN") + "] ")
					case driverutil.ValidationError_Error:
						cmd.Print("[" + chalk.Red.Color("ERROR") + "] ")
					}
					cmd.Println(e.GetMessage())
				}
			}

			var opts jsondiff.Options
			switch diffFormat {
			case "console":
				opts = jsondiff.DefaultConsoleOptions()
			case "json":
				opts = jsondiff.DefaultJSONOptions()
			case "html":
				opts = jsondiff.DefaultHTMLOptions()
			default:
				return fmt.Errorf("invalid diff format: %s", diffFormat)
			}
			opts.SkipMatches = !diffFull

			str, anyChanges := driverutil.RenderJsonDiff(response.Current, response.Modified, opts)
			if !anyChanges {
				cmd.Println("no changes")
			} else {
				cmd.Println(str)
			}
			return nil
		},
	}
	dryRunCmd.PersistentFlags().BoolVar(&diffFull, "diff-full", false, "show full diff, including all unchanged fields")
	dryRunCmd.PersistentFlags().StringVar(&diffFormat, "diff-format", "console", "diff format (console, json, html)")

	dryRunnableCmds := []*cobra.Command{
		BuildCortexOpsSetConfigurationCmd(),
		BuildCortexOpsSetDefaultConfigurationCmd(),
		BuildCortexOpsResetConfigurationCmd(),
		BuildCortexOpsResetDefaultConfigurationCmd(),
		BuildCortexOpsInstallCmd(),
		BuildCortexOpsUninstallCmd(),
	}

	for _, cmd := range dryRunnableCmds {
		cmd.Use = strings.TrimPrefix(cmd.Use, "config ")
		cmd.Short = fmt.Sprintf("[dry-run] %s", cmd.Short)
		dryRunCmd.AddCommand(cmd)
	}
	return dryRunCmd
}

type DryRunClient struct {
	Client   CortexOpsClient
	Request  *DryRunRequest
	Response *DryRunResponse
}

// ResetConfiguration implements CortexOpsClient.
func (dc *DryRunClient) ResetConfiguration(ctx context.Context, req *ResetRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	dc.Request = &DryRunRequest{
		Target: driverutil.Target_ActiveConfiguration,
		Action: driverutil.Action_Reset,
		Patch:  req.GetPatch(),
		Mask:   req.GetMask(),
	}
	var err error
	dc.Response, err = dc.Client.DryRun(ctx, dc.Request, opts...)
	if err != nil {
		return nil, fmt.Errorf("[dry-run] error: %w", err)
	}
	return &emptypb.Empty{}, nil
}

// ResetDefaultConfiguration implements CortexOpsClient.
func (dc *DryRunClient) ResetDefaultConfiguration(ctx context.Context, _ *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	dc.Request = &DryRunRequest{
		Target: driverutil.Target_DefaultConfiguration,
		Action: driverutil.Action_Reset,
	}
	var err error
	dc.Response, err = dc.Client.DryRun(ctx, dc.Request, opts...)
	if err != nil {
		return nil, fmt.Errorf("[dry-run] error: %w", err)
	}
	return &emptypb.Empty{}, nil
}

// SetConfiguration implements CortexOpsClient.
func (dc *DryRunClient) SetConfiguration(ctx context.Context, in *CapabilityBackendConfigSpec, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	in.Enabled = nil
	dc.Request = &DryRunRequest{
		Target: driverutil.Target_ActiveConfiguration,
		Action: driverutil.Action_Set,
		Spec:   in,
	}
	var err error
	dc.Response, err = dc.Client.DryRun(ctx, dc.Request, opts...)
	if err != nil {
		return nil, fmt.Errorf("[dry-run] error: %w", err)
	}
	return &emptypb.Empty{}, nil
}

// SetDefaultConfiguration implements CortexOpsClient.
func (dc *DryRunClient) SetDefaultConfiguration(ctx context.Context, in *CapabilityBackendConfigSpec, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	in.Enabled = nil
	dc.Request = &DryRunRequest{
		Target: driverutil.Target_DefaultConfiguration,
		Action: driverutil.Action_Set,
		Spec:   in,
	}
	var err error
	dc.Response, err = dc.Client.DryRun(ctx, dc.Request, opts...)
	if err != nil {
		return nil, fmt.Errorf("[dry-run] error: %w", err)
	}
	return &emptypb.Empty{}, nil
}

// ListPresets implements CortexOpsClient.
func (dc *DryRunClient) ListPresets(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PresetList, error) {
	return dc.Client.ListPresets(ctx, in, opts...)
}

// GetConfiguration implements CortexOpsClient.
func (dc *DryRunClient) GetConfiguration(ctx context.Context, in *driverutil.GetRequest, opts ...grpc.CallOption) (*CapabilityBackendConfigSpec, error) {
	return dc.Client.GetConfiguration(ctx, in, opts...)
}

// GetDefaultConfiguration implements CortexOpsClient.
func (dc *DryRunClient) GetDefaultConfiguration(ctx context.Context, in *driverutil.GetRequest, opts ...grpc.CallOption) (*CapabilityBackendConfigSpec, error) {
	return dc.Client.GetDefaultConfiguration(ctx, in, opts...)
}

// DryRun implements CortexOpsClient.
func (dc *DryRunClient) DryRun(_ context.Context, _ *DryRunRequest, _ ...grpc.CallOption) (*DryRunResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "[dry-run] method DryRun not implemented")
}

func (dc *DryRunClient) ConfigurationHistory(_ context.Context, _ *driverutil.ConfigurationHistoryRequest, _ ...grpc.CallOption) (*ConfigurationHistoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "[dry-run] method ConfigurationHistory not implemented")
}

// Install implements CortexOpsClient.
func (dc *DryRunClient) Install(ctx context.Context, _ *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	dc.Request = &DryRunRequest{
		Action: driverutil.Action_Set,
		Target: driverutil.Target_ActiveConfiguration,
		Spec: &CapabilityBackendConfigSpec{
			Enabled: lo.ToPtr(true),
		},
	}
	var err error
	dc.Response, err = dc.Client.DryRun(ctx, dc.Request, opts...)
	if err != nil {
		return nil, fmt.Errorf("[dry-run] error: %w", err)
	}
	return &emptypb.Empty{}, nil
}

// Status implements CortexOpsClient.
func (dc *DryRunClient) Status(_ context.Context, _ *emptypb.Empty, _ ...grpc.CallOption) (*driverutil.InstallStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "[dry-run] method Status not implemented")
}

// Uninstall implements CortexOpsClient.
func (dc *DryRunClient) Uninstall(ctx context.Context, _ *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	dc.Request = &DryRunRequest{
		Action: driverutil.Action_Set,
		Target: driverutil.Target_ActiveConfiguration,
		Spec: &CapabilityBackendConfigSpec{
			Enabled: lo.ToPtr(false),
		},
	}
	var err error
	dc.Response, err = dc.Client.DryRun(ctx, dc.Request, opts...)
	if err != nil {
		return nil, fmt.Errorf("[dry-run] error: %w", err)
	}
	return &emptypb.Empty{}, nil
}
