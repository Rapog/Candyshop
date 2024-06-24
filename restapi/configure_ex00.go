// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"ex01/restapi/operations"
)

//go:generate swagger generate server --target ../../Exercise00 --name Ex00 --spec ../docs.yaml --principal interface{}

func configureFlags(api *operations.Ex00API) {
	//api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{}
}

func configureAPI(api *operations.Ex00API) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError
	api.Context()

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	//if api.BuyCandyHandler == nil {
	api.BuyCandyHandler = operations.BuyCandyHandlerFunc(func(params operations.BuyCandyParams) middleware.Responder {
		var candyMap = map[string]int64{
			"CE": 10, "AA": 15, "NT": 17, "DE": 21, "YR": 23,
		}
		_, exist := candyMap[*params.Order.CandyType]
		if !exist || *params.Order.Money < 0 || *params.Order.CandyCount < 0 {
			var d operations.BuyCandyBadRequestBody = operations.BuyCandyBadRequestBody{
				Error: http.StatusText(operations.BuyCandyBadRequestCode),
			}
			var req = operations.NewBuyCandyBadRequest().WithPayload(&d)
			//req.WriteResponse(w http.ResponseWriter)
			return req
		}

		var sumValue = *params.Order.CandyCount * candyMap[*params.Order.CandyType]
		var c = *params.Order.Money - sumValue
		if sumValue <= *params.Order.Money {
			var okResp operations.BuyCandyCreatedBody = operations.BuyCandyCreatedBody{
				Change: c,
				Thanks: "Thank you!",
			}
			return operations.NewBuyCandyCreated().WithPayload(&okResp)
		} else {
			var okResp operations.BuyCandyPaymentRequiredBody = operations.BuyCandyPaymentRequiredBody{
				Error: fmt.Sprintf("You need %d more money!", c*-1),
			}
			return operations.NewBuyCandyPaymentRequired().WithPayload(&okResp)
		}
	})
	//}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
	certFile := "localhost/cert.pem"
	keyFile := "localhost/key.pem"
	caFile := "minica.pem"

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to load server certificate and key: %v", err)
	}

	caCert, err := os.ReadFile(caFile)
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %v", err)
	}
	rootCAPool := x509.NewCertPool()
	if ok := rootCAPool.AppendCertsFromPEM(caCert); !ok {
		log.Fatalf("Failed to append CA certificate to pool")
	}

	tlsConfig.Certificates = []tls.Certificate{cert}
	tlsConfig.ClientCAs = rootCAPool
	tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert

}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
