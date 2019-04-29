// Converts between protobuf structs and Go representation of packets

package main

import (
	"encoding/json"
	"log"
	"time"
	"chatbot_go/pbx"
	"chatbot_go/types"
)

func pbServCtrlSerialize(ctrl *MsgServerCtrl) *pbx.ServerMsg_Ctrl {
	var params map[string][]byte
	if ctrl.Params != nil {
		if in, ok := ctrl.Params.(map[string]interface{}); ok {
			params = interfaceMapToByteMap(in)
		}
	}

	return &pbx.ServerMsg_Ctrl{Ctrl: &pbx.ServerCtrl{
		Id:     ctrl.Id,
		Topic:  ctrl.Topic,
		Code:   int32(ctrl.Code),
		Text:   ctrl.Text,
		Params: params}}
}

func pbServPresSerialize(pres *MsgServerPres) *pbx.ServerMsg_Pres {
	var what pbx.ServerPres_What
	switch pres.What {
	case "on":
		what = pbx.ServerPres_ON
	case "off":
		what = pbx.ServerPres_OFF
	case "ua":
		what = pbx.ServerPres_UA
	case "upd":
		what = pbx.ServerPres_UPD
	case "gone":
		what = pbx.ServerPres_GONE
	case "acs":
		what = pbx.ServerPres_ACS
	case "term":
		what = pbx.ServerPres_TERM
	case "msg":
		what = pbx.ServerPres_MSG
	case "read":
		what = pbx.ServerPres_READ
	case "recv":
		what = pbx.ServerPres_RECV
	case "del":
		what = pbx.ServerPres_DEL
	default:
		log.Fatal("Unknown pres.what value", pres.What)
	}
	return &pbx.ServerMsg_Pres{Pres: &pbx.ServerPres{
		Topic:        pres.Topic,
		Src:          pres.Src,
		What:         what,
		UserAgent:    pres.UserAgent,
		SeqId:        int32(pres.SeqId),
		DelId:        int32(pres.DelId),
		DelSeq:       pbDelQuerySerialize(pres.DelSeq),
		TargetUserId: pres.AcsTarget,
		ActorUserId:  pres.AcsActor,
		Acs:          pbAccessModeSerialize(pres.Acs)}}
}

func pbServInfoSerialize(info *MsgServerInfo) *pbx.ServerMsg_Info {
	return &pbx.ServerMsg_Info{Info: &pbx.ServerInfo{
		Topic:      info.Topic,
		FromUserId: info.From,
		What:       pbInfoNoteWhatSerialize(info.What),
		SeqId:      int32(info.SeqId),
	}}
}

func pbServMetaSerialize(meta *MsgServerMeta) *pbx.ServerMsg_Meta {
	return &pbx.ServerMsg_Meta{Meta: &pbx.ServerMeta{
		Id:    meta.Id,
		Topic: meta.Topic,
		Desc:  pbTopicDescSerialize(meta.Desc),
		Sub:   pbTopicSubSliceSerialize(meta.Sub),
		Del:   pbDelValuesSerialize(meta.Del),
	}}
}

