package avatar

const (
	avatarURL = "https://ui-avatars.com/api/?name="
)

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

func (c *avatarService) New(userName string) string {
	return avatarURL + userName
}
