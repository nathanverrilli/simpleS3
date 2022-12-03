# simpleS3 for Backblaze B2
This is demonstration code to connect to Backblaze B2 object 
storage using the Amazon AWS S3 API.

It demonstrates connection, listing buckets, 
listing the objects in buckets, and downloading
an object (file). It assumes there is at least
one available bucket for the provided credentials,
and at least one object to download.

This program has been tested with Backblaze B2 storage,
and uses only capabilities available via that interface.

## Credential File
A sample credential file is included to 
show what it should look like,

## Possible Enhancements

_In other words, capabilities that are **not** demonstrated in this sample:_

* Uploading objects (files)
* Deleting objects (files)
* Creating buckets
* Deleting buckets
* Copying objects (files) within a bucket
* Copying object (files) from bucket to bucket

At this time, no further work on this project is planned.
