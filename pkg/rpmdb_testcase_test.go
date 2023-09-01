package rpmdb

func intRef(i ...int) *int {
	if len(i) == 0 {
		return nil
	}
	return &i[0]
}

type commonPackageInfo struct {
	Epoch           *int
	Name            string
	Version         string
	Release         string
	Arch            string
	SourceRpm       string
	Size            int
	License         string
	Vendor          string
	Modularitylabel string
	Summary         string
	SigMD5          string
}

func toPackageInfo(pkgs []*commonPackageInfo) []*PackageInfo {
	pkgList := make([]*PackageInfo, 0, len(pkgs))
	for _, p := range pkgs {
		pkgList = append(pkgList, &PackageInfo{
			Epoch:           p.Epoch,
			Name:            p.Name,
			Version:         p.Version,
			Release:         p.Release,
			Arch:            p.Arch,
			SourceRpm:       p.SourceRpm,
			Size:            p.Size,
			License:         p.License,
			Vendor:          p.Vendor,
			Modularitylabel: p.Modularitylabel,
			Summary:         p.Summary,
			SigMD5:          p.SigMD5,
		})
	}

	return pkgList
}
