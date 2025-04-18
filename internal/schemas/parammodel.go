package schemas

type UserLogin struct {
	Username string `form:"username" binding:"max=15"`
	Password string `form:"password" binding:"min=5,max=20"`
	Email    string `form:"email"`
}

type UserRegister struct {
	Username string `form:"username" binding:"max=15"`
	Password string `form:"password" binding:"max=20"`
	Email    string `form:"email"`
}

type FileInfoResponse struct {
	ID   string
	Name string
}

type FileDownloadResponse struct {
	Path string
	Name string
}

type FolderInfoResponse struct {
	ID   string
	Name string
}
