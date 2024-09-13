package dto

import (
	"encoding/xml"
)

// ErrorResponse est utilisé pour formater les erreurs S3 en XML
type ErrorResponse struct {
	XMLName   xml.Name `xml:"Error"`
	Code      string   `xml:"Code"`
	Message   string   `xml:"Message"`
	Resource  string   `xml:"Resource"`
	RequestID string   `xml:"RequestId"`
}

// ListBucketsResponse est utilisé pour retourner la liste des buckets existants
type ListBucketsResponse struct {
	XMLName xml.Name     `xml:"ListAllMyBucketsResult"`
	Owner   Owner        `xml:"Owner"`
	Buckets []BucketInfo `xml:"Buckets>Bucket"`
}

type Owner struct {
	ID          string `xml:"ID"`
	DisplayName string `xml:"DisplayName"`
}

type BucketInfo struct {
	Name         string `xml:"Name"`
	CreationDate string `xml:"CreationDate"`
}

// ListObjectsResponse représente la réponse S3 pour la liste des objets dans un bucket
type ListObjectsResponse struct {
	XMLName     xml.Name     `xml:"ListBucketResult"`
	Name        string       `xml:"Name"`
	Prefix      string       `xml:"Prefix"`
	Marker      string       `xml:"Marker"`
	MaxKeys     int          `xml:"MaxKeys"`
	IsTruncated bool         `xml:"IsTruncated"`
	Contents    []ObjectInfo `xml:"Contents"`
}

type ObjectInfo struct {
	Key          string `xml:"Key"`
	LastModified string `xml:"LastModified"`
	ETag         string `xml:"ETag"`
	Size         int64  `xml:"Size"`
	StorageClass string `xml:"StorageClass"`
}
