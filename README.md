# genddl
Generate RDB DDL by go struct

**THIS IS A ALPHA QUALITY RELEASE. API MAY CHANGE WITHOUT NOTICE.**

## Install

```
$ go get install github.com/mackee/go-genddl/cmd/genddl
```

## Example

Look [_example](https://github.com/mackee/go-genddl/blob/master/_example) sources.

## Usage

### 1. Write schema struct.

```go

//go:generate genddl -outpath=./mysql.sql -driver=mysql

//+table: person
type Person struct { //=> CREATE TABLE `person` (
	ID uint64 `db:"id,primarykey,autoincrement"` //=> `id`   BIGINT unsigned NOT NULL PRIMARY KEY AUTO_INCREMENT,
	Name string `db:"name,unique"`               //=> `name` VARCHAR(191) NOT NULL UNIQUE,
	Age uint32 `db:"age,null"`                   //=> `age`  INTEGER unsigned NULL
}
```

Default `NOT NULL`. You can add tag `null` if you want nullable column.

### 2. Run `go generate`

```
$ ls
person.go
$ go generate
$ person.go mysq.sql
```
