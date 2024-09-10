package cmd

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/robfig/cron/v3"
	_ "github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"tesodev-korpes/OrderService/client"
	config3 "tesodev-korpes/OrderService/config"
	"tesodev-korpes/OrderService/internal"
	"tesodev-korpes/pkg"
)

func BootOrderService(client *mongo.Client, h_client *client.CustomerClient, e *echo.Echo) {
	config := config3.GetOrderConfig("dev")
	orderCol, err := pkg.GetMongoCollection(client, config.DbConfig.DBName, config.DbConfig.ColName)
	if err != nil {
		panic(err)
	}

	repo := internal.NewRepository(orderCol)
	service := internal.NewService(repo, h_client)
	internal.NewHandler(e, service)

	c := cron.New()

	// Test için her dakika çalışacak cron job
	_, err = c.AddFunc("* * * * *", func() {
		err := service.CalculateDailyTotalOrders(context.Background())
		if err != nil {
			fmt.Println("Failed to calculate daily total orders:", err)
		} else {
			fmt.Println("Daily total orders calculated successfully!")
		}
	})

	if err != nil {
		fmt.Println("Error adding cron job:", err)
	}

	// Cron başlatılır
	c.Start()

	// Sunucu başlatılır
	e.Logger.Fatal(e.Start(config.Port))
}

//daily repositorye ulaşmak
//günlük toplam siparişi hesaplamak
//bunları daily repositorye kaydet
