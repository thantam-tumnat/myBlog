package entities

var Topics = []string{
	UserCreated{}.TopicName(),
}

type Event interface {
	TopicName() string
}

type EventHandler interface {
	Handle(topic string, evenBytes []byte)
}

type EventProducer interface {
	Produce(event Event) error
}

type UserCreated struct {
	UserID      int    `json:"userId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UserImage   string `json:"userImage"`
}

//produce 
func (UserCreated) TopicName() string {
	return "myblogs.user.created"
}