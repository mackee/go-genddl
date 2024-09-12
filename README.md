# genddl
Generate RDB DDL by go struct

## Install

```
$ go install github.com/mackee/go-genddl/cmd/genddl@latest
```

## Example

Look [example](https://github.com/mackee/go-genddl/blob/master/_example) sources.

## Usage

### 1. Write schema struct.

```go

//go:generate genddl -outpath=./mysql.sql -driver=mysql

//genddl:table person
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

//genddl:table team
type Team struct { ... }

//genddl:table person
type Person struct { ... }

func (s Person) _schemaIndex(methods index.Methods) []index.Definition {
	return []index.Definition{
		methods.PrimaryKey(s.ID, s.CreatedAt),  //=> PRIMARY KEY (`id`, `created_at`),
		methods.Unique(s.UserCode, s.Type),     //=> UNIQUE (`usercode`, `type`),
		methods.ForeignKey(s.TeamID, Team{}.ID, index.ForeignKeyDeleteCascade, index.ForeignKeyUpdateSetDefault),
		    //=> FOREGIN KEY (`team_id`) REFERENCES team(`id`) ON DELETE CASCADE ON UPDATE SET DEFAULT
		methods.Complex(s.Age, s.Name),         //=> CREATE INDEX person_age_name (`age`, `name`);
	}
}
```

### CLI Options

```
Usage of genddl:
  -driver string
        target driver name. support mysql, pg, sqlite3 (default "mysql")
  -foreignkeyname
        Provides a name for the definition of a foreign-key.
  -innerindex create table
        Placement of index definition. If this specified, the definition was placement inner of create table
  -outerforeignkey
        Placement of foreign key definition. If this specified, the definition was placement end of DDL file.
  -outeruniquekey
        Placement of unique key definition. If this specified, the definition was placement outer of CREATE TABLE.
  -outpath string
        schema target path
  -schemadir string
        schema declaretion directory
  -tablecollate string
        Provides a collate for the definition of tables.
  -uniquename
        Provides a name for the definition of a unique index.
  -withoutdroptable
        If this specified, the DDL file does not contain DROP TABLE statement.
```
