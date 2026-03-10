// Copyright IBM Corp. 2013, 2023
// SPDX-License-Identifier: MIT

package clause

type Delete struct {
	Modifier string
}

func (d Delete) Name() string {
	return "DELETE"
}

func (d Delete) Build(builder Builder) {
	builder.WriteString("DELETE")

	if d.Modifier != "" {
		builder.WriteByte(' ')
		builder.WriteString(d.Modifier)
	}
}

func (d Delete) MergeClause(clause *Clause) {
	clause.Name = ""
	clause.Expression = d
}
