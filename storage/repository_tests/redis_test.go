package repository_tests

import (
	"testing"
	"time"
)

func TestRedisCache(t *testing.T) {
	cache, err := getRedisFromArgs()

	if err != nil {
		t.Error(err)
		return
	}

	err = cache.Put("test_a", "val1", time.Second*10)
	if err != nil {
		t.Error("failed to put string in cache: ", err)
	}

	val := ""
	err = cache.Get("test_a", &val)
	if err != nil {
		t.Error("failed to retrieve string from cache: ", err)
	}

	err = cache.Delete("test_a")
	if err != nil {
		t.Error("failed to delete string from cache: ", err)
	}

}

func TestRedisCacheTimeouts(t *testing.T) {
	cache, err := getRedisFromArgs()
	if err != nil {
		t.Error(err)
		return
	}

	err = cache.Put("test_b", "val2", time.Millisecond*10)
	time.Sleep(time.Millisecond * 15)

	val := ""
	err = cache.Get("test_b", &val)

	if val != "" {
		t.Errorf("timeout not working, time = 10ms")
	}

	err = cache.Put("test_c", "val2", time.Millisecond*500)
	time.Sleep(time.Millisecond * 505)

	val = ""
	err = cache.Get("test_b", &val)

	if val != "" {
		t.Errorf("timeout not working, time = 500ms")
	}
}
