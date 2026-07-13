package middleware

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

type CacheMiddleware struct {
	keyPrefix   string
	redisClient *redis.Client
	cacheTime   time.Duration
}

type cacheResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (rw *cacheResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *cacheResponseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

// canonicalQueryKey 将 query 参数按 key、value 排序后规范化，避免传参顺序不同导致缓存 key 不一致。
func canonicalQueryKey(rawQuery string) string {
	if rawQuery == "" {
		return ""
	}

	values, err := url.ParseQuery(rawQuery)
	if err != nil {
		query := strings.ReplaceAll(rawQuery, "&", ":")
		return strings.ReplaceAll(query, "=", ":")
	}

	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys)*2)
	for _, k := range keys {
		vs := append([]string(nil), values[k]...)
		sort.Strings(vs)
		for _, v := range vs {
			parts = append(parts, k, v)
		}
	}
	return strings.Join(parts, ":")
}

func NewCacheMiddleware(redisClient *redis.Client, cacheTime time.Duration, keyPrefix string) *CacheMiddleware {
	if cacheTime == 0 {
		cacheTime = 1 * time.Minute
	}
	if keyPrefix == "" {
		keyPrefix = "cache"
	}
	return &CacheMiddleware{
		redisClient: redisClient,
		cacheTime:   cacheTime,
		keyPrefix:   keyPrefix,
	}
}

func (m *CacheMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			next(w, r)
			return
		}

		rpath := strings.TrimLeft(r.URL.Path, "/")
		pathKey := strings.ReplaceAll(rpath, "/", ":")
		if shop := strings.TrimSpace(r.Header.Get("shop")); shop != "" {
			pathKey += fmt.Sprintf(":%s", shop)
		}
		if lang := strings.TrimSpace(r.Header.Get("language")); lang != "" {
			pathKey += fmt.Sprintf(":%s", lang)
		} else if lang = strings.TrimSpace(r.Header.Get("lang")); lang != "" {
			pathKey += fmt.Sprintf(":%s", lang)
		}

		if query := canonicalQueryKey(r.URL.RawQuery); query != "" {
			pathKey += fmt.Sprintf(":%s", query)
		}

		uid := ""
		key := fmt.Sprintf("%s:%s:%s", m.keyPrefix, r.Method, pathKey)
		tokenStr := strings.TrimSpace(strings.TrimLeft(r.Header.Get("Authorization"), "Bearer"))
		if parsed, _ := ParseBindEmailToken(tokenStr, "YnVkb25n"); parsed != nil {
			uid = parsed.UserId
		}
		if uid != "" {
			key += fmt.Sprintf(":%s", uid)
		}

		key = fmt.Sprintf("%s:%x", m.keyPrefix, sha1.Sum([]byte(key)))
		val, err := m.redisClient.Get(key).Bytes()
		if err == nil && len(val) > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.Write(val)
			fmt.Println("cache hit ->", key, uid)
			return
		}

		fmt.Println("cache miss", key)
		crw := &cacheResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next(crw, r)

		if crw.statusCode == http.StatusOK && crw.body.Len() > 0 {
			respBytes := crw.body.Bytes()
			_ = m.redisClient.Set(key, respBytes, m.cacheTime)
		}
	}
}
