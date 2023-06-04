package avatar

func NewAvatarUseCases() UseCases {
	return &avatarService{}
}

type UseCases interface {
	AddAvatarUseCase
}

type AddAvatarUseCase interface {
	New(userName string) string
}

type avatarService struct {
}

const (
	avatarURL = "https://ui-avatars.com/api/?name="
)

func (c *avatarService) New(userName string) string {
	return avatarURL + userName
}
