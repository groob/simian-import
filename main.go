package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/groob/plist"

	"cloud.google.com/go/datastore"
)

type DSPackageInfo struct {
	Plist        string    `datastore:"_plist,noindex"`
	BlobKey      string    `datastore:"blobstore_key,omitempty"`
	Catalogs     []string  `datastore:"catalogs"`
	MModAccess   []string  `datastore:"manifest_mod_access"`
	Created      time.Time `datastore:"created"`
	Filename     string    `datastore:"filename"`
	Manifests    string    `datastore:"manifests,omitempty"`
	MTime        time.Time `datastore:"mtime"`
	MunkiName    string    `datastore:"munki_name"`
	Name         string    `datastore:"name"`
	PkgDataSHA   string    `datastore:"pkgdata_sha256"`
	User         string    `datastore:"user"`
	InstallTypes []string  `datastore:"install_types"`
}

func main() {
	var (
		flProject  = flag.String("gcp.project", "", "GCP Project Name")
		flPkgsinfo = flag.String("pkgsinfo", "", "path to pkgsinfo")
	)

	flag.Parse()

	ctx := context.Background()
	dsClient, err := datastore.NewClient(ctx, *flProject)
	if err != nil {
		log.Fatal(err)
	}

	if *flPkgsinfo == "" {
		flag.Usage()
		log.Fatal("must provide pkgsinfo file path")
	}

	info, err := ioutil.ReadFile(*flPkgsinfo)
	if err != nil {
		log.Fatal(err)
	}
	p, err := pkgInfoFromPlist(info)
	if err != nil {
		log.Fatal(err)
	}

	k := datastore.NameKey("PackageInfo", p.Filename, nil)
	if _, err := dsClient.Put(ctx, k, p); err != nil {
		log.Fatal(err)
	}

}

func pkgInfoFromPlist(plistData []byte) (*DSPackageInfo, error) {
	var p PkgsInfo
	if err := plist.Unmarshal(plistData, &p); err != nil {
		return nil, err
	}
	pkgName := p.InstallerItemLocation
	dsp := &DSPackageInfo{
		Name:         p.Name,
		Plist:        string(plistData),
		InstallTypes: []string{"managed_updates"},
		MModAccess:   []string{"support", "security"},
		PkgDataSHA:   p.InstallerItemHash,
		Created:      time.Now(),
		MTime:        time.Now(),
		MunkiName:    strings.TrimSuffix(filepath.Base(pkgName), filepath.Ext(pkgName)),
		Filename:     pkgName,
	}
	return dsp, nil

}
