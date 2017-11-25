go-xmp
===========

[![Build Status](https://travis-ci.org/trimmer-io/go-xmp.svg?branch=master)](https://travis-ci.org/trimmer-io/go-xmp)
[![GoDoc](https://godoc.org/trimmer.io/go-xmp?status.svg)](https://godoc.org/trimmer.io/go-xmp/xmp)


go-xmp is a native [Go](http://golang.org/) SDK for the [Extensible Metadata Platform](http://www.adobe.com/devnet/xmp.html) (XMP) as defined by the Adobe XMP Specification [Part 1](http://wwwimages.adobe.com/content/dam/Adobe/en/devnet/xmp/pdfs/XMP%20SDK%20Release%20cc-2016-08/XMPSpecificationPart1.pdf), [Part 2](http://wwwimages.adobe.com/content/dam/Adobe/en/devnet/xmp/pdfs/XMP%20SDK%20Release%20cc-2016-08/XMPSpecificationPart2.pdf) and [Part 3](http://wwwimages.adobe.com/content/dam/Adobe/en/devnet/xmp/pdfs/XMP%20SDK%20Release%20cc-2016-08/XMPSpecificationPart3.pdf), a.k.a ISO 16684-1:2011(E).

Features
--------

### Included metadata models

* XMP DublinCore (dc)
* XMP Media Management (xmpMM)
* XMP Dynamic Media (xmpDM)
* XMP Rights (xmpRights)
* XMP Jobs (xmpBJ)
* XMP Paged Text (xmpTPg)
* EXIF v2.3.1 (exif, exifEX)
* Adobe Camera Raw (crs)
* Creative Commons (cc)
* DJI Drones (dji)
* ID3 v2.2, v2.3, v2.4 (id3)
* iXML audio recorder (ixml)
* iTunes/MP4 (itunes)
* ISO/MP4 (mp4)
* Quicktime (qt)
* PhotoMechanic (pm)
* Tiff (tiff)
* Riff (riff)
* Photoshop (ps)
* PDF (pdf)

### Metadata models available under commercial license

* ACES Image Metadata
* AEScart
* ARRI Camera Metadata
* ASC CDL
* EBU Broadcast WAV
* Getty Images
* IPTC Core 1.2, IPTC Extension 1.3, IPTC Video Metadata 1.0
* Plus Licensing Metadata
* SMPTE DPX Image Metadata
* SMPTE MXF Metadata
* OpenEXR Image Header Metadata
* XMP Media Production SDK (Universal Metadata Container)


Documentation
-------------

- [API Reference](http://godoc.org/trimmer.io/go-xmp/xmp)
- [FAQ](https://github.com/trimmer-io/go-xmp/wiki/FAQ)

Installation
------------

Install go-xmp using the "go get" command:

    go get trimmer.io/go-xmp

The Go distribution is go-xmp's only dependency.

Examples
--------




Benchmarks
----------

```
go test ./test/ -bench=. -benchmem

goos: darwin
goarch: amd64
pkg: trimmer.io/go-xmp/test
BenchmarkUnmarshalXMP_5kB-8       5000      321524 ns/op     58071 B/op     1056 allocs/op
BenchmarkMarshalXMP_5kB-8         5000      270981 ns/op     61384 B/op      758 allocs/op
BenchmarkMarshalJSON_5kB-8        5000      338354 ns/op     91855 B/op     1023 allocs/op
BenchmarkUnmarshalJSON_5kB-8      5000      382196 ns/op     60387 B/op     1022 allocs/op
BenchmarkUnmarshalXMP_85kB-8       300     5152080 ns/op    902794 B/op    17779 allocs/op
BenchmarkMarshalXMP_85kB-8         300     4292143 ns/op    966356 B/op    12209 allocs/op
BenchmarkMarshalJSON_85kB-8        300     5378268 ns/op   1453004 B/op    16535 allocs/op
BenchmarkUnmarshalJSON_85kB-8      200     5512114 ns/op    880161 B/op    14497 allocs/op
```

XMP Marshal Benchmark using `premiere-cc.xmp`, a rather large xmpDM file with history, xmpMM:Pantry etc.

```
  Compression Results 417        mean               min                  max
  -----------------------------------------------------------------------------
        Original sizes       4013 (100.0)        918 (100.0)      86940 (100.0)
             XMP sizes       3545 ( 90.4)        723 ( 52.3)      78325 (147.0)
        XMP Gzip sizes       1177 ( 32.0)        369 (  8.4)       8086 ( 61.3)
      XMP Snappy sizes       1195 ( 32.5)        387 (  8.4)       8104 ( 62.9)
            JSON sizes       2147 ( 52.8)        389 ( 31.6)      65127 ( 91.7)
       JSON Gzip sizes        889 ( 23.7)        209 (  6.9)       7714 ( 50.5)
     JSON Snappy sizes        907 ( 24.3)        227 (  7.0)       7732 ( 50.6)
  -----------------------------------------------------------------------------
        XML->XMP times          371.674µs           74.091µs         4.816883ms
       XMP->JSON times          234.785µs            36.91µs         4.517591ms
        XMP->XML times          259.084µs           20.004µs         4.505254ms
        XMP Gzip times          214.342µs          105.685µs         1.036886ms
      XMP Gunzip times           59.948µs           19.924µs          285.113µs
      XMP Snappy times           28.325µs            7.975µs          234.161µs
    XMP Unsnappy times           25.841µs            7.899µs          195.265µs
       JSON Gzip times          197.268µs            93.86µs          968.985µs
     JSON Gunzip times           53.622µs           17.655µs         2.913516ms
     JSON Snappy times           23.581µs            7.864µs          221.674µs
   JSON Unsnappy times           37.215µs            7.856µs          398.361µs
```

Size matters when storing XMP in a database or sending documents over a network. Above is a quick comparison between common compression methods gzip and snappy regarding runtime and size for documents in the samples/ directory. What's also included is a comparison of the uncompressed documents in XMP/XML and XMP/JSON format. Original means the initial XMP document as stored in .xmp sidecar files. To be fair, some originals use padding, so the mean size distribution is larger than what go-xmp generated here because padding was turned off during write.

Contributing
------------

See [CONTRIBUTING.md](https://github.com/trimmer-io/go-xmp/blob/master/.github/CONTRIBUTING.md).


License
-------

go-xmp is available under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).

