package commands

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type msgCallback func(m msgQueue)

var msgPool map[int]*msgQueue
var msgQueueWait = 800 * time.Millisecond

type msgQueue struct {
	UserID   int
	Messages []*tgbotapi.Message
	Timer    *time.Timer
	Callback msgCallback
}

func enqueueMessage(q *msgQueue, m *tgbotapi.Message) {
	q.Messages = append(q.Messages, m)
	q.Timer.Reset(msgQueueWait)
}

func EvalMessageToQueue(update tgbotapi.Update) bool {
	var q *msgQueue
	var exists bool

	m := update.Message
	// Its not an update
	if m.ForwardDate == 0 {
		return false
	}

	uid := m.From.ID
	// There is not a Queue for that ID
	if q, exists = msgPool[uid]; !exists {
		return false
	}

	enqueueMessage(q, m)
	return true
}

func StartMsgQueue(msg *tgbotapi.Message, cb msgCallback) {
	if msgPool == nil {
		msgPool = make(map[int]*msgQueue)
	}

	var q *msgQueue
	var exists bool

	uid := msg.From.ID

	if q, exists = msgPool[uid]; !exists {
		end := func() {
			endMsgQueue(uid)
		}

		q = &msgQueue{
			UserID:   uid,
			Timer:    time.AfterFunc(msgQueueWait, end),
			Callback: cb,
		}
		msgPool[uid] = q
	}
}

func endMsgQueue(uid int) {
	if q, exists := msgPool[uid]; exists {
		delete(msgPool, uid)
		q.Callback(*q)
	}
}
