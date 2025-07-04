//go:build !minimal

package commands

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"

	"github.com/aity-cloud/monty/pkg/logger"
	"github.com/aity-cloud/monty/pkg/supportagent/dateparser"
	"github.com/aity-cloud/monty/pkg/supportagent/filereader"
	"github.com/aity-cloud/monty/pkg/supportagent/shipper"
	"google.golang.org/grpc"
)

const (
	rancherRegex = `cattle-system-rancher-[0-9a-f]+-[0-9a-z]+(-previous)?$`
)

var reRancher = regexp.MustCompile(rancherRegex)

type Distribution string

const (
	RKE  Distribution = "rke"
	RKE2 Distribution = "rke2"
	K3S  Distribution = "k3s"
)

func shipRKELogs(ctx context.Context, cc grpc.ClientConnInterface, lg *slog.Logger) error {
	lg.Info("shipping rke etcd logs")
	if err := shipRKEEtcd(ctx, cc, lg); err != nil {
		return err
	}
	lg.Info("shipping rke kube-api logs")
	if err := shipRKEKubeApi(ctx, cc, lg); err != nil {
		return err
	}
	lg.Info("shipping rke kubelet logs")
	if err := shipRKEKubelet(ctx, cc, lg); err != nil {
		return err
	}
	lg.Info("shipping rke kube-controller-manager logs")
	if err := shipRKEKubeControllerManager(ctx, cc, lg); err != nil {
		return err
	}
	lg.Info("shipping rke kube-scheduler logs")
	if err := shipRKEKubeScheduler(ctx, cc, lg); err != nil {
		return err
	}
	lg.Info("shipping rke kube-proxy logs")
	if err := shipRKEKubeProxy(ctx, cc, lg); err != nil {
		return err
	}
	lg.Info("shipping rke rancher logs")
	if err := shipRKERancher(ctx, cc, lg); err != nil {
		return err
	}
	return nil
}

func shipRKEEtcd(ctx context.Context, cc grpc.ClientConnInterface, lg *slog.Logger) error {
	reader, err := filereader.NewFileReader("k8s/containerlogs/etcd")
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	defer reader.Close()
	p := &dateparser.DockerParser{
		TimestampRegex: dateparser.EtcdRegex,
	}
	s := shipper.NewOTLPShipper(cc, p, lg,
		shipper.WithComponent("etcd"),
		shipper.WithLogType("controlplane"),
	)
	return s.Publish(ctx, reader.Scan())
}

func shipRKEKubeApi(ctx context.Context, cc grpc.ClientConnInterface, lg *slog.Logger) error {
	reader, err := filereader.NewFileReader("k8s/containerlogs/kube-apiserver")
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	defer reader.Close()
	p := &dateparser.DockerParser{
		TimestampRegex: dateparser.KlogRegex,
	}
	s := shipper.NewOTLPShipper(cc, p, lg,
		shipper.WithComponent("kube-apiserver"),
		shipper.WithLogType("controlplane"),
	)
	return s.Publish(ctx, reader.Scan())
}

func shipRKEKubelet(ctx context.Context, cc grpc.ClientConnInterface, lg *slog.Logger) error {
	reader, err := filereader.NewFileReader("k8s/containerlogs/kubelet")
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	defer reader.Close()
	p := &dateparser.DockerParser{
		TimestampRegex: dateparser.KlogRegex,
	}
	s := shipper.NewOTLPShipper(cc, p, lg,
		shipper.WithComponent("kubelet"),
		shipper.WithLogType("controlplane"),
	)
	return s.Publish(ctx, reader.Scan())
}

func shipRKEKubeControllerManager(ctx context.Context, cc grpc.ClientConnInterface, lg *slog.Logger) error {
	reader, err := filereader.NewFileReader("k8s/containerlogs/kube-controller-manager")
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	defer reader.Close()
	p := &dateparser.DockerParser{
		TimestampRegex: dateparser.KlogRegex,
	}
	s := shipper.NewOTLPShipper(cc, p, lg,
		shipper.WithComponent("kube-controller-manager"),
		shipper.WithLogType("controlplane"),
	)
	return s.Publish(ctx, reader.Scan())
}

