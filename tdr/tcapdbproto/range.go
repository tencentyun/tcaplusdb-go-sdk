package tcapdbproto

import (
	"sort"
	"sync"

	pref "google.golang.org/protobuf/reflect/protoreflect"
)

type messageField struct {
	fd pref.FieldDescriptor
	v  pref.Value
}

var messageFieldPool = sync.Pool{
	New: func() interface{} { return new([]messageField) },
}

type (
	// FieldRnger is an interface for visiting all fields in a message.
	// The protoreflect.Message type implements this interface.
	FieldRanger interface{ Range(VisitField) }
	// VisitField is called everytime a message field is visited.
	VisitField = func(pref.FieldDescriptor, pref.Value) bool
)

// RangeFields iterates over the fields of fs according to the specified order.
func RangeFields(fs FieldRanger, less FieldOrder, fn VisitField) {
	if less == nil {
		fs.Range(fn)
		return
	}

	// Obtain a pre-allocated scratch buffer.
	p := messageFieldPool.Get().(*[]messageField)
	fields := (*p)[:0]
	defer func() {
		if cap(fields) < 1024 {
			*p = fields
			messageFieldPool.Put(p)
		}
	}()

	// Collect all fields in the message and sort them.
	fs.Range(func(fd pref.FieldDescriptor, v pref.Value) bool {
		fields = append(fields, messageField{fd, v})
		return true
	})
	sort.Slice(fields, func(i, j int) bool {
		return less(fields[i].fd, fields[j].fd)
	})

	// Visit the fields in the specified ordering.
	for _, f := range fields {
		if !fn(f.fd, f.v) {
			return
		}
	}
}

type mapEntry struct {
	k pref.MapKey
	v pref.Value
}

var mapEntryPool = sync.Pool{
	New: func() interface{} { return new([]mapEntry) },
}

type (
	// EntryRanger is an interface for visiting all fields in a message.
	// The protoreflect.Map type implements this interface.
	EntryRanger interface{ Range(VisitEntry) }
	// VisitEntry is called everytime a map entry is visited.
	VisitEntry = func(pref.MapKey, pref.Value) bool
)

// RangeEntries iterates over the entries of es according to the specified order.
func RangeEntries(es EntryRanger, less KeyOrder, fn VisitEntry) {
	if less == nil {
		es.Range(fn)
		return
	}

	// Obtain a pre-allocated scratch buffer.
	p := mapEntryPool.Get().(*[]mapEntry)
	entries := (*p)[:0]
	defer func() {
		if cap(entries) < 1024 {
			*p = entries
			mapEntryPool.Put(p)
		}
	}()

	// Collect all entries in the map and sort them.
	es.Range(func(k pref.MapKey, v pref.Value) bool {
		entries = append(entries, mapEntry{k, v})
		return true
	})
	sort.Slice(entries, func(i, j int) bool {
		return less(entries[i].k, entries[j].k)
	})

	// Visit the entries in the specified ordering.
	for _, e := range entries {
		if !fn(e.k, e.v) {
			return
		}
	}
}
