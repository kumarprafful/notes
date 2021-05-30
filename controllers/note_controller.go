package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"notes/auth"
	"notes/formaterror"
	"notes/models"
	"notes/responses"

	"github.com/valyala/fasthttp"
)

func (s *Server) CreateANote(ctx *fasthttp.RequestCtx) {
	note := models.Note{}
	err := json.Unmarshal(ctx.Request.Body(), &note)
	if err != nil {
		responses.ERROR(ctx, http.StatusBadRequest, err)
		return
	}
	userID, err := auth.ExtractTokenID(ctx)
	if err != nil {
		responses.ERROR(ctx, http.StatusForbidden, err)
		return
	}
	note.UserID = uint32(userID)
	note.Prepare()
	err = note.Validate()
	if err != nil {
		responses.ERROR(ctx, http.StatusBadRequest, err)
		return
	}
	_, err = note.SaveNote(s.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(ctx, http.StatusBadRequest, formattedError)
		return
	}
	responses.JSON(ctx, http.StatusCreated, map[string]string{
		"note_id": fmt.Sprint(note.ID),
	})
}

func (s *Server) CreateContent(ctx *fasthttp.RequestCtx) {
	content := models.Content{}
	err := json.Unmarshal(ctx.Request.Body(), &content)
	if err != nil {
		responses.ERROR(ctx, http.StatusBadRequest, err)
		return
	}
	content.Prepare()
	err = content.Validate()
	if err != nil {
		responses.ERROR(ctx, http.StatusBadRequest, err)
		return
	}
	_, err = content.SaveContent(s.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(ctx, http.StatusBadRequest, formattedError)
		return
	}
	responses.JSON(ctx, http.StatusCreated, map[string]string{
		"content_id": fmt.Sprint(content.ID),
	})
}

func (s *Server) FetchANote(ctx *fasthttp.RequestCtx) {
	note_id, err := ctx.QueryArgs().GetUint("note_id")
	if err != nil {
		responses.ERROR(ctx, http.StatusBadRequest, errors.New("note_id is invalid or not provided"))
		return
	}
	n := models.Note{}
	note, err := n.GetNoteByID(s.DB, uint32(note_id))
	if err != nil {
		responses.ERROR(ctx, http.StatusBadRequest, err)
		return
	}
	responses.JSON(ctx, http.StatusOK, note)
}

func (s *Server) FetchContentOfNote(ctx *fasthttp.RequestCtx) {
	note_id, err := ctx.QueryArgs().GetUint("note_id")
	if err != nil {
		responses.ERROR(ctx, http.StatusBadRequest, errors.New("note_id is invalid or not provided"))
		return
	}
	n := models.Note{}
	contents, err := n.ContentsOfNotes(s.DB, uint32(note_id))
	if err != nil {
		responses.ERROR(ctx, http.StatusBadRequest, err)
		return
	}
	responses.JSON(ctx, http.StatusOK, contents)
}