func shipRKEKubeScheduler(ctx context.Context, cc grpc.ClientConnInterface, lg *slog.Logger) error {
	reader, err := filereader.NewFileReader("k8s/containerlogs/kube-scheduler")
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	defer reader.Close()
	p := &dateparser.DockerParser{
		TimestampRegex: dateparser.KlogRegex,
	}
	s := shipper.NewOTLPShipper(cc, p, lg,
		shipper.WithComponent("kube-scheduler"),
		shipper.WithLogType("controlplane"),
	)
	return s.Publish(ctx, reader.Scan())
}

func shipRKEKubeProxy(ctx context.Context, cc grpc.ClientConnInterface, lg *slog.Logger) error {
	reader, err := filereader.NewFileReader("k8s/containerlogs/kube-proxy")
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	defer reader.Close()
	p := &dateparser.DockerParser{
		TimestampRegex: dateparser.KlogRegex,
	}
	s := shipper.NewOTLPShipper(cc, p, lg,
		shipper.WithComponent("kube-proxy"),
		shipper.WithLogType("controlplane"),
	)
	return s.Publish(ctx, reader.Scan())
}

func shipRKERancher(ctx context.Context, cc grpc.ClientConnInterface, lg *slog.Logger) error {
	files, err := filepath.Glob("rancher/containerlogs/server-*")
	if err != nil {
		lg.Error("failed to glob rancher logs", logger.Err(err))
		return err
	}

	p := &dateparser.MultipleParser{
		Dateformats: []dateparser.Dateformat{
			{
				DateRegex: dateparser.RancherRegex,
				Layout:    dateparser.RancherLayout,
			},
			{
				DateRegex: dateparser.KlogRegex,
				Layout:    dateparser.KlogLayout,
			},
		},
	}

	for _, file := range files {
		if !reRancher.MatchString(file) {
			continue
		}
		s := shipper.NewOTLPShipper(cc, p, lg,
			shipper.WithLogType("rancher"),
		)
		reader, err := filereader.NewFileReader(file)
		if err != nil {
			if os.IsExist(err) {
				return nil
			}
			return err
		}
		err = s.Publish(ctx, reader.Scan())
		if err != nil {
			lg.With(logger.Err(err)).Error(fmt.Sprintf("failed to publish rancher logs %s", file))
		}
		reader.Close()
	}
	return nil
}

func shipK3sLogs(ctx context.Context, cc grpc.ClientConnInterface, lg *slog.Logger) error {
	lg.Info("shipping k3s service logs")
	if err := shipK3sControlplane(ctx, cc, lg); err != nil {
		return err
	}
	lg.Info("shipping k3s rancher logs")
	if err := shipK3sRancher(ctx, cc, lg); err != nil {
		return err
	}
	return nil
}

func zoneAndYearFromDatefile() (string, string) {
	if _, err := os.Stat("systeminfo/date"); err == nil {
		dateFile, err := os.Open("systeminfo/date")
		if err != nil {
			return "", ""
		}
		defer dateFile.Close()
		scanner := bufio.NewScanner(dateFile)
		scanner.Scan()
		dateline := scanner.Text()
		re := regexp.MustCompile(`^[A-Z][a-z]{2} [A-Z][a-z]{2} \d{1,2} \d{2}:\d{2}:\d{2} ([A-Z]{3}) (\d{4})`)
		matches := re.FindStringSubmatch(dateline)
		if len(matches) != 0 {
			return matches[1], matches[2]
		}
	}

	return "", ""
}

