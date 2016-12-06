#!/usr/bin/env ruby

require 'yaml'
require 'active_support'


release_names=[
  'capi',
  'consul',
  'diego',
  'etcd',
  'loggregator',
  'nats',
  'cf-mysql',
  'cf-routing',
  'uaa',
  'garden-runc',
  'cflinuxfs2-rootfs',
  'binary-buildpack',
  'dotnet-core-buildpack',
  'go-buildpack',
  'java-buildpack',
  'java-offline-buildpack',
  'nodejs-buildpack',
  'php-buildpack',
  'python-buildpack',
  'ruby-buildpack',
  'staticfile-buildpack',
]

deployment_configuration_path = ENV.fetch('DEPLOYMENT_CONFIGURATION_PATH')
deployment_manifest_path = ENV.fetch("DEPLOYMENT_MANIFEST_PATH")

deployment_manifest = File.read("deployment-configuration/#{deployment_configuration_path}")
properties_section = deployment_manifest.split(/^releases:$/).first

release_array = release_names.map do |release_name|
  release_resource = "#{release_name}-release"

  url = File.read("#{release_resource}/url").strip
  version = File.read("#{release_resource}/version").strip
  sha1 = File.read("#{release_resource}/sha1").strip

  {
    'name' => release_name,
    'url' => url,
    'version' => version,
    'sha1' => sha1
  }
end

stemcell_version = File.read("stemcell/version").strip

bosh_artifacts_section = {
  'releases' => release_array,
  'stemcells' => [
    {
      'alias' => "default",
      'os' => "ubuntu-trusty",
      'version' => stemcell_version
    }
  ]
}.to_yaml[4..-1]


File.open("deployment-manifest/#{deployment_manifest_path}", 'w') do |file|
  file.write([
    properties_section,
    bosh_artifacts_section
  ].join("\n"))
end
