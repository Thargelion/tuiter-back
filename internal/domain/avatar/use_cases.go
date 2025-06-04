package avatar

const (
	avatarURL = "https://ui-avatars.com/api/?name="
)

func NewAvatarUseCases() *Service {
	return &Service{}
}

type UseCases interface {
	AddAvatarUseCase
}

type AddAvatarUseCase interface {
	New(userName string) string
}

type Service struct{}

func (c *Service) New(userName string) string {
	return avatarURL + userName
}
