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

package main

import (
	"context"
	"fmt"
	"github.com/CocaineCong/tangseng/pkg/prometheus"
	"net"

	logs "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/CocaineCong/tangseng/app/search_engine/analyzer"
	"github.com/CocaineCong/tangseng/app/search_engine/repository/storage"
	"github.com/CocaineCong/tangseng/app/search_engine/rpc"
	"github.com/CocaineCong/tangseng/app/search_engine/service"
	"github.com/CocaineCong/tangseng/config"
	"github.com/CocaineCong/tangseng/consts"
	pb "github.com/CocaineCong/tangseng/idl/pb/search_engine"
	"github.com/CocaineCong/tangseng/loading"
	"github.com/CocaineCong/tangseng/pkg/discovery"
)

func main() {
	ctx := context.Background()
	loading.Loading()
	// bi_dao.InitDB() // TODO starrocks完善才开启
	analyzer.InitSeg()
	storage.InitStorageDB(ctx)
	rpc.Init()

	// etcd 地址
	etcdAddress := []string{config.Conf.Etcd.Address}
	// 服务注册
	etcdRegister := discovery.NewRegister(etcdAddress, logrus.New())
	grpcAddress := config.Conf.Services[consts.SearchServiceName].Addr[0]
	defer etcdRegister.Stop()
	node := discovery.Server{
		Name: config.Conf.Domain[consts.SearchServiceName].Name,
		Addr: grpcAddress,
	}
	server := grpc.NewServer(
		grpc.UnaryInterceptor(prometheus.UnaryServerInterceptor),
		grpc.StreamInterceptor(prometheus.StreamServerInterceptor),
	)
	defer server.Stop()
	// 绑定service
	pb.RegisterSearchEngineServiceServer(server, service.GetSearchEngineSrv())
	lis, err := net.Listen("tcp", grpcAddress)
	fmt.Println("ready to register")
	prometheus.RegisterServer(server, config.Conf.Services[consts.SearchServiceName].AddrMetrics[0], consts.SearchServiceName)
	if err != nil {
		panic(err)
	}
	if _, err := etcdRegister.Register(node, 10); err != nil {
		logs.LogrusObj.Errorf("start service failed, original error: %T %v", errors.Cause(err), errors.Cause(err))
		logs.LogrusObj.Panicf("stack trace: \n%+v\n", err)
	}
	logrus.Info("service started listen on ", grpcAddress)
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}
