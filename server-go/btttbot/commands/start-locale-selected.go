package commands

import (
	"github.com/strongo/bots-framework/core"
	"net/url"
)

var StartLanguageSelectedCommand = bots.Command{
	Code: "start-locale",
	CallbackAction: func(whc bots.WebhookContext, callbackURL *url.URL) (m bots.MessageFromBot, err error) {
		code5 := callbackURL.Query().Get("code5")
		chatEntity := whc.ChatEntity()
		chatEntity.SetPreferredLanguage(code5)      // TODO: Should be called 1st - BAD bug
		if err = whc.SetLocale(code5); err != nil { // TODO: Should work without calling chatEntity.SetPreferredLanguage(code5)
			return
		}
		return welcomeUser(whc)
	},
}
