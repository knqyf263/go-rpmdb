package rpmdb

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
)

// updateTestCase is a helper function to update a given test case
// by running
func updateTestCase(t *testing.T, tc testCase) error {
	t.Helper()
	outputFile := filepath.Join(tc.testDir, "packages.csv")

	f, err := os.OpenFile(outputFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	cmdArgs := []string{
		"run",
		"-i",
		"-v", fmt.Sprintf("./%s:/mnt/export", tc.testDir),
		"--rm",
		"--entrypoint",
		"bash",
		tc.image,
		"/mnt/export/command.sh",
	}
	cmd := exec.Command("docker", cmdArgs...)
	fmt.Printf("Executing: %s\n", cmd.String())
	cmd.Stdout = f
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	fmt.Printf("Generated test data for %s\n", tc.name)
	return nil
}

// convert ¶-delimited CSV to a *commonPackageInfo
//
// ¶ is used because it is unlikely to be used in package metadata.
func csvToCommonPackageInfo(line []string) *commonPackageInfo {
	if len(line) < 11 {
		return nil
	}
	cpi := &commonPackageInfo{}
	if i, err := strconv.Atoi(line[0]); err == nil {
		cpi.Epoch = intRef(i)
	} else {
		cpi.Epoch = intRef()
	}
	cpi.Name = line[1]
	cpi.Version = line[2]
	cpi.Release = line[3]
	cpi.Arch = line[4]
	if line[4] == "(none)" {
		cpi.Arch = ""
	}
	cpi.SourceRpm = line[5]
	if line[5] == "(none)" {
		cpi.SourceRpm = ""
	}

	if i, err := strconv.Atoi(line[6]); err == nil {
		cpi.Size = i
	}
	cpi.License = line[7]
	cpi.Vendor = line[8]
	if line[8] == "(none)" {
		cpi.Vendor = ""
	}
	cpi.Summary = line[9]
	cpi.SigMD5 = line[10]
	if line[10] == "(none)" {
		cpi.SigMD5 = ""
	}
	// CentOS 5 zlib is missing the 2 trailing 0s?
	if len(cpi.SigMD5) == 32 && cpi.SigMD5[len(cpi.SigMD5)-2:len(cpi.SigMD5)] == "00" {
		cpi.SigMD5 = cpi.SigMD5[0 : len(cpi.SigMD5)-2]
	}

	if len(line) > 11 {
		if line[11] == "(none)" {
			cpi.Modularitylabel = ""
		} else {
			cpi.Modularitylabel = line[11]
		}
	}

	return cpi
}

// readTestCase reads a test case from a CSV file and returns
// a slice of *PackageInfo or an error
func readTestCase(t *testing.T, tc testCase) ([]*PackageInfo, error) {
	t.Helper()
	outputFile := filepath.Join(tc.testDir, "packages.csv")
	f, err := os.Open(outputFile)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(f)
	reader.Comma = '¶'
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	pkgs := []*commonPackageInfo{}
	for _, rec := range records {
		pkg := csvToCommonPackageInfo(rec)
		if pkg != nil {
			pkgs = append(pkgs, pkg)
		}
	}

	return toPackageInfo(pkgs), nil
}

func updatePackageFile(t *testing.T, tc packageContentTestCase, pi *PackageInfo) error {
	t.Helper()
	pfile := filepath.Join(tc.testDir, tc.wantPackageFile)
	f, err := os.OpenFile(pfile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(pi)
}

func readPackageFile(t *testing.T, tc packageContentTestCase) (*PackageInfo, error) {
	t.Helper()
	pfile := filepath.Join(tc.testDir, tc.wantPackageFile)
	f, err := os.Open(pfile)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()
	pi := &PackageInfo{}
	err = json.NewDecoder(f).Decode(pi)
	if err != nil {
		return nil, err
	}
	return pi, nil
}

func updateInstalledFiles(t *testing.T, tc packageContentTestCase, fi []FileInfo) error {
	t.Helper()
	pfile := filepath.Join(tc.testDir, tc.wantInstalledFilesFile)
	f, err := os.OpenFile(pfile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(fi)
}

func readInstalledFiles(t *testing.T, tc packageContentTestCase) ([]FileInfo, error) {
	t.Helper()
	pfile := filepath.Join(tc.testDir, tc.wantInstalledFilesFile)
	f, err := os.Open(pfile)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()
	fi := &[]FileInfo{}
	err = json.NewDecoder(f).Decode(fi)
	if err != nil {
		return nil, err
	}
	return *fi, nil
}

func updateInstalledFileNames(t *testing.T, tc packageContentTestCase, fn []string) error {
	t.Helper()
	pfile := filepath.Join(tc.testDir, tc.wantInstalledFileNamesFile)
	f, err := os.OpenFile(pfile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(fn)
}

func readInstalledFileNames(t *testing.T, tc packageContentTestCase) ([]string, error) {
	t.Helper()
	pfile := filepath.Join(tc.testDir, tc.wantInstalledFileNamesFile)
	f, err := os.Open(pfile)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()
	fn := &[]string{}
	err = json.NewDecoder(f).Decode(fn)
	if err != nil {
		return nil, err
	}
	return *fn, nil
}