func shipK3sControlplane(ctx context.Context, cc grpc.ClientConnInterface, lg *slog.Logger) error {
	// Extract timezone and year from the date output
	timezone, year := zoneAndYearFromDatefile()

	reader, err := filereader.NewFileReader("journald/k3s")
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	defer reader.Close()

	p := dateparser.NewDayMonthParser(
		dateparser.JournaldRegex,
		dateparser.JournaldLayout,
		dateparser.WithTimezone(timezone),
		dateparser.WithYear(year),
	)

	s := shipper.NewOTLPShipper(cc, p, lg,
		shipper.WithComponent("k3s"),
		shipper.WithLogType("controlplane"),
	)
	return s.Publish(ctx, reader.Scan())
}

func shipK3sRancher(ctx context.Context, cc grpc.ClientConnInterface, lg *slog.Logger) error {
	timezone, year := zoneAndYearFromDatefile()

	files, err := filepath.Glob("k3s/podlogs/cattle-system-rancher-*")
	if err != nil {
		lg.Error("failed to glob rancher logs", logger.Err(err))
		return err
	}

	p := &dateparser.MultipleParser{
		Dateformats: []dateparser.Dateformat{
			{
				DateRegex: dateparser.RancherRegex,
				Layout:    dateparser.RancherLayout,
			},
			{
				DateRegex:  dateparser.KlogRegex,
				Layout:     dateparser.KlogLayout,
				DateSuffix: fmt.Sprintf(" %s %s", timezone, year),
			},
		},
	}

	for _, file := range files {
		if !reRancher.MatchString(file) {
			continue
		}
		s := shipper.NewOTLPShipper(cc, p, lg,
			shipper.WithLogType("rancher"),
		)
		reader, err := filereader.NewFileReader(file)
		if err != nil {
			if os.IsExist(err) {
				return nil
			}
			return err
		}
		err = s.Publish(ctx, reader.Scan())
		if err != nil {
			lg.With(logger.Err(err)).Error(fmt.Sprintf("failed to publish rancher logs %s", file))
		}
		reader.Close()
	}

	return nil
}

func shipRKE2Logs(ctx context.Context, cc grpc.ClientConnInterface, lg *slog.Logger) error {
	lg.Info("shipping rke2 etcd logs")
	if err := shipRKE2Etcd(ctx, cc, lg); err != nil {
		return err
	}
	lg.Info("shipping rke2 kubelet logs")
	if err := shipRKE2Kubelet(ctx, cc, lg); err != nil {
		return err
	}
	lg.Info("shipping rke2 kube-api logs")
	if err := shipRKE2KubeAPI(ctx, cc, lg); err != nil {
		return err
	}
	lg.Info("shipping rke2 kube-controller-manager logs")
	if err := shipRKE2KubeControllerManager(ctx, cc, lg); err != nil {
		return err
	}
	lg.Info("shipping rke2 kube-scheduler logs")
	if err := shipRKE2KubeScheduler(ctx, cc, lg); err != nil {
		return err
	}
	lg.Info("shipping rke2 kube-proxy logs")
	if err := shipRKE2KubeProxy(ctx, cc, lg); err != nil {
		return err
	}
	lg.Info("shipping rke2 journald logs")
	if err := shipRKE2Journald(ctx, cc, lg); err != nil {
		return err
	}
	lg.Info("shipping rke2 rancher logs")
	if err := shipRKE2Rancher(ctx, cc, lg); err != nil {
		return err
	}
	return nil
}

func shipRKE2Etcd(ctx context.Context, cc grpc.ClientConnInterface, lg *slog.Logger) error {
	files, err := filepath.Glob("rke2/podlogs/kube-system-etcd-*")
	if err != nil {
		lg.Error("failed to glob rke2 etcd logs", logger.Err(err))
		return err
	}

	p := &dateparser.RKE2EtcdParser{}

	for _, file := range files {
		s := shipper.NewOTLPShipper(cc, p, lg,
			shipper.WithComponent("etcd"),
			shipper.WithLogType("controlplane"),
		)
		reader, err := filereader.NewFileReader(file)
		if err != nil {
			if os.IsExist(err) {
				return nil
			}
			return err
		}
		err = s.Publish(ctx, reader.Scan())
		if err != nil {
			lg.With(logger.Err(err)).Error(fmt.Sprintf("failed to publish rke2 etcd logs %s", file))
		}
		reader.Close()
	}

	return nil
}

