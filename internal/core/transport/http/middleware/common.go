package core_http_middleware

import (
	"net/http"
	"time"

	core_logger "github.com/Martaller1999/martaller-todo/internal/core/logger"
	core_http_response "github.com/Martaller1999/martaller-todo/internal/core/transport/http/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	requestIDHeader = "X-Request-ID"
)

func RequestID() Middleware { //зашиваем в миддлвэр код, который нужно выполнить до и после обработки HTTP handlerов
	// когда приходит http запрос, он попадает на эту мидлварю
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader) // либо получаем request id из входящего запроса
			if requestID == "" {
				requestID = uuid.NewString() //либо дообогащаем хэдеры входящего запроса своим request id
			}
			r.Header.Set(requestIDHeader, requestID)
			w.Header().Set(requestIDHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log *core_logger.Logger) Middleware { //входящий запрос, обогатившись request id, он попадает сюда
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			l := log.With( // и обогатится логгером на автоматич. логгирование доп. полей
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)

			ctx := core_logger.ToContext(r.Context(), l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			rw := core_http_response.NewResponseWriter(w)

			before := time.Now()
			log.Debug( //до выполнения хэндлера логируется входящий запрос
				">>> incoming HTTP request",
				zap.String("http method", r.Method),
				zap.Time("time", before.UTC()),
			)

			next.ServeHTTP(rw, r) //передача на осн. http handler

			log.Debug( //после выполнения хэндлера логируется информация обработки: результат и время
				"<<< обработан HTTP запрос",
				zap.Int("status_code", rw.GetStatusCode()),
				zap.Duration("latency", time.Now().Sub(before)),
			)
		})
	}
}

func Panic() Middleware { //если не получилось проставить request id и логер, то приложение падает с паникой
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(
						p,
						"во время HTTP-запроса получена неожиданная паника",
					)
				}
			}()

			next.ServeHTTP(w, r)

		})
	}
}

// на 5:34 разбор
