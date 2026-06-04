package util

import (
	"reflect"
	"sync"
)

// MergePartial applies non-zero/non-nil fields from partial onto old and
// returns the result. Pointer and interface fields are only applied when
// non-nil; struct fields with exported members are merged recursively, while
// atomic value structs such as time.Time (no exported fields) are replaced
// wholesale when non-zero; all other fields (scalars, slices, maps) are applied
// only when not the zero value.
//
// This is the right tool for Discord partial-update payloads (MESSAGE_UPDATE,
// GUILD_MEMBER_UPDATE, etc.) where Discord omits unchanged fields, leaving them
// as zero values in the unmarshaled struct.
func MergePartial[T any](old, partial T) T {
	mergeValues(reflect.ValueOf(&old).Elem(), reflect.ValueOf(&partial).Elem())
	return old
}

func mergeValues(dst, src reflect.Value) {
	if dst.Kind() != reflect.Struct {
		return
	}
	for i := range dst.NumField() {
		d := dst.Field(i)
		s := src.Field(i)
		if !d.CanSet() {
			continue
		}
		switch d.Kind() {
		case reflect.Ptr, reflect.Interface:
			if !s.IsNil() {
				d.Set(s)
			}
		case reflect.Struct:
			// Value structs with no exported (settable) fields — e.g. time.Time —
			// cannot be merged field-by-field: the recursion would set nothing and a
			// non-zero source value would be silently dropped (GuildMember.JoinedAt,
			// Message.Timestamp, …). Treat them as atomic and replace wholesale when
			// the source is non-zero.
			if hasExportedFields(d.Type()) {
				mergeValues(d, s)
			} else if !s.IsZero() {
				d.Set(s)
			}
		default:
			if !s.IsZero() {
				d.Set(s)
			}
		}
	}
}

// exportedFieldsCache memoises hasExportedFields per struct type.
var exportedFieldsCache sync.Map // reflect.Type -> bool

// hasExportedFields reports whether t (a struct type) has at least one exported
// field, i.e. whether mergeValues can meaningfully recurse into it.
func hasExportedFields(t reflect.Type) bool {
	if v, ok := exportedFieldsCache.Load(t); ok {
		return v.(bool)
	}
	has := false
	for i := range t.NumField() {
		if t.Field(i).IsExported() {
			has = true
			break
		}
	}
	exportedFieldsCache.Store(t, has)
	return has
}
