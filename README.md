# genddl
Generate RDB DDL by go struct

**THIS IS A ALPHA QUALITY RELEASE. API MAY CHANGE WITHOUT NOTICE.**

## Install

```
$ go get install github.com/mackee/go-genddl/cmd/genddl
```

## Example

Look [example](https://github.com/mackee/go-genddl/blob/master/_example) sources.

## Usage

### 1. Write schema struct.

```go

//go:generate genddl -outpath=./mysql.sql -driver=mysql

//+table: person
type Person struct { //=> CREATE TABLE `person` (
	ID uint64 `db:"id,primarykey,autoincrement"` //=> `id`            BIGINT unsigned NOT NULL PRIMARY KEY AUTO_INCREMENT,
	Name string `db:"name,unique"`               //=> `name`          VARCHAR(191) NOT NULL UNIQUE,
	Age uint32 `db:"age,null"`                   //=> `age`           INTEGER unsigned NULL,
	UserCode string `db:"usercode"`              //=> `usercode`      VARCHAR(191) NOT NULL,
	Type uint32 `db:"type"`                      //=> `type`          INTEGER unsigned NOT NULL,
	TeamID uint64 `db:"team_id"`                 //=> `team_id`       BIGINT unsigned NOT NULL,
	CreatedAt time.Time `db:"created_at"`        //=> `created_at`    DATETIME NOT NULL
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

## Other Features

### Indexes

If you want to set indexes, write method for schema struct. It name must be `_schemaIndex`.

Example:
```go
import (
	"github.com/mackee/go-genddl/index"
)

//+table:team
type Team struct { ... }

//+table:person
type Person struct { ... }

func (s Person) _schemaIndex(methods index.Methods) []index.Definition {
	return []index.Definition{
		methods.PrimaryKey(s.ID, s.CreatedAt),  //=> PRIMARY KEY (`id`, `created_at`),
		methods.Unique(s.UserCode, s.Type),     //=> UNIQUE (`usercode`, `type`),
		methods.Complex(s.Age, s.Name),         //=> INDEX (`age`, `name`),
		methods.ForeignKey(s.TeanID, Team{}.ID, index.ForeignKeyDeleteCascade, index.ForeignKeyUpdateSetDefault),
		    //=> FOREGIN KEY (`team_id`) REFERENCES team(`id`) ON DELETE CASCADE ON UPDATE SET DEFAULT
	}
}
```
