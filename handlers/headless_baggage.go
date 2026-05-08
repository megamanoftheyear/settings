package handlers

import (
	"errors"
	"example/settings/data"
)

type HeadlessBaggageHandler struct {
	dict  data.AddressDict
	rules []*data.Rule
}

func (h *HeadlessBaggageHandler) Handle(tg *data.Telegram) (*data.OutTelegram, error) {
	receivers := make([]string, 0, len(tg.Segments))

	for _, rule := range h.rules {
		if _, ok := rule.Exec(tg.Segments...); !ok {
			continue
		}

		for _, link := range rule.Links {
			address, err := h.findAddress(link)
			if err != nil {
				return nil, err
			}

			receivers = append(receivers, address.Address)
		}

		break
	}

	if len(receivers) == 0 {
		return nil, errors.New("no telegram receivers found")
	}

	return &data.OutTelegram{Receivers: receivers}, nil
}

func (h *HeadlessBaggageHandler) findAddress(link string) (*data.Address, error) {
	addr, ok := h.dict.Get(link)
	if !ok {
		return nil, errors.New("no such address")
	}

	if addr.Status.IsBlocked() {
		return nil, errors.New("address is blocked")
	}

	return addr, nil
}
