package main

import (
	"fmt"
	"kitchen-api/internal/app/addon"
	menucategory "kitchen-api/internal/app/menu_category"
	menuitem "kitchen-api/internal/app/menu_item"
	"kitchen-api/internal/app/order"
	orderitem "kitchen-api/internal/app/order_item"
	"kitchen-api/internal/app/payment"
	"kitchen-api/internal/app/restaurant"
	"kitchen-api/internal/app/user"
	"kitchen-api/internal/server"
	"log"
	"os"

	"kitchen-api/internal/config"
	"kitchen-api/internal/database"

	"gorm.io/driver/postgres"
)

func main() {
	config := config.Default()

	stripeKey := config.GetString("stripe.secretKey")
	if key := os.Getenv("STRIPE_SECRET_KEY"); key != "" {
		stripeKey = key
	}

	fmt.Println("stipe key", stripeKey)

	database, err := openPostgresDb()
	if err != nil {
		log.Fatal("Error while connecting to database", err)
	}

	// Fix GORM driver detection for lib/pq to avoid extra pgx protocol parameters
	if dialector, ok := database.OrmInstance.Dialector.(*postgres.Dialector); ok {
		dialector.DriverName = "postgres"
		if dialector.Config != nil {
			dialector.Config.DriverName = "postgres"
		}
	}
	database.OrmInstance.PrepareStmt = true

	// Ensure pgcrypto extension and enum types exist
	database.OrmInstance.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`)
	database.OrmInstance.Exec(`
		DO $$ BEGIN
			CREATE TYPE user_role AS ENUM ('SUPER_ADMIN', 'ADMIN', 'USER');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
	`)
	database.OrmInstance.Exec(`
		DO $$ BEGIN
			CREATE TYPE menu_type AS ENUM ('BAR', 'KITCHEN');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
	`)
	database.OrmInstance.Exec(`
		DO $$ BEGIN
			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'orders' AND column_name = 'table_id') 
			   AND EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'tables') THEN
				UPDATE orders o
				SET table_name = COALESCE(t.number::text, 'N/A')
				FROM tables t
				WHERE o.table_id = t.id AND (o.table_name IS NULL OR o.table_name = '');
			END IF;
			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'orders' AND column_name = 'table_name') THEN
				UPDATE orders SET table_name = 'N/A' WHERE table_name IS NULL;
			END IF;

			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'payment_sessions' AND column_name = 'table_id') 
			   AND EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'tables') THEN
				UPDATE payment_sessions p
				SET table_name = COALESCE(t.number::text, 'N/A')
				FROM tables t
				WHERE p.table_id = t.id AND (p.table_name IS NULL OR p.table_name = '');
			END IF;
			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'payment_sessions' AND column_name = 'table_name') THEN
				UPDATE payment_sessions SET table_name = 'N/A' WHERE table_name IS NULL;
			END IF;
		END $$;
	`)

	err = database.OrmInstance.AutoMigrate(
		&restaurant.Restaurants{},
		&user.Users{},
		&menucategory.MenuCategory{},
		&menuitem.MenuItem{},
		&order.Order{},
		&orderitem.OrderItem{},
		&orderitem.OrderItemAddon{},
		&addon.Addon{},
		&payment.PaymentSession{},
	)
	if err != nil {
		log.Fatal("Error while running auto migrations: ", err)
	}

	database.OrmInstance.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_orders_stripe_session_id ON orders(stripe_session_id) WHERE stripe_session_id IS NOT NULL;")

	router := server.NewRouter(database)
	port := config.GetString("server.port")

	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	router.Run(":" + port)
}

func openPostgresDb() (*database.OrmDb, error) {
	host := config.Default().GetString("db.host")
	user := config.Default().GetString("db.user")
	password := config.Default().GetString("db.password")
	dbName := config.Default().GetString("db.database")
	port := config.Default().GetInt("db.port")
	orm, err := database.OpenPostgresORM(host, port, user, password, dbName)
	if err != nil {
		return nil, err
	}
	return &orm, nil
}