// Convert ClientComMessage to pbx.ClientMsg
func pbCliSerialize(msg *ClientComMessage) *pbx.ClientMsg {
	var pkt pbx.ClientMsg

	switch {
	case msg.Hi != nil:
		pkt.Message = &pbx.ClientMsg_Hi{Hi: &pbx.ClientHi{
			Id:        msg.Hi.Id,
			UserAgent: msg.Hi.UserAgent,
			Ver:       msg.Hi.Version,
			DeviceId:  msg.Hi.DeviceID,
			Platform:  msg.Hi.Platform,
			Lang:      msg.Hi.Lang}}
	case msg.Acc != nil:
		pkt.Message = &pbx.ClientMsg_Acc{Acc: &pbx.ClientAcc{
			Id:     msg.Acc.Id,
			UserId: msg.Acc.User,
			Token:  msg.Acc.Token,
			Scheme: msg.Acc.Scheme,
			Secret: msg.Acc.Secret,
			Login:  msg.Acc.Login,
			Tags:   msg.Acc.Tags,
			Cred:   pbCredentialsSerialize(msg.Acc.Cred),
			Desc:   pbSetDescSerialize(msg.Acc.Desc)}}
	case msg.Login != nil:
		pkt.Message = &pbx.ClientMsg_Login{Login: &pbx.ClientLogin{
			Id:     msg.Login.Id,
			Scheme: msg.Login.Scheme,
			Secret: msg.Login.Secret,
			Cred:   pbCredentialsSerialize(msg.Login.Cred)}}
	case msg.Sub != nil:
		pkt.Message = &pbx.ClientMsg_Sub{Sub: &pbx.ClientSub{
			Id:       msg.Sub.Id,
			Topic:    msg.Sub.Topic,
			SetQuery: pbSetQuerySerialize(msg.Sub.Set),
			GetQuery: pbGetQuerySerialize(msg.Sub.Get)}}
	case msg.Leave != nil:
		pkt.Message = &pbx.ClientMsg_Leave{Leave: &pbx.ClientLeave{
			Id:    msg.Leave.Id,
			Topic: msg.Leave.Topic,
			Unsub: msg.Leave.Unsub}}
	case msg.Pub != nil:
		pkt.Message = &pbx.ClientMsg_Pub{Pub: &pbx.ClientPub{
			Id:      msg.Pub.Id,
			Topic:   msg.Pub.Topic,
			NoEcho:  msg.Pub.NoEcho,
			Head:    interfaceMapToByteMap(msg.Pub.Head),
			Content: interfaceToBytes(msg.Pub.Content)}}
	case msg.Get != nil:
		pkt.Message = &pbx.ClientMsg_Get{Get: &pbx.ClientGet{
			Id:    msg.Get.Id,
			Topic: msg.Get.Topic,
			Query: pbGetQuerySerialize(&msg.Get.MsgGetQuery)}}
	case msg.Set != nil:
		pkt.Message = &pbx.ClientMsg_Set{Set: &pbx.ClientSet{
			Id:    msg.Set.Id,
			Topic: msg.Set.Topic,
			Query: pbSetQuerySerialize(&msg.Set.MsgSetQuery)}}
	case msg.Del != nil:
		var what pbx.ClientDel_What
		switch msg.Del.What {
		case "msg":
			what = pbx.ClientDel_MSG
		case "topic":
			what = pbx.ClientDel_TOPIC
		case "sub":
			what = pbx.ClientDel_SUB
		}
		pkt.Message = &pbx.ClientMsg_Del{Del: &pbx.ClientDel{
			Id:     msg.Del.Id,
			Topic:  msg.Del.Topic,
			What:   what,
			DelSeq: pbDelQuerySerialize(msg.Del.DelSeq),
			UserId: msg.Del.User,
			Hard:   msg.Del.Hard}}
	case msg.Note != nil:
		pkt.Message = &pbx.ClientMsg_Note{Note: &pbx.ClientNote{
			Topic: msg.Note.Topic,
			What:  pbInfoNoteWhatSerialize(msg.Note.What),
			SeqId: int32(msg.Note.SeqId)}}
	}

	if pkt.Message == nil {
		return nil
	}

	pkt.OnBehalfOf = msg.from
	pkt.AuthLevel = pbx.AuthLevel(msg.authLvl)

	return &pkt
}

