package factory

import (
	"github.com/jinzhu/gorm"
	"github.com/keyo-oliveira/codepix/application/usecases"
	"github.com/keyo-oliveira/codepix/infra/repository"
)

func TransactionUseCaseFactory(database *gorm.DB) usecases.TransactionUseCase {
	pixRepository := repository.PixKeyRepositoryDb{
		Db: database,
	}
	transactionRepository := repository.TransactionRepositoryDb{Db: database}

	transactionUseCase := usecases.TransactionUseCase{
		TransactionRepository: &transactionRepository,
		PixRepository:         pixRepository,
	}

	return transactionUseCase

}
