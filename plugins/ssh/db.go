package main

type Server struct {
	ID             int
	Name           string `validate:"required"`
	Host           string `validate:"required"`
	Port           int    `validate:"required"`
	CredentialID   int    `validate:"required"`
	CredentialName string `gorm:"-"`
}

func (Server) TableName() string {
	return "ssh_server"
}

type Credential struct {
	ID         int
	Name       string `validate:"required"`
	Username   string `validate:"required"`
	Password   string
	PrivateKey string
}

func (Credential) TableName() string {
	return "ssh_credential"
}

func (s SSH) migrateDB() {
	s.db.AutoMigrate(&Credential{}, &Server{})
}
