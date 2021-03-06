package lang

import (
	"go/types"
	"os"
	"bytes"
	"bufio"

	"fmt"
)

/*
	NOTE: I've made the design decision for now to mimic static methods as best I can,
	which in this case means creating a private class and a single public object for that class.

	In practice I think that RT.java is more or less a catchall for a host of static methods
	that could just as easily be generally pure functions. I find having an empty struct
	associated with them to be easier to read, semantically. I don't think this is idiomatic,
	but then, neither is having such a massive package.

	Also, in case it wasn't obvious, RT stands for runtime.
*/

type rt struct{}
var RT = rt{}

var LOADER_SUFFIX string = "__init"


var DEFAULT_IMPORTS IPersistentMap = RT.Map(
	InternSymbol("Boolean"), types.Bool,
	InternSymbol("Byte"), types.Byte,
	InternSymbol("Rune"), types.Rune, // NOTE "Character" in JVM Clojure
	// TODO ... there's a lot of stuff in here. It basically maps all names to underlying types.
)

func (_ *rt) ReadTrueFalseUnknown(s string) interface{} {
	if s == "true" {
		return true
	} else if s == "false" {
		return false
	}
	return InternKeywordByNsName("unknown")
}

// GetEnvWithDefault is like System.getProperty in Java
func (_ *rt) GetEnvWithDefault(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	return defaultVal
}

var CLOJURE_NS = FindOrCreateNamespace(InternSymbol("clojure.core"))
var OUT *Var = InternVar(CLOJURE_NS, InternSymbol("*out*"), os.Stdout).SetDynamic()
// TODO: This should actually be a LineNumberingPushbackReader for os.Stdin
var IN *Var = InternVar(CLOJURE_NS, InternSymbol("*in*"), os.Stdin).SetDynamic()
var ERR *Var = InternVar(CLOJURE_NS, InternSymbol("*err*"), os.Stderr).SetDynamic()
var TAG_KEY *Keyword = InternKeywordByNsName("tag")
var CONST_KEY *Keyword = InternKeywordByNsName("const")
var AGENT *Var = InternVar(CLOJURE_NS, InternSymbolByNsname("*agent*"), nil).SetDynamic()
var readeval interface{} = RT.ReadTrueFalseUnknown(RT.GetEnvWithDefault("clojure.read.eval", "true"))
var READEVAL = InternVar(CLOJURE_NS, InternSymbol("*read-eval*"), readeval).SetDynamic()
var DATA_READERS = InternVar(CLOJURE_NS, InternSymbol("*data-readers*"), RT.Map()).SetDynamic()
var DEFAULT_DATA_READERS_FN = InternVar(CLOJURE_NS, InternSymbol("*default-data-readers-fn*"), RT.Map()).SetDynamic()
var DEFAULT_DATA_READERS = InternVar(CLOJURE_NS, InternSymbol("*default-data-readers*"), RT.Map()).SetDynamic()

// TODO...more
var LINE_KEY *Keyword = InternKeywordByNsName("line")
var COLUMN_KEY *Keyword = InternKeywordByNsName("column")
var FILE_KEY *Keyword = InternKeywordByNsName("file")
var DECLARE_KEY *Keyword = InternKeywordByNsName("declared")
var DOC_KEY *Keyword = InternKeywordByNsName("doc")
// TODO...more

var LOAD_FILE *Symbol = InternSymbolByNsname("load-file")
var IN_NAMESPACE *Symbol = InternSymbolByNsname("in-ns")
var NAMESPACE *Symbol = InternSymbolByNsname("ns")
var CMD_LINE_ARGS *Var = InternVar(CLOJURE_NS, InternSymbolByNsname("*command-line-args*"), nil).SetDynamic()
var CURRENT_NS *Var = InternVar(CLOJURE_NS, InternSymbolByNsname("*ns*"), CLOJURE_NS).SetDynamic()

var PRINT_META *Var = InternVar(CLOJURE_NS, InternSymbolByNsname("*print-meta*"), false).SetDynamic()
var PRINT_READABLY *Var = InternVar(CLOJURE_NS, InternSymbolByNsname("*print-readably*"), true).SetDynamic()
var PRINT_DUP *Var = InternVar(CLOJURE_NS, InternSymbolByNsname("*print-dup*"), false).SetDynamic()
// TODO...more

