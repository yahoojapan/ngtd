//
// Copyright (C) 2018 Yahoo Japan Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package handler

import (
	"io"

	"github.com/yahoojapan/gongt"
	pb "github.com/yahoojapan/ngtd/proto"
	"github.com/yahoojapan/ngtd/service"
	"golang.org/x/net/context"
)

type GRPC struct{}

func (g *GRPC) Search(ctx context.Context, in *pb.SearchRequest) (*pb.SearchResponse, error) {
	result, err := service.Search(in.Vector, int(in.Size_), in.Epsilon)
	if err != nil {
		return nil, err
	}

	return toSearchResponse(result), nil
}

func (g *GRPC) SearchByID(ctx context.Context, in *pb.SearchRequest) (*pb.SearchResponse, error) {
	result, err := service.SearchByID(in.Id, int(in.Size_), in.Epsilon)
	if err != nil {
		return nil, err
	}

	return toSearchResponse(result), nil
}

func (g *GRPC) StreamSearch(srv pb.NGTD_StreamSearchServer) error {
	for {
		in, err := srv.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		result, err := service.Search(in.Vector, int(in.Size_), in.Epsilon)
		if err != nil {
			srv.Send(&pb.SearchResponse{Error: err.Error()})
		} else {
			srv.Send(toSearchResponse(result))
		}
	}
}

func (g *GRPC) StreamSearchByID(srv pb.NGTD_StreamSearchByIDServer) error {
	for {
		in, err := srv.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		result, err := service.SearchByID(in.Id, int(in.Size_), in.Epsilon)
		if err != nil {
			srv.Send(&pb.SearchResponse{Error: err.Error()})
		} else {
			srv.Send(toSearchResponse(result))
		}
	}
}

func (g *GRPC) Insert(ctx context.Context, in *pb.InsertRequest) (*pb.InsertResponse, error) {
	if err := service.Insert(in.Vector, in.Id); err != nil {
		return nil, err
	}
	return &pb.InsertResponse{}, nil
}

func (g *GRPC) StreamInsert(srv pb.NGTD_StreamInsertServer) error {
	for {
		in, err := srv.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		if err := service.Insert(in.Vector, in.Id); err != nil {
			srv.Send(&pb.InsertResponse{Error: err.Error()})
		}
	}
}

func (g *GRPC) Remove(ctx context.Context, in *pb.RemoveRequest) (*pb.RemoveResponse, error) {
	if err := service.Remove(in.Id); err != nil {
		return nil, err
	}
	return &pb.RemoveResponse{}, nil
}

func (g *GRPC) StreamRemove(srv pb.NGTD_StreamRemoveServer) error {
	for {
		in, err := srv.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		if err := service.Remove(in.Id); err != nil {
			srv.Send(&pb.RemoveResponse{Error: err.Error()})
		}
	}
}

func (g *GRPC) CreateIndex(ctx context.Context, in *pb.CreateIndexRequest) (*pb.Empty, error) {
	if err := gongt.CreateIndex(int(in.PoolSize)); err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (g *GRPC) SaveIndex(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	if err := gongt.SaveIndex(); err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (g *GRPC) GetDimension(ctx context.Context, in *pb.Empty) (*pb.GetDimensionResponse, error) {
	dim := gongt.GetDim()
	return &pb.GetDimensionResponse{Dimension: int32(dim)}, nil
}

func toSearchResponse(s []service.SearchResult) *pb.SearchResponse {
	ret := make([]*pb.ObjectDistance, len(s))
	for i, r := range s {
		if r.Error == nil {
			ret[i] = &pb.ObjectDistance{
				Id:       r.Id,
				Distance: r.Distance,
			}
		} else {
			ret[i] = &pb.ObjectDistance{
				Error: r.Error.Error(),
			}
		}
	}
	return &pb.SearchResponse{
		Result: ret,
	}
}