// Convert pbx.ClientMsg to ClientComMessage
func pbCliDeserialize(pkt *pbx.ClientMsg) *ClientComMessage {
	var msg ClientComMessage
	if hi := pkt.GetHi(); hi != nil {
		msg.Hi = &MsgClientHi{
			Id:        hi.GetId(),
			UserAgent: hi.GetUserAgent(),
			Version:   hi.GetVer(),
			DeviceID:  hi.GetDeviceId(),
			Platform:  hi.GetPlatform(),
			Lang:      hi.GetLang(),
		}
	} else if acc := pkt.GetAcc(); acc != nil {
		msg.Acc = &MsgClientAcc{
			Id:     acc.GetId(),
			User:   acc.GetUserId(),
			Scheme: acc.GetScheme(),
			Secret: acc.GetSecret(),
			Login:  acc.GetLogin(),
			Tags:   acc.GetTags(),
			Desc:   pbSetDescDeserialize(acc.GetDesc()),
			Cred:   pbCredentialsDeserialize(acc.GetCred()),
		}
	} else if login := pkt.GetLogin(); login != nil {
		msg.Login = &MsgClientLogin{
			Id:     login.GetId(),
			Scheme: login.GetScheme(),
			Secret: login.GetSecret(),
			Cred:   pbCredentialsDeserialize(login.GetCred()),
		}
	} else if sub := pkt.GetSub(); sub != nil {
		msg.Sub = &MsgClientSub{
			Id:    sub.GetId(),
			Topic: sub.GetTopic(),
			Get:   pbGetQueryDeserialize(sub.GetGetQuery()),
			Set:   pbSetQueryDeserialize(sub.GetSetQuery()),
		}
	} else if leave := pkt.GetLeave(); leave != nil {
		msg.Leave = &MsgClientLeave{
			Id:    leave.GetId(),
			Topic: leave.GetTopic(),
			Unsub: leave.GetUnsub(),
		}
	} else if pub := pkt.GetPub(); pub != nil {
		msg.Pub = &MsgClientPub{
			Id:      pub.GetId(),
			Topic:   pub.GetTopic(),
			NoEcho:  pub.GetNoEcho(),
			Head:    byteMapToInterfaceMap(pub.GetHead()),
			Content: bytesToInterface(pub.GetContent()),
		}
	} else if get := pkt.GetGet(); get != nil {
		msg.Get = &MsgClientGet{
			Id:          get.GetId(),
			Topic:       get.GetTopic(),
			MsgGetQuery: *pbGetQueryDeserialize(get.GetQuery()),
		}
	} else if set := pkt.GetSet(); set != nil {
		msg.Set = &MsgClientSet{
			Id:          set.GetId(),
			Topic:       set.GetTopic(),
			MsgSetQuery: *pbSetQueryDeserialize(set.GetQuery()),
		}
	} else if del := pkt.GetDel(); del != nil {
		msg.Del = &MsgClientDel{
			Id:     del.GetId(),
			Topic:  del.GetTopic(),
			DelSeq: pbDelQueryDeserialize(del.GetDelSeq()),
			User:   del.GetUserId(),
			Hard:   del.GetHard(),
		}
		switch del.GetWhat() {
		case pbx.ClientDel_MSG:
			msg.Del.What = "msg"
		case pbx.ClientDel_TOPIC:
			msg.Del.What = "topic"
		case pbx.ClientDel_SUB:
			msg.Del.What = "sub"
		}
	} else if note := pkt.GetNote(); note != nil {
		msg.Note = &MsgClientNote{
			Topic: note.GetTopic(),
			SeqId: int(note.GetSeqId()),
		}
		switch note.GetWhat() {
		case pbx.InfoNote_READ:
			msg.Note.What = "read"
		case pbx.InfoNote_RECV:
			msg.Note.What = "recv"
		case pbx.InfoNote_KP:
			msg.Note.What = "kp"
		}
	}

	msg.from = pkt.GetOnBehalfOf()
	msg.authLvl = int(pkt.GetAuthLevel())

	return &msg
}

func interfaceMapToByteMap(in map[string]interface{}) map[string][]byte {
	out := make(map[string][]byte, len(in))
	for key, val := range in {
		if val != nil {
			out[key], _ = json.Marshal(val)
		}
	}
	return out
}

func byteMapToInterfaceMap(in map[string][]byte) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for key, raw := range in {
		if val := bytesToInterface(raw); val != nil {
			out[key] = val
		}
	}
	return out
}

func interfaceToBytes(in interface{}) []byte {
	if in != nil {
		out, _ := json.Marshal(in)
		return out
	}
	return nil
}

func bytesToInterface(in []byte) interface{} {
	var out interface{}
	if len(in) > 0 {
		err := json.Unmarshal(in, &out)
		if err != nil {
			log.Println("pbx: failed to parse bytes", string(in), err)
		}
	}
	return out
}

func intSliceToInt32(in []int) []int32 {
	out := make([]int32, len(in))
	for i, v := range in {
		out[i] = int32(v)
	}
	return out
}

func int32SliceToInt(in []int32) []int {
	out := make([]int, len(in))
	for i, v := range in {
		out[i] = int(v)
	}
	return out
}

func timeToInt64(ts *time.Time) int64 {
	if ts != nil {
		return ts.UnixNano() / int64(time.Millisecond)
	}
	return 0
}

func int64ToTime(ts int64) *time.Time {
	if ts > 0 {
		res := time.Unix(ts/1000, ts%1000).UTC()
		return &res
	}
	return nil
}

