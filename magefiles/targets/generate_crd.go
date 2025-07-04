package targets

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
	"github.com/samber/lo"
)

type CRD mg.Namespace

var runReplace bool

func (CRD) All(ctx context.Context) {
	ctx, tr := Tracer.Start(ctx, "target.crd")
	defer tr.End()
	mg.SerialCtxDeps(ctx, CRD.CRDGen, CRD.ReplaceCRDText)
}

func (CRD) CRDGen(ctx context.Context) error {
	_, tr := Tracer.Start(ctx, "target.crd.crdgen")
	defer tr.End()

	ok, err := checkCRDContstraints()
	if err != nil {
		return err
	}
	if ok {
		runReplace = true
		var commands []*exec.Cmd

		commands = append(commands, exec.Command(mg.GoCmd(), "run", "sigs.k8s.io/kustomize/kustomize/v5@v5.0.3",
			"build", "./config/chart-crds", "-o", "./packages/monty/monty/charts/crds/crds.yaml",
		))
		for _, cmd := range commands {
			buf := new(bytes.Buffer)
			cmd.Stderr = buf
			cmd.Stdout = buf
			err := cmd.Run()
			if err != nil {
				if ex, ok := err.(*exec.ExitError); ok {
					if ex.ExitCode() != 1 {
						return errors.New(buf.String())
					}
					bufStr := buf.String()
					lines := strings.Split(bufStr, "\n")
					for _, line := range lines {
						if strings.TrimSpace(line) == "" {
							continue
						}
						fmt.Fprintln(os.Stderr, line)
						return err
					}
				}
			}
		}

		expr := `del(.. | select(has("description")).description | select(has("type") | not )) | .. style="flow"`

		e1 := lo.Async(func() error {
			if yq, err := exec.LookPath("yq"); err == nil {
				return sh.Run(yq, "-i", expr, "./packages/monty/monty/charts/crds/crds.yaml")
			} else {
				return sh.Run(mg.GoCmd(), "run", "github.com/mikefarah/yq/v4@latest", "-i", expr, "./packages/monty/monty/charts/crds/crds.yaml")
			}
		})

		if err := <-e1; err != nil {
			return err
		}

		// prepend "---" to each file, otherwise kubernetes will think it's json
		for _, f := range []string{"./packages/monty/monty/charts/crds/crds.yaml"} {
			if err := prependDocumentSeparator(f); err != nil {
				return err
			}
		}
	}

	return nil
}

func prependDocumentSeparator(path string) error {
	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	i, err := f.Stat()
	if err != nil {
		return err
	}

	buf := make([]byte, i.Size()+4)
	copy(buf[:4], "---\n")

	_, err = f.Read(buf[4:])
	if err != nil {
		return err
	}

	f.Seek(0, 0)
	f.Truncate(0)
	_, err = f.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

func (CRD) ReplaceCRDText(ctx context.Context) error {
	_, tr := Tracer.Start(ctx, "target.crd.replacecrdtext")
	defer tr.End()

	if !runReplace {
		ok, err := checkCRDContstraints()
		if err != nil {
			return err
		}
		runReplace = ok
	}
	if runReplace {
		files := []string{
			"./packages/monty/monty/charts/crds/crds.yaml",
		}

		for _, file := range files {
			input, err := os.ReadFile(file)
			if err != nil {
				return err
			}

			firstReplace := bytes.Replace(input, []byte("replace-me/monty-serving-cert"), []byte(`"replace-me/monty-serving-cert"`), -1)
			output := bytes.Replace(firstReplace, []byte("replace-me"), []byte("{{ .Release.Namespace }}"), -1)

			if err := os.WriteFile(file, output, 0644); err != nil {
				return err
			}
		}
	}
	return nil
}

func checkCRDContstraints() (bool, error) {
	var crdSources []string
	err := filepath.WalkDir("config/crd", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(entry.Name(), ".yaml") {
			crdSources = append(crdSources, path)
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	err = filepath.WalkDir("config/chart-crds", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(entry.Name(), ".yaml") {
			crdSources = append(crdSources, path)
		}
		return nil
	})
	if err != nil {
		return false, err
	}

	return target.Dir("packages/monty/monty/charts/crds", crdSources...)
}
