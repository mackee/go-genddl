package index

//go:generate stringer -type=ForeignKeyOption
type ForeignKeyOption int

const (
	ForeignKeyDeleteRestrict ForeignKeyOption = iota
	ForeignKeyDeleteCascade
	ForeignKeyDeleteSetNull
	ForeignKeyDeleteSetDefault
	ForeignKeyDeleteNoAction
	ForeignKeyUpdateRestrict
	ForeignKeyUpdateCascade
	ForeignKeyUpdateSetNull
	ForeignKeyUpdateSetDefault
	ForeignKeyUpdateNoAction
	None
)

var ForeignKeyOptionConstMax = None

type Methods interface {
	PrimaryKey(columns ...interface{}) Definition
	Unique(columns ...interface{}) Definition
	Complex(columns ...interface{}) Definition
	ForeignKey(myColumn interface{}, foreignColumn interface{}, options ...ForeignKeyOption) Definition
	Spatial(columns ...interface{}) Definition
	Fulltext(columns ...string) Definition
}

type Definition string

// Columns groups columns for a composite foreign key:
// methods.ForeignKey(index.Columns(s.A, s.B), index.Columns(Parent{}.A, Parent{}.B)).
func Columns(columns ...interface{}) []interface{} {
	return columns
}
