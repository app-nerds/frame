package frame

// func (fa *FrameApplication) WithGoogleAuth(scopes ...string) *FrameApplication {
// 	fa.externalAuths = append(fa.externalAuths, google.New(fa.Config.GoogleClientID, fa.Config.GoogleClientSecret, fa.Config.GoogleRedirectURI, scopes...))
// 	return fa
// }

// func (fa *FrameApplication) SetupExternalAuth(pathsExcludedFromAuth, htmlPaths []string) *FrameApplication {
// 	if fa.sessionStore == nil {
// 		fa.Logger.Fatal("Please setup a session storage before calling SetupExternalAuth()")
// 	}

// 	gothic.Store = fa.sessionStore

// 	goth.UseProviders(
// 		fa.externalAuths...,
// 	)

// 	fa.router.HandleFunc("/auth/{provider}/callback", func(w http.ResponseWriter, r *http.Request) {
// 		var (
// 			err     error
// 			user    goth.User
// 			session *sessions.Session
// 			member  Member
// 		)

// 		user, err = gothic.CompleteUserAuth(w, r)

// 		if err != nil {
// 			http.Redirect(w, r, UnexpectedErrorPath, http.StatusTemporaryRedirect)
// 			return
// 		}

// 		/*
// 		 * If this member doesn't exist yet, create them as an unapproved member
// 		 */
// 		member, err = fa.MemberService.GetMemberByEmail(user.Email, true)

// 		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
// 			member = Member{
// 				AvatarURL:  user.AvatarURL,
// 				Email:      user.Email,
// 				ExternalID: user.UserID,
// 				FirstName:  user.FirstName,
// 				LastName:   user.LastName,
// 				Status: MembersStatus{
// 					ID:     MemberPendingApprovalID,
// 					Status: MemberPendingApproval,
// 				},
// 			}

// 			if err = fa.MemberService.CreateMember(&member); err != nil {
// 				fa.Logger.WithError(err).Error("error creating new member")
// 				http.Redirect(w, r, UnexpectedErrorPath, http.StatusTemporaryRedirect)
// 				return
// 			}

// 			http.Redirect(w, r, SiteAuthAccountPendingPath, http.StatusTemporaryRedirect)
// 			return
// 		}

// 		if err != nil {
// 			fa.Logger.WithError(err).Error("error getting member information in external auth")
// 			http.Redirect(w, r, UnexpectedErrorPath, http.StatusTemporaryRedirect)
// 			return
// 		}

// 		/*
// 		 * If we have an existing member, but they aren't approved yet,
// 		 * redirect them.
// 		 */
// 		if member.Status.Status != MemberActive {
// 			http.Redirect(w, r, SiteAuthAccountPendingPath, http.StatusTemporaryRedirect)
// 			return
// 		}

// 		/*
// 		 * Otherwise, we are good to go!
// 		 */
// 		if session, err = fa.sessionStore.Get(r, fa.sessionName); err != nil {
// 			fa.Logger.WithError(err).Error("error geting session")
// 			http.Redirect(w, r, UnexpectedErrorPath, http.StatusTemporaryRedirect)
// 			return
// 		}

// 		session.Values["email"] = user.Email
// 		session.Values["firstName"] = user.FirstName
// 		session.Values["lastName"] = user.LastName
// 		session.Values["avatarURL"] = user.AvatarURL
// 		// session.Values["status"] =

// 		if err = fa.sessionStore.Save(r, w, session); err != nil {
// 			fa.Logger.WithError(err).Error("error saving session")
// 			http.Redirect(w, r, UnexpectedErrorPath, http.StatusTemporaryRedirect)
// 			return
// 		}

// 		if fa.OnAuthSuccess != nil {
// 			fa.OnAuthSuccess(w, r, member)
// 		}
// 	})

// 	fa.router.HandleFunc("/auth/{provider}", func(w http.ResponseWriter, r *http.Request) {
// 		gothic.BeginAuthHandler(w, r)
// 	})

// 	fa.setupMiddleware(pathsExcludedFromAuth, htmlPaths)
// 	return fa
// }
