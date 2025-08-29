package model

import "github.com/finfix/go-server-grpc/proto"

type GetIconsRes struct {
	Icons []Icon
}

func (s *GetIconsRes) ConvertToProto() (res *proto.GetIconsResponse, err error) {

	protoIcons := make([]*proto.Icon, 0, len(s.Icons))
	for _, icon := range s.Icons {
		protoIcon, err := icon.ConvertToProto()
		if err != nil {
			return res, err
		}
		protoIcons = append(protoIcons, protoIcon)
	}

	return &proto.GetIconsResponse{
		Error: nil,
		Icons: protoIcons,
	}, nil
}
