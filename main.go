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
		flPkgName  = flag.String("pkgname", "", "name of pkg to retrieve")
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
	p, err := pkgInfoFromPlist(info, *flPkgName)
	if err != nil {
		log.Fatal(err)
	}

	k := datastore.NameKey("PackageInfo", *flPkgName, nil)
	if _, err := dsClient.Put(ctx, k, p); err != nil {
		log.Fatal(err)
	}

}

func pkgInfoFromPlist(plistData []byte, pkgName string) (*DSPackageInfo, error) {
	var p PkgsInfo
	if err := plist.Unmarshal(plistData, &p); err != nil {
		return nil, err
	}
	dsp := &DSPackageInfo{
		Name:       p.Name,
		Plist:      string(plistData),
		Catalogs:   p.Catalogs,
		PkgDataSHA: p.InstallerItemHash,
		Created:    time.Now(),
		MTime:      time.Now(),
		MunkiName:  strings.TrimSuffix(filepath.Base(pkgName), filepath.Ext(pkgName)),
		Filename:   pkgName,
	}
	return dsp, nil

}