func shipRKE2Kubelet(ctx context.Context, cc grpc.ClientConnInterface, lg *slog.Logger) error {
	reader, err := filereader.NewFileReader("rke2/agent-logs/kubelet.log")
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	defer reader.Close()

	timezone, year := zoneAndYearFromDatefile()
	p := dateparser.NewDayMonthParser(
		dateparser.KlogRegex,
		dateparser.KlogLayout,
		dateparser.WithTimezone(timezone),
		dateparser.WithYear(year),
	)

	s := shipper.NewOTLPShipper(cc, p, lg,
		shipper.WithComponent("kubelet"),
		shipper.WithLogType("controlplane"),
	)
	return s.Publish(ctx, reader.Scan())
}

func shipRKE2KubeAPI(ctx context.Context, cc grpc.ClientConnInterface, lg *slog.Logger) error {
	files, err := filepath.Glob("rke2/podlogs/kube-system-kube-apiserver-*")
	if err != nil {
		lg.Error("failed to glob rke2 kube-api logs", logger.Err(err))
		return err
	}

	timezone, year := zoneAndYearFromDatefile()
	p := dateparser.NewDayMonthParser(
		dateparser.KlogRegex,
		dateparser.KlogLayout,
		dateparser.WithTimezone(timezone),
		dateparser.WithYear(year),
	)

	for _, file := range files {
		s := shipper.NewOTLPShipper(cc, p, lg,
			shipper.WithComponent("kube-apiserver"),
			shipper.WithLogType("controlplane"),
		)
		reader, err := filereader.NewFileReader(file)
		if err != nil {
			if os.IsExist(err) {
				return nil
			}
			return err
		}
		err = s.Publish(ctx, reader.Scan())
		if err != nil {
			lg.With(logger.Err(err)).Error(fmt.Sprintf("failed to publish rke2 kube-api logs %s", file))
		}
		reader.Close()
	}

	return nil
}

func shipRKE2KubeControllerManager(ctx context.Context, cc grpc.ClientConnInterface, lg *slog.Logger) error {
	files, err := filepath.Glob("rke2/podlogs/kube-system-kube-controller-manager-*")
	if err != nil {
		lg.Error("failed to glob rke2 kube-controller-manager logs", logger.Err(err))
		return err
	}

	timezone, year := zoneAndYearFromDatefile()
	p := dateparser.NewDayMonthParser(
		dateparser.KlogRegex,
		dateparser.KlogLayout,
		dateparser.WithTimezone(timezone),
		dateparser.WithYear(year),
	)

	for _, file := range files {
		s := shipper.NewOTLPShipper(cc, p, lg,
			shipper.WithComponent("kube-controller-manager"),
			shipper.WithLogType("controlplane"),
		)
		reader, err := filereader.NewFileReader(file)
		if err != nil {
			if os.IsExist(err) {
				return nil
			}
			return err
		}
		err = s.Publish(ctx, reader.Scan())
		if err != nil {
			lg.With(logger.Err(err)).Error(fmt.Sprintf("failed to publish rke2 kube-controller-manager logs %s", file))
		}
		reader.Close()
	}

	return nil
}

