package godis

import (
	"testing"
	"time"
)

// Test setting key-values to the DB
func TestSET(t *testing.T) {

	db := setUp()
	for _, c := range cases {
		got := db.SET(c.key, c.value)
		if got != c.key {
			t.Errorf("SET(%q) == %q, want %q", c.key, got, c.key)
		}
	}
}

// Test getting key-values from the DB
func TestGET(t *testing.T) {
	db := setUp()
	for _, c := range cases {
		db.SET(c.key, c.value)
		got := db.GET(c.key)
		if got != c.value {
			t.Errorf("GET(%q) == %q, want %q", c.key, got, c.value)
		}
	}
}

// Test getting non-existent key from the DB
func TestGETNotEXISTS(t *testing.T) {
	db := setUp()
	key := "not exists"
	got := db.GET(key)
	if got != nil {
		t.Errorf("GET(%q) == %v, want %v", key, got, nil)
	}
}

// Test MGETting key-values from the DB
func TestMGET(t *testing.T) {
	testKeys := []string{"key1", "key2", "key 3"}
	want := []string{"new value 1", "value2", "value 3"}
	db := setUp()
	for _, c := range cases {
		db.SET(c.key, c.value)
	}
	got := db.MGET(testKeys...)
	for i, key := range testKeys {
		if got[i] != want[i] {
			t.Errorf("MGET(%q) == %q, want %q", key, got[i], want[i])
		}
	}
}

// Test MGETting non-existent key-values from the DB
func TestMGETNotExists(t *testing.T) {
	testKeys := []string{"non-key1", "non-key2", "non-key3"}
	db := setUp()
	for _, c := range cases {
		db.SET(c.key, c.value)
	}
	got := db.MGET(testKeys...)
	for i, key := range testKeys {
		if got[i] != nil {
			t.Errorf("MGET(%q) == %q, want %p", key, got[i], nil)
		}
	}

}

// Test MGETting few non-existent key-values from the DB
func TestMGETFewNotExists(t *testing.T) {
	testKeys := []string{"key1", "non-key2", "key 3"}
	want := []string{"new value 1", "value2", "value 3"}
	db := setUp()
	for _, c := range cases {
		db.SET(c.key, c.value)
	}
	got := db.MGET(testKeys...)
	for i, key := range testKeys {
		if i == 1 && got[i] != nil {
			t.Errorf("MGET(%q) == %q, want %p", key, got[i], nil)
		}
		if i != 1 && got[i] != want[i] {
			t.Errorf("MGET(%q) == %q, want %q", key, got[i], want[i])
		}
	}

}

// Test MSETing key-values pairs into the DB
func TestMSET(t *testing.T) {
	tests := []string{"key1", "value1", "key2", "value2", "key3", "value3"}
	db := setUp()
	got := db.MSET(tests...)
	if got != true {
		t.Errorf("MSET(%q) == %t, want %t", tests, got, true)
	}
}

// Test incrementing values for given key by 1
func TestINCR(t *testing.T) {
	db := setUp()
	for _, c := range integers {
		db.SET(c.key, c.value)
		got := db.INCR(c.key)
		if got != c.value.(int)+1 {
			t.Errorf("INCR(%q) == %d, want %d", c.key, got, c.value.(int)+1)
		}
	}
}

// Test incrementing non-existent keys
func TestINCRNonExists(t *testing.T) {
	db := setUp()
	got := db.INCR("non-incr-key")
	if got != 1 {
		t.Errorf("INCR(%q) == %d, want %d", "non-incr-key", got, 1)
	}
}

// Test decrementing values for given key by 1
func TestDECR(t *testing.T) {
	db := setUp()
	for _, c := range integers {
		db.SET(c.key, c.value)
		got := db.DECR(c.key)
		if got != c.value.(int)-1 {
			t.Errorf("DECR(%q) == %d, want %d", c.key, got, c.value.(int)-1)
		}
	}
}

