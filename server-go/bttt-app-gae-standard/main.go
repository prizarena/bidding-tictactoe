package bttt_app_gae_standard

import (
	"github.com/julienschmidt/httprouter"
	"github.com/strongo-games/bidding-tictactoe/server-go/btttapp"
	"github.com/strongo-games/bidding-tictactoe/server-go/btttbot-secrets"
	"github.com/strongo/bots-framework/hosts/appengine"
	"github.com/strongo/log"
	"google.golang.org/appengine"
	"html/template"
	"net/http"
)

const templatesPath = "templates/"

func init() {
	// Add log adapter for Google AppEngine
	log.AddLogger(gaehost.GaeLogger)

	httpRouter := httprouter.New()

	// Register bot HTTP handlers
	btttapp.InitApp(httpRouter, gaehost.GaeBotHost{})

	httpRouter.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		indexPage(w, r, "en")
	})

	httpRouter.GET("/ru", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		indexPage(w, r, "ru")
	})
}

var indexByLocale = make(map[string]*template.Template, 2)

func indexPage(w http.ResponseWriter, r *http.Request, locale string) {
	indexTmpl, ok := indexByLocale[locale]
	if !ok {
		indexTmpl = template.Must(template.ParseFiles(templatesPath + "index." + locale + ".html"))
		indexByLocale[locale] = indexTmpl
	}

	if err := indexTmpl.Execute(w, map[string]string{"GaTrackingID": btttbot_secrets.GaTrackingID}); err != nil {
		c := appengine.NewContext(r)
		log.Errorf(c, err.Error())
	}
}
