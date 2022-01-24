package main

import (
	"bytes"
	"fmt"
	"net/smtp"
	"time"
)

var ACTIONS = map[string]string{
	"completeChecklist":    "a complété une checklist",
	"removedChecklistItem": "a supprimé un élément de checklist",
	"addChecklistItem":     "a ajouté un élément de checklist",
	"addChecklist":         "a ajouté une checklist",
	"a-dueAt":              "a fixé une date d'échéance",
	"removeChecklist":      "a supprimé une checklist",
	"archivedList":         "a archivé une liste",
	"unsetCustomField":     "a modifié un champ",
	"deleteComment":        "a supprimé un commentaire",
	"uncheckedItem":        "a décoché un élément de checklist sur la carte",
	"setCustomField":       "a ",
	"editComment":          "a modifié un commentaire sur la carte",
	"addComment":           "a ajouté un commentaire sur la carte",
	"createList":           "a ajouté une liste",
	"addedLabel":           "a ajouté une étiquette sur la carte ",
	"createCard":           "a créé la carte ",
	"checkedItem":          "a coché un élément de checklist",
	"restoredCard":         "a restauré la carte",
	"moveCard":             "a déplacé la carte",
}

type messages map[string][]activity

type mail struct {
	Destinataire string
	From         time.Time
	To           time.Time
	Boards       map[string]boardInfo
}

type boardInfo struct {
	Title string
	Id    string
	Slug  string
	Cards map[string]cardInfo
}

type userInfo struct {
	Name  string
	Email string
}

type cardInfo struct {
	Id            string
	Siret         string
	RaisonSociale string
	Actions       int
	Utilisateurs  map[userInfo]struct{}
}

func loadMessages(activities []activity, users map[string]user) messages {
	msgs := make(messages)
	for _, activity := range activities {
		var destinataires = []string{users[activity.Card.UserID].Services.OIDC.Email}
		for _, member := range activity.Card.Members {
			destinataires = append(destinataires, users[member].Services.OIDC.Email)
		}

		for _, destinataire := range destinataires {
			if destinataire != "" {
				activities := msgs[destinataire]
				activities = append(activities, activity)
				msgs[destinataire] = activities
			}
		}
	}

	return msgs
}

func getMail(msgs messages, from time.Time, to time.Time, users map[string]user, boards map[string]board) map[string]string {
	var mails = make(map[string]string)
	for destinataire, messages := range msgs {
		if includes(WHITELIST, destinataire) {
			var m mail
			m.Destinataire = destinataire
			m.From = from
			m.To = to
			m.Boards = make(map[string]boardInfo)
			m.group(messages, users, boards)
			m.send()
			fmt.Printf("envoi pour %s réalisé\n", m.Destinataire)
			time.Sleep(5 * time.Second)
		}
	}
	return mails
}

func includes(array []string, elem string) bool {
	for _, i := range array {
		if elem == i {
			return true
		}
	}
	return false
}

func (m *mail) group(activities []activity, users map[string]user, boards map[string]board) {
	for _, a := range activities {
		boardInfo := m.Boards[a.BoardId]
		boardInfo.Slug = boards[a.BoardId].Slug
		boardInfo.Title = boards[a.BoardId].Title
		boardInfo.Id = a.BoardId
		if boardInfo.Cards == nil {
			boardInfo.Cards = make(map[string]cardInfo)
		}
		cardInfo := boardInfo.Cards[a.Card.ID]
		user := userInfo{
			Name:  users[a.UserId].Profile.Fullname,
			Email: users[a.UserId].Services.OIDC.Email,
		}
		if cardInfo.Utilisateurs == nil {
			cardInfo.Utilisateurs = make(map[userInfo]struct{})
		}
		cardInfo.Utilisateurs[user] = struct{}{}
		cardInfo.Actions += 1
		cardInfo.RaisonSociale = a.Card.Title
		boardInfo.Cards[a.Card.ID] = cardInfo
		m.Boards[a.BoardId] = boardInfo
	}
}

func (m *mail) send() {
	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Rapport d'activité Wekan \n%s\n\n", mimeHeaders)))
	TEMPLATE.Execute(&body, m)
	err := smtp.SendMail(SMTPHOST+":"+SMTPPORT, nil, SMTPFROM, []string{m.Destinataire}, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
}
