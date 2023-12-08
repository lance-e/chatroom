package logic

import (
	"container/ring"
	"github.com/spf13/viper"
)

type offlineProcessor struct {
	//表示保存的信息数量
	n int

	//使用链表来存储连续的信息
	//保存所有用户的最近的n条消息
	recentRing *ring.Ring

	//保存某个用户的n条离线信息
	userRing map[string]*ring.Ring
}

// OfflineProcessor 对外提供一个offlineProcessor单例
var OfflineProcessor = newOfflineProcessor()

// newOfflineProcessor 生成一个offlineProcessor实例
func newOfflineProcessor() *offlineProcessor {
	n := viper.GetInt("offline-number") //读取配置文件中设置的保存信息数
	return &offlineProcessor{
		n:          n,
		recentRing: ring.New(n),
		userRing:   make(map[string]*ring.Ring),
	}
}

// Save 进行消息的存储
func (o *offlineProcessor) Save(message *Message) {
	if message.Type != MsgTypeNormal {
		return
	}
	//消息直接存在recentRing中，并后移一位
	o.recentRing.Value = message
	o.recentRing = o.recentRing.Next()
	//接下来进行对@的用户的消息进行单独的保存
	for _, user := range message.Ats {
		user = user[1:]
		var (
			r  *ring.Ring
			ok bool
		)
		if r, ok = o.userRing[user]; !ok {
			r = ring.New(o.n)
		}
		r.Value = message
		o.userRing[user] = r.Next()
	}

}

// Send 进行消息的取出,使用ring.Do()来进行链表的遍历
func (o *offlineProcessor) Send(user *User) {
	//取出所有用户的消息
	o.recentRing.Do(func(value any) {
		if value != nil {
			user.MessageChanel <- value.(*Message)
		}
	})
	//判断用户是否是新用户
	if user.IsBool {
		return
	}
	//取出某个用户的信息
	if r, ok := o.userRing[user.NickName]; ok {
		r.Do(func(value any) {
			if value != nil {
				user.MessageChanel <- value.(*Message)
			}
		})
		delete(o.userRing, user.NickName)
	}
}