// Test incrementing non-existent keys
func TestDECRNonExists(t *testing.T) {
	db := setUp()
	got := db.DECR("non-decr-key")
	if got != -1 {
		t.Errorf("DECR(%q) == %d, want %d", "non-decr-key", got, -1)
	}
}

// Test incrementing values for given key by n
func TestINCRBY(t *testing.T) {
	db := setUp()
	n := 3
	for _, c := range integers {
		db.SET(c.key, c.value)
		got := db.INCRBY(c.key, n)
		if got != c.value.(int)+n {
			t.Errorf("INCRBY(%q) == %d, want %d", c.key, got, c.value.(int)+n)
		}
	}
}

// Test decrementing values for given key by n
func TestDECRBY(t *testing.T) {
	db := setUp()
	n := 3
	for _, c := range integers {
		db.SET(c.key, c.value)
		got := db.DECRBY(c.key, n)
		if got != c.value.(int)-n {
			t.Errorf("DECRBY(%q) == %d, want %d", c.key, got, c.value.(int)-n)
		}
	}
}

/* Test incrementing values for given key by 1
func TestINCRmismatchs(t *testing.T) {
    db := setUp()
    n := 1
    for _, c := range strings {
        db.SET(c.key, c.value)
        got := db.INCR(c.key)
        _, ok := c.value.(int)
        //if ok && got != value+n {
        //  t.Errorf("INCR(%q) == %d, want %d", c.key, got, value+n)
        //}
        if !ok {
            t.Errorf("INCR(%q) got type %v, want type %v", c.key, reflect.TypeOf(got), reflect.TypeOf(n))
        }
    }
}*/

func TestSETEXWithinExp(t *testing.T) {
	// One second before expiry time
	key := "mykey"
	val := 25
	exp := 2
	db := setUp()
	db.SETEX(key, int64(exp), val)
	time.Sleep(time.Duration(exp-1) * time.Second)
	got := db.GET(key)
	if got != val {
		t.Errorf("SETEX(%q) == %d, want %d", key, got, val)
	}
}

func TestSETEXAfterExp(t *testing.T) {
	// One second before expiry time
	key := "mykey"
	val := 25
	exp := 1
	db := setUp()
	db.SETEX(key, int64(exp), val)
	time.Sleep(time.Duration(exp+1) * time.Second)
	got := db.GET(key)
	if got != nil {
		t.Errorf("SETEX(%q) == %d, want nil", key, got)
	}
}

func TestSETEXWithZero(t *testing.T) {
	// Zero as expiry time
	key := "mykey"
	val := 25
	exp := 0
	db := setUp()
	res := db.SETEX(key, int64(exp), val)
	if res != false {
		t.Errorf("SETEX(%q) == %d, want %d", key, res, val)
	}
}

func TestSETEXNegative(t *testing.T) {
	// negative valuse as expiry time
	key := "mykey"
	val := 25
	exp := -1
	db := setUp()
	res := db.SETEX(key, int64(exp), val)
	if res != false {
		t.Errorf("SETEX(%q) == %d, want %d", key, res, val)
	}
}

func TestPSETEXWithinExp(t *testing.T) {
	// One second before expiry time
	key := "mykey"
	val := 25
	exp := 1000
	db := setUp()
	db.PSETEX(key, int64(exp), val)
	time.Sleep(time.Duration(exp-1) * time.Millisecond)
	got := db.GET(key)
	if got != val {
		t.Errorf("SETEX(%q) == %d, want %d", key, got, val)
	}
}

func TestPSETEXAfterExp(t *testing.T) {
	// One second before expiry time
	key := "mykey"
	val := 25
	exp := 1000
	db := setUp()
	db.PSETEX(key, int64(exp), val)
	time.Sleep(time.Duration(exp+1) * time.Millisecond)
	got := db.GET(key)
	if got != nil {
		t.Errorf("SETEX(%q) == %d, want nil", key, got)
	}
}



