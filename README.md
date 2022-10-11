# godbf
[![GoDoc](https://godoc.org/github.com/NovikovRoman/godbf?status.svg)](https://godoc.org/github.com/NovikovRoman/godbf)
[![Go Report Card](https://goreportcard.com/badge/github.com/NovikovRoman/godbf)](https://goreportcard.com/report/github.com/NovikovRoman/godbf)
[![Build Status](https://travis-ci.com/NovikovRoman/godbf.svg?branch=master)](https://travis-ci.com/NovikovRoman/godbf)

> This fork from [LindsayBradford/go-dbf](https://github.com/LindsayBradford/go-dbf) has been heavily refactored and changed.

A pure Go library for reading and writing [dBase/xBase](http://en.wikipedia.org/wiki/DBase#File_formats) database files.

You can incorporate the library into your local workspace with the following 'go get' command:

```go
go get github.com/NovikovRoman/godbf
```

Code needing to call into the library needs to include the following import statement:
```go
import (
  "github.com/NovikovRoman/godbf"
)
```

Here is a very simple snippet of example 'load' code to get you going:
```go
  dbfTable, err := godbf.NewFromFile("exampleFile.dbf", nil)

  exampleList := make(ExampleList, dbfTable.NumberOfRecords())

  for i := 0; i < dbfTable.NumberOfRecords(); i++ {
    exampleList[i] = new(ExampleListEntry)

    exampleList[i].someColumnId, err = dbfTable.FieldValueByName(i, "SOME_COLUMN_ID")
  }
```

With encoding:
```go
  dbfTable, err := godbf.NewFromFile("exampleFileCp866.dbf", charmap.CodePage866)

  exampleList := make(ExampleList, dbfTable.NumberOfRecords())

  for i := 0; i < dbfTable.NumberOfRecords(); i++ {
    exampleList[i] = new(ExampleListEntry)

    exampleList[i].someColumnId, err = dbfTable.FieldValueByName(i, "SOME_COLUMN_ID")
  }
```

Further examples can be found by browsing the library's test suite. 
