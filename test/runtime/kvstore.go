// Copyright 2017 Authors of Cilium
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

package RuntimeTest

import (
	"context"

	"github.com/cilium/cilium/test/helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("RuntimeKVStoreTest", func() {

	var initialized bool
	var logger *logrus.Entry
	var vm *helpers.SSHMeta

	initialize := func() {
		if initialized == true {
			return
		}
		logger = log.WithFields(logrus.Fields{"testName": "RuntimeKVStoreTest"})
		logger.Info("Starting")
		vm = helpers.CreateNewRuntimeHelper(helpers.Runtime, logger)
		logger.Info("done creating Cilium and Docker helpers")
		initialized = true
	}
	containers := func(option string) {
		switch option {
		case helpers.Create:
			vm.NetworkCreate(helpers.CiliumDockerNetwork, "")
			vm.ContainerCreate(helpers.Client, helpers.NetperfImage, helpers.CiliumDockerNetwork, "-l id.client")
		case helpers.Delete:
			vm.ContainerRm(helpers.Client)

		}
	}

	BeforeEach(func() {
		initialize()
		By("Stopping cilium service")
		vm.Exec("sudo systemctl stop cilium")
	}, 150)

	AfterEach(func() {
		if CurrentGinkgoTestDescription().Failed {
			vm.ReportFailed(
				"sudo docker ps -a",
				"sudo cilium endpoint list")
		}
		containers(helpers.Delete)
		By("Starting cilium service")
		vm.Exec("sudo systemctl start cilium")
	})

	It("Consul KVStore", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		By("Starting Cilium using Consul as backing key-value store")
		vm.ExecContext(
			ctx,
			"sudo cilium-agent --kvstore consul --kvstore-opt consul.address=127.0.0.1:8500 --debug")
		err := vm.WaitUntilReady(150)
		Expect(err).Should(BeNil())

		vm.Exec("sudo systemctl restart cilium-docker")
		helpers.Sleep(2)
		containers(helpers.Create)
		vm.WaitEndpointsReady()
		eps, err := vm.GetEndpointsNames()
		Expect(err).Should(BeNil())
		Expect(len(eps)).To(Equal(1))
	})

	It("Etcd KVStore", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		By("Starting Cilium using etcd as backing key-value store")
		vm.ExecContext(
			ctx,
			"sudo cilium-agent --kvstore etcd --kvstore-opt etcd.address=127.0.0.1:9732")
		err := vm.WaitUntilReady(150)
		Expect(err).Should(BeNil())

		By("Restarting cilium-docker service")
		vm.Exec("sudo systemctl restart cilium-docker")
		helpers.Sleep(2)
		containers(helpers.Create)

		vm.WaitEndpointsReady()

		eps, err := vm.GetEndpointsNames()
		Expect(err).Should(BeNil())
		Expect(len(eps)).To(Equal(1))
	})
})