func pbGetQuerySerialize(in *MsgGetQuery) *pbx.GetQuery {
	if in == nil {
		return nil
	}

	out := &pbx.GetQuery{
		What: in.What,
	}

	if in.Desc != nil {
		out.Desc = &pbx.GetOpts{
			IfModifiedSince: timeToInt64(in.Desc.IfModifiedSince),
			User:            in.Desc.User,
			Topic:           in.Desc.Topic,
			Limit:           int32(in.Desc.Limit)}
	}
	if in.Sub != nil {
		out.Sub = &pbx.GetOpts{
			IfModifiedSince: timeToInt64(in.Sub.IfModifiedSince),
			User:            in.Sub.User,
			Topic:           in.Sub.Topic,
			Limit:           int32(in.Sub.Limit)}
	}
	if in.Data != nil {
		out.Data = &pbx.GetOpts{
			BeforeId: int32(in.Data.BeforeId),
			SinceId:  int32(in.Data.SinceId),
			Limit:    int32(in.Data.Limit)}
	}
	return out
}

func pbGetQueryDeserialize(in *pbx.GetQuery) *MsgGetQuery {
	msg := MsgGetQuery{}

	if in != nil {
		msg.What = in.GetWhat()

		if desc := in.GetDesc(); desc != nil {
			msg.Desc = &MsgGetOpts{
				IfModifiedSince: int64ToTime(desc.GetIfModifiedSince()),
				Limit:           int(desc.GetLimit()),
			}
		}
		if sub := in.GetSub(); sub != nil {
			msg.Desc = &MsgGetOpts{
				IfModifiedSince: int64ToTime(sub.GetIfModifiedSince()),
				Limit:           int(sub.GetLimit()),
			}
		}
		if data := in.GetData(); data != nil {
			msg.Data = &MsgGetOpts{
				BeforeId: int(data.GetBeforeId()),
				SinceId:  int(data.GetSinceId()),
				Limit:    int(data.GetLimit()),
			}
		}
	}

	return &msg
}

func pbSetDescSerialize(in *MsgSetDesc) *pbx.SetDesc {
	if in == nil {
		return nil
	}

	return &pbx.SetDesc{
		DefaultAcs: pbDefaultAcsSerialize(in.DefaultAcs),
		Public:     interfaceToBytes(in.Public),
		Private:    interfaceToBytes(in.Private),
	}
}

func pbSetDescDeserialize(in *pbx.SetDesc) *MsgSetDesc {
	if in == nil {
		return nil
	}

	return &MsgSetDesc{
		DefaultAcs: pbDefaultAcsDeserialize(in.GetDefaultAcs()),
		Public:     bytesToInterface(in.GetPublic()),
		Private:    bytesToInterface(in.GetPrivate()),
	}
}

func pbSetQuerySerialize(in *MsgSetQuery) *pbx.SetQuery {
	if in == nil {
		return nil
	}

	out := &pbx.SetQuery{
		Desc: pbSetDescSerialize(in.Desc),
	}

	if in.Sub != nil {
		out.Sub = &pbx.SetSub{
			UserId: in.Sub.User,
			Mode:   in.Sub.Mode,
		}
	}
	return out
}

func pbSetQueryDeserialize(in *pbx.SetQuery) *MsgSetQuery {
	msg := MsgSetQuery{}

	if in != nil {
		if desc := in.GetDesc(); desc != nil {
			msg.Desc = pbSetDescDeserialize(desc)
		}
		if sub := in.GetSub(); sub != nil {
			msg.Sub = &MsgSetSub{
				User: sub.GetUserId(),
				Mode: sub.GetMode(),
			}
		}
	}

	return &msg
}

func pbInfoNoteWhatSerialize(what string) pbx.InfoNote {
	var out pbx.InfoNote
	switch what {
	case "kp":
		out = pbx.InfoNote_KP
	case "read":
		out = pbx.InfoNote_READ
	case "recv":
		out = pbx.InfoNote_RECV
	default:
		log.Fatal("unknown info-note.what", what)
	}
	return out
}

func pbInfoNoteWhatDeserialize(what pbx.InfoNote) string {
	var out string
	switch what {
	case pbx.InfoNote_KP:
		out = "kp"
	case pbx.InfoNote_READ:
		out = "read"
	case pbx.InfoNote_RECV:
		out = "recv"
	default:
		log.Fatal("unknown info-note.what", what)
	}
	return out
}

func pbAccessModeSerialize(acs *MsgAccessMode) *pbx.AccessMode {
	if acs == nil {
		return nil
	}

	return &pbx.AccessMode{
		Want:  acs.Want,
		Given: acs.Given,
	}
}

