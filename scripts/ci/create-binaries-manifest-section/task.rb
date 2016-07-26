#!/usr/bin/env ruby

require 'yaml'
require 'active_support'


release_names=[
  'consul',
  'diego',
  'etcd',
  'loggregator',
  'cf-mysql',
  'uaa',
  'garden-linux',
  'cflinuxfs2-rootfs',
  'routing',
  'cf',
]

deployment_configuration_path = ENV.fetch('DEPLOYMENT_CONFIGURATION_PATH')
deployment_manifest_path = ENV.fetch("DEPLOYMENT_MANIFEST_PATH")

deployment_manifest = YAML.load_file("deployment-configuration/#{deployment_configuration_path}")

release_names.each do |release_name|
  release_resource = "#{release_name}-release"

  url = File.read("#{release_resource}/url").strip
  version = File.read("#{release_resource}/version").strip
  sha1 = File.read("#{release_resource}/sha1").strip

  deployment_manifest.deep_merge!('releases' => {
    release_resource =>
    {
      'name' => release_name,
      'url' => url,
      'version' => version,
      'sha1' => sha1
    }
  })
end

stemcell_version = File.read("stemcell/version").strip

deployment_manifest.deep_merge!('stemcells' => [
  {
    'alias' => "default",
    'os' => "ubuntu-trusty",
    'version' => stemcell_version
  }
])

File.open("deployment-manifest/#{deployment_manifest_path}", 'w') do |file|
  YAML.dump(deployment_manifest, file)
end
