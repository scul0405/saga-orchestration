package postgres_repo

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/scul0405/saga-orchestration/internal/account/domain/entity"
	"github.com/scul0405/saga-orchestration/internal/account/domain/valueobject"
	"github.com/scul0405/saga-orchestration/pkg/sonyflake"
	"github.com/scul0405/saga-orchestration/pkg/utils"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"testing"
)

func NewAuthMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Printf("an error '%s' was not expected when opening a stub database connection", err)
		return nil, nil, err
	}

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	}), &gorm.Config{})
	if err != nil {
		log.Printf("an error '%s' was not expected when opening a stub database connection", err)
		return nil, nil, err
	}

	return gdb, mock, nil
}

func TestCheckCustomer(t *testing.T) {
	t.Parallel()

	gdb, mock, err := NewAuthMockDB()
	require.NoError(t, err)
	require.NotNil(t, gdb)
	require.NotNil(t, mock)

	authRepo := NewJwtAuthRepositoryImpl(gdb)

	sf, err := sonyflake.NewSonyFlake()
	require.NoError(t, err)

	t.Run("CheckCustomer", func(t *testing.T) {
		customerID, err := sf.NextID()
		require.NoError(t, err)

		testCustomer := entity.Customer{
			ID:     customerID,
			Active: true,
		}

		rows := sqlmock.NewRows([]string{"id", "active"}).
			AddRow(testCustomer.ID,
				testCustomer.Active,
			)

		mock.ExpectQuery(
			"SELECT \"active\" FROM \"accounts\" WHERE id = $1 ORDER BY \"accounts\".\"id\" LIMIT 1").
			WithArgs(testCustomer.ID).WillReturnRows(rows)

		exists, active, err := authRepo.CheckCustomer(context.Background(), testCustomer.ID)
		require.NoError(t, err)
		require.True(t, exists)
		require.True(t, active)
	})
}

func TestCreateCustomer(t *testing.T) {
	t.Parallel()

	gdb, mock, err := NewAuthMockDB()
	require.NoError(t, err)
	require.NotNil(t, gdb)
	require.NotNil(t, mock)

	authRepo := NewJwtAuthRepositoryImpl(gdb)

	sf, err := sonyflake.NewSonyFlake()
	require.NoError(t, err)

	t.Run("CreateCustomer", func(t *testing.T) {
		customerID, err := sf.NextID()
		require.NoError(t, err)

		testCustomer := entity.Customer{
			ID:     customerID,
			Active: true,
			PersonalInfo: &valueobject.CustomerPersonalInfo{
				FirstName: "dep",
				LastName:  "trai",
				Email:     "deptrai@gmail.com",
			},
			DeliveryInfo: &valueobject.CustomerDeliveryInfo{
				Address:     "123 abc",
				PhoneNumber: "123456789",
			},
			Password: "secret",
		}

		hashedPassword, err := utils.HashPassword(testCustomer.Password)
		require.NoError(t, err)

		sqlmock.NewRows([]string{"id", "active", "first_name", "last_name", "email", "address", "phone_number", "password"}).
			AddRow(testCustomer.ID,
				testCustomer.Active,
				testCustomer.PersonalInfo.FirstName,
				testCustomer.PersonalInfo.LastName,
				testCustomer.PersonalInfo.Email,
				testCustomer.DeliveryInfo.Address,
				testCustomer.DeliveryInfo.PhoneNumber,
				hashedPassword,
			)

		mock.ExpectBegin()
		mock.ExpectQuery("INSERT INTO \"accounts\" (\"active\",\"first_name\",\"last_name\",\"email\",\"address\",\"phone_number\",\"password\",\"updated_at\",\"created_at\",\"id\") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING \"id\"").
			WithArgs(
				testCustomer.Active,
				testCustomer.PersonalInfo.FirstName,
				testCustomer.PersonalInfo.LastName,
				testCustomer.PersonalInfo.Email,
				testCustomer.DeliveryInfo.Address,
				testCustomer.DeliveryInfo.PhoneNumber,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				testCustomer.ID,
			).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testCustomer.ID))
		mock.ExpectCommit()

		err = authRepo.CreateCustomer(context.Background(), &testCustomer)
		require.NoError(t, err)
	})
}

func TestGetCustomerCredentials(t *testing.T) {
	t.Parallel()

	gdb, mock, err := NewAuthMockDB()
	require.NoError(t, err)
	require.NotNil(t, gdb)
	require.NotNil(t, mock)

	authRepo := NewJwtAuthRepositoryImpl(gdb)

	sf, err := sonyflake.NewSonyFlake()
	require.NoError(t, err)

	t.Run("GetCustomerCredentials", func(t *testing.T) {
		customerID, err := sf.NextID()
		require.NoError(t, err)

		testCustomer := entity.Customer{
			ID:     customerID,
			Active: true,
			PersonalInfo: &valueobject.CustomerPersonalInfo{
				Email: "deptrai@123.com",
			},
			Password: "secret",
		}

		hashedPassword, err := utils.HashPassword(testCustomer.Password)
		require.NoError(t, err)

		rows := sqlmock.NewRows([]string{"id", "active", "email", "password"}).
			AddRow(testCustomer.ID,
				testCustomer.Active,
				testCustomer.PersonalInfo.Email,
				hashedPassword,
			)

		mock.ExpectQuery(
			"SELECT \"id\",\"active\",\"password\" FROM \"accounts\" WHERE email = $1 ORDER BY \"accounts\".\"id\" LIMIT 1").
			WithArgs(testCustomer.PersonalInfo.Email).WillReturnRows(rows)

		exists, creds, err := authRepo.GetCustomerCredentials(context.Background(), testCustomer.PersonalInfo.Email)
		require.NoError(t, err)
		require.True(t, exists)
		require.NotNil(t, creds)
		require.Equal(t, testCustomer.ID, creds.CustomerID)
		require.Equal(t, testCustomer.Active, creds.Active)
		require.Equal(t, hashedPassword, creds.Password)
	})
}
