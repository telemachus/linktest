package cli

// App stores information about the application's state.
type App struct {
	ExitValue     int
	HelpWanted    bool
	Verbose       bool
	VersionWanted bool
}

// NoOp determines whether an App should bail out.
func (app *App) NoOp() bool {
	return app.ExitValue != exitSuccess || app.HelpWanted || app.VersionWanted
}
