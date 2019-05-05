package translations

import (
	"context"
	"net/http"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var I18nCtxKey = &ContextKey{name: "translation"}

type ContextKey struct {
	name string
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bundle := r.Context().Value("bundle").(*i18n.Bundle)

		// Specify language by request
		lang := r.FormValue("lang")
		accept := r.Header.Get("Accept-Language")

		// Init localizer
		localizer := i18n.NewLocalizer(bundle, lang, accept)

		// Set context
		ctx := context.WithValue(r.Context(), I18nCtxKey, localizer)

		// and call the next with our new context
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})

}

func T(ctx context.Context, key string) string {
	l := ctx.Value(I18nCtxKey).(*i18n.Localizer)

	return l.MustLocalize(&i18n.LocalizeConfig{
		MessageID: key,
	})
}

func TWithTemplateData(ctx context.Context, key string, data map[string]interface{}) string {
	l := ctx.Value(I18nCtxKey).(*i18n.Localizer)

	return l.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: data,
	})
}
