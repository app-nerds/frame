package frame

func (fa *FrameApplication) getAccountPendingTemplate() string {
	return `{{template "` + fa.primaryLayoutName + `"}}
{{define "Title"}}Account Pending{{end}}

{{define "content"}}
  <div class="account-pending-page">
    <h2>Account Pending</h2>
    
    <p>
      Your account has been created but it has not been approved by an administrator yet.
      Once your account has been approved try logging in again.
    </p>
  </div>
{{end}}
  `
}

func (fa *FrameApplication) getLoginTemplate() string {
	return `{{template "` + fa.primaryLayoutName + `" .}}
{{define "Title"}}Login{{end}}

{{define "content"}}
  <div class="login-page">
    <h2>Log In</h2>

    {{if .ErrorMessage}}
      <p class="error-message">{{.ErrorMessage}}</p>
    {{else}}
      <p>
        Please enter your user name and password to log in.
      </p>
    {{end}}

    <form method="post">
      <label for="email">Email</label>
      <input type="email" id="email" name="email" required autofocus />

      <label for="password">Password</label>
      <input type="password" id="password" name="password" required />

      <footer>
        <button id="login" class="action-button">Log In</button>
      </footer>
    </form>
  </div>
{{end}}`
}

func (fa *FrameApplication) getUnexpectedErrorTemplate() string {
	return `{{template "` + fa.primaryLayoutName + `"}}
{{define "Title"}}Unexpected Error{{end}}

{{define "content"}}
  <div class="unexpected-error-page">
    <h2>Unexpected Error</h2>

    <p>
      We are sorry but the site has experienced an unexpected error. Please 
      try again later.
    </p>
  </div>
{{end}}
  `
}
