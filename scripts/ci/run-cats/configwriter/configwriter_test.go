package configwriter_test

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/cloudfoundry/runtime-ci/scripts/ci/run-cats/configwriter"
	"github.com/cloudfoundry/runtime-ci/scripts/ci/run-cats/configwriter/configwriterfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configwriter", func() {
	var env *configwriterfakes.FakeEnvironment
	BeforeEach(func() {
		env = &configwriterfakes.FakeEnvironment{}
		env.GetIncludeRoutingReturns(true, nil)
		env.GetIncludeDetectReturns(true, nil)
		env.GetIncludeAppsReturns(true, nil)
	})

	Context("when env vars are not set", func() {
		It("returns an empty config object", func() {
			configFile, err := configwriter.NewConfigFile("/dir/name", env)

			Expect(err).NotTo(HaveOccurred())
			Expect(configFile).NotTo(BeNil())
			Expect(configFile.Config.Api).To(Equal(""))
			Expect(configFile.Config.AdminUser).To(Equal(""))
			Expect(configFile.Config.AdminPassword).To(Equal(""))
			Expect(configFile.Config.AppsDomain).To(Equal(""))
			Expect(configFile.Config.SkipSslValidation).To(BeFalse())
			Expect(configFile.Config.IncludeSSO).To(BeFalse())
			Expect(configFile.Config.UseHttp).To(BeFalse())
			Expect(configFile.Config.ExistingUser).To(Equal(""))
			Expect(configFile.Config.ExistingUserPassword).To(Equal(""))
			Expect(configFile.Config.Backend).To(Equal(""))
			Expect(configFile.Config.PersistentAppHost).To(Equal(""))
			Expect(configFile.Config.PersistentAppSpace).To(Equal(""))
			Expect(configFile.Config.PersistentAppOrg).To(Equal(""))
			Expect(configFile.Config.PersistentAppQuotaName).To(Equal(""))
			Expect(configFile.Config.StaticBuildpackName).To(Equal(""))
			Expect(configFile.Config.JavaBuildpackName).To(Equal(""))
			Expect(configFile.Config.RubyBuildpackName).To(Equal(""))
			Expect(configFile.Config.NodeJsBuildpackName).To(Equal(""))
			Expect(configFile.Config.GoBuildpackName).To(Equal(""))
			Expect(configFile.Config.PythonBuildpackName).To(Equal(""))
			Expect(configFile.Config.PhpBuildpackName).To(Equal(""))
			Expect(configFile.Config.BinaryBuildpackName).To(Equal(""))
			Expect(configFile.Config.IncludeApps).To(BeTrue())
			Expect(configFile.Config.IncludeDiegoSSH).To(BeFalse())
			Expect(configFile.Config.IncludeV3).To(BeFalse())
			Expect(configFile.Config.IncludeDiegoDocker).To(BeFalse())
			Expect(configFile.Config.IncludeSecurityGroups).To(BeFalse())
			Expect(configFile.Config.IncludeBackendCompatibility).To(BeFalse())
			Expect(configFile.Config.IncludeInternetDependent).To(BeFalse())
			Expect(configFile.Config.IncludeServices).To(BeFalse())
			Expect(configFile.Config.IncludeRouteServices).To(BeFalse())
			Expect(configFile.Config.IncludeRouting).To(BeTrue())
			Expect(configFile.Config.IncludeTasks).To(BeFalse())
			Expect(configFile.Config.IncludeDetect).To(BeTrue())
			Expect(configFile.DestinationDir).To(Equal("/dir/name"))
		})

		It("Sets 'KeepUserAtSuiteEnd' and 'UseExistingUser' to false", func() {
			configFile, err := configwriter.NewConfigFile("", env)
			Expect(err).NotTo(HaveOccurred())
			Expect(configFile).NotTo(BeNil())
			Expect(configFile.Config.UseExistingUser).To(Equal(false))
			Expect(configFile.Config.KeepUserAtSuiteEnd).To(Equal(false))
		})
	})

	Context("When envvars are set", func() {
		var (
			expectedApi                               string
			expectedAdminUser                         string
			expectedPassword                          string
			expectedAppsDomain                        string
			expectedSkipSslValidation                 bool
			expectedIncludeSSO                        bool
			expectedUseHttp                           bool
			expectedIncludePrivilegedContainerSupport bool
			expectedExistingUser                      string
			expectedKeepUserAtSuiteEnd                bool
			expectedUseExistingUser                   bool
			expectedExistingUserPassword              string
			expectedBackend                           string
			expectedPersistedAppHost                  string
			expectedPersistedAppSpace                 string
			expectedPersistedAppOrg                   string
			expectedPersistedAppQuotaName             string
			expectedDefaultTimeout                    int
			expectedCfPushTimeout                     int
			expectedLongCurlTimeout                   int
			expectedBrokerStartTimeout                int
			expectedAsyncServiceOperationTimeout      int
			expectedStaticBuildpackName               string
			expectedJavaBuildpackName                 string
			expectedRubyBuildpackName                 string
			expectedNodeJsBuildpackName               string
			expectedGoBuildpackName                   string
			expectedPythonBuildpackName               string
			expectedPhpBuildpackName                  string
			expectedBinaryBuildpackName               string
			expectedIncludeApps                       bool
			expectedIncludeDiegoSSH                   bool
			expectedIncludeV3                         bool
			expectedIncludeDiegoDocker                bool
			expectedIncludeSecurityGroups             bool
			expectedIncludeBackendCompatibility       bool
			expectedIncludeInternetDependent          bool
			expectedIncludeServices                   bool
			expectedIncludeRouteServices              bool
			expectedIncludeRouting                    bool
			expectedIncludeTasks                      bool
			expectedIncludeDetect                     bool
		)

		BeforeEach(func() {
			expectedApi = "api.example.com"
			expectedAdminUser = "admin_user"
			expectedPassword = "admin_password"
			expectedAppsDomain = "apps_domain"
			expectedSkipSslValidation = true
			expectedIncludeSSO = true
			expectedUseHttp = true
			expectedIncludePrivilegedContainerSupport = true
			expectedExistingUser = "existing_user"
			expectedKeepUserAtSuiteEnd = true
			expectedUseExistingUser = true
			expectedExistingUserPassword = "expected_existing_user_password"
			expectedBackend = "diego"
			expectedPersistedAppHost = "PERSISTENT_APP_HOST"
			expectedPersistedAppSpace = "PERSISTENT_APP_SPACE"
			expectedPersistedAppOrg = "PERSISTENT_APP_ORG"
			expectedPersistedAppQuotaName = "PERSISTENT_APP_QUOTA_NAME"

			expectedDefaultTimeout = 1
			expectedCfPushTimeout = 2
			expectedLongCurlTimeout = 3
			expectedBrokerStartTimeout = 4
			expectedAsyncServiceOperationTimeout = 5

			expectedStaticBuildpackName = "STATICFILE_BUILDPACK_NAME"
			expectedJavaBuildpackName = "JAVA_BUILDPACK_NAME"
			expectedRubyBuildpackName = "Ruby_BUILDPACK_NAME"
			expectedNodeJsBuildpackName = "NODEJS_BUILDPACK_NAME"
			expectedGoBuildpackName = "GO_BUILDPACK_NAME"
			expectedPythonBuildpackName = "PYTHON_BUILDPACK_NAME"
			expectedPhpBuildpackName = "PHP_BUILDPACK_NAME"
			expectedBinaryBuildpackName = "BINARY_BUILDPACK_NAME"

			expectedIncludeApps = false
			expectedIncludeDiegoSSH = true
			expectedIncludeV3 = true
			expectedIncludeDiegoDocker = true
			expectedIncludeSecurityGroups = true
			expectedIncludeBackendCompatibility = true
			expectedIncludeInternetDependent = true
			expectedIncludeServices = true
			expectedIncludeRouteServices = true
			expectedIncludeRouting = false
			expectedIncludeTasks = true
			expectedIncludeDetect = false

			env.GetSkipSSLValidationReturns(expectedSkipSslValidation, nil)
			env.GetIncludeSSOReturns(expectedIncludeSSO, nil)
			env.GetUseHTTPReturns(expectedUseHttp, nil)
			env.GetIncludePrivilegedContainerSupportReturns(expectedIncludePrivilegedContainerSupport, nil)

			env.GetCFAPIReturns(expectedApi)
			env.GetCFAdminUserReturns(expectedAdminUser)
			env.GetCFAdminPasswordReturns(expectedPassword)
			env.GetCFAppsDomainReturns(expectedAppsDomain)
			env.GetExistingUserReturns(expectedExistingUser)
			env.GetExistingUserPasswordReturns(expectedExistingUserPassword)
			env.UseExistingUserReturns(expectedUseExistingUser)
			env.KeepUserAtSuiteEndReturns(expectedKeepUserAtSuiteEnd)
			env.GetPersistentAppHostReturns(expectedPersistedAppHost)
			env.GetPersistentAppSpaceReturns(expectedPersistedAppSpace)
			env.GetPersistentAppOrgReturns(expectedPersistedAppOrg)
			env.GetPersistentAppQuotaNameReturns(expectedPersistedAppQuotaName)
			env.GetStaticBuildpackNameReturns(expectedStaticBuildpackName)
			env.GetJavaBuildpackNameReturns(expectedJavaBuildpackName)
			env.GetRubyBuildpackNameReturns(expectedRubyBuildpackName)
			env.GetNodeJSBuildpackNameReturns(expectedNodeJsBuildpackName)
			env.GetGoBuildpackNameReturns(expectedGoBuildpackName)
			env.GetPythonBuildpackNameReturns(expectedPythonBuildpackName)
			env.GetPHPBuildpackNameReturns(expectedPhpBuildpackName)
			env.GetBinaryBuildpackNameReturns(expectedBinaryBuildpackName)

			env.GetDefaultTimeoutInSecondsReturns(expectedDefaultTimeout, nil)
			env.GetCFPushTimeoutInSecondsReturns(expectedCfPushTimeout, nil)
			env.GetLongCurlTimeoutInSecondsReturns(expectedLongCurlTimeout, nil)
			env.GetBrokerStartTimeoutInSecondsReturns(expectedBrokerStartTimeout, nil)
			env.GetAsyncServiceOperationTimeoutInSecondsReturns(expectedAsyncServiceOperationTimeout, nil)

			env.GetIncludeAppsReturns(expectedIncludeApps, nil)
			env.GetIncludeDiegoSSHReturns(expectedIncludeDiegoSSH, nil)
			env.GetIncludeV3Returns(expectedIncludeV3, nil)
			env.GetIncludeDiegoDockerReturns(expectedIncludeDiegoDocker, nil)
			env.GetIncludeSecurityGroupsReturns(expectedIncludeSecurityGroups, nil)
			env.GetIncludeBackendCompatibilityReturns(expectedIncludeBackendCompatibility, nil)
			env.GetIncludeInternetDependentReturns(expectedIncludeInternetDependent, nil)
			env.GetIncludeServicesReturns(expectedIncludeServices, nil)
			env.GetIncludeRouteServicesReturns(expectedIncludeRouteServices, nil)
			env.GetIncludeRoutingReturns(expectedIncludeRouting, nil)
			env.GetIncludeTasksReturns(expectedIncludeTasks, nil)
			env.GetIncludeDetectReturns(expectedIncludeDetect, nil)

			env.GetBackendReturns(expectedBackend, nil)
		})

		It("Generates a config object with the correct CF env variables set", func() {
			configFile, err := configwriter.NewConfigFile("/some/dir", env)
			Expect(err).NotTo(HaveOccurred())
			Expect(configFile).NotTo(BeNil())
			Expect(configFile.Config.Api).To(Equal(expectedApi))
			Expect(configFile.Config.AdminUser).To(Equal(expectedAdminUser))
			Expect(configFile.Config.AdminPassword).To(Equal(expectedPassword))
			Expect(configFile.Config.AppsDomain).To(Equal(expectedAppsDomain))
			Expect(configFile.Config.SkipSslValidation).To(Equal(expectedSkipSslValidation))
			Expect(configFile.Config.IncludeSSO).To(Equal(expectedIncludeSSO))
			Expect(configFile.Config.UseHttp).To(Equal(expectedUseHttp))
			Expect(configFile.Config.IncludePrivilegedContainerSupport).To(Equal(expectedIncludePrivilegedContainerSupport))
			Expect(configFile.Config.ExistingUser).To(Equal(expectedExistingUser))
			Expect(configFile.Config.ExistingUserPassword).To(Equal(expectedExistingUserPassword))
			Expect(configFile.Config.Backend).To(Equal(expectedBackend))
			Expect(configFile.Config.PersistentAppHost).To(Equal(expectedPersistedAppHost))
			Expect(configFile.Config.PersistentAppSpace).To(Equal(expectedPersistedAppSpace))
			Expect(configFile.Config.PersistentAppOrg).To(Equal(expectedPersistedAppOrg))
			Expect(configFile.Config.PersistentAppQuotaName).To(Equal(expectedPersistedAppQuotaName))
			Expect(configFile.Config.DefaultTimeout).To(Equal(expectedDefaultTimeout))
			Expect(configFile.Config.CfPushTimeout).To(Equal(expectedCfPushTimeout))
			Expect(configFile.Config.LongCurlTimeout).To(Equal(expectedLongCurlTimeout))
			Expect(configFile.Config.BrokerStartTimeout).To(Equal(expectedBrokerStartTimeout))
			Expect(configFile.Config.AsyncServiceOperationTimeout).To(Equal(expectedAsyncServiceOperationTimeout))
			Expect(configFile.Config.StaticBuildpackName).To(Equal(expectedStaticBuildpackName))
			Expect(configFile.Config.JavaBuildpackName).To(Equal(expectedJavaBuildpackName))
			Expect(configFile.Config.RubyBuildpackName).To(Equal(expectedRubyBuildpackName))
			Expect(configFile.Config.NodeJsBuildpackName).To(Equal(expectedNodeJsBuildpackName))
			Expect(configFile.Config.GoBuildpackName).To(Equal(expectedGoBuildpackName))
			Expect(configFile.Config.PythonBuildpackName).To(Equal(expectedPythonBuildpackName))
			Expect(configFile.Config.PhpBuildpackName).To(Equal(expectedPhpBuildpackName))
			Expect(configFile.Config.BinaryBuildpackName).To(Equal(expectedBinaryBuildpackName))
			Expect(configFile.Config.IncludeApps).To(Equal(expectedIncludeApps))
			Expect(configFile.Config.IncludeDiegoSSH).To(Equal(expectedIncludeDiegoSSH))
			Expect(configFile.Config.IncludeV3).To(Equal(expectedIncludeV3))
			Expect(configFile.Config.IncludeDiegoDocker).To(Equal(expectedIncludeDiegoDocker))
			Expect(configFile.Config.IncludeSecurityGroups).To(Equal(expectedIncludeSecurityGroups))
			Expect(configFile.Config.IncludeBackendCompatibility).To(Equal(expectedIncludeBackendCompatibility))
			Expect(configFile.Config.IncludeInternetDependent).To(Equal(expectedIncludeInternetDependent))
			Expect(configFile.Config.IncludeServices).To(Equal(expectedIncludeServices))
			Expect(configFile.Config.IncludeRouteServices).To(Equal(expectedIncludeRouteServices))
			Expect(configFile.Config.IncludeRouting).To(Equal(expectedIncludeRouting))
			Expect(configFile.Config.IncludeTasks).To(Equal(expectedIncludeTasks))
			Expect(configFile.Config.IncludeDetect).To(Equal(expectedIncludeDetect))
			Expect(configFile.DestinationDir).To(Equal("/some/dir"))
		})

		It("Sets 'KeepUserAtSuiteEnd' and 'UseExistingUser' to true", func() {
			configFile, err := configwriter.NewConfigFile("/some/dir", env)
			Expect(err).NotTo(HaveOccurred())
			Expect(configFile.Config.UseExistingUser).To(Equal(true))
			Expect(configFile.Config.KeepUserAtSuiteEnd).To(Equal(true))
		})
	})

	Describe("marshaling the struct", func() {
		It("does not render optional keys if their values are empty", func() {
			var configJson []byte
			configFile, err := configwriter.NewConfigFile("", env)
			Expect(err).NotTo(HaveOccurred())
			configJson, err = json.Marshal(configFile.Config)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(configJson)).To(MatchJSON(`{
                                              "api": "",
                                              "admin_user": "",
                                              "admin_password": "",
                                              "apps_domain": "",
                                              "skip_ssl_validation": false,
                                              "include_sso": false,
                                              "use_http": false,
                                              "use_existing_user": false,
                                              "keep_user_at_suite_end": false,
                                              "include_ssh": false,
                                              "include_v3": false,
                                              "include_docker": false,
                                              "include_security_groups": false,
                                              "include_backend_compatibility": false,
                                              "include_internet_dependent": false,
                                              "include_services": false,
                                              "include_route_services": false,
                                              "include_apps": true,
                                              "include_routing": true,
                                              "include_tasks": false,
                                              "include_detect": true
                                              }`))
		})

		Context("when any env variables are provided", func() {
			BeforeEach(func() {
				env.GetSkipSSLValidationReturns(true, nil)
				env.GetIncludeSSOReturns(true, nil)
				env.GetUseHTTPReturns(true, nil)

				env.GetBackendReturns("diego", nil)

				env.GetCFAPIReturns("non-empty-value")
				env.GetCFAdminUserReturns("non-empty-value")
				env.GetCFAdminPasswordReturns("non-empty-value")
				env.GetCFAppsDomainReturns("non-empty-value")
				env.GetExistingUserReturns("non-empty-value")
				env.UseExistingUserReturns(true)
				env.KeepUserAtSuiteEndReturns(true)
				env.GetExistingUserPasswordReturns("non-empty-value")
				env.GetPersistentAppHostReturns("non-empty-value")
				env.GetPersistentAppSpaceReturns("non-empty-value")
				env.GetPersistentAppOrgReturns("non-empty-value")
				env.GetPersistentAppQuotaNameReturns("non-empty-value")
				env.GetDefaultTimeoutInSecondsReturns(1, nil)
				env.GetCFPushTimeoutInSecondsReturns(1, nil)
				env.GetLongCurlTimeoutInSecondsReturns(1, nil)
				env.GetBrokerStartTimeoutInSecondsReturns(1, nil)
				env.GetAsyncServiceOperationTimeoutInSecondsReturns(1, nil)
				env.GetStaticBuildpackNameReturns("non-empty-value")
				env.GetJavaBuildpackNameReturns("non-empty-value")
				env.GetRubyBuildpackNameReturns("non-empty-value")
				env.GetNodeJSBuildpackNameReturns("non-empty-value")
				env.GetGoBuildpackNameReturns("non-empty-value")
				env.GetPythonBuildpackNameReturns("non-empty-value")
				env.GetPHPBuildpackNameReturns("non-empty-value")
				env.GetBinaryBuildpackNameReturns("non-empty-value")
				env.GetIncludeAppsReturns(false, nil)
				env.GetIncludeDiegoSSHReturns(true, nil)
				env.GetIncludeV3Returns(true, nil)
				env.GetIncludeDiegoDockerReturns(true, nil)
				env.GetIncludeSecurityGroupsReturns(true, nil)
				env.GetIncludeBackendCompatibilityReturns(true, nil)
				env.GetIncludeInternetDependentReturns(true, nil)
				env.GetIncludeServicesReturns(true, nil)
				env.GetIncludeRouteServicesReturns(true, nil)
				env.GetIncludeRoutingReturns(false, nil)
				env.GetIncludeTasksReturns(true, nil)
				env.GetIncludeDetectReturns(false, nil)

			})

			It("renders the variables in the integration_config", func() {
				var configJson []byte
				configFile, err := configwriter.NewConfigFile("", env)
				Expect(err).NotTo(HaveOccurred())
				configJson, err = json.Marshal(configFile.Config)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(configJson)).To(MatchJSON(`{
                                              "api": "non-empty-value",
                                              "admin_user": "non-empty-value",
                                              "admin_password": "non-empty-value",
                                              "apps_domain": "non-empty-value",
                                              "skip_ssl_validation": true,
                                              "include_sso": true,
                                              "use_http": true,
                                              "existing_user": "non-empty-value",
                                              "use_existing_user": true,
                                              "keep_user_at_suite_end": true,
                                              "existing_user_password": "non-empty-value",
                                              "backend": "diego",
                                              "staticfile_buildpack_name": "non-empty-value",
                                              "java_buildpack_name": "non-empty-value",
                                              "ruby_buildpack_name": "non-empty-value",
                                              "nodejs_buildpack_name": "non-empty-value",
                                              "go_buildpack_name": "non-empty-value",
                                              "python_buildpack_name": "non-empty-value",
                                              "php_buildpack_name": "non-empty-value",
                                              "binary_buildpack_name": "non-empty-value",
                                              "persistent_app_host": "non-empty-value",
                                              "persistent_app_space": "non-empty-value",
                                              "persistent_app_org": "non-empty-value",
                                              "persistent_app_quota_name": "non-empty-value",
                                              "default_timeout": 1,
                                              "cf_push_timeout": 1,
                                              "long_curl_timeout": 1,
                                              "broker_start_timeout": 1,
                                              "async_service_operation_timeout": 1,
                                              "include_ssh": true,
                                              "include_v3": true,
                                              "include_docker": true,
                                              "include_security_groups": true,
                                              "include_backend_compatibility": true,
                                              "include_internet_dependent": true,
                                              "include_services": true,
                                              "include_route_services": true,
                                              "include_apps": false,
                                              "include_routing": false,
                                              "include_tasks": true,
                                              "include_detect": false
                                              }`))

			})
		})

		Context("when only required env variables are provided", func() {
			BeforeEach(func() {
				env.GetSkipSSLValidationReturns(false, nil)
				env.GetIncludeSSOReturns(false, nil)
				env.GetUseHTTPReturns(false, nil)

				env.GetBackendReturns("", nil)

				env.GetCFAPIReturns("non-empty-value")
				env.GetCFAdminUserReturns("non-empty-value")
				env.GetCFAdminPasswordReturns("non-empty-value")
				env.GetCFAppsDomainReturns("non-empty-value")
				env.GetExistingUserReturns("")
				env.GetExistingUserPasswordReturns("")
				env.GetPersistentAppHostReturns("")
				env.GetPersistentAppSpaceReturns("")
				env.GetPersistentAppOrgReturns("")
				env.GetPersistentAppQuotaNameReturns("")
				env.GetDefaultTimeoutInSecondsReturns(0, nil)
				env.GetCFPushTimeoutInSecondsReturns(0, nil)
				env.GetLongCurlTimeoutInSecondsReturns(0, nil)
				env.GetBrokerStartTimeoutInSecondsReturns(0, nil)
				env.GetAsyncServiceOperationTimeoutInSecondsReturns(0, nil)
				env.GetStaticBuildpackNameReturns("")
				env.GetJavaBuildpackNameReturns("")
				env.GetRubyBuildpackNameReturns("")
				env.GetNodeJSBuildpackNameReturns("")
				env.GetGoBuildpackNameReturns("")
				env.GetPythonBuildpackNameReturns("")
				env.GetPHPBuildpackNameReturns("")
				env.GetBinaryBuildpackNameReturns("")
				env.GetIncludeDiegoSSHReturns(false, nil)
				env.GetIncludeV3Returns(false, nil)
				env.GetIncludeDiegoDockerReturns(false, nil)
				env.GetIncludeSecurityGroupsReturns(false, nil)
				env.GetIncludeBackendCompatibilityReturns(false, nil)
				env.GetIncludeInternetDependentReturns(false, nil)
				env.GetIncludeServicesReturns(false, nil)
				env.GetIncludeRouteServicesReturns(false, nil)
				env.GetIncludeAppsReturns(true, nil)
				env.GetIncludeRoutingReturns(true, nil)
				env.GetIncludeTasksReturns(false, nil)
				env.GetIncludeDetectReturns(true, nil)
			})

			It("renders the default variables in the integration_config", func() {
				var configJson []byte
				configFile, err := configwriter.NewConfigFile("", env)
				Expect(err).NotTo(HaveOccurred())
				configJson, err = json.Marshal(configFile.Config)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(configJson)).To(MatchJSON(`{
                                              "api": "non-empty-value",
                                              "admin_user": "non-empty-value",
                                              "admin_password": "non-empty-value",
                                              "apps_domain": "non-empty-value",
                                              "skip_ssl_validation": false,
                                              "include_sso": false,
                                              "use_http": false,
                                              "use_existing_user": false,
                                              "keep_user_at_suite_end": false,
                                              "include_ssh": false,
                                              "include_v3": false,
                                              "include_docker": false,
                                              "include_security_groups": false,
                                              "include_backend_compatibility": false,
                                              "include_internet_dependent": false,
                                              "include_services": false,
                                              "include_route_services": false,
                                              "include_apps": true,
                                              "include_routing": true,
                                              "include_tasks": false,
                                              "include_detect": true
                                              }`))
			})
		})

	})

	Describe("writing the integration_config.json file", func() {
		var tempDir string
		var err error

		BeforeEach(func() {
			tempDir, err = ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())
			env.GetCFAPIReturns("cf-api-value")
		})

		AfterEach(func() {
			err := os.RemoveAll(tempDir)
			Expect(err).NotTo(HaveOccurred())
		})

		It("writes the marshalled config to the file", func() {
			configFile, err := configwriter.NewConfigFile(tempDir, env)
			Expect(err).NotTo(HaveOccurred())

			var file *os.File
			file, err = configFile.WriteConfigToFile()
			Expect(err).NotTo(HaveOccurred())

			Expect(tempDir + "/integration_config.json").To(BeARegularFile())
			Expect(tempDir + "/integration_config.json").To(Equal(file.Name()))

			contents, err := ioutil.ReadFile(file.Name())
			Expect(err).NotTo(HaveOccurred())
			Expect(contents).To(MatchJSON(`{
                                              "api": "cf-api-value",
                                              "admin_user": "",
                                              "admin_password": "",
                                              "apps_domain": "",
                                              "skip_ssl_validation": false,
                                              "include_sso": false,
                                              "use_http": false,
                                              "use_existing_user": false,
                                              "keep_user_at_suite_end": false,
                                              "include_ssh": false,
                                              "include_v3": false,
                                              "include_docker": false,
                                              "include_security_groups": false,
                                              "include_backend_compatibility": false,
                                              "include_internet_dependent": false,
                                              "include_services": false,
                                              "include_route_services": false,
                                              "include_apps": true,
                                              "include_routing": true,
                                              "include_tasks": false,
                                              "include_detect": true
                                              }`))
		})
	})

	Context("when the destinationDir doesn't end with a trailing slash", func() {
		BeforeEach(func() {
			env.GetCFAPIReturns("cf-api-value")
		})

		It("should successfully write integration_config.json", func() {
			configFile, err := configwriter.NewConfigFile("/tmp", env)
			Expect(err).NotTo(HaveOccurred())

			_, err = configFile.WriteConfigToFile()
			Expect(err).NotTo(HaveOccurred())

			var file *os.File
			file, err = os.Open("/tmp/integration_config.json")
			Expect(err).NotTo(HaveOccurred())
			contents, err := ioutil.ReadFile(file.Name())

			Expect(contents).To(MatchJSON(`{
                                    "api": "cf-api-value",
                                    "admin_user": "",
                                    "admin_password": "",
                                    "apps_domain": "",
                                    "skip_ssl_validation": false,
                                    "include_sso": false,
                                    "use_http": false,
                                    "use_existing_user": false,
                                    "keep_user_at_suite_end": false,
                                    "include_ssh": false,
                                    "include_v3": false,
                                    "include_docker": false,
                                    "include_security_groups": false,
                                    "include_backend_compatibility": false,
                                    "include_internet_dependent": false,
                                    "include_services": false,
                                    "include_route_services": false,
                                    "include_apps": true,
                                    "include_routing": true,
                                    "include_tasks": false,
                                    "include_detect": true
                                    }`))
		})
	})

	Describe("exporting the config environment variable", func() {
		Context("when no env vars are set", func() {
			var tempDir string
			var err error

			BeforeEach(func() {
				tempDir, err = ioutil.TempDir("", "")
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				err := os.RemoveAll(tempDir)
				Expect(err).NotTo(HaveOccurred())

				os.Unsetenv("CONFIG")
			})

			It("exports the location of the integration_config.json file", func() {
				configFile, err := configwriter.NewConfigFile("/some/path", env)
				Expect(err).NotTo(HaveOccurred())

				configFile.ExportConfigFilePath()

				Expect(os.Getenv("CONFIG")).To(Equal("/some/path/integration_config.json"))
			})
		})
	})
})
