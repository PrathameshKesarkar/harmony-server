package event

import (
	"database/sql"
	"github.com/kataras/golog"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/time/rate"
	"harmony-server/authentication"
	"harmony-server/globals"
	"harmony-server/harmonydb"
)

type GetMessagesData struct {
	Token string `mapstructure:"token"`
	Guild string `mapstructure:"guild"`
	Channel string `mapstructure:"channel"`
	LastMessage string `mapstructure:"lastmessage"`
}

type Message struct {
	Userid    string `json:"userid"`
	Createdat int    `json:"createdat"`
	Message   string `json:"message"`
	Attachment *string `json:"attachment"`
	Messageid string `json:"messageid"`
}

func OnGetMessages(ws *globals.Client, rawMap map[string]interface{}, limiter *rate.Limiter) {
	var data GetMessagesData
	if err := mapstructure.Decode(rawMap, &data); err != nil {
		sendErr(ws, "Something was wrong with your request. Please try again")
		return
	}
	if data.Guild == "" || data.Token == "" || data.Channel == "" {
		sendErr(ws, "Something was wrong with your request. Please try again")
		golog.Warnf("Error decoding getmessages request")
		return
	}
	userid ,err := authentication.VerifyToken(data.Token)
	if err != nil { // token is invalid! Get outta here!
		deauth(ws)
		return
	}
	if globals.Guilds[data.Guild] == nil || globals.Guilds[data.Guild].Clients[userid] == nil {
		return
	}
	var res *sql.Rows
	if data.LastMessage == "" {
		res, err = harmonydb.DBInst.Query("SELECT messageid, author, createdat, message, attachment FROM messages WHERE guildid=$1 AND channelid=$2 ORDER BY createdat DESC LIMIT 30", data.Guild, data.Channel)
	} else {
		res, err = harmonydb.DBInst.Query("SELECT messageid, author, createdat, message, attachment FROM messages WHERE guildid=$1 AND channelid=$2 AND createdat < (SELECT createdat FROM messages WHERE messageid=$3) ORDER BY createdat DESC LIMIT 30", data.Guild, data.Channel, data.LastMessage)
	}
	if err != nil {
		sendErr(ws, "We weren't able to get a list of messages. Please try again")
		golog.Warnf("Error getting recent messages : %v", err)
		return
	}
	var returnMsgs [] Message
	for res.Next() {
		var msg Message
		err := res.Scan(&msg.Messageid, &msg.Userid, &msg.Createdat, &msg.Message, &msg.Attachment)
		if err != nil {
			sendErr(ws, "We weren't able to get a list of messages. Please try again")
			golog.Warnf("Error scanning next row getting messages. Reason: %v", err)
			return
		}
		returnMsgs = append(returnMsgs, msg)
	}
	ws.Send(&globals.Packet{
		Type: "getmessages",
		Data: map[string]interface{}{
			"messages": returnMsgs,
		},
	})
}
