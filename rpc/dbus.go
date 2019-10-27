package rpc

import (
	"fmt"
	"github.com/fabiofalci/sconsify/sconsify"
	"github.com/guelfey/go.dbus"
	"github.com/guelfey/go.dbus/introspect"
)

const intro = `
<node>
	<interface name="org.mpris.MediaPlayer2.Player">
		<method name="PlayPause"/>
		<method name="Next"/>
		<method name="Previous"/>
		<method name="Pause"/>
		<method name="Stop"/>
		<method name="Play"/>
	</interface>` + introspect.IntrospectDataString + `</node> `

type DbusMethods struct {
	publisher *sconsify.Publisher
}

func StartDbus(publisher *sconsify.Publisher, fallbackOnServer bool) {
	if !tryToStartDbusSession(publisher) {
		if fallbackOnServer {
			StartServer(publisher)
		}
	}
}

func tryToStartDbusSession(publisher *sconsify.Publisher) bool {
	conn, err := dbus.SessionBus()
	if err != nil {
		fmt.Println("Cannot access dbus, ignoring...")
		return false
	}
	reply, err := conn.RequestName("org.mpris.MediaPlayer2.sconsify", dbus.NameFlagDoNotQueue)
	if err != nil {
		fmt.Println("Cannot request dbus name, ignoring...")
		return false
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		fmt.Println("org.mpris.MediaPlayer2.sconsify name already taken, ignoring...")
		return false
	}
	dbusMethods := new(DbusMethods)
	dbusMethods.publisher = publisher
	err = conn.Export(dbusMethods, "/org/mpris/MediaPlayer2", "org.mpris.MediaPlayer2.Player")
	if err != nil {
		fmt.Println("Cannot export mpris Player")
		return false
	}
	err = conn.Export(introspect.Introspectable(intro), "/org/mpris/MediaPlayer2", "org.freedesktop.DBus.Introspectable")
	if err != nil {
		fmt.Println("Cannot export mpris Introspectable, ignoring...")
	}
	return true
}

func (dbus DbusMethods) PlayPause() *dbus.Error {
	dbus.publisher.PlayPauseToggle()
	return nil
}

func (dbus DbusMethods) Next() *dbus.Error {
	dbus.publisher.NextPlay()
	return nil
}

func (dbus DbusMethods) Previous() *dbus.Error {
	dbus.publisher.PreviousPlay()
	return nil
}

func (dbus DbusMethods) Pause() *dbus.Error {
	dbus.publisher.Pause()
	return nil
}

func (dbus DbusMethods) Stop() *dbus.Error {
	dbus.publisher.Stop()
	return nil
}

func (dbus DbusMethods) Play() *dbus.Error {
	dbus.publisher.Play(nil)
	return nil
}

