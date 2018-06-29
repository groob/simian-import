package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"time"

	"cloud.google.com/go/datastore"
)

func main() {
	var (
		flProject  = flag.String("gcp.project", "", "GCP Project Name")
		flTagname  = flag.String("tag.name", "", "name of tag")
		flComputer = flag.String("tag.computer", "", "add a computer")
		flUser     = flag.String("user", "", "your email")
	)

	flag.Parse()

	ctx := context.Background()
	dsClient, err := datastore.NewClient(ctx, *flProject)
	if err != nil {
		log.Fatal(err)
	}

	tag := NewTag(*flTagname, *flUser, *flProject)
	tag.addKind("Computer", *flComputer)

	if err := tag.marshalKeys(); err != nil {
		log.Fatal(err)
	}

	k := datastore.NameKey("Tag", *flTagname, nil)
	if _, err := dsClient.Put(ctx, k, tag); err != nil {
		log.Fatal(err)
	}
}

type Tag struct {
	Keys    []string  `datastore:"keys"`
	Time    time.Time `datastore:"mrtime"`
	User    string    `datastore:"user"`
	key     *key      `datastore:"-"`
	tagName string    `datastore:"-"`
}

func NewTag(tagName, user, projectName string) *Tag {
	tag := Tag{
		key: &key{
			Values: []value{value{KeyValue: keyValue{PartitionID: partitionId{ProjectID: projectName}}}},
		},
		User:    user,
		Time:    time.Now().UTC(),
		tagName: tagName,
	}
	return &tag
}

func (t *Tag) marshalKeys() error {
	out, err := json.MarshalIndent(t.key, "", "  ")
	t.Keys = []string{string(out)}
	return err
}

func (t *Tag) addKind(kind, name string) {
	// check if present
	for _, kv := range t.key.Values[0].KeyValue.Path {
		k, ok := kv["kind"]
		if !ok || k != kind {
			continue
		}

		value, ok := kv["name"]
		if !ok {
			panic("expected name next to kind")
		}
		if name == value {
			log.Printf("kind=%s name=%s already exists in tag", name, kind)
			return
		}
	}

	t.key.Values[0].KeyValue.Path = append(t.key.Values[0].KeyValue.Path, map[string]string{
		"kind": kind,
		"name": name,
	})
}

type key struct {
	Values []value `json:"values"`
}

type value struct {
	KeyValue keyValue `json:"keyValue"`
}

type keyValue struct {
	PartitionID partitionId         `json:"partitionId"`
	Path        []map[string]string `json:"path"`
}

type partitionId struct {
	ProjectID string `json:"projectId"`
}
