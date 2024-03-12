// GithubRelease.go

// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    githubRelease, err := UnmarshalGithubRelease(bytes)
//    bytes, err = githubRelease.Marshal()

package model

import "github.com/goccy/go-json"

func Unmarshal(data []byte) (GithubRelease, error) {
	var r GithubRelease
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GithubRelease) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type GithubRelease struct {
	URL             string  `json:"url"`
	AssetsURL       string  `json:"assets_url"`
	UploadURL       string  `json:"upload_url"`
	HTMLURL         string  `json:"html_url"`
	ID              int64   `json:"id"`
	Author          Author  `json:"author"`
	NodeID          string  `json:"node_id"`
	TagName         string  `json:"tag_name"`
	TargetCommitish string  `json:"target_commitish"`
	Name            string  `json:"name"`
	Draft           bool    `json:"draft"`
	Prerelease      bool    `json:"prerelease"`
	CreatedAt       string  `json:"created_at"`
	PublishedAt     string  `json:"published_at"`
	Assets          []Asset `json:"assets"`
	TarballURL      string  `json:"tarball_url"`
	ZipballURL      string  `json:"zipball_url"`
	Body            string  `json:"body"`
}
