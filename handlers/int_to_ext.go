package handlers

import (
	"errors"
	"example/settings/data"
)

type IntToExtHandler struct {
	rules []*data.Rule // map[client_id + context][]*data.Rule
	dict  data.AddressDict
}

func (h *IntToExtHandler) Handle(telegram *data.Telegram) ([]*data.OutTelegram, error) {
	matches, err := h.findData(telegram.ClientID, telegram.Segments...)
	if err != nil {
		return nil, err
	}

	out := make([]*data.OutTelegram, 0, len(matches))
	for sender, addresses := range matches {
		outRcv := make([]string, 0, len(addresses))
		for _, address := range addresses {
			outRcv = append(outRcv, address.Address)
		}

		out = append(out, &data.OutTelegram{
			Receivers: outRcv,
			Sender:    sender.Address,
			ClientID:  telegram.ClientID,
		})
	}

	return out, nil
}

// не работает корректно, но это сейчас не нужно, т.к. цель отразить работу с настройками
// проигнорирован дефолтный адрес, а также делегированные
func (h *IntToExtHandler) findData(clientID string, segments ...*data.Segment) (map[*data.Address][]*data.Address, error) {
	result := make(map[*data.Address][]*data.Address, len(segments))

	for _, rule := range h.rules {
		matchSegments, ok := rule.Exec(segments...)

		// должна быть ошибка, если матчей(получателей) больше 1 для 1 правила
		if !ok || len(matchSegments) > 0 || len(rule.Links) != 1 {
			continue
		}

		receiver, err := h.findAddress(clientID, matchSegments[0].Value)
		if err != nil {
			return nil, err
		}

		sender, err := h.findAddress(clientID, rule.Links[0])
		if err != nil {
			return nil, err
		}

		result[sender] = append(result[sender], receiver)
	}

	return result, nil
}

func (h *IntToExtHandler) findReceivers(clientID string, segments ...*data.Segment) ([]*data.Address, error) {
	receivers := make([]*data.Address, 0, len(segments))

	for _, segment := range segments {
		receiver, err := h.findAddress(clientID, segment.Value)
		if err != nil {
			return nil, err
		}

		receivers = append(receivers, receiver)
	}

	return receivers, nil
}

func (h *IntToExtHandler) findAddress(clientID string, address string) (*data.Address, error) {
	found, exists := h.dict.Get(address)
	if !exists {
		return nil, errors.New("address not found")
	}

	if found.Status.IsBlocked() {
		return nil, errors.New("address is blocked")
	}

	if found.ClientID != clientID {
		return nil, errors.New("client id not match")
	}

	return found, nil
}
