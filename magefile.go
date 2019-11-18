//+build mage

package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var Default = Build

var GOPATH = os.Getenv("GOPATH")

func gopath(path string) string {
	return GOPATH + path
}

func getPackages() ([]string, error) {
	out, err := sh.Output("go", "list", "./...")
	if err != nil {
		return nil, err
	}

	return strings.Split(out, "\n"), nil
}

func getPackageFiles(name string) ([]string, error) {
	out, err := sh.Output("go", "list", "-f", "{{range .GoFiles}}{{$.Dir}}/{{.}} {{end}}", name)
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.Trim(out, " "), " "), nil
}

// Builds the mmctl binary
func Build() {
	mg.SerialDeps(Vendor, Check)

	sh.RunV("go", "build", "-mod=vendor")
}

// Installs the mmctl binary in the GOPATH
func Install() {
	mg.SerialDeps(Vendor, Check)

	sh.RunV("go", "install", "-mod=vendor")
}

// Packages mmctl in the build directory
func Package() error {
	mg.SerialDeps(Vendor, Check)

	sh.RunV("mkdir", "-p", "build")

	fmt.Println("Building Linux amd64")
	if err := sh.RunWith(map[string]string{"GOOS": "linux", "GOARCH": "amd64"}, "go", "build", "-mod=vendor"); err != nil {
		return err
	}
	sh.RunV("tar", "cf", "build/linux_amd64.tar", "mmctl")

	fmt.Println("Building OSX amd64")
	if err := sh.RunWith(map[string]string{"GOOS": "darwin", "GOARCH": "amd64"}, "go", "build", "-mod=vendor"); err != nil {
		return err
	}
	sh.RunV("tar", "cf", "build/darwin_amd64.tar", "mmctl")

	fmt.Println("Build Windows amd64")
	if err := sh.RunWith(map[string]string{"GOOS": "windows", "GOARCH": "amd64"}, "go", "build", "-mod=vendor"); err != nil {
		return err
	}
	sh.RunV("zip", "build/windows_amd64.zip", "mmctl.exe")

	sh.RunV("rm", "mmctl", "mmctl.exe")
	return nil
}

// Runs gofmt through the mmctl codebase
func Gofmt() error {
	fmt.Println("Running gofmt")

	packages, err := getPackages()
	if err != nil {
		return err
	}

	for _, p := range packages {
		fmt.Printf("Checking %s\n", p)
		files, err := getPackageFiles(p)
		if err != nil {
			return err
		}

		if len(files) > 0 {
			args := append([]string{"-d", "-s"}, files...)
			out, err := sh.Output("gofmt", args...)
			if err != nil {
				return err
			}

			if out != "" {
				fmt.Println(out)
				return errors.New("Gofmt failure")
			}
		}
	}
	fmt.Println("Gofmt success")
	return nil
}

// Runs govet through the mmctl codebase
func Govet() error {
	fmt.Println("Running govet")
	packages, err := getPackages()
	if err != nil {
		return err
	}

	sh.RunV("go", "get", "golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow")

	args := append([]string{"vet"}, packages...)
	sh.RunV("go", args...)

	args = append([]string{"vet", "-vettool=" + gopath("/bin/shadow")}, packages...)
	sh.RunV("go", args...)

	fmt.Println("Govet success")
	return nil
}

// Runs the unit test suite
func UnitTest() error {
	fmt.Println("Running unit tests")
	packages, err := getPackages()
	if err != nil {
		return err
	}
	args := append([]string{"test", "-race", "-v", "-count", "1"}, packages...)
	sh.RunV("go", args...)
	return nil
}

// Runs all the test suites
func Test() error {
	fmt.Println("Running all tests")
	packages, err := getPackages()
	if err != nil {
		return err
	}
	args := append([]string{"test", "-race", "-v", "-count", "1", "-tags", "e2e"}, packages...)
	sh.RunV("go", args...)
	return nil
}

// Runs all checks on the mmctl codebase
func Check() {
	mg.SerialDeps(Gofmt, Govet)
}

// Downloads all the dependencies to the vendor folder
func Vendor() {
	sh.RunV("go", "mod", "vendor")
}
