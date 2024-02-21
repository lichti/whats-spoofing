package main

import (
	"fmt"
	"strings"

	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func cmdGetGroup(args []string) (output string) {
	if len(args) < 1 {
		output = "\n[getgroup] Usage: getgroup <jid>"
		log.Errorf("%s", output)
		return
	}
	group, ok := parseJID(args[0])
	if !ok {
		output = "\n[getgroup] You need to specify a valid group JID"
		log.Errorf("%s", output)
		return
	} else if group.Server != types.GroupServer {
		output = fmt.Sprintf("\n[getgroup] Input must be a group JID (@%s)", types.GroupServer)
		log.Errorf("%s", output)
		return
	}
	resp, err := cli.GetGroupInfo(group)
	if err != nil {
		output = fmt.Sprintf("\n[getgroup] Failed to get group info: %v", err)
		log.Errorf("%s", output)
		return
	} else {
		output = fmt.Sprintf("\n[getgroup] Group info: %+v", resp)
		log.Infof("%s", output)
		return
	}
}

func cmdListGroups(args []string) (output string) {
	groups, err := cli.GetJoinedGroups()
	if err != nil {
		output = fmt.Sprintf("\n[listgroup] Failed to get group list: %v", err)
		log.Errorf("%s", output)
		return
	} else {
		for _, group := range groups {
			output = fmt.Sprintf("%s \n[listgroup] %+v: %+v", output, group.GroupName.Name, group.JID)
			log.Infof("%s", output)
		}
		return
	}
}

func cmdSendSpoofedReply(args []string) (output string) {
	if len(args) < 4 {
		output = "\n[send-spoofed-reply] Usage: send-spoofed-reply <chat_jid> <msgID:!|#ID> <spoofed_jid> <spoofed_text>|<text>"
		log.Errorf("%s", output)
		return
	}

	chat_jid, ok := parseJID(args[0])
	if !ok {
		output = "\n[send-spoofed-reply] You need to specify a valid Chat ID (Group or User)"
		log.Errorf("%s", output)
		return
	}

	msgID := args[1]
	if msgID[0] == '!' {
		msgID = cli.GenerateMessageID()
	}

	spoofed_jid, ok2 := parseJID(args[2])
	if !ok2 {
		output = "\n[send-spoofed-reply] You need to specify a valid User ID to spoof"
		log.Errorf("%s", output)
		return
	}

	parameters := strings.SplitN(strings.Join(args[3:], " "), "|", 2)
	spoofed_text := parameters[0]
	text := parameters[1]

	_, resp, err := sendSpoofedReplyMessage(chat_jid, spoofed_jid, msgID, spoofed_text, text)
	if err != nil {
		output = fmt.Sprintf("\n[send-spoofed-reply] Error on sending spoofed msg: %v", err)
		log.Errorf("%s", output)
		return
	} else {
		output = fmt.Sprintf("\n[send-spoofed-reply] spoofed msg sended: %+v", resp)
		log.Infof("%s", output)
		return
	}
}

func cmdSendSpoofedImgReply(args []string) (output string) {
	if len(args) < 5 {
		output = "\n[send-spoofed-img-reply] Usage: send-spoofed-img-reply <chat_jid> <msgID:!|#ID> <spoofed_jid> <spoofed_file> <spoofed_text>|<text>"
		log.Errorf("%s", output)
		return
	}
	chat_jid, ok := parseJID(args[0])
	if !ok {
		output = "\n[send-spoofed-img-reply] You need to specify a valid Chat ID (Group or User)"
		log.Errorf("%s", output)
		return
	}

	msgID := args[1]
	if msgID[0] == '!' {
		msgID = cli.GenerateMessageID()
	}

	spoofed_jid, ok2 := parseJID(args[2])
	if !ok2 {
		output = "\n[send-spoofed-img-reply] You need to specify a valid User ID to spoof"
		log.Errorf("%s", output)
		return
	}

	spoofed_file := args[3]

	parameters := strings.SplitN(strings.Join(args[4:], " "), "|", 2)
	spoofed_text := parameters[0]
	text := parameters[1]

	_, resp, err := sendSpoofedReplyImg(chat_jid, spoofed_jid, msgID, spoofed_file, spoofed_text, text)
	if err != nil {
		output = fmt.Sprintf("\n[send-spoofed-img-reply] Error on sending spoofed msg: %v", err)
		log.Errorf("%s", output)
		return
	} else {
		output = fmt.Sprintf("\n[send-spoofed-img-reply] spoofed msg sended: %+v", resp)
		log.Infof("%s", output)
		return
	}
}

func cmdSendSpoofedDemo(args []string) (output string) {
	if len(args) < 4 {
		output = "\n[send-spoofed-demo] Usage: send-spoofed-demo <toGender:boy|girl> <language:br|en|ar> <chat_jid> <spoofed_jid>"
		log.Errorf("%s", output)
		return
	}

	var toGender string
	if args[0] != "boy" && args[0] != "girl" {
		output = "\n[send-spoofed-demo] Error: <boy|girl>"
		log.Errorf("%s", output)
		return
	} else {
		toGender = args[0]
	}

	var language string
	if args[1] != "br" && args[1] != "en" && args[1] != "ar" {
		output = "\n[send-spoofed-demo] Error: <br|en|ar>"
		log.Errorf("%s", output)
		return
	} else {
		language = args[1]
	}
	
	chat_jid, ok := parseJID(args[2])
	if !ok {
		output = "\n[send-spoofed-demo] You need to specify a valid Chat ID (Group or User)"
		log.Errorf("%s", output)
		return
	}
	spoofed_jid, ok2 := parseJID(args[3])
	if !ok2 {
		output = "\n[send-spoofed-demo] You need to specify a valid User ID to spoof"
		log.Errorf("%s", output)
		return
	}
	sendSpoofedTalkDemo(chat_jid, spoofed_jid, toGender, language, "")
	output = fmt.Sprintf("\n[send-spoofed-demo] spoofed msg sended to %s as %s", chat_jid, spoofed_jid)
	return

}

func cmdSendSpoofedDemoImg(args []string) (output string) {
	if len(args) < 5 {
		log.Errorf("\n[send-spoofed-demo-img] Usage: send-spoofed-demo-img <toGender:boy|girl> <language:br|en> <chat_jid> <spoofed_jid> <spoofed_img>")
		return
	}

	var toGender string
	if args[0] != "boy" && args[0] != "girl" {
		output = "\n[send-spoofed-demo-img] Error: <boy|girl>"
		log.Errorf("%s", output)
		return
	} else {
		toGender = args[0]
	}

	var language string
	if args[1] != "br" && args[1] != "en" && args[1] != "ar" {
		output = "\n[send-spoofed-demo-img] Error: <br|en|ar>"
		log.Errorf("%s", output)
		return
	} else {
		language = args[1]
	}
	

	chat_jid, ok := parseJID(args[2])
	if !ok {
		output = "\n[send-spoofed-demo-img] You need to specify a valid Chat ID (Group or User)"
		log.Errorf("%s", output)
		return
	}
	spoofed_jid, ok2 := parseJID(args[3])
	if !ok2 {
		output = "\n[send-spoofed-demo-img] You need to specify a valid User ID to spoof"
		log.Errorf("%s", output)
		return
	}

	spoofed_img := args[4]

	sendSpoofedTalkDemo(chat_jid, spoofed_jid, toGender, language, spoofed_img)
	output = fmt.Sprintf("\n[send-spoofed-demo-img] send-spoofed-demo-img: spoofed msg sended to %s as %s", chat_jid, spoofed_jid)
	return
}

func cmdSpoofedReplyThis(args []string, msg *waProto.Message) (output string) {
	if len(args) < 4 {
		output = "\n[spoofed-reply-this] Usage: spoofed-reply-this <chat_jid> <msgID:!|#ID> <spoofed_jid> <text>"
		log.Errorf("%s", output)
		return
	}

	chat_jid, ok := parseJID(args[0])
	if !ok {
		output = "\n[send-spoofed-reply] You need to specify a valid Chat ID (Group or User)"
		log.Errorf("%s", output)
		return
	}

	msgID := args[1]
	if msgID[0] == '!' {
		msgID = cli.GenerateMessageID()
	}

	spoofed_jid, ok2 := parseJID(args[2])
	if !ok2 {
		output = "\n[send-spoofed-reply] You need to specify a valid User ID to spoof"
		log.Errorf("%s", output)
		return
	}

	text := strings.Join(args[3:], " ")

	_, resp, err := sendSpoofedReplyThis(chat_jid, spoofed_jid, msgID, text, msg)
	if err != nil {
		output = fmt.Sprintf("\n[reply-spoofed-this] Error on sending spoofed msg: %v", err)
		log.Errorf("%s", output)
		return
	} else {
		output = fmt.Sprintf("\n[reply-spoofed-this] spoofed msg sended: %+v", resp)
		log.Infof("%s", output)
		return
	}
}
