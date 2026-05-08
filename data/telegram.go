package data

type Telegram struct {
	TrackID  string
	ClientID string
	FromSita bool
	Segments []*Segment
}

type Segments []*Segment
type Segment struct {
	Source string
	Type   string
	Value  string
}

func (s Segments) Find(source, segmentType string) []*Segment {
	segments := make([]*Segment, 0, len(s))

	for _, segment := range s {
		if segment.Source == source && segment.Type == segmentType {
			segments = append(segments, segment)
		}
	}

	return segments
}

type OutTelegram struct {
	Receivers []string
	Sender    string
	ClientID  string
}

/*
	Владимир Василёнок (АЭРОН) - настройки
*/
