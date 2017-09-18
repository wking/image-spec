## OCI Image Layout Specification

This is version 1.1.0 of this specification.

* The OCI Image Layout is a slash separated layout of OCI content-addressable blobs and [location-addressable](https://en.wikipedia.org/wiki/Content-addressable_storage#Content-addressed_vs._location-addressed) references (refs).
* This layout MAY be used in a variety of different transport mechanisms: archive formats (e.g. tar, zip), shared filesystem environments (e.g. nfs), or networked file fetching (e.g. http, ftp, rsync).

Given an image layout and a ref, a tool can create an [OCI Runtime Specification bundle](https://github.com/opencontainers/runtime-spec/blob/v1.0.0/bundle.md) by:

* Following the ref to find a [manifest](manifest.md#image-manifest), possibly via an [image index](image-index.md)
* [Applying the filesystem layers](layer.md#applying) in the specified order
* Converting the [image configuration](config.md) into an [OCI Runtime Specification `config.json`](https://github.com/opencontainers/runtime-spec/blob/v1.0.0/config.md)

# Content

An image layout consists of an `oci-layout` file which MUST exist and be a valid [layout object](#layout-object).

## Layout object

This section defines the `application/vnd.oci.layout.header.v1+json` [media type](media-types.md).

Content of this type MUST be a JSON object.

The object MUST include an `imageLayoutVersion` entry to provide the version of the image-layout in use.
The value will align with the OCI Image Specification version at the time changes to the layout are made, and will pin a given version until changes to the image layout are required.

The object MUST include a `refEngines` entry, [as defined by version 0.1 of the OCI Ref-Engine Discovery specification][ref-engines-objects].

The object MAY include a `casEngines` entry, [as defined by version 0.1 of the OCI Ref-Engine Discovery specification][ref-engines-objects].

### Example Layout Object

```json,title=OCI%20Layout&mediatype=application/vnd.oci.layout.header.v1%2Bjson
{
    "imageLayoutVersion": "1.1.0",
    "refEngines": [
        {
            "protocol": "oci-index-template-v1",
            "uri": "index.json"
        }
    ],
    "casEngines": [
        {
            "protocol": "oci-cas-template-v1",
            "uri": "blobs/{algorithm}/{encoded}"
        }
    ]
}
```

Where the above `refEngines` and `casEngines` entries are used, the remainder of a 1.1.0 layout will be identical to [a 1.0.0 layout][layout-1.0.0].

### Example Layout

This is an example image layout corresponding to the [example layout object](#example-layout-object)

```
$ cd example.com/app/
$ find . -type f
./blobs/sha256/3588d02542238316759cbf24502f4344ffcc8a60c803870022f335d1390c13b4
./blobs/sha256/4b0bc1c4050b03c95ef2a8e36e25feac42fd31283e8c30b3ee5df6b043155d3c
./blobs/sha256/7968321274dc6b6171697c33df7815310468e694ac5be0ec03ff053bb135e768
./index.json
./oci-layout
```

[layout-1.0.0]: https://github.com/opencontainers/image-spec/blob/v1.0.0/image-layout.md
[ref-engines-objects]: https://github.com/xiekeyang/oci-discovery/blob/44ec3cf3113e29a743ad04220ccb7ff5197dab2a/ref-engine-discovery.md#ref-engines-objects
