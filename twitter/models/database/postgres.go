package database

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
	conn *gorm.DB
}

func NewPostgresDB() *PostgresDB {
	return &PostgresDB{}
}

type User struct {
	ID       uint
	Username string
	Mail     string
	FullName string
	Created  []User `gorm:"foreignkey:UserID"`
	UserID   *uint
}

func (db *PostgresDB) Connect() error {
	name := viper.GetString("postgres.name")
	host := viper.GetString("postgres.host")
	port := viper.GetString("postgres.port")
	user := viper.GetString("postgres.user")
	password := viper.GetString("postgres.password")

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, name, password, port)
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
		return err
	} else {
		log.Infof("postgres: successfully connected to database on %s %s", host, port)
	}
	db.conn = conn
	db.setup()
	return nil
}

func (db *PostgresDB) setup() {
	err := db.conn.AutoMigrate(
		&User{},
	)
	log.Info("postgres: migrated tables")
	if err == nil {
		db.conn.Exec("DELETE FROM users")
		user1 := User{Username: "usuario1", Mail: "usuario1@Test.com", FullName: "NombreUsuario1"}
		db.conn.Create(&user1)

		user2 := User{Username: "usuario2", Mail: "usuario2@Test.com", FullName: "NombreUsuario2", UserID: &(user1.ID)}
		db.conn.Create(&user2)

		user3 := User{Username: "usuario3", Mail: "usuario3@Test.com", FullName: "NombreUsuario3", UserID: &(user1.ID)}
		db.conn.Create(&user3)

		user4 := User{Username: "usuario4", Mail: "usuario4@Test.com", FullName: "NombreUsuario4", UserID: &(user1.ID)}
		db.conn.Create(&user4)

	}
}

func (db *PostgresDB) GetSeguidores(args ...interface{}) ([]map[string]interface{}, error) {

	fmt.Printf("args %+v\n\n", args)

	//err := db.conn.Preload("Created").Where(" username = ? ", args[0]).Find(&result).Error
	var result []map[string]interface{}
	err := db.conn.Model(User{}).Where(" username = ? ", args[0]).Find(&result).Error
	if err != nil {
		log.Error("Error recuperando seguidores ", err)
		return nil, err

	}

	fmt.Printf("Resultado %+v\n", result)
	return result, nil
}
