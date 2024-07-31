package main

import (
	"github.com/labstack/echo/v4"
	_ "github.com/labstack/echo/v4"
	"os"
	"sync"
	"tesodev-korpes/CustomerService/cmd"
	orderCmd "tesodev-korpes/OrderService/cmd"
	"tesodev-korpes/pkg"
	"tesodev-korpes/pkg/middlewares"
	"tesodev-korpes/shared/config"
)

type RequestProcessor struct {
	counter int
	mutex   sync.Mutex
}

// Increment increments the counter
func (rp *RequestProcessor) Increment() {
	rp.mutex.Lock()
	defer rp.mutex.Unlock()
	rp.counter++
}

// GetCounter returns the current value of the counter
func (rp *RequestProcessor) GetCounter() int {
	rp.mutex.Lock()
	defer rp.mutex.Unlock()
	return rp.counter
}

func main() {
	//todo : what is dev,qa,prod ? explain why we are using them in the lecture
	dbConf := config.GetDBConfig("dev")

	client, err := pkg.GetMongoClient(dbConf.MongoDuration, dbConf.MongoClientURI)
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(pkg.CorrelationIDMiddleware)

	e.Use(middlewares.Logger())
	/*stats := middlewares.NewStats()
	e.Use(middleware.Logger())
	e.Use(stats.Process)
	e.Use(middlewares.ScopedServiceMiddleware)*/

	h_client := pkg.NewRestClient()

	if len(os.Args) < 2 {
		panic("Please provide a service to start: customer, order, or both")
	}
	input := os.Args[1] // This argument specifies which service to start. This setup allows
	// the program to determine its behavior based on user input provided at runtime.

	switch input {
	case "customer":
		cmd.BootCustomerService(client, e)
	case "order":
		orderCmd.BootOrderService(client, h_client, e)
	case "both":
		cmd.BootCustomerService(client, e)
		//  allowing both cmd.BootCustomerService(client, e)
		// and BootOrderService(client, e) functions to run simultaneously in the 'both' case
		go orderCmd.BootOrderService(client, h_client, e)
	default:
		panic("Invalid input. Use 'customer', 'order', or 'both'.")
	}

	// Keep the main function alive
	select {}
}

//challenge : after you create a func boot order service, manage somehow to run specific project //switch
// cases are used above.

//description : when you give an input here it should look that input and boot THAT specific project
//if the input says "both" it should.

//PS : do not forget to create and call a different column for order service and do not forget to boot order service
//from another port different than customer service---> .. OrderService config.go

//orderCol, err := pkg.GetMongoCollection(client, "tesodev", "order")
//if err != nil {
//	panic(err)