func pbAccessModeDeserialize(acs *pbx.AccessMode) *MsgAccessMode {
	if acs == nil {
		return nil
	}

	return &MsgAccessMode{
		Want:  acs.Want,
		Given: acs.Given,
	}
}

func pbDefaultAcsSerialize(defacs *MsgDefaultAcsMode) *pbx.DefaultAcsMode {
	if defacs == nil {
		return nil
	}

	return &pbx.DefaultAcsMode{
		Auth: defacs.Auth,
		Anon: defacs.Anon}
}

func pbDefaultAcsDeserialize(defacs *pbx.DefaultAcsMode) *MsgDefaultAcsMode {
	if defacs == nil {
		return nil
	}

	return &MsgDefaultAcsMode{
		Auth: defacs.GetAuth(),
		Anon: defacs.GetAnon(),
	}
}

func pbTopicDescSerialize(desc *MsgTopicDesc) *pbx.TopicDesc {
	if desc == nil {
		return nil
	}
	return &pbx.TopicDesc{
		CreatedAt: timeToInt64(desc.CreatedAt),
		UpdatedAt: timeToInt64(desc.UpdatedAt),
		TouchedAt: timeToInt64(desc.TouchedAt),
		Defacs:    pbDefaultAcsSerialize(desc.DefaultAcs),
		Acs:       pbAccessModeSerialize(desc.Acs),
		SeqId:     int32(desc.SeqId),
		ReadId:    int32(desc.ReadSeqId),
		RecvId:    int32(desc.RecvSeqId),
		DelId:     int32(desc.DelId),
		Public:    interfaceToBytes(desc.Public),
		Private:   interfaceToBytes(desc.Private),
	}
}

func pbTopicDescDeserialize(desc *pbx.TopicDesc) *MsgTopicDesc {
	if desc == nil {
		return nil
	}
	return &MsgTopicDesc{
		CreatedAt:  int64ToTime(desc.GetCreatedAt()),
		UpdatedAt:  int64ToTime(desc.GetUpdatedAt()),
		TouchedAt:  int64ToTime(desc.GetTouchedAt()),
		DefaultAcs: pbDefaultAcsDeserialize(desc.GetDefacs()),
		Acs:        pbAccessModeDeserialize(desc.GetAcs()),
		SeqId:      int(desc.SeqId),
		ReadSeqId:  int(desc.ReadId),
		RecvSeqId:  int(desc.RecvId),
		DelId:      int(desc.DelId),
		Public:     bytesToInterface(desc.Public),
		Private:    bytesToInterface(desc.Private),
	}
}


func pbTopicSubSliceSerialize(subs []MsgTopicSub) []*pbx.TopicSub {
	if subs == nil || len(subs) == 0 {
		return nil
	}

	out := make([]*pbx.TopicSub, len(subs))
	for i := 0; i < len(subs); i++ {
		out[i] = pbTopicSubSerialize(&subs[i])
	}
	return out
}

func pbTopicSubSerialize(sub *MsgTopicSub) *pbx.TopicSub {
	out := &pbx.TopicSub{
		UpdatedAt: timeToInt64(sub.UpdatedAt),
		DeletedAt: timeToInt64(sub.DeletedAt),
		Online:    sub.Online,
		Acs:       pbAccessModeSerialize(&sub.Acs),
		ReadId:    int32(sub.ReadSeqId),
		RecvId:    int32(sub.RecvSeqId),
		Public:    interfaceToBytes(sub.Public),
		Private:   interfaceToBytes(sub.Private),
		UserId:    sub.User,
		Topic:     sub.Topic,
		TouchedAt: timeToInt64(sub.TouchedAt),
		SeqId:     int32(sub.SeqId),
		DelId:     int32(sub.DelId),
	}
	if sub.LastSeen != nil {
		out.LastSeenTime = timeToInt64(sub.LastSeen.When)
		out.LastSeenUserAgent = sub.LastSeen.UserAgent
	}
	return out
}

