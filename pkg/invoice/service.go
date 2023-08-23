package invoice

const LogService = "pkg/invoice/service"

type Service interface {
}

type service struct {
	repository Repository
}

func NewService(repository Repository) service {
	return service{repository}
}
