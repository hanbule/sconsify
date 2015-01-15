package sconsify

func InitialiseEngine(events *Events, ui UserInterface, askForFirstTrack bool) error {
	select {
	case playlists := <-events.PlaylistsUpdates():
		err := ui.NewPlaylists(playlists)
		if err != nil {
			return err
		}
	case <-events.ShutdownEngineUpdates():
		// TODO it is an error
		return nil
	}

	getNextToPlay := func() {
		if track := ui.GetNextToPlay(); track != nil {
			events.Play(track)
		}
	}

	if askForFirstTrack {
		getNextToPlay()
	}

	running := true
	for running {
		select {
		case track := <-events.TrackPausedUpdates():
			ui.TrackPaused(track)
		case track := <-events.TrackPlayingUpdates():
			ui.TrackPlaying(track)
		case track := <-events.TrackNotAvailableUpdates():
			ui.TrackNotAvailable(track)
		case <-events.PlayTokenLostUpdates():
			if err := ui.PlayTokenLost(); err != nil {
				running = false
			}
		case <-events.NextPlayUpdates():
			getNextToPlay()
		case newPlaylist := <-events.PlaylistsUpdates():
			ui.NewPlaylists(newPlaylist)
		case <-events.ShutdownEngineUpdates():
			running = false
		}
	}

	events.ShutdownSpotify()
	// wait for spotify shutdown
	<-events.ShutdownEngineUpdates()

	return nil
}
