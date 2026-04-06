// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type libvaccelPaths struct {
	libDir    string
	imagesDir string
	modelsDir string
	inputDir  string
	labelsDir string
}

func pkgConfigVar(t *testing.T, pkg, variable string) string {
	t.Helper()
	cmd := exec.Command("pkg-config", "--variable="+variable, pkg) //nolint:gosec
	cmd.Env = os.Environ()
	out, err := cmd.Output()
	if err != nil {
		t.Skipf("pkg-config %s not found: %v", pkg, err)
	}
	return strings.TrimSpace(string(out))
}

func resolveLibvaccelPaths(t *testing.T) libvaccelPaths {
	t.Helper()
	prefix := pkgConfigVar(t, "vaccel", "prefix")
	return libvaccelPaths{
		libDir:    pkgConfigVar(t, "vaccel", "libdir"),
		imagesDir: filepath.Join(prefix, "share", "vaccel", "images"),
		modelsDir: filepath.Join(prefix, "share", "vaccel", "models"),
		inputDir:  filepath.Join(prefix, "share", "vaccel", "input"),
		labelsDir: filepath.Join(prefix, "share", "vaccel", "labels"),
	}
}

func TestExamples(t *testing.T) {
	paths := resolveLibvaccelPaths(t)
	env := append(
		os.Environ(),
		"LD_LIBRARY_PATH="+paths.libDir,
		"VACCEL_PLUGINS=libvaccel-noop.so",
	)

	tests := []struct {
		name    string
		dir     string
		args    []string
		env     []string
		wantOut string
	}{
		{
			name: "noop",
		},
		{
			name: "classify",
			args: []string{filepath.Join(paths.imagesDir, "example.jpg")},
			wantOut: `Output(1):  This is a dummy classification tag!
Output(2):  This is a dummy classification tag!`,
		},
		{
			name: "detect",
			args: []string{filepath.Join(paths.imagesDir, "example.jpg")},
			wantOut: `Output(1):  This is a dummy imgname!
Output(2):  This is a dummy imgname!`,
		},
		{
			name: "exec",
			args: []string{filepath.Join(paths.libDir, "libmytestlib.so"), "10"},
			wantOut: `Output(1):  10
Output(2):  10
Output(3):  10`,
		},
		{
			name: "nonser",
			args: []string{filepath.Join(paths.libDir, "libmytestlib.so")},
			wantOut: `Input: 10 20 30 40 50 
Output: 10 20 30 40 50 `,
		},
		{
			name: "tf",
			args: []string{filepath.Join(paths.modelsDir, "tf")},
			wantOut: `Success!
Output tensor => type:1 nr_dims:2
dim[0]: 1
dim[1]: 30
Result Tensor:
Tensor shape: [1 30]
Values:
  [1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000]`,
		},
		{
			name: "tflite",
			args: []string{filepath.Join(paths.modelsDir, "tf/lstm2.tflite")},
			wantOut: `Success, TFLite status:  0
Output tensor => type:1 nr_dims:2
dim[0]: 1
dim[1]: 30
Result Tensor:
Tensor shape: [1 30]
Values:
  [1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000, 1.0000]`,
		},
		{
			name: "torch",
			args: []string{
				filepath.Join(paths.imagesDir, "example.jpg"),
				filepath.Join(paths.modelsDir, "torch", "cnn_trace.pt"),
				filepath.Join(paths.labelsDir, "imagenet.txt"),
			},
			wantOut: `Success!
Prediction: tench, Tinca tinca`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main := filepath.Join(tt.name, "main.go")

			cmd := exec.Command( //nolint:gosec
				"stdbuf",
				append([]string{"-oL", "-eL", "go", "run", main}, tt.args...)...,
			)
			cmdEnv := make([]string, len(env)+len(tt.env))
			copy(cmdEnv, env)
			copy(cmdEnv[len(env):], tt.env)
			cmd.Env = cmdEnv

			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("example failed: %v\n%s", err, out)
			}
			if tt.wantOut != "" && !strings.Contains(string(out), tt.wantOut) {
				t.Errorf("got %q, want %q", out, tt.wantOut)
			}
		})
	}
}
