package crsfmiddleware

import (
    "context"
    "net/http"
    "github.com/traefik/traefik/v2/pkg/middlewares"
    "github.com/traefik/traefik/v2/pkg/middlewares/emptybackendhandler"
    "github.com/traefik/traefik/v2/pkg/server/middleware"
)

type Config struct{}

func CreateConfig() *Config {
    return &Config{}
}

type CSRFTokenMiddleware struct {
    next http.Handler
    name string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
    return &CSRFTokenMiddleware{
        next: next,
        name: name,
    }, nil
}

func (c *CSRFTokenMiddleware) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
    csrfToken := req.Header.Get("X-CSRFToken")
    if csrfToken != "" {
        http.SetCookie(rw, &http.Cookie{
            Name:  "csrftoken",
            Value: csrfToken,
            Path:  "/",
        })
    }

    c.next.ServeHTTP(rw, req)
}

func init() {
    middleware.NewRegistry().Add("crsfmiddleware", func(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
        return New(ctx, next, config, name)
    })
}
