package broker

func (b *B) TopicMetadata(topics ...string) (*TopicMetadataResponse, error) {
	reqMsg := TopicMetadataRequest(topics)
	req := Request{
		RequestMessage: &reqMsg,
	}
	respMsg := &TopicMetadataResponse{}
	resp := Response{ResponseMessage: respMsg}
	if err := b.Do(&req, &resp); err != nil {
		return nil, err
	}
	return respMsg, nil
}

func (b *B) GroupCoordinator(group string) (*GroupCoordinatorResponse, error) {
	reqMsg := GroupCoordinatorRequest(group)
	req := Request{
		RequestMessage: &reqMsg,
	}
	respMsg := &GroupCoordinatorResponse{}
	resp := Response{ResponseMessage: respMsg}
	if err := b.Do(&req, &resp); err != nil {
		return nil, err
	}
	if respMsg.ErrorCode.HasError() {
		return nil, respMsg.ErrorCode
	}
	return respMsg, nil
}
