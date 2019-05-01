package pipelines

type subscription struct {
	// Device is optional
	Device string
	// Pipeline
	PipelineId string
}

type subscriptionStore struct {
	// map[group][measurement][]subscriptions
	subscriptions map[string]map[string][]subscription
}

func NewSubscriptionStore() *subscriptionStore {
	s := &subscriptionStore{
		subscriptions: make(map[string]map[string][]subscription, 0),
	}
	return s
}

// Return all pipelines that are subscribed to given parameters
// Groups can be any number, device is optional
func (s *subscriptionStore) getSubscriptedPipelines(groups *[]string, device string, measurement string) (*[]string, error) {
	subscriptions := make([]string, 0)

	for _, v := range *groups {
		// Iterate over groups
		g := s.subscriptions[v]
		// if group matches
		if g != nil {
			m := g[measurement]
			// if measurement matches
			if m != nil {
				for _, subs := range m {
					// iterate over subscriptions
					if subs.Device == "" || subs.Device == device {
						subscriptions = append(subscriptions, subs.PipelineId)
					}
				}
			}
		}
	}
	return &subscriptions, nil
}

func (s *subscriptionStore) subscribePipeline(group string, device string, measurement string, pipeline string) error {
	sub := &subscription{
		Device:     device,
		PipelineId: pipeline,
	}

	existingGroup := s.subscriptions[group]
	if existingGroup == nil {
		s.subscriptions[group] = map[string][]subscription{}
	}

	existingMeasurement := s.subscriptions[group][measurement]
	if existingMeasurement == nil {
		s.subscriptions[group][measurement] = []subscription{}
	}

	s.subscriptions[group][measurement] = append(s.subscriptions[group][measurement], *sub)
	return nil
}
