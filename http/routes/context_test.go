package routes

import (
	"testing"

	httpTest "github.com/saharsh-samples/go-mux-sql-starter/http/test"
	"github.com/saharsh-samples/go-mux-sql-starter/test"
)

func TestBootstrap(t *testing.T) {

	// Arrange
	responseWriter := &httpTest.MockResponseWriter{}

	// Act
	out := Bootstrap(&ContextIn{})

	// Assert
	test.AssertEquals("", 1, len(out.RoutesToRegister), t)

	livenessCheckRoute, _ := out.RoutesToRegister[0].(*LivenessCheck)
	livenessCheckRoute.Get(responseWriter, nil)
	test.AssertEquals("", 200, responseWriter.WriteHeaderStatusCodeArg, t)
}
