package db

import "github.com/lib/pq"

const driverName = "postgres"

func doAppend(buf []byte, prefix, value string) []byte {
	if value != "" {
		buf = append(buf, prefix...)
		buf = append(buf, value...)
	}
	return buf
}

func appendDriverErrorDetails(buf []byte, err error, stmt *Stmt, tx *Tx) []byte {
	pe, ok := err.(*pq.Error)
	if !ok {
		return buf
	}

	buf = append(buf, "\n\nDriver details\nError code: "...)
	buf = append(buf, pe.Code...)
	buf = append(buf, " ("...)
	buf = append(buf, pe.Code.Name()...)
	buf = append(buf, ")"...)

	buf = doAppend(buf, "\nMessage: ", pe.Message)
	buf = doAppend(buf, "\nDetail: ", pe.Detail)
	buf = doAppend(buf, "\nHint: ", pe.Hint)
	buf = doAppend(buf, "\nWhere:\n", pe.Where)
	buf = doAppend(buf, "\nSchema: ", pe.Schema)
	buf = doAppend(buf, "\nTable: ", pe.Table)
	buf = doAppend(buf, "\nColumn: ", pe.Column)
	buf = doAppend(buf, "\nConstraint: ", pe.Constraint)
	buf = doAppend(buf, "\nData type: ", pe.DataTypeName)
	return buf
}

func constraintFromDriverError(err error) (string, bool) {
	if pe, ok := err.(*pq.Error); ok && pe.Code.Class() == "23" {
		return pe.Constraint, true
	}

	return "", false
}
