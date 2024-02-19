package postgres_repo

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/scul0405/saga-orchestration/pkg/sonyflake"
	"github.com/scul0405/saga-orchestration/services/account/internal/domain/entity"
	"github.com/scul0405/saga-orchestration/services/account/internal/domain/valueobject"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"testing"
)

func NewCustomerMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
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

func TestGetCustomerPersonalInfo(t *testing.T) {
	t.Parallel()

	gdb, mock, err := NewCustomerMockDB()
	require.NoError(t, err)
	require.NotNil(t, gdb)
	require.NotNil(t, mock)

	customerRepo := NewCustomerRepositoryImpl(gdb)

	sf, err := sonyflake.NewSonyFlake()
	require.NoError(t, err)

	t.Run("GetCustomerPersonalInfo", func(t *testing.T) {
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
		}

		rows := sqlmock.NewRows([]string{"id", "active", "first_name", "last_name", "email"}).
			AddRow(testCustomer.ID,
				testCustomer.Active,
				testCustomer.PersonalInfo.FirstName,
				testCustomer.PersonalInfo.LastName,
				testCustomer.PersonalInfo.Email)

		mock.ExpectQuery(
			"SELECT \"first_name\",\"last_name\",\"email\" FROM \"accounts\" WHERE id = $1 AND active = TRUE ORDER BY \"accounts\".\"id\" LIMIT 1").
			WithArgs(testCustomer.ID).WillReturnRows(rows)

		personalInfo, err := customerRepo.GetCustomerPersonalInfo(context.Background(), testCustomer.ID)
		require.NoError(t, err)
		require.NotNil(t, personalInfo)
		require.Equal(t, testCustomer.PersonalInfo.FirstName, personalInfo.FirstName)
		require.Equal(t, testCustomer.PersonalInfo.LastName, personalInfo.LastName)
		require.Equal(t, testCustomer.PersonalInfo.Email, personalInfo.Email)
	})
}

func TestGetCustomerDeliveryInfo(t *testing.T) {
	t.Parallel()

	gdb, mock, err := NewCustomerMockDB()
	require.NoError(t, err)
	require.NotNil(t, gdb)
	require.NotNil(t, mock)

	customerRepo := NewCustomerRepositoryImpl(gdb)

	sf, err := sonyflake.NewSonyFlake()
	require.NoError(t, err)

	t.Run("GetCustomerDeliveryInfo", func(t *testing.T) {
		customerID, err := sf.NextID()
		require.NoError(t, err)

		testCustomer := entity.Customer{
			ID:     customerID,
			Active: true,
			DeliveryInfo: &valueobject.CustomerDeliveryInfo{
				Address:     "test address",
				PhoneNumber: "test phone number",
			},
		}

		rows := sqlmock.NewRows([]string{"id", "active", "address", "phone_number"}).
			AddRow(testCustomer.ID,
				testCustomer.Active,
				testCustomer.DeliveryInfo.Address,
				testCustomer.DeliveryInfo.PhoneNumber)

		mock.ExpectQuery(
			"SELECT \"address\",\"phone_number\" FROM \"accounts\" WHERE id = $1 AND active = TRUE ORDER BY \"accounts\".\"id\" LIMIT 1").
			WithArgs(testCustomer.ID).WillReturnRows(rows)

		personalInfo, err := customerRepo.GetCustomerDeliveryInfo(context.Background(), testCustomer.ID)
		require.NoError(t, err)
		require.NotNil(t, personalInfo)
		require.Equal(t, testCustomer.DeliveryInfo.Address, personalInfo.Address)
		require.Equal(t, testCustomer.DeliveryInfo.PhoneNumber, personalInfo.PhoneNumber)
	})
}

func TestUpdateCustomerPersonalInfo(t *testing.T) {
	t.Parallel()

	gdb, mock, err := NewCustomerMockDB()
	require.NoError(t, err)
	require.NotNil(t, gdb)
	require.NotNil(t, mock)

	customerRepo := NewCustomerRepositoryImpl(gdb)

	sf, err := sonyflake.NewSonyFlake()
	require.NoError(t, err)

	t.Run("UpdateCustomerPersonalInfo", func(t *testing.T) {
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
		}

		mock.ExpectBegin()
		mock.ExpectExec(
			"UPDATE \"accounts\" SET \"first_name\"=$1,\"last_name\"=$2,\"email\"=$3,\"updated_at\"=$4 WHERE id = $5 AND active = TRUE").
			WithArgs(testCustomer.PersonalInfo.FirstName,
				testCustomer.PersonalInfo.LastName,
				testCustomer.PersonalInfo.Email,
				sqlmock.AnyArg(),
				testCustomer.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = customerRepo.UpdateCustomerPersonalInfo(context.Background(), customerID, testCustomer.PersonalInfo)
		require.NoError(t, err)
	})
}

func TestUpdateCustomerDeliveryInfo(t *testing.T) {
	t.Parallel()

	gdb, mock, err := NewCustomerMockDB()
	require.NoError(t, err)
	require.NotNil(t, gdb)
	require.NotNil(t, mock)

	customerRepo := NewCustomerRepositoryImpl(gdb)

	sf, err := sonyflake.NewSonyFlake()
	require.NoError(t, err)

	t.Run("UpdateCustomerDeliveryInfo", func(t *testing.T) {
		customerID, err := sf.NextID()
		require.NoError(t, err)

		testCustomer := entity.Customer{
			ID:     customerID,
			Active: true,
			DeliveryInfo: &valueobject.CustomerDeliveryInfo{
				Address:     "test address",
				PhoneNumber: "0123456789",
			},
		}

		mock.ExpectBegin()
		mock.ExpectExec(
			"UPDATE \"accounts\" SET \"address\"=$1,\"phone_number\"=$2,\"updated_at\"=$3 WHERE id = $4 AND active = TRUE").
			WithArgs(testCustomer.DeliveryInfo.Address,
				testCustomer.DeliveryInfo.PhoneNumber,
				sqlmock.AnyArg(),
				testCustomer.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = customerRepo.UpdateCustomerDeliveryInfo(context.Background(), customerID, testCustomer.DeliveryInfo)
		require.NoError(t, err)
	})
}