func shipRKE2KubeScheduler(ctx context.Context, cc grpc.ClientConnInterface, lg *slog.Logger) error {
	files, err := filepath.Glob("rke2/podlogs/kube-system-kube-scheduler-*")
	if err != nil {
		lg.Error("failed to glob rke2 kube-scheduler logs", logger.Err(err))
		return err
	}

	timezone, year := zoneAndYearFromDatefile()
	p := dateparser.NewDayMonthParser(
		dateparser.KlogRegex,
		dateparser.KlogLayout,
		dateparser.WithTimezone(timezone),
		dateparser.WithYear(year),
	)

	for _, file := range files {
		s := shipper.NewOTLPShipper(cc, p, lg,
			shipper.WithComponent("kube-scheduler"),
			shipper.WithLogType("controlplane"),
		)
		reader, err := filereader.NewFileReader(file)
		if err != nil {
			if os.IsExist(err) {
				return nil
			}
			return err
		}
		err = s.Publish(ctx, reader.Scan())
		if err != nil {
			lg.With(logger.Err(err)).Error(fmt.Sprintf("failed to publish rke2 kube-scheduler logs %s", file))
		}
		reader.Close()
	}

	return nil
}

func shipRKE2KubeProxy(ctx context.Context, cc grpc.ClientConnInterface, lg *slog.Logger) error {
	files, err := filepath.Glob("rke2/podlogs/kube-system-kube-proxy-*")
	if err != nil {
		lg.Error("failed to glob rke2 kube-proxy logs", logger.Err(err))
		return err
	}

	timezone, year := zoneAndYearFromDatefile()
	p := dateparser.NewDayMonthParser(
		dateparser.KlogRegex,
		dateparser.KlogLayout,
		dateparser.WithTimezone(timezone),
		dateparser.WithYear(year),
	)

	for _, file := range files {
		s := shipper.NewOTLPShipper(cc, p, lg,
			shipper.WithComponent("kube-proxy"),
			shipper.WithLogType("controlplane"),
		)
		reader, err := filereader.NewFileReader(file)
		if err != nil {
			if os.IsExist(err) {
				return nil
			}
			return err
		}
		err = s.Publish(ctx, reader.Scan())
		if err != nil {
			lg.With(logger.Err(err)).Error(fmt.Sprintf("failed to publish rke2 kube-proxy logs %s", file))
		}
		reader.Close()
	}

	return nil
}

func shipRKE2Journald(ctx context.Context, cc grpc.ClientConnInterface, lg *slog.Logger) error {
	reader, err := filereader.NewFileReader("journald/rke2-server")
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	defer reader.Close()

	timezone, year := zoneAndYearFromDatefile()
	p := dateparser.NewDayMonthParser(
		dateparser.JournaldRegex,
		dateparser.JournaldLayout,
		dateparser.WithTimezone(timezone),
		dateparser.WithYear(year),
	)

	s := shipper.NewOTLPShipper(cc, p, lg,
		shipper.WithComponent("rke2"),
		shipper.WithLogType("controlplane"),
	)
	return s.Publish(ctx, reader.Scan())
}

func shipRKE2Rancher(ctx context.Context, cc grpc.ClientConnInterface, lg *slog.Logger) error {
	files, err := filepath.Glob("rke2/podlogs/cattle-system-rancher-*")
	if err != nil {
		lg.Error("failed to glob rke2 rancher logs", logger.Err(err))
		return err
	}

	timezone, year := zoneAndYearFromDatefile()
	p := &dateparser.MultipleParser{
		Dateformats: []dateparser.Dateformat{
			{
				DateRegex: dateparser.RancherRegex,
				Layout:    dateparser.RancherLayout,
			},
			{
				DateRegex:  dateparser.KlogRegex,
				Layout:     dateparser.KlogLayout,
				DateSuffix: fmt.Sprintf(" %s %s", timezone, year),
			},
		},
	}

	for _, file := range files {
		if !reRancher.MatchString(file) {
			continue
		}
		s := shipper.NewOTLPShipper(cc, p, lg,
			shipper.WithLogType("rancher"),
		)
		reader, err := filereader.NewFileReader(file)
		if err != nil {
			if os.IsExist(err) {
				return nil
			}
			return err
		}
		err = s.Publish(ctx, reader.Scan())
		if err != nil {
			lg.With(logger.Err(err)).Error(fmt.Sprintf("failed to publish rke2 rancher logs %s", file))
		}
		reader.Close()
	}

	return nil
}
