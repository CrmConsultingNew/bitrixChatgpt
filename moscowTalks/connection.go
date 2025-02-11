package moscowTalks

import (
	"bitrix_app/backend/bitrix/authorize"
	"fmt"
	"io"
	"log"
	"net/http"
)

var GlobalAuthIdMoscowTalks string

func ConnectionBitrixLocalAppMoscowTalks(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection Moscow talks is starting...")
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading request body:", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
	log.Println("Moscow talks response:", string(bs))
	defer r.Body.Close()

	authValues := authorize.ParseValues(w, bs) //todo here we must to add this data in dbase?
	fmt.Printf("Moscow talks -> authValues.AuthID : %s, authValues.MemberID: %s", authValues.AuthID, authValues.MemberID)

	//w.Write([]byte(authValues.AuthID))
	redirectURL := "https://crmconsulting-api.ru/moscowTalks" // TODO

	// Use http.Redirect to redirect the client
	// The http.StatusFound status code is commonly used for redirects
	http.Redirect(w, r, redirectURL, http.StatusFound)

	fmt.Println("redirect is done...")
	GlobalAuthIdMoscowTalks = authValues.AuthID

	//events.OnCrmDealAddEventRegistration(authValues.AuthID) //todo return this method
}
