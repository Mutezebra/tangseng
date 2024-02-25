// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package rpc

import (
	"context"
	"fmt"
	"github.com/CocaineCong/tangseng/pkg/prometheus"
	"log"
	"time"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"

	"github.com/CocaineCong/tangseng/config"
	"github.com/CocaineCong/tangseng/consts"
	"github.com/CocaineCong/tangseng/idl/pb/favorite"
	"github.com/CocaineCong/tangseng/idl/pb/index_platform"
	"github.com/CocaineCong/tangseng/idl/pb/search_engine"
	"github.com/CocaineCong/tangseng/idl/pb/search_vector"
	"github.com/CocaineCong/tangseng/idl/pb/user"
	"github.com/CocaineCong/tangseng/pkg/discovery"
)

var (
	Register   *discovery.Resolver
	ctx        context.Context
	CancelFunc context.CancelFunc

	UserClient          user.UserServiceClient
	FavoriteClient      favorite.FavoritesServiceClient
	SearchEngineClient  search_engine.SearchEngineServiceClient
	IndexPlatformClient index_platform.IndexPlatformServiceClient
	SearchVectorClient  search_vector.SearchVectorServiceClient
)

// Init 初始化所有的rpc请求
func Init() {
	Register = discovery.NewResolver([]string{config.Conf.Etcd.Address}, logrus.New())
	resolver.Register(Register)
	ctx, CancelFunc = context.WithTimeout(context.Background(), 3*time.Second)

	defer Register.Close()
	initClient(config.Conf.Domain[consts.UserServiceName].Name, &UserClient)
	initClient(config.Conf.Domain[consts.FavoriteServiceName].Name, &FavoriteClient)
	initClient(config.Conf.Domain[consts.SearchServiceName].Name, &SearchEngineClient)
	initClient(config.Conf.Domain[consts.IndexPlatformName].Name, &IndexPlatformClient)
	initClient(config.Conf.Domain[consts.SearchVectorName].Name, &SearchVectorClient)
}

// initClient 初始化所有的rpc客户端
func initClient(serviceName string, client interface{}) {
	conn, err := connectServer(serviceName)

	if err != nil {
		panic(err)
	}

	switch c := client.(type) {
	case *user.UserServiceClient:
		*c = user.NewUserServiceClient(conn)
	case *favorite.FavoritesServiceClient:
		*c = favorite.NewFavoritesServiceClient(conn)
	case *search_engine.SearchEngineServiceClient:
		*c = search_engine.NewSearchEngineServiceClient(conn)
	case *index_platform.IndexPlatformServiceClient:
		*c = index_platform.NewIndexPlatformServiceClient(conn)
	case *search_vector.SearchVectorServiceClient:
		*c = search_vector.NewSearchVectorServiceClient(conn)
	default:
		panic("unsupported worker type")
	}
}

func connectServer(serviceName string) (conn *grpc.ClientConn, err error) {
	prometheus.EnableHandlingTimeHistogram()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	addr := fmt.Sprintf("%s:///%s", Register.Scheme(), serviceName)

	// Load balance
	if config.Conf.Services[serviceName].LoadBalance {
		log.Printf("load balance enabled for %s\n", serviceName)
		opts = append(opts, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")))
	}

	conn, err = grpc.DialContext(ctx, addr, opts...)
	err = errors.Wrapf(err, "failed to connect to gRPC service,address is %v", addr)
	return
}
