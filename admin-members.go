package frame

import "net/http"

func (fa *FrameApplication) handleAdminMembersManage(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"appName": fa.appName,
		"scripts": JavascriptIncludes{
			{Type: "module", Src: "/pages/admin-members-manage.js"},
		},
	}

	fa.RenderTemplate(w, "admin-members-manage.tmpl", data)
}

func (fa *FrameApplication) handleAdminApiGetMembers(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		members []Member
	)

	if members, err = fa.MemberService.GetMembers(r, 1); err != nil {
		fa.Logger.WithError(err).Error("error getting members")
		fa.WriteJSON(w, http.StatusInternalServerError, fa.CreateGenericErrorResponse("There was a problem retrieving members", err.Error(), ""))
		return
	}

	fa.WriteJSON(w, http.StatusOK, members)
}
