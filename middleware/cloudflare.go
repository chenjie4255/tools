package middleware

import (
	"context"
	"net/http"
)

var cfRay = http.CanonicalHeaderKey("Cf-Ray")           //CloudFlare Ray
var cfCountry = http.CanonicalHeaderKey("Cf-Ipcountry") // CF Country

type contextKey int

const (
	cfRayKey contextKey = iota
	cfCountryKey
)

// CF middleware
func CF(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if ray := r.Header.Get(cfRay); ray != "" {
			ctx = context.WithValue(ctx, cfRayKey, ray)
		}

		if country := r.Header.Get(cfCountry); country != "" {
			ctx = context.WithValue(ctx, cfCountryKey, country)
		}

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func IsFromCloudFlare(ctx context.Context) bool {
	return CFRay(ctx) != ""
}

func GetCloudFlareCountry(ctx context.Context) string {
	return CFCountry(ctx)
}

//IsFromOutseaRequest 判断是否请求的来源是海外
func IsFromOutseaRequest(ctx context.Context) bool {
	return IsFromCloudFlare(ctx) && GetCloudFlareCountry(ctx) != "CN"
}

func CFCountry(ctx context.Context) string {
	country := ctx.Value(cfCountryKey)
	if country != nil {
		return country.(string)
	}

	return ""
}

func CFRay(ctx context.Context) string {
	ray := ctx.Value(cfRayKey)
	if ray != nil {
		return ray.(string)
	}

	return ""
}