var IN_NS_VAR *Var = InternVar(CLOJURE_NS, IN_NAMESPACE, false)
var NS_VAR *Var = InternVar(CLOJURE_NS, NAMESPACE, false)
var FN_LOADER_VAR *Var = InternVar(CLOJURE_NS, InternSymbolByNsname("*fn-loader*"), nil).SetDynamic()
var PRINT_INITIALIZED *Var = InternVarByNsAndSym(CLOJURE_NS, InternSymbolByNsname("print-initialized"))
// TOOD...more


func (_ *rt) EMPTY_ARRAY() []interface{} {
	return make([]interface{}, 1)
}


// Overloaded with optional init argument.
func (_ *rt) Var(ns string, name string, init interface{}) *Var {
	return InternVar(FindOrCreateNamespace(InternSymbolByNsname(ns)), InternSymbolByNsname(name), init)
}

func (_ *rt) Map(init... interface{}) IPersistentMap {
	if init == nil {
		return EMPTY_PERSISTENT_ARRAY_MAP
	} else if len(init) <= HASHTABLE_THRESHOLD {
		return CreatePersistentArrayMapWithCheck(init)
	}
	return CreatePersistentHashMapWithCheck(init...)
}

func (_ *rt) IsReduced(r interface{}) bool {
	switch r.(type) {
	case Reduced:
		return true
	default:
		return false
	}
}

// TODO....so much

func (_ *rt) Seq(coll interface{}) ISeq {
	// TODO: This definitely doesn't work at the moment
	switch c := coll.(type) {
	case ASeq:
		return c
	case LazySeq:
		return c.Seq()
	}

	return RT.seqFrom(coll)
}

func (_ *rt) seqFrom(coll interface{}) ISeq {
	switch c := coll.(type) {
	case Seqable:
		return c.Seq()
		// TODO: the rest of this
	}
	// TODO
	return nil
}

func (_ *rt) SubVec(v IPersistentVector, start int, end int) IPersistentVector {
	if end < start || start < 0 || end > v.Count() {
		panic(IndexOutOfBoundsException)
	}
	if start == end {
		return nil
		// return EMPTY_PERSISTENT_VECTOR TODO
	}
	return &SubVector{} // TODO
}



// unordered

// TODO
func (_ *rt) Count(o interface{}) int {
	return 0
}

func (_ *rt) PrintString(o interface{}) string {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	RT.Print(o, w)
	w.Flush()
	return b.String()

}

// TODO
func (_ *rt) ReadString(s string, opts interface{}) interface{} {
	return nil
}

// Figure out what the relevant string representation of the object is and write it into
// the writer (expecting a buffer)
func (_ *rt) Print(x interface{}, w *bufio.Writer) {

	// TODO a lot of stuff

	switch obj := x.(type) {
	// MORE STUFF
	case ISeq:
		w.WriteRune('(')
		RT.printInnerSeq(RT.Seq(x), w)
		w.WriteRune(')')
	case IPersistentList:
		w.WriteRune('(')
		RT.printInnerSeq(RT.Seq(x), w)
		w.WriteRune(')')
	case IPersistentMap:
		w.WriteRune('{')
		for s := RT.Seq(obj); s != nil; s = s.Next() {
			e := s.First().(IMapEntry)
			RT.Print(e.Key(), w)
			w.WriteRune(' ')
			RT.Print(e.Val(), w)
			if s.Next() != nil {
				w.WriteString(", ")
			}
		}
		w.WriteRune('}')
	case fmt.Stringer:
		w.WriteString(obj.String())
	default:

		w.WriteString(fmt.Sprintf("%v", obj))
	}

}

func (_ *rt) printInnerSeq(x ISeq, w *bufio.Writer) {
	for s := x; s != nil; s = s.Next() {
		RT.Print(s.First(), w)
		if s.Next() != nil {
			w.WriteRune(' ')
		}
	}
}


// TODO
func (_ *rt) ToArray(coll interface{}) []interface{} {
	return nil
}

func (_ *rt) SeqToArray(seq ISeq) []interface{} {
	len := RT.Length(seq)
	ret := make([]interface{}, len)
	for i := 0; seq != nil; seq = seq.Next() {
		ret[i] = seq.First()
		i++
	}
	return ret
}

func (_ *rt) SeqToPassedArray(seq ISeq, passed []interface{}) []interface{} {
	dest := passed
	length := RT.Count(seq)
	if length > len(dest) {
		dest = make([]interface{}, length) // NOTE: This does some reflection in JVM Clojure.
	}
	for i := 0; seq != nil; seq = seq.Next() {
		dest[i] = seq.First()
		i++
	}
	if length < len(passed) {
		dest[length] = nil
	}
	return dest
}

