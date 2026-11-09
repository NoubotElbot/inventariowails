package sqlite

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

// whereBuilder accumulates SQL WHERE conditions and their positional arguments.
// The zero value is ready to use; conditions are joined with AND.
//
// Example:
//
//	var wb whereBuilder
//	wb.And("deleted_at IS NULL")
//	if search != "" {
//	    wb.And("name LIKE ?", "%"+search+"%")
//	}
//	query := "SELECT ... FROM t " + wb.SQL() + " ORDER BY id LIMIT ? OFFSET ?"
//	err   := db.Select(&rows, query, append(wb.Args(), limit, offset)...)
type whereBuilder struct {
	conds []string
	args  []any
}

// And appends a condition and its positional ? arguments.
func (w *whereBuilder) And(cond string, args ...any) {
	w.conds = append(w.conds, cond)
	w.args = append(w.args, args...)
}

// In appends a "col IN (...)" condition for a non-empty int64 slice.
// It uses sqlx.In to expand the single ? placeholder. Guard with len(ids) > 0.
func (w *whereBuilder) In(col string, ids []int64) error {
	cond, args, err := sqlx.In(col+" IN (?)", ids)
	if err != nil {
		return err
	}
	w.conds = append(w.conds, cond)
	w.args = append(w.args, args...)
	return nil
}

// SQL returns "WHERE cond1 AND cond2 ..." or "" when no conditions were added.
func (w *whereBuilder) SQL() string {
	if len(w.conds) == 0 {
		return ""
	}
	return "WHERE " + strings.Join(w.conds, " AND ")
}

// Args returns a copy of the accumulated arguments so callers can safely
// append pagination args (LIMIT, OFFSET) without aliasing internal state.
func (w *whereBuilder) Args() []any {
	out := make([]any, len(w.args))
	copy(out, w.args)
	return out
}
