package main

import (
	"context"
	"log"
	"net/http"

	"github.com/deVamshi/golang_food_delivery_api/internal/fooditem"
	"github.com/deVamshi/golang_food_delivery_api/internal/hotel"
	"github.com/deVamshi/golang_food_delivery_api/internal/user"
	"github.com/deVamshi/golang_food_delivery_api/pkg/dbcontext"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var(
	userservice user.UserService
	usercontroller user.UserController
	ctx context.Context
)

var(
	fooditemservice fooditem.FoodItemService
	fooditemcontroller fooditem.FoodItemController
)

func main() {
	// load env vars
	log.Println("checkb")
	ctx = context.Background()
	APP_ENV, err := godotenv.Read("../.env")
	if err != nil {
		log.Fatal(err)
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	dbClient, err := dbcontext.ConnectDB(APP_ENV["MONGOURI"])
	defer func() {
		if err := dbcontext.DisconnectDB(); err != nil {
			log.Fatal(err)
		}
		log.Print("Disconnected to DB\nBye!")
	}()

	if err != nil {
		log.Fatal(err)
		return
	}

	v1 := r.Group("/v1")
	{
		hotel.RegisterHandlers(v1, hotel.NewService(hotel.NewRepository(dbClient)))
		var usercollection = dbClient.Collection("users")
		userservice = user.NewUserService(usercollection,ctx)
		usercontroller = user.New(userservice)
		usercontroller.RegisterUserRoutes(v1) 
		var itemscollection = dbClient.Collection("fooditems")
		fooditemservice = fooditem.NewFoodItemService(itemscollection,ctx)
		fooditemcontroller = fooditem.New(fooditemservice)
		fooditemcontroller.RegisterFoodItemRoutes(v1)

	}


	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// func() {
	// 	time.Sleep(3 * time.Second)
	// 	server.Shutdown(context.Background())
	// }()

	server.ListenAndServe()

}