func (_ *rt) Length (list ISeq) int {
	i := 0
	for c := list; c != nil; c = c.Next() {
		i++
	}
	return i
}

func (_ *rt) Keys(coll interface{}) ISeq {
	switch c := coll.(type) {
	case IPersistentMap:
		return CreateKeySeqFromMap(c)
	default:
		return CreateKeySeq(RT.Seq(coll))
	}
}

func (_ *rt) Cons(x interface{}, coll interface{}) ISeq {
	if coll == nil {
		return &PersistentList{
			_first: x,
			_rest:  nil,
			_count: 1,
		}
	}
	switch c := coll.(type) {
	case ISeq:
		return &Cons{
			_first: x,
			_more:  c,
		}
	default:
		return &Cons{
			_first: x,
			_more:  RT.Seq(coll),
		}
	}

}

func (_ *rt) First(x interface{}) interface{} {
	switch o := x.(type) {
	case ISeq:
		return o.First()
	}
	seq := RT.Seq(x)
	if seq == nil {
		return nil
	}
	return seq.First()
}

func (_ *rt) Second(x interface{}) interface{} {
	return RT.First(RT.Next(x))
}

func (_ *rt) Third(x interface{}) interface{} {
	return RT.First(RT.Next(RT.Next(x)))
}

func (_ *rt) Fourth(x interface{}) interface{} {
	return RT.First(RT.Next(RT.Next(RT.Next(x))))
}

func (_ *rt) Next(x interface{}) ISeq {
	switch s := x.(type) {
	case ISeq:
		return s.Next()
	}
	seq := RT.Seq(x)
	if seq == nil {
		return nil
	}
	return seq.Next()
}

func (_ *rt) More(x interface{}) ISeq {
	switch s := x.(type) {
	case ISeq:
		return s.More()
	}
	seq := RT.Seq(x)
	if seq == nil {
		return EMPTY_PERSISTENT_LIST
	}
	return seq.More()
}

func (_ *rt) Peek(x interface{}) interface{} {
	if x == nil {
		return nil
	}
	return x.(IPersistentStack).Peek()
}

func (_ *rt) Pop(x interface{}) interface{} {
	if x == nil {
		return nil
	}
	return x.(IPersistentStack).Pop()
}

func (_ *rt) Get(coll interface{}, key interface{}, notFound interface{}) interface{} {
	switch coll.(type) {
	case ILookup:
		return coll.(ILookup).ValAt(key, notFound)
	}
	return RT.getFrom(coll, key, notFound)
}

func (_ *rt) getFrom(coll interface{}, key interface{}, notFound interface{}) interface{} {
	if coll == nil {
		return nil
	}
	switch c := coll.(type) {
	// TODO: Figure out map type checking
	case IPersistentSet:
		if c.Contains(key) {
			return c.Get(key)
		}
		return notFound
	}
	// TODO: This implementation is incomplete
	return nil
}


func (_ *rt) Assoc(coll interface{}, key interface{}, val interface{}) Associative {
	if coll == nil {
		array := []interface{}{key, val}
		return &PersistentArrayMap{array: array}
	}
	return coll.(Associative).Assoc(key, val)
}

func (_ *rt) Contains(coll interface{}, key interface{}) interface{} {
	if coll == nil {
		return false
	}
	// TODO: figure out map type checking
	switch c := coll.(type) {
	case Associative:
		return c.ContainsKey(key)
	case IPersistentSet:
		return c.Contains(key)
	}

	// TODO...more of this. Sigh.
	return nil
}


/*
	List (Persistent) support
*/

func (_ *rt) List(args ...interface{}) ISeq {
	if len(args) == 0 {
		return nil
	} else if len(args) == 1 {
		return &PersistentList{
			_first: args[0],
			_rest:  nil,
			_count: 1,
		}
	} else {
		newarray := make([]interface{}, len(args) + 1)
		copy(newarray, args)
		return RT.ListStar(newarray...)
	}
}

func (_ *rt) ListStar(args ...interface{}) ISeq {
	if len(args) == 1 {
		return args[0].(ISeq)
	} else {
		return RT.Cons(args[0], RT.ListStar(args[1:]...))
	}
}