func pbTopicSubSliceDeserialize(subs []*pbx.TopicSub) []MsgTopicSub {
	if subs == nil || len(subs) == 0 {
		return nil
	}

	out := make([]MsgTopicSub, len(subs))
	for i := 0; i < len(subs); i++ {
		out[i] = MsgTopicSub{
			UpdatedAt: int64ToTime(subs[i].GetUpdatedAt()),
			DeletedAt: int64ToTime(subs[i].GetDeletedAt()),
			Online:    subs[i].GetOnline(),
			ReadSeqId: int(subs[i].GetReadId()),
			RecvSeqId: int(subs[i].GetRecvId()),
			Public:    bytesToInterface(subs[i].GetPublic()),
			Private:   bytesToInterface(subs[i].GetPrivate()),
			User:      subs[i].GetUserId(),
			Topic:     subs[i].GetTopic(),
			TouchedAt: int64ToTime(subs[i].GetTouchedAt()),
			SeqId:     int(subs[i].GetSeqId()),
			DelId:     int(subs[i].GetDelId()),
		}
		if acs := subs[i].GetAcs(); acs != nil {
			out[i].Acs = *pbAccessModeDeserialize(acs)
		}
		if subs[i].GetLastSeenTime() > 0 {
			out[i].LastSeen = &MsgLastSeenInfo{
				When:      int64ToTime(subs[i].GetLastSeenTime()),
				UserAgent: subs[i].GetLastSeenUserAgent(),
			}
		}
	}
	return out
}

func pbSubSliceDeserialize(subs []*pbx.TopicSub) []types.Subscription {
	if subs == nil || len(subs) == 0 {
		return nil
	}

	out := make([]types.Subscription, len(subs))
	for i := 0; i < len(subs); i++ {
		out[i] = types.Subscription{
			ObjHeader: types.ObjHeader{
				UpdatedAt: *int64ToTime(subs[i].GetUpdatedAt()),
				DeletedAt: int64ToTime(subs[i].GetDeletedAt()),
			},
			User:    subs[i].GetUserId(),
			Topic:   subs[i].GetTopic(),
			DelId:   int(subs[i].GetDelId()),
			Private: bytesToInterface(subs[i].GetPrivate()),
		}
		out[i].SetPublic(bytesToInterface(subs[i].GetPublic()))
		if acs := subs[i].GetAcs(); acs != nil {
			out[i].ModeGiven.UnmarshalText([]byte(acs.GetGiven()))
			out[i].ModeWant.UnmarshalText([]byte(acs.GetWant()))
		}
		if subs[i].GetLastSeenTime() > 0 {
			out[i].SetLastSeenAndUA(int64ToTime(subs[i].GetLastSeenTime()),
				subs[i].GetLastSeenUserAgent())
		}
	}
	return out
}

func pbDelQuerySerialize(in []MsgDelRange) []*pbx.SeqRange {
	if in == nil {
		return nil
	}

	out := make([]*pbx.SeqRange, len(in))
	for i, dq := range in {
		out[i] = &pbx.SeqRange{Low: int32(dq.LowId), Hi: int32(dq.HiId)}
	}

	return out
}

func pbDelQueryDeserialize(in []*pbx.SeqRange) []MsgDelRange {
	if in == nil {
		return nil
	}

	out := make([]MsgDelRange, len(in))
	for i, sr := range in {
		out[i].LowId = int(sr.GetLow())
		out[i].HiId = int(sr.GetHi())
	}

	return out
}

func pbDelValuesSerialize(in *MsgDelValues) *pbx.DelValues {
	if in == nil {
		return nil
	}

	return &pbx.DelValues{
		DelId:  int32(in.DelId),
		DelSeq: pbDelQuerySerialize(in.DelSeq),
	}
}

func pbDelValuesDeserialize(in *pbx.DelValues) *MsgDelValues {
	if in == nil {
		return nil
	}

	return &MsgDelValues{
		DelId:  int(in.GetDelId()),
		DelSeq: pbDelQueryDeserialize(in.GetDelSeq()),
	}
}

func pbCredentialsSerialize(in []MsgAccCred) []*pbx.Credential {
	if in == nil {
		return nil
	}

	out := make([]*pbx.Credential, len(in))
	for i := range in {
		cr := &in[i]
		out[i] = &pbx.Credential{
			Method:   cr.Method,
			Value:    cr.Value,
			Response: cr.Response,
			Params:   interfaceToBytes(cr.Params)}
	}

	return out
}

func pbCredentialsDeserialize(in []*pbx.Credential) []MsgAccCred {
	if in == nil {
		return nil
	}

	out := make([]MsgAccCred, len(in))
	for i, cr := range in {
		out[i].Method = cr.GetMethod()
		out[i].Value = cr.GetValue()
		out[i].Response = cr.GetResponse()
		out[i].Params = bytesToInterface(cr.GetParams())
	}

	return out
}
