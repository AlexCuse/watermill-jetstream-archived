package jetstream

import (
	"fmt"
	"github.com/nats-io/nats.go"
)

type SubjectCalculator func(topic string) *Subjects

type Subjects struct {
	Primary    string
	Additional []string
}

func (s *Subjects) All() []string {
	return append([]string{s.Primary}, s.Additional...)
}

type topicInterpreter struct {
	js                nats.JetStreamContext
	subjectCalculator SubjectCalculator
	autoProvision     bool
}

func defaultSubjectCalculator(topic string) *Subjects {
	return &Subjects{
		Primary: fmt.Sprintf("%s.*", topic),
	}
}

func newTopicInterpreter(js nats.JetStreamContext, formatter SubjectCalculator, autoProvision bool) *topicInterpreter {
	if formatter == nil {
		formatter = defaultSubjectCalculator
	}

	return &topicInterpreter{
		js:                js,
		subjectCalculator: formatter,
		autoProvision:     autoProvision,
	}
}

func (b *topicInterpreter) ensureStream(topic string) error {
	var err error

	if b.autoProvision {
		_, err = b.js.StreamInfo(topic)

		if err != nil {
			_, err = b.js.AddStream(&nats.StreamConfig{
				Name:        topic,
				Description: "",
				Subjects:    b.subjectCalculator(topic).All(),
			})

			if err != nil {
				return err
			}
		}
	}

	return err
}
