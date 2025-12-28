package handler

import "github.com/secure-notes/internal/service"

type NoteHandler struct {
	svc *service.NoteService
}
type createNoteReq struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func NewHandler(svc *service.NoteService) *NoteHandler {
	return &NoteHandler{svc: svc}
}

type UserAuthHandler struct {
	svc *service.UserAuth
}
type createUsereAuthReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserAuthHandler(svc *service.UserAuth) *UserAuthHandler {
	return &UserAuthHandler{svc: svc}
}
