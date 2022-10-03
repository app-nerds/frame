package frame

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func (fa *FrameApplication) handleMemberCurrent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	member := map[string]interface{}{
		"id":        ctx.Value("memberID"),
		"email":     ctx.Value("email"),
		"firstName": ctx.Value("firstName"),
		"lastName":  ctx.Value("lastName"),
		"avatarURL": ctx.Value("avatarURL"),
	}

	fa.WriteJSON(w, http.StatusOK, member)
}

func (fa *FrameApplication) handleMemberLogout(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		session *sessions.Session
	)

	if session, err = fa.sessionStore.Get(r, fa.sessionName); err != nil {
		fa.Logger.WithError(err).Error("error getting session information")
		fa.WriteJSON(w, http.StatusInternalServerError, fa.CreateGenericErrorResponse("error getting session information", err.Error(), ""))
		return
	}

	session.Options.MaxAge = -1

	if err = fa.sessionStore.Save(r, w, session); err != nil {
		fa.Logger.WithError(err).Error("error deleting session")
		fa.WriteJSON(w, http.StatusInternalServerError, fa.CreateGenericErrorResponse("error deleting session", err.Error(), ""))
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
